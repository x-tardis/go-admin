package actions

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thinkgos/gin-middlewares/requestid"
	"github.com/thinkgos/sharp/iorm"
	"gorm.io/gorm"

	"github.com/thinkgos/sharp/core/paginator"

	"github.com/x-tardis/go-admin/common/dto"
	"github.com/x-tardis/go-admin/common/models"
	"github.com/x-tardis/go-admin/pkg/deployed"
	"github.com/x-tardis/go-admin/pkg/izap"
	"github.com/x-tardis/go-admin/pkg/servers"
)

// IndexAction 通用查询动作
func IndexAction(m models.ActiveRecord, d dto.Index, f func() interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		msgID := requestid.FromRequestID(c)
		list := f()
		object := m.Generate()
		req := d.Generate()
		var count int64

		//查询列表
		err := req.Bind(c)
		if err != nil {
			servers.Fail(c, http.StatusUnprocessableEntity, servers.WithMsg("参数验证失败"))
			return
		}

		//数据权限检查
		p := GetPermissionFromContext(c)

		pageParam := req.GetPaginatorParam()
		err = deployed.DB.WithContext(c).Model(object).
			Scopes(
				dto.MakeCondition(req.GetNeedSearch()),
				iorm.Paginate(pageParam),
				Permission(object.TableName(), p),
			).
			Find(list).Limit(-1).Offset(-1).
			Count(&count).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			izap.Sugar.Errorf("MsgID[%s] Index error: %s", msgID, err)
			servers.Fail(c, http.StatusInternalServerError, servers.WithMsg("查询失败"))
			return
		}

		servers.OK(c, servers.WithData(&paginator.Page{
			List:      list,
			Total:     count,
			PageIndex: pageParam.PageIndex,
			PageSize:  pageParam.PageSize,
		}), servers.WithMsg("查询成功"))
		c.Next()
	}
}
