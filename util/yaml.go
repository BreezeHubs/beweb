package util

import (
	"errors"
	"github.com/BreezeHubs/beweb"
	"gopkg.in/yaml.v3"
	"strconv"
)

// ResponseYAML 通用响应
func ResponseYAML(c *beweb.Context, code int, value any) error {
	bytes, err := yaml.Marshal(value)
	if err != nil {
		return errors.New("ResponseXML: " + err.Error())
	}

	c.Resp.Header().Set("Content-Type", "application/yaml;charset=utf-8")
	c.Resp.Header().Set("Content-Length", strconv.Itoa(len(bytes)))
	c.Response(code, bytes)
	return nil
}
