package servers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/x-tardis/go-admin/pkg/gcontext"
)

type Response struct {
	Code      int         `json:"code" example:"200"`
	Data      interface{} `json:"data"`
	Msg       string      `json:"msg,omitempty"`
	RequestId string      `json:"requestId,omitempty"`
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

func WithMessage(msg string) Option {
	return func(r *Response) {
		r.Msg = msg
	}
}

func WithRequestId(requestId string) Option {
	return func(r *Response) {
		r.RequestId = requestId
	}
}

func JSON(c *gin.Context, httpCode int, opts ...Option) {
	r := Response{Code: http.StatusOK, Data: "{}"}

	for _, opt := range opts {
		opt(&r)
	}
	c.JSON(httpCode, r)
}

func Success(c *gin.Context, opts ...Option) {
	JSON(c, http.StatusOK, opts...)
}

// 通常成功数据处理
func OKWithRequestID(c *gin.Context, data interface{}, msg string) {
	Success(c, WithData(data), WithMessage(msg), WithRequestId(gcontext.GenerateMsgIDFromContext(c)))
}

func Fail(c *gin.Context, code int, msg string) {
	Success(c, WithCode(code), WithMessage(msg))
}

// 失败数据处理
func FailWithRequestID(c *gin.Context, code int, msg string) {
	Success(c, WithCode(code), WithMessage(msg))
}
