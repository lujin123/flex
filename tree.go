package flex

import (
	"fmt"
	"regexp"
	"strings"
)

//
//todo 参数和正则现在还不支持
const (
	root byte = iota
	static
	named
	pattern
)

const (
	start     = '{'
	end       = '}'
	separator = ':'
)

const (
	defaultKey  = "key"
	defaultExpr = "[0-9a-zA-Z_.]+"
)

type node struct {
	key      string
	typ      byte
	expr     *regexp.Regexp //编译后的正则表达式
	regNode  *node
	children map[string]*node
	params   Param
	handler  HandlerFunc
}

type tree struct {
	root *node
}

func newTree() *tree {
	return &tree{root: &node{
		key:      "/",
		typ:      root,
		children: make(map[string]*node),
	}}
}

func (t *tree) insert(path string, handler HandlerFunc) {
	n := t.root
	if path != n.key {
		path = trimPathPrefix(path)
		keys := splitPath(path)
		for _, s := range keys {
			typ := parseNodeType(s)
			var (
				key string
				reg string
				tn  *node
				ok  bool
			)

			if typ == static {
				key = s
				tn, ok = n.children[key]
				if !ok {
					tn = &node{
						key:      key,
						typ:      static,
						children: make(map[string]*node),
					}
					n.children[key] = tn
				}
			} else {
				//如果是正则节点，那么这部分就只能有一个节点，否则就会冲突
				if len(n.children) > 0 {
					panic(fmt.Sprintf("<%s> conflict with static path", path))
				}
				ln := len(s)
				if typ == pattern {
					i := strings.IndexByte(s, separator)

					if i < 0 {
						// /post/{id}/hello => /post/{id:[0-9a-zA-Z_.]+}/hello
						key = s[1 : ln-1]
						reg = defaultExpr
					} else if i == 1 {
						// /post/{:\\d+}/hello => /post/{id:\\d+}/hello
						key = defaultKey
						reg = s[i+1 : ln-1]
					} else {
						// /post/{id:\\d+}/hello => /post/{id:\\d+}/hello
						key = s[1:i]
						reg = s[i+1 : ln-1]
					}
				} else {
					// /post/:id/hello => /post/{id:[0-9a-zA-Z_.]+}/hello
					key = s[1 : ln-1]
					reg = defaultExpr
				}
				tn = n.regNode
				if tn == nil {
					tn = &node{
						key:      key,
						typ:      typ,
						expr:     regexp.MustCompile(reg),
						children: make(map[string]*node),
					}
					n.regNode = tn
				}
			}
			n = tn
		}
	}

	if n.handler != nil {
		panic(fmt.Sprintf("<%s> route exists", path))
	}

	n.handler = handler
}

func (t *tree) find(path string) *node {
	n := t.root
	if path == n.key {
		return n
	}

	params := make(map[string][]byte)

	path = trimPathPrefix(path)
	ln := len(path)
	if path[ln-1] == '/' {
		path = path[:ln-1]
	}

	keys := splitPath(path)
	for _, key := range keys {
		var (
			child *node
			ok    bool
		)
		if n.regNode != nil {
			child = n.regNode
			val := child.expr.Find([]byte(key))
			//如果不匹配就是没有找到对应的路由
			if val == nil {
				return nil
			}
			params[child.key] = val
		} else {
			child, ok = n.children[key]
			if !ok {
				return nil
			}
		}
		n = child
	}
	n.params = params
	return n
}

func parseNodeType(s string) byte {
	n := len(s)
	if n == 0 {
		panic("route key empty")
	}
	if s[0] == start {
		if s[n-1] == end {
			return pattern
		} else {
			panic("router syntax error")
		}
	} else if s[0] == separator {
		return named
	}
	return static
}
