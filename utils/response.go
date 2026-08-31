package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 统一API响应结构
type Response struct {
	Code string      `json:"code"`
	Data interface{} `json:"data,omitempty"`
	Msg  string      `json:"msg"`
}

const (
	SUCCESS = "00000"
	ERROR   = "50000"
)

// Success 成功响应
func Success(c *gin.Context, data interface{}, msg string) {
	if msg == "" {
		msg = "操作成功"
	}
	c.JSON(http.StatusOK, Response{
		Code: SUCCESS,
		Data: data,
		Msg:  msg,
	})
}

// Error 错误响应
func Error(c *gin.Context, code string, msg string) {
	if code == "" {
		code = ERROR
	}
	c.JSON(http.StatusBadRequest, Response{
		Code: code,
		Data: nil,
		Msg:  msg,
	})
}

// ErrorWithStatus 带状态码的错误响应
func ErrorWithStatus(c *gin.Context, status int, code string, msg string) {
	if code == "" {
		code = ERROR
	}
	c.JSON(status, Response{
		Code: code,
		Data: nil,
		Msg:  msg,
	})
}

// Fail 快捷失败响应 (使用默认错误码)
func Fail(c *gin.Context, msg string) {
	Error(c, ERROR, msg)
}
