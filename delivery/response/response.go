package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Response 响应结构体
type Response struct {
	ErrorCode int         `json:"error_code"` // 自定义错误码
	Data      interface{} `json:"data"`       // 数据
	Message   string      `json:"message"`    // 信息
}

// Success 响应成功 ErrorCode 为 0 表示成功
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		0,
		data,
		"ok",
	})
}

// Fail 响应失败 ErrorCode 不为 0 表示失败
func Fail(c *gin.Context, errorCode int, msg string) {
	c.JSON(http.StatusOK, Response{
		errorCode,
		nil,
		msg,
	})
}

// FailByError 失败响应 返回自定义错误的错误码、错误信息
func FailByError(c *gin.Context, error CustomError) {
	Fail(c, error.ErrorCode, error.ErrorMsg)
}

// TokenFail 登录授权失效
func TokenFail(c *gin.Context) {
	FailByError(c, Errors.TokenError)
}

// ValidateFail 请求参数验证失败
func ValidateFail(c *gin.Context, msg string) {
	Fail(c, 10000, msg)
}

// SqlFail 写入数据库失败
func SqlFail(c *gin.Context, msg string) {
	Fail(c, 10001, msg)
}
