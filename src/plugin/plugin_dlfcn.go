// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build linux darwin

package plugin

// #include <stdlib.h>
// #include <string.h>
//
// #include <dlfcn.h>
// #cgo LDFLAGS: -ldl
import "C"
import (
	"fmt"
	"reflect"
	"unsafe"
)

type libHandle struct {
	addr unsafe.Pointer
}

func dlopen(name string) (libHandle, error) {
	cstr := C.CString(name)
	defer C.free(unsafe.Pointer(cstr))
	h := C.dlopen(cstr, C.RTLD_NOW|C.RTLD_GLOBAL)
	if h == nil {
		return libHandle{}, dlerror()
	}
	return libHandle{h}, nil
}

func (lib libHandle) dlclose() error {
	o := C.dlclose(lib.addr)
	if o != 0 {
		return dlerror()
	}
	return nil
}

func dlerror() error {
	cmsg := C.dlerror()
	msg := C.GoString(cmsg)
	return fmt.Errorf("plugin: %s", msg)
}

func (lib libHandle) dlsym(name string) (unsafe.Pointer, error) {
	sym := C.CString(name)
	defer C.free(unsafe.Pointer(sym))

	addr := C.dlsym(lib.addr, sym)
	if addr == nil {
		return nil, dlerror()
	}
	return addr, nil
}

func (lib libHandle) lookupC(name string, valptr interface{}) error {
	addr, err := lib.dlsym(name)
	if err != nil {
		return err
	}
	if addr == nil {
		return fmt.Errorf("plugin: nil pointer to symbol %q", name)
	}

	rv := reflect.ValueOf(valptr)
	if !rv.IsValid() || rv.Kind() != reflect.Ptr {
		return errNotPtr
	}

	if rv.IsNil() {
		return errNilPtr
	}

	val := rv.Elem()
	switch val.Kind() {
	case reflect.Int:
		val.SetInt(int64(*(*int)(addr)))
	case reflect.Int8:
		val.SetInt(int64(*(*int8)(addr)))
	case reflect.Int16:
		val.SetInt(int64(*(*int16)(addr)))
	case reflect.Int32:
		val.SetInt(int64(*(*int32)(addr)))
	case reflect.Int64:
		val.SetInt(int64(*(*int64)(addr)))
	case reflect.Uint:
		val.SetUint(uint64(*(*uint)(addr)))
	case reflect.Uint8:
		val.SetUint(uint64(*(*uint8)(addr)))
	case reflect.Uint16:
		val.SetUint(uint64(*(*uint16)(addr)))
	case reflect.Uint32:
		val.SetUint(uint64(*(*uint32)(addr)))
	case reflect.Uint64:
		val.SetUint(uint64(*(*uint64)(addr)))
	case reflect.Uintptr:
		val.SetUint(uint64(*(*uintptr)(addr)))
	case reflect.Float32:
		val.SetFloat(float64(*(*float32)(addr)))
	case reflect.Float64:
		val.SetFloat(float64(*(*float64)(addr)))
	case reflect.Complex64:
		val.SetComplex(complex128(*(*complex64)(addr)))
	case reflect.Complex128:
		val.SetComplex(complex128(*(*complex128)(addr)))
	case reflect.String:
		val.SetString(C.GoString(*(**C.char)(addr)))
	case reflect.UnsafePointer:
		val.SetPointer(addr)
	default:
		return fmt.Errorf("plugin: invalid type %T", valptr)
	}

	return nil
}
