package main

import (
	"fmt"
	"github.com/BreezeHubs/beweb"
	"github.com/BreezeHubs/beweb/util"
	"net/http"
	"time"
)

func main() {
	//创建服务
	h := beweb.NewHTTPServer()

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

	h.Get("response", func(ctx *beweb.Context) {
		ctx.Response(200, []byte("success"))
	})

	h.Get("response/json", func(ctx *beweb.Context) {
		util.ResponseJSON(ctx, 200, "SUCCESS", "请求成功", map[string]string{
			"id":   "1",
			"name": "haha",
		})
	})

	h.Get("response/xml", func(ctx *beweb.Context) {
		util.ResponseXML(ctx, 200, map[string]string{
			"id":   "1",
			"name": "haha",
		})
	})

	h.Get("/xml/:id", func(ctx *beweb.Context) {
		type xml struct {
			Id   int    `xml:"id"`
			Name string `xml:"name"`
		}
		id, _ := ctx.PathParam("id").Int()
		util.ResponseXML(ctx, 200, &xml{
			Id:   id,
			Name: "haha",
		})
	})

	h.Get("cookie", func(ctx *beweb.Context) {
		ck := &http.Cookie{
			Name:    "test",
			Value:   "test",
			Expires: time.Now().Add(1 * time.Hour),
		}
		ctx.SetCookie(ck)
	})

	h.Start(":8080")
}
