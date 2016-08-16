// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package plugin provides programmatic access to symbols in shared libraries.
//
// More detailed informations available at:
//  https://docs.google.com/document/d/1nr-TQHw_er6GOQRsF6T43GGhFDelrAP0NqSS_00RgZQ/edit#
package plugin

// TODO: implement stubs for Windows and others
// TODO: implement LookupC for *C.funcs

import (
	"errors"
)

var (
	errNilPtr = errors.New("plugin: nil pointer")
	errNotPtr = errors.New("plugin: expected a pointer to a value")
)

// Plugin is an opaque handle to a plugin library opened at runtime.
type Plugin struct {
	handle libHandle
}

// Open opens a plugin by name.
func Open(name string) (Plugin, error) {
	lib, err := dlopen(name)
	if err != nil {
		return Plugin{}, err
	}
	return Plugin{lib}, nil
}

// Close closes the plugin.
func (p Plugin) Close() error {
	return p.handle.dlclose()
}

// LookupC looks up a symbol in a C style plugin, passing in a pointer to
// a value with the type it is expected to have.
// The value must be a C variable type.
//
// TODO: implement LookupC for function types with a C style API.
func (p Plugin) LookupC(name string, valptr interface{}) error {
	return p.handle.lookupC(name, valptr)
}
