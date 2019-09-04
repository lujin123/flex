// Copyright 2017 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package binding

import "net/http"

type queryBinding struct {
}

func (queryBinding) Name() string {
	return "query"
}

func (queryBinding) Bind(r *http.Request, v interface{}) error {
	values := r.URL.Query()

	if err := mapForm(v, values); err != nil {
		return err
	}

	//todo validate

	return nil
}
