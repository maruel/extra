// Copyright 2018 The Periph Authors. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package procextra

/*
#include <pthread.h>
*/
import "C"
import (
	"fmt"
)

func setHighPriority() error {
	// https://developer.apple.com/library/content/documentation/Darwin/Conceptual/KernelProgramming/scheduler/scheduler.html
	sp := C.struct_sched_param{sched_priority: -1}
	if ret := C.pthread_setschedparam(C.pthread_self(), C.SCHED_FIFO, &sp); ret != 0 {
		return fmt.Errorf("procextra: failed to set process priority to SCHED_FIFO: %d", ret)
	}
	return nil
}
