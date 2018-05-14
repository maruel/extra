// Copyright 2018 The Periph Authors. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// Package procextra contains code to mutate the process behavior with code that
// requires cgo.
package procextra

// SetHighPriority sets the process as high priority.
//
// It is implemented for Linux, macOS and Windows. This requires super user
// rights (root or administrator) to succeed.
//
// SystemD
//
// For process running as a systemd service, do not forget to set
// LimitRTPRIO=infinity.
func SetHighPriority() error {
	return setHighPriority()
}
