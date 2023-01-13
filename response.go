package beweb

import "unsafe"

func (c *Context) Response(code int, value []byte) error {
	c.Resp.WriteHeader(code)
	_, err := c.Resp.Write(value)

	//存储响应结果用于log
	if err == nil {
		c.ResponseData.Status = code
		c.ResponseData.Content = *(*string)(unsafe.Pointer(&value))
	}
	return err
}
