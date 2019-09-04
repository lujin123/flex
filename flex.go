package flex

import (
	"log"
	"net/http"
	"sync"
)

type (
	HandlerFunc    func(ctx *Context) error
	MiddlewareFunc func(h HandlerFunc) HandlerFunc
	HttpErrHandler func(ctx *Context, err error)

	Param map[string][]byte
	M     map[string]interface{}
)

type Flex struct {
	*Router

	errHandler      HttpErrHandler
	notFoundHandler HandlerFunc
	pool            sync.Pool
}

func New() *Flex {
	flex := &Flex{
		Router:          newRouter(),
		errHandler:      defaultErrHandler,
		notFoundHandler: defaultNotFoundHandler,
	}

	flex.pool.New = func() interface{} {
		return flex.allocateContext()
	}

	return flex
}

func (flex *Flex) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	ctx := flex.pool.Get().(*Context)
	ctx.reset(req, rw)

	flex.handleConn(ctx)
	flex.pool.Put(ctx)
}

func (flex *Flex) handleConn(ctx *Context) {
	h := flex.notFoundHandler
	if nd, err := flex.findRouter(ctx.Method(), ctx.Path()); err == nil {
		h = nd.handler
		ctx.params = nd.params
	}

	if err := h(ctx); err != nil {
		flex.errHandler(ctx, err)
	}
}

func (flex *Flex) allocateContext() *Context {
	return &Context{
		flex: flex,
		Resp: NewResponse(nil),
	}
}

func defaultErrHandler(ctx *Context, err error) {
	log.Printf("error: %+v", err)
	ctx.Write(400, []byte(err.Error()))
}

func defaultNotFoundHandler(ctx *Context) error {
	return ctx.Write(404, []byte("page not found"))
}
