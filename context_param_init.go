package beweb

import (
	"bytes"
	"io"
)

// 初始化form数据
func (c *Context) initFormValue() error {
	//不必要form cache，因为 c.Req.ParseForm() 会完成是否需要重复parse的动作，所以该调用的幂等的
	if err := c.Req.ParseForm(); err != nil {
		return err
	}

	c.FormParams = make(map[string]string, len(c.Req.Form))
	for k, v := range c.Req.Form {
		if len(v) > 0 {
			c.FormParams[k] = v[len(v)-1]
		}
	}
	return nil
}

// 初始化url get数据
func (c *Context) initQueryValue() {
	//不存在query缓存，则初始化缓存
	c.QueryParams = make(map[string]string, len(c.Req.URL.Query()))
	for k, v := range c.Req.URL.Query() {
		if len(v) > 0 {
			c.QueryParams[k] = v[len(v)-1]
		}
	}
}

// 初始化body数据
func (c *Context) initBody() {
	//不存在body缓存，则初始化缓存
	if c.Body == nil {
		bodyData, _ := io.ReadAll(c.Req.Body)
		c.Req.Body = io.NopCloser(bytes.NewBuffer(bodyData))
		c.Body = bodyData
	}
}

// 初始化header数据
func (c *Context) initHeaderValue() {
	//不存在header缓存，则初始化缓存
	c.HeaderParams = make(map[string]string, len(c.Req.Header))
	for k, v := range c.Req.Header {
		if len(v) > 0 {
			c.HeaderParams[k] = v[len(v)-1]
		}
	}
}
