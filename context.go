package flex

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/lujin123/flex/binding"
)

type Context struct {
	Req  *http.Request
	Resp *Response
	flex *Flex

	params Param
}

func (ctx *Context) Path() string {
	path := ctx.Req.URL.RawPath
	if path == "" {
		path = ctx.Req.URL.Path
	}
	return path
}

func (ctx *Context) Method() string {
	return ctx.Req.Method
}

func (ctx *Context) ContentType() string {
	return filterHeader(ctx.reqHeader(HeaderContentType))
}

func (ctx *Context) Params() Param {
	return ctx.params
}

func (ctx *Context) Param(key string) string {
	v, ok := ctx.params[key]
	if ok {
		return string(v)
	}
	return ""
}

func (ctx *Context) SetCookie(name, value string, maxAge int, path, domain string, secure, httpOnly bool) {
	if path == "" {
		path = "/"
	}
	http.SetCookie(ctx.Resp, &http.Cookie{
		Name:     name,
		Value:    url.QueryEscape(value),
		MaxAge:   maxAge,
		Path:     path,
		Domain:   domain,
		Secure:   secure,
		HttpOnly: httpOnly,
	})
}

func (ctx *Context) Write(code int, b []byte) error {
	ctx.Resp.WriteHeader(code)
	_, err := ctx.Resp.Write(b)
	return err
}

func (ctx *Context) JSON(code int, data interface{}) error {
	ctx.writeContentType(MIMEApplicationJSONCharsetUTF8)
	ctx.Resp.WriteHeader(code)
	return json.NewEncoder(ctx.Resp).Encode(data)
}

///////////////////////////binding request methods///////////////////////////
func (ctx *Context) ShouldBind(ptr interface{}) error {
	return ctx.ShouldBindWith(ptr, ctx.binding())
}

func (ctx *Context) ShouldBindWith(ptr interface{}, b binding.Binding) error {
	return b.Bind(ctx.Req, ptr)
}

func (ctx *Context) ShouldBindJSON(ptr interface{}) error {
	return ctx.ShouldBindWith(ptr, binding.JSON)
}

func (ctx *Context) ShouldBindQuery(ptr interface{}) error {
	return ctx.ShouldBindWith(ptr, binding.Query)
}

func (ctx *Context) ShouldBindXML(ptr interface{}) error {
	return ctx.ShouldBindWith(ptr, binding.XML)
}

//////////////////////////private methods////////////////////////////
func (ctx *Context) writeContentType(value string) {
	header := ctx.Resp.Header()
	header.Set(HeaderContentType, value)
}

func (ctx *Context) reqHeader(key string) string {
	return ctx.Req.Header.Get(key)
}

func (ctx *Context) reset(r *http.Request, w http.ResponseWriter) {
	ctx.Req = r
	ctx.Resp.reset(w)

	ctx.params = nil
}

func (ctx *Context) binding() binding.Binding {
	if ctx.Method() == http.MethodGet {
		return binding.Form
	}

	contentType := ctx.ContentType()

	switch contentType {
	case MIMEApplicationJSON:
		return binding.JSON
	case MIMEApplicationXML, MIMETextXML:
		return binding.XML
	case MIMEApplicationForm:
		return binding.FormPost
	case MIMEMultipartForm:
		return binding.FormMultipart
	default:
		return binding.Form
	}
}
