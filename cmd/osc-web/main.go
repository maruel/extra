// Copyright 2017 The Periph Authors. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// oscilloscope runs a web based oscilloscope!
package main

import (
	"errors"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"periph.io/x/periph"
	"periph.io/x/periph/host"
)

const cacheControl30d = "Cache-Control:public, max-age=259200" // 30d
const cacheControl5m = "Cache-Control:public, max-age=300"     // 5m

var rootTmpl = `<!DOCTYPE html>
<html>
<head>
	<meta charset="utf-8" />
	<title>{{.Hostname}}</title>
</head>
<body>
{{.State}}
<br>
</body>
</html>`

type webServer struct {
	ln       net.Listener
	server   http.Server
	state    *periph.State
	hostname string
	rootTmpl *template.Template
}

func (s *webServer) rootHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Ugh", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("Cache-Control", cacheControl5m)
	keys := map[string]interface{}{
		"Hostname": s.hostname,
		"State":    s.state,
	}
	if err := s.rootTmpl.Execute(w, keys); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *webServer) faviconHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Ugh", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "image/png")
	//w.Header().Set("Cache-Control", cacheControl30d)
	// TODO(maruel): Add nice icon.
	w.Write([]byte{})
}

func (s *webServer) Close() error {
	return s.ln.Close()
}

func newWebServer(port string, state *periph.State) (*webServer, error) {
	s := &webServer{state: state}
	var err error
	if s.hostname, err = os.Hostname(); err != nil {
		return nil, err
	}
	if s.rootTmpl, err = template.New("name").Parse(rootTmpl); err != nil {
		return nil, err
	}
	http.HandleFunc("/", s.rootHandler)
	http.HandleFunc("/favicon.ico", s.faviconHandler)
	s.ln, err = net.Listen("tcp", port)
	if err != nil {
		return nil, err
	}
	s.server = http.Server{
		Addr:           s.ln.Addr().String(),
		Handler:        loggingHandler{http.DefaultServeMux},
		ReadTimeout:    60 * time.Second,
		WriteTimeout:   60 * time.Second,
		MaxHeaderBytes: 1 << 16,
	}
	go s.server.Serve(s.ln)
	return s, nil
}

type loggingHandler struct {
	handler http.Handler
}

type loggingResponseWriter struct {
	http.ResponseWriter
	length int
	status int
}

func (l *loggingResponseWriter) Write(data []byte) (size int, err error) {
	size, err = l.ResponseWriter.Write(data)
	l.length += size
	return
}

func (l *loggingResponseWriter) WriteHeader(status int) {
	l.ResponseWriter.WriteHeader(status)
	l.status = status
}

func (l loggingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	lrw := &loggingResponseWriter{ResponseWriter: w}
	l.handler.ServeHTTP(lrw, r)
	log.Printf("%s - %3d %6db %4s %s\n", r.RemoteAddr, lrw.status, lrw.length, r.Method, r.RequestURI)
}

func mainImpl() error {
	port := flag.String("p", "127.0.0.1:6060", "IP and port to bind to")
	verbose := flag.Bool("v", false, "verbose logging")
	flag.Parse()
	if flag.NArg() != 0 {
		return errors.New("unsupported arguments")
	}
	if !*verbose {
		log.SetOutput(ioutil.Discard)
	}
	state, err := host.Init()
	if err != nil {
		return err
	}
	s, err := newWebServer(*port, state)
	if err != nil {
		return err
	}
	select {}
	s.Close()
	return nil
}

func main() {
	if err := mainImpl(); err != nil {
		fmt.Fprintf(os.Stderr, "oscilloscope: %s.\n", err)
		os.Exit(1)
	}
}
