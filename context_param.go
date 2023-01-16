package beweb

import (
	"bytes"
	"errors"
)

type stringValue struct {
	value string
	err   error
}

/* -------------------------------------- Form -------------------------------------- */

// FormValue 获取form值
func (c *Context) FormValue(key string) (string, error) {
	if c.FormParams == nil {
		if err := c.initFormValue(); err != nil {
			return "", err
		}
	}

	value, ok := c.FormParams[key]
	if !ok {
		return "", errors.New("key：" + key + "对应数据不存在")
	}
	return value, nil
}

// FormValueAll 获取form所有值
func (c *Context) FormValueAll() (map[string]string, error) {
	if c.FormParams == nil {
		if err := c.initFormValue(); err != nil {
			return nil, err
		}
	}
	return c.FormParams, nil
}

/* -------------------------------------- Form -------------------------------------- */

/* -------------------------------------- Query -------------------------------------- */

// QueryValue 获取url get的值
func (c *Context) QueryValue(key string) (string, error) {
	if c.QueryParams == nil {
		c.initQueryValue()
	}

	value, ok := c.QueryParams[key]
	if !ok {
		return "", errors.New("key：" + key + "对应数据不存在")
	}
	return value, nil
}

// QueryValueAll 获取url get的所有值
func (c *Context) QueryValueAll() map[string]string {
	if c.QueryParams == nil {
		c.initQueryValue()
	}
	return c.QueryParams
}

/* -------------------------------------- Query -------------------------------------- */

/* -------------------------------------- Path -------------------------------------- */

// PathValue 获取路由参数
func (c *Context) PathValue(key string) (string, error) {
	value, ok := c.PathParams[key]
	if !ok {
		return "", errors.New("key：" + key + "对应数据不存在")
	}
	return value, nil
}

// PathValueAll 获取路由所有参数
func (c *Context) PathValueAll() map[string]string {
	return c.PathParams
}

/* -------------------------------------- Path -------------------------------------- */

/* -------------------------------------- Json -------------------------------------- */

// BindJSON 解析JSON
func (c *Context) BindJSON(value any) error {
	if c.Req.Body == nil {
		return errors.New("http body不能为nil")
	}
	if value == nil {
		return errors.New("value不能为nil")
	}

	d := json.NewDecoder(bytes.NewBuffer(c.Body))
	return d.Decode(value)
}

/* -------------------------------------- Json -------------------------------------- */

/* -------------------------------------- Header -------------------------------------- */

// HeaderValue 获取header参数
func (c *Context) HeaderValue(key string) (string, error) {
	value, ok := c.HeaderParams[key]
	if !ok {
		return "", errors.New("key：" + key + "对应数据不存在")
	}
	return value, nil
}

// HeaderValueAll 获取header所有参数
func (c *Context) HeaderValueAll() map[string]string {
	return c.HeaderParams
}

/* -------------------------------------- Header -------------------------------------- */
