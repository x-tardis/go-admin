package servers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/x-tardis/go-admin/pkg/xxfield"
)

// Code code interface
type Code interface {
	fmt.Stringer
	Value() int
}

// Response 应答
type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg,omitempty"`
	Data interface{} `json:"data"`
	err  error       `json:"-"` // for context error
}

// Option option
type Option func(r *Response)

// WithData 设置data
func WithData(data interface{}) Option {
	return func(r *Response) {
		r.Data = data
	}
}

// WithCode 设置code
func WithCode(code int) Option {
	return func(r *Response) {
		r.Code = code
	}
}

// WithMsg 设置msg
func WithMsg(msg string) Option {
	return func(r *Response) {
		r.Msg = msg
	}
}

// WithICode 设置code和msg
func WithICode(c Code) Option {
	return func(r *Response) {
		r.Code = c.Value()
		r.Msg = c.String()
	}
}

// WithError 设置err
func WithError(err error) Option {
	return func(r *Response) {
		r.err = err
	}
}

// WithDError 设置err和data为err.Error()
func WithDError(err error) Option {
	return func(r *Response) {
		r.err = err
		r.Data = err.Error()
	}
}

// OK 通常成功数据处理,永不设置err
func OK(c *gin.Context, opts ...Option) {
	r := Response{
		http.StatusOK,
		http.StatusText(http.StatusOK),
		"{}",
		nil,
	}

	for _, opt := range opts {
		opt(&r)
	}
	c.JSON(http.StatusOK, r)
}

func Fail(c *gin.Context, code int, opts ...Option) {
	r := Response{
		code,
		http.StatusText(code),
		"{}",
		nil,
	}

	for _, opt := range opts {
		opt(&r)
	}

	if r.err != nil {
		xxfield.ContextError(c, r.err)
	}
	c.JSON(http.StatusOK, r)
}
