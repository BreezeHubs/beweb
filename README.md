"# beweb" 

# web框架主体与router
![](./resource/beweb-route.png)

```go
h := NewHTTPServer()

h.Get("/user", func(ctx *Context) {
    ctx.Resp.Write([]byte("hello world"))
})

h.Start(":8080")
```

