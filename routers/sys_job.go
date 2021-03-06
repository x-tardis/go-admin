package routers

import (
	"github.com/gin-gonic/gin"

	"github.com/x-tardis/go-admin/apis/actions"
	"github.com/x-tardis/go-admin/apis/system"
	"github.com/x-tardis/go-admin/app/service"
	"github.com/x-tardis/go-admin/models"
)

// 需认证的路由代码
func SysJobRouter(v1 gin.IRouter) {
	r := v1.Group("/jobs")
	{
		sysJob := &models.Job{}
		r.GET("", actions.PermissionAction(), actions.IndexAction(sysJob, new(service.SysJobSearch), func() interface{} {
			list := make([]models.Job, 0)
			return &list
		}))
		r.GET("/:id", actions.PermissionAction(), actions.ViewAction(new(service.SysJobById), func() interface{} {
			return &service.SysJobItem{}
		}))
		r.POST("", actions.CreateAction(new(service.SysJobControl)))
		r.PUT("", actions.PermissionAction(), actions.UpdateAction(new(service.SysJobControl)))
		r.DELETE("", actions.PermissionAction(), actions.DeleteAction(new(service.SysJobById)))
	}

	ctl := new(system.Job)
	rx := v1.Group("/job")
	{
		rx.GET("/remove/:id", ctl.RemoveJobForService)
		rx.GET("/start/:id", ctl.StartJobForService)
	}
}
