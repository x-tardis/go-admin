package servers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/x-tardis/go-admin/pkg/middleware"
)

// Code code interface
type Code interface {
	fmt.Stringer
	Value() int
}
type Response struct {
	Code   int         `json:"code" example:"200"`
	Msg    string      `json:"msg,omitempty"`
	Detail string      `json:"detail,omitempty"` // 错误携带的信息, 用于开发者调试
	Data   interface{} `json:"data"`
}

type Option func(r *Response)

func WithData(data interface{}) Option {
	return func(r *Response) {
		r.Data = data
	}
}

func WithCode(code int) Option {
	return func(r *Response) {
		r.Code = code
	}
}

func WithMsg(msg string) Option {
	return func(r *Response) {
		r.Msg = msg
	}
}

// WithICode Code interface 使应答修改code和msg,用于显示
func WithICode(code Code) Option {
	return func(r *Response) {
		r.Code = code.Value()
		r.Msg = code.String()
	}
}

// WithDetail detail 开发调试使用
func WithDetail(detail string) Option {
	return func(r *Response) {
		r.Detail = detail
	}
}

// WithError err detail为err的stringer
func WithError(err error) Option {
	return func(r *Response) {
		r.Detail = err.Error()
	}
}

func JSON(c *gin.Context, httpCode int, opts ...Option) {
	r := Response{
		Code: httpCode,
		Msg:  http.StatusText(httpCode),
		Data: "{}",
	}

	for _, opt := range opts {
		opt(&r)
	}
	c.JSON(httpCode, r)
}

// 通常成功数据处理
func OKWithRequestID(c *gin.Context, data interface{}, msg string) {
	c.Header("X-Request-Id", middleware.GenerateMsgIDFromContext(c))
	JSON(c, http.StatusOK, WithData(data), WithMsg(msg))
}

func Fail(c *gin.Context, code int, msg string) {
	JSON(c, http.StatusOK, WithCode(code), WithMsg(msg))
}

// 失败数据处理
func FailWithRequestID(c *gin.Context, code int, msg string) {
	c.Header("X-Request-Id", middleware.GenerateMsgIDFromContext(c))
	JSON(c, http.StatusOK, WithCode(code), WithMsg(msg))
}
