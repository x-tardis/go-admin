package xxfiled

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/thinkgos/gin-middlewares/requestid"
	"go.uber.org/zap"
)

type errorKey struct{}

func ContextError(c *gin.Context, err error) {
	if err != nil {
		ctx := c.Request.Context()
		ctx = context.WithValue(ctx, errorKey{}, err)
		c.Request = c.Request.WithContext(ctx)
	}
}

func Error(c *gin.Context) zap.Field {
	err, ok := c.Request.Context().Value(errorKey{}).(error)
	if !ok || err == nil {
		return zap.Skip()
	}
	return zap.Error(err)
}

func RequestId(c *gin.Context) zap.Field {
	return zap.String("requestID", requestid.FromRequestID(c))
}
