// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !cgo !linux,!darwin
// +build !windows

package plugin

import "unsafe"

type libHandle int

func dlopen(name string) (libHandle, error) {
	panic("not implemented")
}

func (lib libHandle) dlclose() error {
	panic("not implemented")
}

func (lib libHandle) dlsym(name string) (unsafe.Pointer, error) {
	panic("not implemented")
}

func (lib libHandle) lookupC(name string, valptr interface{}) error {
	panic("not implemented")
}
