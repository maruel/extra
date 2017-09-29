// Copyright 2017 The Periph Authors. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// Handheld IR camera!
package main

import (
	"errors"
	"flag"
	"fmt"
	"image"
	"log"
	"os"

	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"

	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/gpio/gpioreg"
	"periph.io/x/periph/conn/i2c/i2creg"
	"periph.io/x/periph/conn/spi/spireg"
	"periph.io/x/periph/devices/lepton"
	"periph.io/x/periph/host"
)

/*
func query(dev *lepton.Dev) error {
	status, err := dev.GetStatus()
	if err != nil {
		return err
	}
	fmt.Printf("Status.CameraStatus: %s\n", status.CameraStatus)
	fmt.Printf("Status.CommandCount: %d\n", status.CommandCount)
	serial, err := dev.GetSerial()
	if err != nil {
		return err
	}
	fmt.Printf("Serial:              0x%x\n", serial)
	uptime, err := dev.GetUptime()
	if err != nil {
		return err
	}
	fmt.Printf("Uptime:              %s\n", uptime)
	temp, err := dev.GetTemp()
	if err != nil {
		return err
	}
	fmt.Printf("Temp:         %s\n", temp)
	temp, err = dev.GetTempHousing()
	if err != nil {
		return err
	}
	fmt.Printf("Temp housing: %s\n", temp)
	pos, err := dev.GetShutterPos()
	if err != nil {
		return err
	}
	fmt.Printf("ShutterPos:     %s\n", pos)
	mode, err := dev.GetFFCModeControl()
	if err != nil {
		return err
	}
	fmt.Printf("FCCMode.FFCShutterMode:          %s\n", mode.FFCShutterMode)
	fmt.Printf("FCCMode.ShutterTempLockoutState: %s\n", mode.ShutterTempLockoutState)
	fmt.Printf("FCCMode.VideoFreezeDuringFFC:    %t\n", mode.VideoFreezeDuringFFC)
	fmt.Printf("FCCMode.FFCDesired:              %t\n", mode.FFCDesired)
	fmt.Printf("FCCMode.ElapsedTimeSinceLastFFC: %s\n", mode.ElapsedTimeSinceLastFFC)
	fmt.Printf("FCCMode.DesiredFFCPeriod:        %s\n", mode.DesiredFFCPeriod)
	fmt.Printf("FCCMode.ExplicitCommandToOpen:   %t\n", mode.ExplicitCommandToOpen)
	fmt.Printf("FCCMode.DesiredFFCTempDelta:     %s\n", mode.DesiredFFCTempDelta)
	fmt.Printf("FCCMode.ImminentDelay:           %d\n", mode.ImminentDelay)
	return nil
}
*/

func captureLoop(dev *lepton.Dev, b screen.Buffer, t screen.Texture) error {
	for {
		_, err := dev.ReadImg()
		if err != nil {
			log.Fatal(err)
			return err
		}
	}
}

func mainUI(s screen.Screen, dev *lepton.Dev) {
	w, err := s.NewWindow(&screen.NewWindowOptions{Title: "PocketLepton"})
	if err != nil {
		log.Fatal(err)
	}
	defer w.Release()

	b, err := s.NewBuffer(dev.Bounds().Size())
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer b.Release()

	tImg, err := s.NewTexture(dev.Bounds().Size())
	if err != nil {
		log.Fatal(err)
	}
	defer tImg.Release()
	tImg.Upload(image.Point{}, b, b.Bounds())
	go captureLoop(dev, b, tImg)
	mainUILoop(w)
}

func mainUILoop(w screen.Window) {
	var sz size.Event
	for {
		switch e := w.NextEvent().(type) {
		case lifecycle.Event:
			if e.To == lifecycle.StageDead {
				return
			}
		case key.Event:
			if e.Code == key.CodeEscape {
				return
			}
		case paint.Event:
			paintFrame(sz)
			w.Publish()
		case size.Event:
			sz = e
		case error:
			log.Print(e)
		}
	}
}

