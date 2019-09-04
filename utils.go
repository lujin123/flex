package flex

import (
	"path"
	"strings"
)

func joinPaths(absPath, relativePath string) string {
	if relativePath == "" {
		return absPath
	}

	finalPath := path.Join(absPath, relativePath)
	n := len(finalPath)
	if n == 0 {
		panic("The length of the string can't be 0")
	}
	if finalPath[n-1] == '/' {
		finalPath = finalPath[:n-1]
	}
	return finalPath
}

func applyMiddleware(h HandlerFunc, middleware ...MiddlewareFunc) HandlerFunc {
	for i := len(middleware) - 1; i >= 0; i-- {
		h = middleware[i](h)
	}
	return h
}

func splitPath(path string) []string {
	return strings.Split(path, "/")
}
func trimPathPrefix(path string) string {
	return strings.TrimPrefix(path, "/")
}

func filterHeader(s string) string {
	for i, c := range s {
		if c == ' ' || c == ';' {
			return s[:i]
		}
	}
	return s
}
