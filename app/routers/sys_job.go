package routers

import (
	"github.com/gin-gonic/gin"

	"github.com/x-tardis/go-admin/app/apis/sysjob"
	"github.com/x-tardis/go-admin/app/models"
	"github.com/x-tardis/go-admin/app/service/dto"
	"github.com/x-tardis/go-admin/common/actions"
)

// 需认证的路由代码
func SysJobRouter(v1 gin.IRouter) {
	r := v1.Group("/sysjob")
	{
		sysJob := &models.Job{}
		r.GET("", actions.PermissionAction(), actions.IndexAction(sysJob, new(dto.SysJobSearch), func() interface{} {
			list := make([]models.Job, 0)
			return &list
		}))
		r.GET("/:id", actions.PermissionAction(), actions.ViewAction(new(dto.SysJobById), func() interface{} {
			return &dto.SysJobItem{}
		}))
		r.POST("", actions.CreateAction(new(dto.SysJobControl)))
		r.PUT("", actions.PermissionAction(), actions.UpdateAction(new(dto.SysJobControl)))
		r.DELETE("", actions.PermissionAction(), actions.DeleteAction(new(dto.SysJobById)))
	}
	sysJob := &sysjob.SysJob{}

	v1.GET("/job/remove/:id", sysJob.RemoveJobForService)
	v1.GET("/job/start/:id", sysJob.StartJobForService)
}