func paintFrame(sz size.Event) {
	/*
		if *meta {
			fmt.Printf("SinceStartup:       %s\n", frame.Metadata.SinceStartup)
			fmt.Printf("FrameCount:         %d\n", frame.Metadata.FrameCount)
			fmt.Printf("Temp:        %s\n", frame.Metadata.Temp)
			fmt.Printf("TempHousing: %s\n", frame.Metadata.TempHousing)
			fmt.Printf("FFCSince:           %s\n", frame.Metadata.FFCSince)
			fmt.Printf("FFCDesired:         %t\n", frame.Metadata.FFCDesired)
			fmt.Printf("Overtemp:           %t\n", frame.Metadata.Overtemp)
		}
	*/
	/*
		const inset = 10
		for _, r := range imageutil.Border(sz.Bounds(), inset) {
			w.Fill(r, blue0, screen.Src)
		}
		w.Fill(sz.Bounds().Inset(inset), blue1, screen.Src)
		w.Upload(image.Point{20, 0}, b, b.Bounds())
		w.Fill(image.Rect(50, 50, 350, 120), red, screen.Over)

		// By default, draw the entirety of the texture using the Over
		// operator. Uncomment one or both of the lines below to see
		// their different effects.
		op := screen.Over
		// op = screen.Src
		t0Rect := t0.Bounds()
		// t0Rect = image.Rect(16, 0, 240, 100)

		// Draw the texture t0 twice, as a 1:1 copy and under the
		// transform src2dst.
		w.Copy(image.Point{150, 100}, t0, t0Rect, op, nil)
		src2dst := f64.Aff3{
			+0.5 * cos30, -1.0 * sin30, 100,
			+0.5 * sin30, +1.0 * cos30, 200,
		}
		w.Draw(src2dst, t0, t0Rect, op, nil)
		w.DrawUniform(src2dst, yellow, t0Rect.Inset(30), screen.Over, nil)

		// Draw crosses at the transformed corners of t0Rect.
		for _, sx := range []int{t0Rect.Min.X, t0Rect.Max.X} {
			for _, sy := range []int{t0Rect.Min.Y, t0Rect.Max.Y} {
				dx := int(src2dst[0]*float64(sx) + src2dst[1]*float64(sy) + src2dst[2])
				dy := int(src2dst[3]*float64(sx) + src2dst[4]*float64(sy) + src2dst[5])
				w.Fill(image.Rect(dx-0, dy-1, dx+1, dy+2), darkGray, screen.Src)
				w.Fill(image.Rect(dx-1, dy-0, dx+2, dy+1), darkGray, screen.Src)
			}
		}

		// Draw t1.
		w.Copy(image.Point{400, 50}, t1, t1.Bounds(), screen.Src, nil)
	*/
}

func initialize(spiName, csName string, spiHz int, i2cName string, i2cHz int) (*lepton.Dev, error) {
	if _, err := host.Init(); err != nil {
		return nil, err
	}
	spiBus, err := spireg.Open(spiName)
	if err != nil {
		return nil, err
	}
	if spiHz != 0 {
		if err := spiBus.LimitSpeed(int64(spiHz)); err != nil {
			return nil, err
		}
	}
	var cs gpio.PinOut
	if len(csName) != 0 {
		if p := gpioreg.ByName(csName); p != nil {
			cs = p
		} else {
			return nil, fmt.Errorf("%s is not a valid pin", csName)
		}
	}

	i2cBus, err := i2creg.Open(i2cName)
	if err != nil {
		return nil, err
	}
	if i2cHz != 0 {
		if err := i2cBus.SetSpeed(int64(i2cHz)); err != nil {
			return nil, err
		}
	}
	return lepton.New(spiBus, i2cBus, cs)
}

func mainImpl() error {
	i2cName := flag.String("i2c", "", "I²C bus to use")
	spiName := flag.String("spi", "", "SPI bus to use")
	csName := flag.String("cs", "", "SPI CS line to use instead of the default")
	i2cHz := flag.Int("i2chz", 0, "I²C bus speed")
	spiHz := flag.Int("spihz", 0, "SPI bus speed")
	flag.Parse()
	if flag.NArg() != 0 {
		return errors.New("unexpected flags")
	}
	dev, err := initialize(*spiName, *csName, *spiHz, *i2cName, *i2cHz)
	if err != nil {
		return err
	}
	driver.Main(func(s screen.Screen) { mainUI(s, dev) })
	return nil
}

func main() {
	if err := mainImpl(); err != nil {
		fmt.Fprintf(os.Stderr, "\npocketlepton: %s.\n", err)
		os.Exit(1)
	}
}
