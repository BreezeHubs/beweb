package beweb

import (
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type fileUploader struct {
	FileField    string
	SavePathFunc func(header *multipart.FileHeader) string
}

func NewFileUploader(fileField string, savePathFunc func(header *multipart.FileHeader) string) *fileUploader {
	if fileField == "" {
		fileField = "file"
	}
	if savePathFunc == nil {
		savePathFunc = func(header *multipart.FileHeader) string {
			return "./upload/" + header.Filename
		}
	}
	return &fileUploader{
		FileField:    fileField,
		SavePathFunc: savePathFunc,
	}
}

func (f fileUploader) Handler() HandleFunc {
	return func(ctx *Context) {
		//上传逻辑

		//读取到文件内容
		file, header, err := ctx.FormFileValue(f.FileField)
		if err != nil {
			ctx.Response(http.StatusInternalServerError, []byte("文件上传失败: "+err.Error()))
			return
		}
		defer file.Close()

		//计算保存路径，将计算逻辑交给用户
		savePath := f.SavePathFunc(header)

		//创建目录
		if err := os.MkdirAll(filepath.Dir(savePath), 0644); err != nil {
			ctx.Response(http.StatusInternalServerError, []byte("文件上传失败: "+err.Error()))
			return
		}

		//保存文件
		//os.O_WRONLY 写入权限
		//os.O_TRUNC 文件存在则清空内容
		//os.O_CREATE 创建文件权限
		copyFile, err := os.OpenFile(savePath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
		if err != nil {
			ctx.Response(http.StatusInternalServerError, []byte("文件上传失败: "+err.Error()))
			return
		}
		defer copyFile.Close()

		_, err = io.CopyBuffer(copyFile, file, nil)
		if err != nil {
			ctx.Response(http.StatusInternalServerError, []byte("文件上传失败: "+err.Error()))
			return
		}

		//返回响应
		ctx.Response(http.StatusOK, []byte("文件上传成功"))
	}
}

type fileDownloader struct {
	Dir string
}

func NewFileDownloader(dir string) *fileDownloader {
	return &fileDownloader{Dir: dir}
}

func (f fileDownloader) Handler() HandleFunc {
	return func(ctx *Context) {
		// xxx?file=xxx
		file, err := ctx.QueryValue("file")
		if err != nil {
			ctx.Response(http.StatusBadRequest, []byte("传递的目标文件参数错误："+err.Error()))
			return
		}

		//校验、安全性处理
		file = filepath.Clean(file) //返回同目录的最短路径
		filePath := filepath.Join(f.Dir, file)
		filePath, err = filepath.Abs(filePath) //返回path相对当前路径的绝对路径
		if err != nil {
			ctx.Response(http.StatusBadRequest, []byte("传递的目标文件参数错误："+err.Error()))
			return
		}
		if !strings.Contains(filePath, f.Dir) {
			ctx.Response(http.StatusBadRequest, []byte("传递的目标文件参数异常"))
			return
		}

		//download必要的header设置
		header := ctx.Resp.Header()
		header.Set("Content-Disposition", "attachment;filename="+filepath.Base(filePath))
		header.Set("Content-Description", "File Transfer")
		header.Set("Content-Type", "application/octet-stream") //代表格式为通用的二进制文件
		header.Set("Content-Transfer-Encoding", "binary")      //binary相当于直接输出
		header.Set("Expires", "0")
		header.Set("Cache-Control", "must-revalidate")
		header.Set("Pragma", "public")

		http.ServeFile(ctx.Resp, ctx.Req, filePath)
	}
}
