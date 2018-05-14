// Copyright 2018 The Periph Authors. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package procextra

import "testing"

func TestSetHighPriority(t *testing.T) {
	// Ignore the result, do not do anything else. It may fail if not running as
	// root, succeed but in this case we want it to terminate ASAP.
	_ = SetHighPriority()
}
