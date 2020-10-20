package servers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thinkgos/gin-middlewares/requestid"
	xrequestid "github.com/thinkgos/http-middlewares/requestid"
)

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg,omitempty"`
	Data interface{} `json:"data"`
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

// WithError data 为err
func WithError(err error) Option {
	return func(r *Response) {
		r.Data = err.Error()
	}
}

// 通常成功数据处理
func OK(c *gin.Context, opts ...Option) {
	r := Response{
		http.StatusOK,
		http.StatusText(http.StatusOK),
		"{}",
	}

	for _, opt := range opts {
		opt(&r)
	}
	c.Header(xrequestid.RequestIDHeader, requestid.FromRequestID(c))
	c.JSON(http.StatusOK, r)
}

func Fail(c *gin.Context, code int, opts ...Option) {
	r := Response{
		code,
		http.StatusText(code),
		"{}",
	}

	for _, opt := range opts {
		opt(&r)
	}
	c.Header(xrequestid.RequestIDHeader, requestid.FromRequestID(c))
	c.JSON(http.StatusOK, r)
}

func JSON(c *gin.Context, httpCode int, obj interface{}) {
	c.Header(xrequestid.RequestIDHeader, requestid.FromRequestID(c))
	c.JSON(httpCode, obj)
}
