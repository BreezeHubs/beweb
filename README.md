# beweb

## web框架主体与router
![](./resource/beweb-route.png)

### 创建服务
```go
//创建服务
h := beweb.NewHTTPServer()
//运行
h.Start(":8080")
```

### 创建路由
path 限制：
- 不支持空字符串
- 必须以 / 开头
- 不能以 / 结尾
- 中间不能有连续的 ///

优先级：  
- route: /test/:id                => http://127.0.0.1:8080/test/1  
- route: /test/*\/user            => http://127.0.0.1:8080/test/abc/user  
- route: /test/Reg(^\d{4}-\d{8}$) => http://127.0.0.1:8080/test/0931-87562388   

互相不能共存，会导致painc，只能和【静态路由】共存 
例如：  
```go
h.Get("/user/:id", func(ctx *beweb.Context) {
    fmt.Println("hello user")
})
h.Get("/user/*/abc", func(ctx *beweb.Context) {
    fmt.Println("hello user")
})
```
> panic: method: GET, path: /user/*/abc, error: 不允许同时存在【参数路由】和【通配符路由】，已存在【参数路由】

```go
//创建静态路由
h.Get("/user", func(ctx *beweb.Context) {
    fmt.Println("hello world")
})

//创建路由参数路由
h.Get("/user/:id", func(ctx *beweb.Context) {
    fmt.Println("hello user")
})

//创建通配符路由
h.Get("/order/*/detail", func(ctx *beweb.Context) {
    fmt.Println("hello order")
})

//创建正则路由
h.Get("/info/Reg(^\\d{4}-\\d{8}$)", func(ctx *beweb.Context) {
    fmt.Println("hello info")
})
```

### 获取参数
```go
h.Get("/param/:name", func(ctx *beweb.Context) {
    //获取路由参数
    value, err := ctx.PathValue("name")
    fmt.Println(value, err)
    value, err = ctx.PathParam("name").String()
    fmt.Println(value, err)
    
    //获取所有路由参数
    all := ctx.PathValueAll() //map[string]string
    fmt.Println(all)
    all = ctx.PathParams //map[string]string
    fmt.Println(all)
    
    //获取get参数
    value, err = ctx.QueryValue("id")
    fmt.Println(value, err)
    id, err := ctx.QueryParam("id").Int64()
    fmt.Println(id, err)
    
    //获取所有get参数
    all = ctx.QueryValueAll() //map[string]string
    fmt.Println(all)
    all = ctx.QueryParams //map[string]string
    fmt.Println(all)
    
    //获取Form参数
    value, err = ctx.FormValue("date")
    fmt.Println(value, err)
    value, err = ctx.FormParam("date").String()
    fmt.Println(value, err)
    
    //获取所有Form参数
    all, err = ctx.FormValueAll() //map[string]string
    fmt.Println(all, err)
    all = ctx.FormParams //map[string]string
    fmt.Println(all)
    
    //获取json参数
    type user struct {
        id   int    `json:"id"`
        name string `json:"name"`
    }
    var u user
    err = ctx.BindJSON(&u)
    fmt.Println(u, err)
})
```

### 响应数据
```go
//返回字符串
h.Get("response", func(ctx *beweb.Context) {
    ctx.Response(200, []byte("success"))
})

//使用扩展包 import "github.com/BreezeHubs/beweb/util"
//json
h.Get("response/json", func(ctx *beweb.Context) {
    util.ResponseJSON(ctx, 200, "SUCCESS", "请求成功", map[string]string{
    "id":   "1",
    "name": "haha",
    })
})

//xml
h.Get("response/xml", func(ctx *beweb.Context) {
    util.ResponseXML(ctx, 200, map[string]string{
    "id":   "1",
    "name": "haha",
    })
})
```