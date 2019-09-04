package flex

import (
	"errors"
	"net/http"
)

var allowMethods = []string{
	http.MethodGet,
	http.MethodPost,
	http.MethodPut,
	http.MethodDelete,
	http.MethodHead,
	http.MethodPatch,
	http.MethodOptions,
}

type Router struct {
	basePath   string
	middleware []MiddlewareFunc
	trees      map[string]*tree
}

func newRouter() *Router {
	trees := make(map[string]*tree, len(allowMethods))

	for _, method := range allowMethods {
		trees[method] = newTree()
	}
	return &Router{
		basePath:   "/",
		middleware: nil,
		trees:      trees,
	}
}

func (r *Router) Use(middleware ...MiddlewareFunc) *Router {
	r.middleware = append(r.middleware, middleware...)
	return r
}

func (r *Router) Group(relativePath string, middleware ...MiddlewareFunc) *Router {
	return &Router{
		basePath:   r.calcAbsPath(relativePath),
		middleware: r.combineMiddleware(middleware),
		trees:      r.trees,
	}
}

func (r *Router) Get(relativePath string, handler HandlerFunc, middleware ...MiddlewareFunc) *Router {
	return r.handle(http.MethodGet, relativePath, handler, middleware)
}

func (r *Router) Post(relativePath string, handler HandlerFunc, middleware ...MiddlewareFunc) *Router {
	return r.handle(http.MethodPost, relativePath, handler, middleware)
}

func (r *Router) Put(relativePath string, handler HandlerFunc, middleware ...MiddlewareFunc) *Router {
	return r.handle(http.MethodPut, relativePath, handler, middleware)
}

func (r *Router) Delete(relativePath string, handler HandlerFunc, middleware ...MiddlewareFunc) *Router {
	return r.handle(http.MethodDelete, relativePath, handler, middleware)
}

func (r *Router) Patch(relativePath string, handler HandlerFunc, middleware ...MiddlewareFunc) *Router {
	return r.handle(http.MethodPatch, relativePath, handler, middleware)
}

func (r *Router) handle(httpMethod, relativePath string, handler HandlerFunc, middleware []MiddlewareFunc) *Router {
	absPath := r.calcAbsPath(relativePath)
	middleware = r.combineMiddleware(middleware)

	r.trees[httpMethod].insert(absPath, func(ctx *Context) error {
		h := applyMiddleware(handler, middleware...)
		return h(ctx)
	})
	return r
}

func (r *Router) calcAbsPath(relativePath string) string {
	return joinPaths(r.basePath, relativePath)
}

func (r *Router) combineMiddleware(middleware []MiddlewareFunc) []MiddlewareFunc {
	size := len(r.middleware) + len(middleware)
	mergeMiddleware := make([]MiddlewareFunc, size)
	copy(mergeMiddleware, r.middleware)
	copy(mergeMiddleware[len(r.middleware):], middleware)
	return mergeMiddleware
}

func (r *Router) findRouter(method, path string) (*node, error) {
	tree, ok := r.trees[method]
	if !ok {
		return nil, errors.New(method + " method not allowed")
	}
	n := tree.find(path)
	if n == nil || n.handler == nil {
		return nil, errors.New("path not found")
	}
	return n, nil
}
