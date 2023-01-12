# beweb

<br>

## 1 web框架主体与router
![](./resource/beweb-route.png)

<br>

## 2 middleware设计

<br>

### 2.1 AOP
横向关注点：与业务不那么密切，但又需要处理的  
常见有：  
- 可观测性：logging、metric和tracing
- 安全相关：登录、鉴权与权限控制
- 错误处理：e.g.错误页面支持
- 可用性保证：熔断、限流和降级等

<br>

### 2.2 access log


<br>

### 2.3 Trace，OpenTelemetry

<br>

### 2.4 Prometheus

<br>

### 2.5 error page

<br>

### 2.6panic recover

<br>

## 3 other part

<br>

### 3.1 页面模板引擎

<br>

### 3.2 文件处理

<br>

### 3.3 Session

<br>

## 4 使用

<br>

### 4.1 创建服务
服务配置后面再安排，先鸽着[狗头]
```go
//创建服务
h := beweb.NewHTTPServer()
//运行
h.Start(":8080")
```

<br>

### 4.2 创建路由
path 限制：
- 不支持空字符串
- 必须以 / 开头
- 不能以 / 结尾
- 中间不能有连续的 ///

优先级：  
- 【静态路由】  route: /test/123                => http://127.0.0.1:8080/test/123  
- 【参数路由】  route: /test/:id                => http://127.0.0.1:8080/test/1  
- 【通配符路由】 route: /test/*\/user            => http://127.0.0.1:8080/test/abc/user  
- 【正则路由】  route: /test/Reg(^\d{4}-\d{8}$) => http://127.0.0.1:8080/test/0931-87562388   

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

<br>

### 4.3 获取参数
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

<br>

### 4.4 响应数据
```go
//返回字符串
h.Get("/response", func(ctx *beweb.Context) {
    ctx.Response(200, []byte("success"))
})

//使用扩展包 import "github.com/BreezeHubs/beweb/util"
//json
h.Get("/response/json", func(ctx *beweb.Context) {
    type xml struct {
        Id   int    `xml:"id"`
        Name string `xml:"name"`
    }
    util.ResponseJSON(ctx, 200, "SUCCESS", "请求成功", &xml{
        Id:   1,
        Name: "haha",
    })
})

//xml
h.Get("/response/xml", func(ctx *beweb.Context) {
    type xml struct {
        Id   int    `xml:"id"`
        Name string `xml:"name"`
    }
    util.ResponseXML(ctx, 200, &xml{
        Id:   1,
        Name: "haha",
    })
})
```

<br>

### 4.5 cookie
```go
h.Get("/cookie", func(ctx *beweb.Context) {
    ck := &http.Cookie{
        Name:    "test",
        Value:   "test",
        Expires: time.Now().Add(1 * time.Hour),
    }
    ctx.SetCookie(ck)
})
```

<br>

### 4.6 服务配置
优雅退出设置
```go
//创建服务
h := beweb.NewHTTPServer(
    beweb.WithGracefullyExit(true, func() {
        fmt.Println("test：进行一些回收动作...")
        time.Sleep(2 * time.Second)
        fmt.Println("test：回收完成")
    }),
)
```
run后使用一次Ctrl+c触发退出，两次则强制退出  
> E:\beweb\cmd> go run .  
test：进行一些回收动作...  
test：回收完成  
exit  
