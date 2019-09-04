# flex
`golang`的路由组件，支持正则形式的路由，快不快不知道，没测过...

写这个库的原因是很不喜欢`gin`那种没有返回值的库，导致用起装饰器很烦，基本都是侵入性很强，
所以写一个带返回值的，这样就可以开心的用装饰器来做各种骚操作了...

原来是要实现`radix tree`的，但是在支持正则的情况下有点麻烦，没有实现好，
就换了一个取巧的方法实现了，利用`URL`都用`/`分隔的特点，直接分割开，一个做
一个节点，这样就简单很多了，越简单越少`Bug`呀...

## 安装
```bash
> go get github.com/lujin123/flex
```

## 使用

用起来很简单，跟其`gin`或者`echo`之类的很像的，例子如下：

```go
func main() {
	router := flex.New()
	v1 := router.Group("/v1", logMiddleware)
	{
		v1.Get("/post/{id:\\d+}", detail)
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
```

## middleware
稍微说下`middleware`的实现，主要就是利用装饰器模式，将所有的`middleware`
包裹住`controller`就行了，很喜欢这样的用法，找回了写`python`的感觉，装饰器实例：

```go
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
```
只需要将它放在一个合适位置即可

## 路由规则

### 静态路由

这个没啥说的，就是不带参数的，这个是必须完全匹配的

### 参数路由

1. `/post/{id:\\d+}`

这个路由会匹配 `/post/123` 这样的地址

2. `/post/{:\\d+}`

这个和上面的一样，区别在于，没有指定参数名字，会默认给个`key`，所以这样的不可以存在多个，不然就会丢失参数

3. `/post/{id}` 或者 `/post/:id`

这个两种只有参数名字的，默认会给个正则，默认的正则是`[0-9a-zA-Z_.]+`，所以可以匹配 `/post/123` 或者 `/post/hello`
