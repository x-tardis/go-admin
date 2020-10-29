package router

import (
	"context"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"

	"github.com/x-tardis/go-admin/app/apis/ws"
)

func RegisterWs(engine *gin.Engine, authMiddleware *jwt.GinJWTMiddleware) {
	go ws.WsHub.Run(context.Background())
	go ws.FileMonitor(context.Background(), "temp/job.log", ws.JobGroup, ws.SendGroup)
	// 需要认证
	wsGroup := engine.Group("/ws")
	wsGroup.Use(authMiddleware.MiddlewareFunc(), OperLog())
	{
		ctl := new(ws.Cws)
		wsjob := wsGroup.Group("/job")
		{
			wsjob.GET("/login/:id", ctl.WsClient)
			wsjob.GET("/logout/:id", ctl.UnWsClient)
		}
	}
}
