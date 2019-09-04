package flex

import (
	"bufio"
	"net"
	"net/http"
)

type Response struct {
	rw         http.ResponseWriter
	statusCode int
	commit     bool
}

func NewResponse(w http.ResponseWriter) *Response {
	return &Response{rw: w}
}

func (r *Response) Header() http.Header {
	return r.rw.Header()
}

func (r *Response) WriteHeader(code int) {
	r.statusCode = code
	r.rw.WriteHeader(code)
	r.commit = true
}

func (r *Response) Write(b []byte) (n int, err error) {
	if !r.commit {
		if r.statusCode == 0 {
			r.statusCode = http.StatusOK
		}
		r.WriteHeader(r.statusCode)
	}
	n, err = r.rw.Write(b)
	return
}

// Flush implements the http.Flusher interface to allow an HTTP handler to flush
// buffered data to the client.
// See [http.Flusher](https://golang.org/pkg/net/http/#Flusher)
func (r *Response) Flush() {
	r.rw.(http.Flusher).Flush()
}

// Hijack implements the http.Hijacker interface to allow an HTTP handler to
// take over the connection.
// See [http.Hijacker](https://golang.org/pkg/net/http/#Hijacker)
func (r *Response) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return r.rw.(http.Hijacker).Hijack()
}

func (r *Response) reset(w http.ResponseWriter) {
	r.rw = w
	r.statusCode = http.StatusOK
	r.commit = false
}
