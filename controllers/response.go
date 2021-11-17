package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// ResponseData 定义响应返回规范
type ResponseData struct {
	Code ResCode `json:"code"`
	Message interface{} `json:"message"`
	Data interface{} `json:"data,omitempty"`
}

// ResponseError定义返回错误信息
func ResponseError(c *gin.Context,code ResCode)  {
	rd := &ResponseData{
		Code: code,
		Message: code.Msg(),
		Data: nil,
	}
	c.JSON(http.StatusOK,rd)
}

// ResponseWithMsg 返回指定响应信息
func ResponseErrorWithMsg(c *gin.Context,code ResCode,msg interface{})  {
	rd := &ResponseData{
		Code: code,
		Message: msg,
		Data: nil,
	}
	c.JSON(http.StatusOK,rd)
}

// ResponseSuccess 返回成功响应
func ResponseSuccess(c *gin.Context,data interface{})  {
	rd := &ResponseData{
		Code: CodeSuccess,
		Message: CodeSuccess.Msg(),
		Data: data,
	}
	c.JSON(http.StatusOK,rd)
}


