// Copyright 2018 The Periph Authors. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package procextra_test

import (
	"log"
	"runtime"
	"runtime/debug"

	"periph.io/x/extra/hostextra/procextra"
)

func ExampleSetHighPriority() {
	// GC one last them and then disable GC.
	runtime.GC()
	debug.SetGCPercent(-1)

	// Disable the Go runtime scheduler for this goroutine.
	runtime.LockOSThread()

	if err := procextra.SetHighPriority(); err != nil {
		log.Fatal(err)
	}

	// Start CPU intensive work that do not do memory allocation.
}
