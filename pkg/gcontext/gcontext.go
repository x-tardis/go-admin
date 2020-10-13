package gcontext

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"

	"github.com/x-tardis/go-admin/pkg/infra"
)

// GenerateMsgIDFromContext 生成msgID
func GenerateMsgIDFromContext(c *gin.Context) string {
	var msgID string
	data, ok := c.Get("msgID")
	if !ok {
		msgID = infra.GenerateUUID()
		c.Set("msgID", msgID)
		return msgID
	}
	msgID = cast.ToString(data)
	return msgID
}

// GetOrm 获取orm连接
func GetOrm(c *gin.Context) (*gorm.DB, error) {
	msgID := GenerateMsgIDFromContext(c)
	idb, exist := c.Get("db")
	if !exist {
		return nil, fmt.Errorf("msgID[%s], db connect not exist", msgID)
	}
	switch idb.(type) {
	case *gorm.DB:
		// 新增操作
		return idb.(*gorm.DB), nil
	default:
		return nil, errors.New(fmt.Sprintf("msgID[%s], db connect not exist", msgID))
	}
}
