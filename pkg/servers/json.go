package servers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/x-tardis/go-admin/tools"
)

type Response struct {
	Code      int         `json:"code" example:"200"`
	Data      interface{} `json:"data"`
	Msg       string      `json:"msg"`
	RequestId string      `json:"requestId"`
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
	Success(c, WithData(data), WithMessage(msg), WithRequestId(tools.GenerateMsgIDFromContext(c)))
}

// 失败数据处理
func FailWithRequestID(c *gin.Context, code int, msg string) {
	Success(c, WithCode(code), WithMessage(msg))
}

type Page struct {
	List      interface{} `json:"list"`
	Count     int         `json:"count"`
	PageIndex int         `json:"pageIndex"`
	PageSize  int         `json:"pageSize"`
}

// 分页数据处理
func PageOK(c *gin.Context, result interface{}, count int, pageIndex int, pageSize int, msg string) {
	OKWithRequestID(c, Page{
		result,
		count,
		pageIndex,
		pageSize,
	}, msg)
}
