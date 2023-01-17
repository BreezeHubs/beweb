package beweb

import (
	"mime/multipart"
	"net/http"
	"path/filepath"
)

/* -------------------------------------- Form File -------------------------------------- */

func (c *Context) FormFileValue(key string) (multipart.File, *multipart.FileHeader, error) {
	file, header, err := c.Req.FormFile(key)
	return file, header, err
}

/* -------------------------------------- Form File -------------------------------------- */

func (c *Context) DownLoadFile(filePath string) {
	//download必要的header设置
	header := c.Resp.Header()
	header.Set("Content-Disposition", "attachment;filename="+filepath.Base(filePath))
	header.Set("Content-Description", "File Transfer")
	header.Set("Content-Type", "application/octet-stream") //代表格式为通用的二进制文件
	header.Set("Content-Transfer-Encoding", "binary")      //binary相当于直接输出
	header.Set("Expires", "0")
	header.Set("Cache-Control", "must-revalidate")
	header.Set("Pragma", "public")

	http.ServeFile(c.Resp, c.Req, filePath)
}
