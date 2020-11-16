package xxfield

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/thinkgos/gin-middlewares/requestid"
	"go.uber.org/multierr"
	"go.uber.org/zap"
)

type errorKey struct{}

// ContextError context error
func ContextError(c *gin.Context, err error) {
	if err == nil {
		panic("xxfield: context error is nil")
	}
	ctx := c.Request.Context()
	if e, ok := ctx.Value(errorKey{}).(error); ok && e != nil {
		multierr.AppendInto(&e, err)
		err = e
	}
	ctx = context.WithValue(ctx, errorKey{}, err)
	c.Request = c.Request.WithContext(ctx)
}

// Error error zap field
func Error(c *gin.Context) zap.Field {
	err, ok := c.Request.Context().Value(errorKey{}).(error)
	if !ok || err == nil {
		return zap.Skip()
	}
	return zap.Error(err)
}

// RequestId request id
func RequestId(c *gin.Context) zap.Field {
	return zap.String("requestID", requestid.FromRequestID(c))
}
