package api

import (
	"github.com/gin-gonic/gin"
)

// Response 基础序列化器
type Response struct {
	Code  int         `json:"code"`
	Data  interface{} `json:"data,omitempty"`
	Msg   string      `json:"msg"`
	Error string      `json:"error,omitempty"`
}

// Ping 状态检查页面
func Ping(c *gin.Context) {
	c.JSON(200, Response{
		Code: 0,
		Msg:  "Pong",
	})
}
