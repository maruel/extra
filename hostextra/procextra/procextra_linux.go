// Copyright 2018 The Periph Authors. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package procextra

/*
#include <sched.h>
*/
import "C"
import (
	"fmt"
)

func setHighPriority() error {
	sp := C.struct_sched_param{__sched_priority: -1}
	if ret := C.sched_setscheduler(0, C.SCHED_FIFO, &sp); ret != 0 {
		return fmt.Errorf("procextra: failed to set process priority to SCHED_FIFO: %d", ret)
	}
	return nil
}
