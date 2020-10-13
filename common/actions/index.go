package actions

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/x-tardis/go-admin/common/dto"
	"github.com/x-tardis/go-admin/common/models"
	"github.com/x-tardis/go-admin/pkg/gcontext"
	"github.com/x-tardis/go-admin/pkg/logger"
	"github.com/x-tardis/go-admin/pkg/paginator"
	"github.com/x-tardis/go-admin/pkg/servers"
)

// IndexAction 通用查询动作
func IndexAction(m models.ActiveRecord, d dto.Index, f func() interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		db, err := gcontext.GetOrm(c)
		if err != nil {
			logger.Error(err)
			return
		}

		msgID := gcontext.GenerateMsgIDFromContext(c)
		list := f()
		object := m.Generate()
		req := d.Generate()
		var count int64

		//查询列表
		err = req.Bind(c)
		if err != nil {
			servers.FailWithRequestID(c, http.StatusUnprocessableEntity, "参数验证失败")
			return
		}

		//数据权限检查
		p := GetPermissionFromContext(c)

		err = db.WithContext(c).Model(object).
			Scopes(
				dto.MakeCondition(req.GetNeedSearch()),
				dto.Paginate(req.GetPageSize(), req.GetPageIndex()),
				Permission(object.TableName(), p),
			).
			Find(list).Limit(-1).Offset(-1).
			Count(&count).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Errorf("MsgID[%s] Index error: %s", msgID, err)
			servers.FailWithRequestID(c, http.StatusInternalServerError, "查询失败")
			return
		}

		servers.Success(c, servers.WithData(&paginator.Page{
			List:      list,
			Count:     int(count),
			PageIndex: req.GetPageIndex(),
			PageSize:  req.GetPageSize(),
		}), servers.WithMessage("查询成功"))
		c.Next()
	}
}
