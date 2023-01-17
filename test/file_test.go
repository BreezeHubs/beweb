package test

import (
	"bytes"
	"fmt"
	"github.com/BreezeHubs/beweb"
	"github.com/BreezeHubs/beweb/util"
	"html/template"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strconv"
	"testing"
	"time"
)

func TestFileUploader_Handle(t *testing.T) {
	s := beweb.NewHTTPServer()

	s.Get("/upload_page", func(ctx *beweb.Context) {
		tpl := template.New("upload")
		tpl, err := tpl.Parse(`<html><body>
<form action="/upload" method="post" enctype="multipart/form-data">
<input type="file" name="myfile"><button type="submit">上传</button>
</form>
</body><html>`)
		if err != nil {
			fmt.Println(err)
			return
		}

		page := &bytes.Buffer{}
		if err = tpl.Execute(page, nil); err != nil {
			fmt.Println(err)
			return
		}
		ctx.Response(http.StatusOK, page.Bytes())
	})

	//方式1：路由注册的file handler
	s.Post("/upload",
		beweb.NewFileUploader(
			"myfile",
			func(header *multipart.FileHeader) string {
				return "./testdata/upload/" + strconv.Itoa(time.Now().Nanosecond()) + filepath.Ext(header.Filename)
			},
		).Handle(),
	)

	//方式2：FormFile
	s.Post("/upload1", func(ctx *beweb.Context) {
		value, header, err := ctx.FormFileValue("myfile")
		fmt.Println(value, header, err)
	})

	s.Start(":8080")
}

func TestFileDownloader_Handle(t *testing.T) {
	s := beweb.NewHTTPServer()

	s.Get("/download", beweb.NewFileDownloader("./testdata/upload").Handle())

	s.Get("/download1", func(ctx *beweb.Context) {
		handler := beweb.NewFileDownloader("./../testdata/upload").Handle()
		handler(ctx)

		if ctx.ResponseStatus == http.StatusBadRequest {
			util.ResponseJSONFail(ctx, "FILE ERROR", string(ctx.ResponseContent))
			return
		}
	})

	s.Get("/download2", func(ctx *beweb.Context) {
		file, _ := ctx.QueryValue("file")
		ctx.DownLoadFile("./../testdata/upload/" + file)
	})

	s.Start(":8080")
}

func TestStaticResourceHandler_Handle(t *testing.T) {
	s := beweb.NewHTTPServer()

	s.Get("/static/:file",
		beweb.NewStaticResourceHandler("./../testdata/static/").
			SetCacheCapSize(1000).
			SetCacheSize(2*1024*1024, 100*1024*1024).
			SetExtTypeMap(map[string]string{}).
			Handle(),
	)
	s.Start(":8080")
}
