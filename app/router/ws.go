package router

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"

	"github.com/x-tardis/go-admin/pkg/ws"
)

func RegisterWs(engine *gin.Engine, authMiddleware *jwt.GinJWTMiddleware) {
	go ws.WebsocketManager.Start()

	// 需要认证
	wsGroup := engine.Group("/ws")
	wsGroup.Use(authMiddleware.MiddlewareFunc(), OperLog())
	{
		wsjob := wsGroup.Group("/job")
		{
			wsjob.GET("/login/:id/:channel", ws.WebsocketManager.WsClient)
			wsjob.GET("/logout/:id/:channel", ws.WebsocketManager.UnWsClient)
		}
	}
}
