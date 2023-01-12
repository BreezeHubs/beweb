package util

import (
	"encoding/xml"
	"errors"
	"github.com/BreezeHubs/beweb"
	"strconv"
)

// ResponseXML 通用响应
func ResponseXML(c *beweb.Context, code int, value any) error {
	bytes, err := xml.Marshal(value)
	if err != nil {
		return errors.New("ResponseXML: " + err.Error())
	}

	c.Resp.Header().Set("Content-Type", "application/xml;charset=utf-8")
	c.Resp.Header().Set("Content-Length", strconv.Itoa(len(bytes)))
	return c.Response(code, bytes)
}
