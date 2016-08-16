// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package plugin

import (
	"fmt"
	"syscall"
)

type libHandle struct {
	dll *syscall.DLL
}

func dlopen(name string) (libHandle, error) {
	dll, err := syscall.LoadDLL(name)
	if err != nil {
		return libHandle{}, fmt.Errorf("plugin: %s", err)
	}
	return libHandle{dll}, nil
}

func (lib libHandle) dlclose() error {
	return lib.dll.Release()
}

func (lib libHandle) dlsym(name string) (*syscall.Proc, error) {
	proc, err := lib.dll.FindProc(name)
	if err != nil {
		return nil, fmt.Errorf("plugin: %s", err)
	}
	return proc, nil
}

func (lib libHandle) lookupC(name string, valptr interface{}) error {
	panic("not implemented")
}
