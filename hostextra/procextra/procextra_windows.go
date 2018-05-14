// Copyright 2018 The Periph Authors. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package procextra

import (
	"fmt"
	"syscall"

	"golang.org/x/sys/windows"
)

func setHighPriority() error {
	dll, err := syscall.LoadDLL("kernel32.dll")
	if err != nil {
		return err
	}
	setPriorityClass, err := dll.FindProc("SetPriorityClass")
	if err != nil {
		return err
	}
	setThreadPriority, err := dll.FindProc("SetThreadPriority")
	if err != nil {
		return err
	}

	// https://msdn.microsoft.com/en-us/library/windows/desktop/ms685100.aspx
	if r1, _, _ := setPriorityClass.Call(0xFFFFFFFF, 0x100); r1 == 0 {
		return fmt.Errorf("procextra: failed to set process priority to REALTIME_PRIORITY_CLASS: %v", windows.GetLastError())
	}

	if r1, _, _ := setThreadPriority.Call(0xFFFFFFFF, 15); r1 == 0 {
		return fmt.Errorf("procextra: failed to set thread priority to THREAD_PRIORITY_TIME_CRITICAL: %v", windows.GetLastError())
	}
	return nil
}
