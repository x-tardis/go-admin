package routers

import (
	"github.com/gin-gonic/gin"

	"github.com/x-tardis/go-admin/apis/ws"
)

func WsJob(wsGroup gin.IRouter) {
	ctl := new(ws.Cws)
	rJob := wsGroup.Group("/job")
	{
		rJob.GET("/login/:id", ctl.WsClient)
		rJob.GET("/logout/:id", ctl.UnWsClient)
	}
}
