// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package binding

import "net/http"

type Binding interface {
	Name() string
	Bind(r *http.Request, v interface{}) error
}

type BindingBody interface {
	Binding
	BindBody(b []byte, v interface{}) error
}

var (
	JSON          = jsonBinding{}
	XML           = xmlBinding{}
	Query         = queryBinding{}
	Form          = formBinding{}
	FormPost      = formPostBinding{}
	FormMultipart = formMultipartBinding{}
)
