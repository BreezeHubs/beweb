package beweb

import (
	"bytes"
	"errors"
	"io"
	"strconv"
	"unsafe"
)

type stringValue struct {
	value string
	err   error
}

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

// FormParam 获取form值
func (c *Context) FormParam(key string) stringValue {
	if c.FormParams == nil {
		if err := c.initFormValue(); err != nil {
			return stringValue{err: err}
		}
	}

	value, ok := c.FormParams[key]
	if !ok {
		return stringValue{err: errors.New("key：" + key + "对应数据不存在")}
	}
	return stringValue{value: value}
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

// QueryParam 获取url get的值
func (c *Context) QueryParam(key string) stringValue {
	if c.QueryParams == nil {
		c.initQueryValue()
	}

	value, ok := c.QueryParams[key]
	if !ok {
		return stringValue{err: errors.New("key：" + key + "对应数据不存在")}
	}
	return stringValue{value: value}
}

// QueryValueAll 获取url get的所有值
func (c *Context) QueryValueAll() map[string]string {
	if c.QueryParams == nil {
		c.initQueryValue()
	}
	return c.QueryParams
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

// PathValue 获取路由参数
func (c *Context) PathValue(key string) (string, error) {
	value, ok := c.PathParams[key]
	if !ok {
		return "", errors.New("key：" + key + "对应数据不存在")
	}
	return value, nil
}

// PathParam 获取路由参数
func (c *Context) PathParam(key string) stringValue {
	value, ok := c.PathParams[key]
	if !ok {
		return stringValue{err: errors.New("key：" + key + "对应数据不存在")}
	}
	return stringValue{value: value}
}

// PathValueAll 获取路由所有参数
func (c *Context) PathValueAll() map[string]string {
	return c.PathParams
}

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

// 初始化body数据
func (c *Context) initBody() {
	//不存在body缓存，则初始化缓存
	if c.Body == nil {
		bodyData, _ := io.ReadAll(c.Req.Body)
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

// HeaderValue 获取header参数
func (c *Context) HeaderValue(key string) (string, error) {
	value, ok := c.HeaderParams[key]
	if !ok {
		return "", errors.New("key：" + key + "对应数据不存在")
	}
	return value, nil
}

// HeaderParam 获取header参数
func (c *Context) HeaderParam(key string) stringValue {
	value, ok := c.HeaderParams[key]
	if !ok {
		return stringValue{err: errors.New("key：" + key + "对应数据不存在")}
	}
	return stringValue{value: value}
}

// HeaderValueAll 获取header所有参数
func (c *Context) HeaderValueAll() map[string]string {
	return c.HeaderParams
}

// 转换格式
func (s stringValue) String() (string, error) {
	if s.err != nil {
		return "", s.err
	}
	return s.value, nil
}

func (s stringValue) Int64() (int64, error) {
	if s.err != nil {
		return 0, s.err
	}
	return strconv.ParseInt(s.value, 10, 64)
}

func (s stringValue) Int() (int, error) {
	if s.err != nil {
		return 0, s.err
	}
	val, err := strconv.ParseInt(s.value, 10, 64)
	return int(val), err
}

func (s stringValue) Bool() (bool, error) {
	if s.err != nil {
		return false, s.err
	}
	return strconv.ParseBool(s.value)
}

func (s stringValue) Bytes() ([]byte, error) {
	if s.err != nil {
		return nil, s.err
	}
	return *(*[]byte)(unsafe.Pointer(&s.value)), nil
}
