// Copyright 2018 The Periph Authors. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// +build !cgo
// +build !windows

package procextra

import "errors"

func setHighPriority() error {
	return errors.New("procextra: needs cgo")
}
