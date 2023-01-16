package beweb

import (
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

type FileUploader struct {
	FileField    string
	SavePathFunc func(header *multipart.FileHeader) string
}

func NewFileUploader(fileField string, savePathFunc func(header *multipart.FileHeader) string) *FileUploader {
	if fileField == "" {
		fileField = "file"
	}
	if savePathFunc == nil {
		savePathFunc = func(header *multipart.FileHeader) string {
			return "./upload/" + header.Filename
		}
	}
	return &FileUploader{
		FileField:    fileField,
		SavePathFunc: savePathFunc,
	}
}

func (f FileUploader) Handler() HandleFunc {
	return func(ctx *Context) {
		//上传逻辑

		//读取到文件内容
		file, header, err := ctx.Req.FormFile(f.FileField)
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
