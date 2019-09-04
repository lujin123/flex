package flex

import (
	"fmt"
	"testing"
)

func detail(ctx *Context) error {
	fmt.Println("post detail...")
	return nil
}

func TestFlex(t *testing.T) {
	router := newRouter()
	v1 := router.Group("/v1")
	{
		v1.Get("/post/{id:\\d+}", detail)
	}
}
