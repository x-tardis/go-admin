package servers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thinkgos/gin-middlewares/requestid"
	xrequestid "github.com/thinkgos/http-middlewares/requestid"
)

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

func WithPrompt(s fmt.Stringer) Option {
	return func(r *Response) {
		r.Msg = s.String()
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

func JSONs(c *gin.Context, httpCode int, obj interface{}) {
	c.Header(xrequestid.RequestIDHeader, requestid.FromRequestID(c))
	c.JSON(httpCode, obj)
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
	c.Header(xrequestid.RequestIDHeader, requestid.FromRequestID(c))
	c.JSON(httpCode, r)
}

// JSONCustom http.StatusBadRequest式应答,自定义prompt码应答,一般给前端显示使用
func JSONCustom(c *gin.Context, prompt fmt.Stringer, opts ...Option) {
	rsp := Response{
		Code: http.StatusBadRequest,
		Msg:  prompt.String(),
		Data: "{}",
	}
	for _, opt := range opts {
		opt(&rsp)
	}
	c.Header(xrequestid.RequestIDHeader, requestid.FromRequestID(c))
	c.JSON(http.StatusBadRequest, rsp)
}

// JSONDetail http.StatusBadRequest式应答,detail为err的stringer
func JSONDetail(c *gin.Context, err error, opts ...Option) {
	rsp := Response{
		Code:   http.StatusBadRequest,
		Msg:    http.StatusText(http.StatusBadRequest),
		Data:   "{}",
		Detail: err.Error(),
	}
	for _, opt := range opts {
		opt(&rsp)
	}
	c.Header(xrequestid.RequestIDHeader, requestid.FromRequestID(c))
	c.JSON(http.StatusBadRequest, rsp)
}

// 通常成功数据处理
func OK(c *gin.Context, opts ...Option) {
	JSON(c, http.StatusOK, opts...)
}

func Fail(c *gin.Context, code int, opts ...Option) {
	r := Response{
		Code: code,
		Msg:  http.StatusText(code),
		Data: "{}",
	}

	for _, opt := range opts {
		opt(&r)
	}
	c.Header(xrequestid.RequestIDHeader, requestid.FromRequestID(c))
	c.JSON(http.StatusOK, r)
}
