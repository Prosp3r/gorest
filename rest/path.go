// Tideland Go REST Server Library - REST - Path
//
// Copyright (C) 2009-2017 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package rest

//--------------------
// IMPORTS
//--------------------

import (
	"net/http"
	"strings"

	"github.com/tideland/golib/stringex"
)

//--------------------
// CONSTANTS
//--------------------

// Path indexes for the different parts.
const (
	PathDomain     = 0
	PathResource   = 1
	PathResourceID = 2
)

//--------------------
// PATH
//--------------------

// Path provides access to the parts of a
// request path interesting for handling a
// job.
type Path interface {
	// Length returns the number of parts of the path.
	Length() int

	// Part returns the parts of the URL path based on the
	// index or an empty string.
	Part(index int) string

	// Domain returns the requests domain.
	Domain() string

	// Resource returns the requests resource.
	Resource() string

	// ResourceID return the requests resource ID.
	ResourceID() string
}

// path implements Path.
type path struct {
	parts []string
}

// newPath returns the analyzed path.
func newPath(env *environment, r *http.Request) *path {
	parts := stringex.SplitMap(r.URL.Path, "/", func(part string) (string, bool) {
		if part == "" {
			return "", false
		}
		return part, true
	})[env.basepartsLen:]
	switch len(parts) {
	case 1:
		parts = append(parts, env.defaultResource)
	case 0:
		parts = append(parts, env.defaultDomain, env.defaultResource)
	}
	return &path{
		parts: parts,
	}
}

// Length implements Path.
func (p *path) Length() int {
	return len(p.parts)
}

// Part implements Path.
func (p *path) Part(index int) string {
	if len(p.parts) <= index {
		return ""
	}
	return p.parts[index]
}

// Domain implements Path.
func (p *path) Domain() string {
	return p.parts[PathDomain]
}

// Resource implements Path.
func (p *path) Resource() string {
	return p.parts[PathResource]
}

// ResourceID implements Path.
func (p *path) ResourceID() string {
	if len(p.parts) > 2 {
		return strings.Join(p.parts[PathResourceID:], "/")
	}
	return ""
}

// EOF
