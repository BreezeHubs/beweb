package beweb

import (
	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func (c *Context) Response(code int, value []byte) error {
	c.Resp.WriteHeader(code)
	_, err := c.Resp.Write(value)
	return err
}
