package beweb

import (
	"errors"
	"strconv"
	"unsafe"
)

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

// PathParam 获取路由参数
func (c *Context) PathParam(key string) stringValue {
	value, ok := c.PathParams[key]
	if !ok {
		return stringValue{err: errors.New("key：" + key + "对应数据不存在")}
	}
	return stringValue{value: value}
}

// HeaderParam 获取header参数
func (c *Context) HeaderParam(key string) stringValue {
	value, ok := c.HeaderParams[key]
	if !ok {
		return stringValue{err: errors.New("key：" + key + "对应数据不存在")}
	}
	return stringValue{value: value}
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
