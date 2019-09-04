package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/lujin123/flex"
)

type Query struct {
	Limit  int64 `form:"limit,default=20"`
	Offset int64 `form:"offset,default=0"`
}

type PostJson struct {
	Id   uint64 `json:"id" form:"id"`
	Name string `json:"name" form:"name"`
	Age  uint64 `json:"age" form:"age"`
}

func logMiddleware(h flex.HandlerFunc) flex.HandlerFunc {
	return func(context *flex.Context) error {
		start := time.Now()
		log.Printf("middleware logs before time: %v\n", start)
		err := h(context)
		end := time.Now()
		log.Printf("middleware logs after time: %v\n", end)
		log.Printf("request latency: %dns\n", end.Sub(start).Nanoseconds())
		return err
	}
}

func detail(ctx *flex.Context) error {
	s := "detail ..."
	fmt.Println(s)
	return ctx.Write(200, []byte(s))
}

func comment(ctx *flex.Context) error {
	s := "comment ..."
	fmt.Println(ctx.Params())
	fmt.Println(s)
	return ctx.JSON(200, ctx.Params())
}

func list(ctx *flex.Context) error {
	var data Query
	err := ctx.ShouldBindQuery(&data)
	if err != nil {
		return ctx.JSON(400, err)
	}
	fmt.Printf("list data: %+v\n", data)
	s := "list ..."
	fmt.Println(s)
	return ctx.Write(200, []byte(s))
}

func main() {
	router := flex.New()
	v1 := router.Group("/v1", logMiddleware)
	{
		v1.Get("/post/{id:\\d+}", detail)
		//v1.Get("/post/{id:\\d+}/comment", comment)
		//v1.Get("/post/{:\\d+}/comment", comment)
		//v1.Get("/post/{id}/comment", comment)
		v1.Get("/post/:id/comment", comment)
		v1.Get("/post/{id:\\d+/hello}", func(ctx *flex.Context) error {
			return ctx.JSON(200, flex.M{
				"id":   ctx.Param("id"),
				"name": "hello",
			})
		})
		v1.Get("/post", list)
		v1.Post("/post", func(ctx *flex.Context) error {
			var data PostJson
			//if err := ctx.ShouldBindJSON(&data); err != nil {
			//	return ctx.JSON(400, err)
			//}
			if err := ctx.ShouldBind(&data); err != nil {
				return ctx.JSON(400, err)
			}
			fmt.Printf("post data: %+v\n", data)
			return ctx.JSON(200, data)
		})
	}

	srv := &http.Server{
		Addr:              ":12345",
		Handler:           router,
		TLSConfig:         nil,
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      5 * time.Second,
		IdleTimeout:       5 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}
