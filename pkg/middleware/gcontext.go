package middleware

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
	"github.com/thinkgos/gin-middlewares/requestid"
)

// GetOrm 获取orm连接
func GetOrm(c *gin.Context) (*gorm.DB, error) {
	msgID := requestid.FromRequestID(c)
	idb, exist := c.Get("db")
	if !exist {
		return nil, fmt.Errorf("msgID[%s], db connect not exist", msgID)
	}
	switch v := idb.(type) {
	case *gorm.DB:
		// 新增操作
		return v, nil
	default:
		return nil, fmt.Errorf("msgID[%s], db connect not exist", msgID)
	}
}
