package flex

import (
	"fmt"
	"testing"
)

func log(h HandlerFunc) HandlerFunc {
	return func(ctx *Context) error {
		fmt.Println("log middleware before...")
		err := h(ctx)
		fmt.Println("log middleware after...")
		return err
	}
}
func hello(ctx *Context) error {
	fmt.Println("hello handler...")
	return nil
}

func TestTreeInsert(t *testing.T) {
	tree := newTree()
	tree.insert("/post/{id:\\d+}/hello", hello)
	t.Log(tree)
}

func TestTreeFind(t *testing.T) {
	tree := newTree()
	tree.insert("/post/{id:\\w+}/hello", hello)
	tree.insert("/post/{id:\\w+}/comment", hello)
	tree.insert("/article/{id:\\d+}", hello)
	n := tree.find("/post/abc/hello")
	t.Log(string(n.params["id"]) == "abc")
	n2 := tree.find("/article/123")
	t.Log(string(n2.params["id"]) == "123")
}
