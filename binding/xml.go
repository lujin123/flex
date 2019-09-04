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

type xmlBinding struct {
}

func (xmlBinding) Name() string {
	return "json"
}

func (xmlBinding) Bind(r *http.Request, v interface{}) error {
	return decodeXML(r.Body, v)
}

func (xmlBinding) BindBody(b []byte, v interface{}) error {
	return decodeXML(bytes.NewBuffer(b), v)
}

func decodeXML(r io.Reader, v interface{}) error {
	if err := json.NewDecoder(r).Decode(v); err != nil {
		return err
	}

	//todo validate xml data

	return nil
}
