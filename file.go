package beweb

import (
	"github.com/golang/groupcache/lru"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

/* ------------------------------------ Uploader ------------------------------------ */

type fileUploader struct {
	fileField    string
	savePathFunc func(header *multipart.FileHeader) string
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
		fileField:    fileField,
		savePathFunc: savePathFunc,
	}
}

func (f *fileUploader) Handle() HandleFunc {
	return func(ctx *Context) {
		//上传逻辑

		//读取到文件内容
		file, header, err := ctx.FormFileValue(f.fileField)
		if err != nil {
			ctx.Response(http.StatusInternalServerError, []byte("文件上传失败: "+err.Error()))
			return
		}
		defer file.Close()

		//计算保存路径，将计算逻辑交给用户
		savePath := f.savePathFunc(header)

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

/* ------------------------------------ Uploader ------------------------------------ */

/* ------------------------------------ Downloader ------------------------------------ */

type fileDownloader struct {
	dir string
}

func NewFileDownloader(dir string) *fileDownloader {
	return &fileDownloader{dir: dir}
}

func (f *fileDownloader) Handle() HandleFunc {
	return func(ctx *Context) {
		// xxx?file=xxx
		file, err := ctx.QueryValue("file")
		if err != nil {
			ctx.Response(http.StatusBadRequest, []byte("传递的目标文件参数错误："+err.Error()))
			return
		}

		//校验、安全性处理
		file = filepath.Clean(file) //返回同目录的最短路径
		filePath := filepath.Join(f.dir, file)
		filePath, err = filepath.Abs(filePath) //返回path相对当前路径的绝对路径
		if err != nil {
			ctx.Response(http.StatusBadRequest, []byte("传递的目标文件参数错误："+err.Error()))
			return
		}
		//filePath = filepath.ToSlash(filePath) // \ 转为 /
		//if !strings.Contains(filePath, strings.ReplaceAll(f.Dir, ".", "")) {
		//	ctx.Response(http.StatusBadRequest, []byte("传递的目标文件参数异常"))
		//	return
		//}

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

/* ------------------------------------ Downloader ------------------------------------ */

/* ------------------------------------ StaticResourceHandler ------------------------------------ */

type staticResourceHandler struct {
	dir        string
	extTypeMap map[string]string

	isCache      bool       //是否开启缓存，cacheMaxSize>即开启
	cache        *lru.Cache //缓存数据
	cacheKVSize  int        //缓存K-V数量
	cacheMinSize int        //最小缓存大小，即>cacheMinSize的文件才缓存
	cacheMaxSize int        //最大缓存大小，即<cacheMinSize的文件才缓存
}

func NewStaticResourceHandler(dir string) *staticResourceHandler {
	if dir == "" {
		dir = "./"
	}
	return &staticResourceHandler{dir: dir, extTypeMap: map[string]string{
		"jpeg": "image/jpeg",
		"jpe":  "image/jpeg",
		"jpg":  "image/jpeg",
		"png":  "image/png",
		"pdf":  "image/pdf",
	},
		cache: lru.New(1000), isCache: true, //缓存数据容量大小 200M
		cacheMinSize: 2 * 1024 * 1024, cacheMaxSize: 100 * 1024 * 1024} //2M - 100M的文件会缓存
}

func (f *staticResourceHandler) SetExtTypeMap(extTypeMap map[string]string) *staticResourceHandler {
	for key, value := range extTypeMap {
		f.extTypeMap[key] = value
	}
	return f
}

func (f *staticResourceHandler) SetCacheSize(min, max int) *staticResourceHandler {
	f.cacheMinSize = min
	f.cacheMaxSize = max
	if f.cacheMaxSize == 0 {
		f.isCache = false
	}
	return f
}

func (f *staticResourceHandler) SetCacheCapSize(size int) *staticResourceHandler {
	f.cacheKVSize = size
	return f
}

func (f *staticResourceHandler) Handle() HandleFunc {
	return func(ctx *Context) {
		//1、获取目标文件名
		file, err := ctx.PathValue("file")
		if err != nil {
			ctx.Response(http.StatusBadRequest, []byte("传递的目标文件参数错误："+err.Error()))
			return
		}

		//校验、安全性处理
		file = filepath.Clean(file) //返回同目录的最短路径
		filePath := filepath.Join(f.dir, file)
		filePath, err = filepath.Abs(filePath) //返回path相对当前路径的绝对路径
		if err != nil {
			ctx.Response(http.StatusBadRequest, []byte("传递的目标文件参数错误："+err.Error()))
			return
		}
		//filePath = filepath.ToSlash(filePath) // \ 转为 /
		//if !strings.Contains(filePath, strings.ReplaceAll(f.Dir, ".", "")) {
		//	ctx.Response(http.StatusBadRequest, []byte("传递的目标文件参数异常"))
		//	return
		//}

		typeString := ""
		ok := true
		//filepath.Ext(filePath) = ".jpg"，需要去掉”.“
		if len([]byte(filepath.Ext(filePath))) > 1 {
			typeString, ok = f.extTypeMap[filepath.Ext(filePath)[1:]]
			if !ok {
				ctx.Response(http.StatusBadRequest, []byte("不支持的文件类型"))
				return
			}
		}

		header := ctx.Resp.Header()

		if f.isCache {
			//2、判断是否有缓存
			if data, ok := f.cache.Get(filePath); ok {
				header := ctx.Resp.Header()
				header.Set("Content-Type", typeString)
				header.Set("Content-Length", strconv.Itoa(len(data.([]byte))))
				ctx.Response(http.StatusOK, data.([]byte))
				return
			}
		}

		//2、定位到目标文件，读取
		data, err := os.ReadFile(filePath)
		if err != nil {
			ctx.Response(http.StatusInternalServerError, []byte("传递的目标文件读取异常"+err.Error()))
			return
		}

		if f.isCache && len(data) >= f.cacheMinSize && len(data) <= f.cacheMaxSize {
			f.cache.Add(filePath, data) //加缓存
		}

		//3、返回文件内容
		header.Set("Content-Type", typeString)
		header.Set("Content-Length", strconv.Itoa(len(data)))

		ctx.Response(http.StatusOK, data)
	}
}

/* ------------------------------------ StaticResourceHandler ------------------------------------ */
