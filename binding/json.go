// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package binding

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type jsonBinding struct {
}

func (jsonBinding) Name() string {
	return "json"
}

func (jsonBinding) Bind(r *http.Request, v interface{}) error {
	return decodeJSON(r.Body, v)
}

func (jsonBinding) BindBody(b []byte, v interface{}) error {
	return decodeJSON(bytes.NewBuffer(b), v)
}

func decodeJSON(r io.Reader, v interface{}) error {
	if err := json.NewDecoder(r).Decode(v); err != nil {
		return err
	}

	//todo validate json data

	return nil
}
