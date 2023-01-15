package util

import (
	"errors"
	"github.com/BreezeHubs/beweb"
	jsoniter "github.com/json-iterator/go"
	"net/http"
	"strconv"
)

var (
	JSONSuccessReason  = "SUCCESS"
	JSONSuccessMessage = "请求成功"
	JSONFailReason     = "BAD_REQUEST"
	JSONFailMessage    = "请求错误"
	JSONErrorReason    = "SERVER_ERROR"
	JSONErrorMessage   = "系统异常"

	json = jsoniter.ConfigCompatibleWithStandardLibrary
)

type Response struct {
	Code    int    `json:"code"`    //状态码
	Reason  string `json:"reason"`  //业务响应码
	Message string `json:"message"` //信息描述
	Data    any    `json:"data"`    //数据
}

// ResponseJSON 通用响应
func ResponseJSON(c *beweb.Context, code int, reason, message string, value any) error {
	resp := &Response{Code: code, Reason: reason, Message: message, Data: value}

	bytes, err := json.Marshal(resp)
	if err != nil {
		return errors.New("ResponseJSON: " + err.Error())
	}

	if c.Resp != nil {
		c.Resp.Header().Set("Content-Type", "application/json;charset=utf-8")
		c.Resp.Header().Set("Content-Length", strconv.Itoa(len(bytes)))
	}

	c.Response(code, bytes)
	return nil
}

// ResponseJSONSuccess 请求成功
func ResponseJSONSuccess(c *beweb.Context, value any) error {
	return ResponseJSON(c, http.StatusOK, JSONSuccessReason, JSONSuccessMessage, value)
}

// ResponseJSONFail 请求错误
func ResponseJSONFail(c *beweb.Context, reason, message string) error {
	if reason == "" {
		reason = JSONFailReason
	}
	if message == "" {
		message = JSONFailMessage
	}
	return ResponseJSON(c, http.StatusBadRequest, reason, message, nil)
}

// ResponseJSONError 服务器错误
func ResponseJSONError(c *beweb.Context, reason, message string) error {
	if reason == "" {
		reason = JSONErrorReason
	}
	if message == "" {
		message = JSONErrorMessage
	}
	return ResponseJSON(c, http.StatusInternalServerError, reason, message, nil)
}
