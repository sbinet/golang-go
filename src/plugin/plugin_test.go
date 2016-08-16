// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build cgo
// +build linux darwin

package plugin

import (
	"runtime"
	"testing"
	"unsafe"
)

func libc() string {
	switch runtime.GOOS {
	case "linux":
		return "libc.so.6"
	case "darwin":
		return "libc.dylib"
	}
	return "N/A"
}

func libm() string {
	switch runtime.GOOS {
	case "linux":
		return "libm.so.6"
	case "darwin":
		return "libm.dylib"
	}
	return "N/A"
}

func libNotThere() string {
	switch runtime.GOOS {
	case "linux":
		return "libc_NOT_THERE.so"
	case "darwin":
		return "libc_NOT_THERE.dylib"
	}
	return "N/A"
}

func TestOpen(t *testing.T) {
	for _, test := range []struct {
		lib string
		ok  bool
	}{
		{
			lib: libc(),
			ok:  true,
		},
		{
			lib: libm(),
			ok:  true,
		},
		{
			lib: libNotThere(),
			ok:  false,
		},
	} {
		p, err := Open(test.lib)
		if !test.ok {
			if err == nil {
				t.Errorf("%s: expected an error loading library\n", test.lib)
				p.Close()
				continue
			}
			continue
		}
		if err != nil {
			t.Error(err)
			continue
		}

		err = p.Close()
		if err != nil {
			t.Error(err)
			continue
		}
	}
}

func TestLookupC(t *testing.T) {
	for _, test := range []struct {
		lib string
		sym string
		ok  bool
	}{
		{
			lib: libc(),
			sym: "puts",
			ok:  true,
		},
		{
			lib: libc(),
			sym: "puts_NOT_THERE",
			ok:  false,
		},
		{
			lib: libm(),
			sym: "fabs",
			ok:  true,
		},
	} {
		p, err := Open(test.lib)
		if err != nil {
			t.Error(err)
			continue
		}
		defer p.Close()

		var val unsafe.Pointer
		err = p.LookupC(test.sym, &val)
		if !test.ok {
			if err == nil {
				t.Errorf("%s: expected an error loading symbol %q\n", test.lib, test.sym)
				continue
			}
			continue
		}
		if err != nil {
			t.Error(err)
			continue
		}

		err = p.Close()
		if err != nil {
			t.Error(err)
			continue
		}
	}
}
