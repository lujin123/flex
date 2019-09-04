// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package binding

import "net/http"

const defaultMemory = 32 * 1024 * 1024

type formBinding struct {
}

func (formBinding) Name() string {
	return "form"
}

func (formBinding) Bind(r *http.Request, v interface{}) error {
	if err := r.ParseForm(); err != nil {
		return err
	}

	if err := r.ParseMultipartForm(defaultMemory); err != nil {
		return err
	}

	if err := mapForm(v, r.Form); err != nil {
		return err
	}
	return nil
}

type formPostBinding struct {
}

func (formPostBinding) Name() string {
	return "form-urlencoded"
}

func (formPostBinding) Bind(r *http.Request, v interface{}) error {
	if err := r.ParseForm(); err != nil {
		return err
	}

	if err := mapForm(v, r.PostForm); err != nil {
		return err
	}
	return nil
}

type formMultipartBinding struct {
}

func (formMultipartBinding) Name() string {
	return "multipart/form-data"
}

func (formMultipartBinding) Bind(r *http.Request, v interface{}) error {
	if err := r.ParseMultipartForm(defaultMemory); err != nil {
		return err
	}

	if err := mapForm(v, r.MultipartForm.Value); err != nil {
		return err
	}
	return nil
}
