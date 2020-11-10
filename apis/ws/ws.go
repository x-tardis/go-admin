package ws

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/thinkgos/easyws"

	"github.com/x-tardis/go-admin/pkg/servers"
)

const JobGroup = "job"

// 初始化 wsManager 管理器
var WsHub = easyws.New()

type Cws struct{}

// gin 处理 websocket handler
func (Cws) WsClient(c *gin.Context) {
	id := c.Param("id")

	upGrader := websocket.Upgrader{
		// cross origin domain
		CheckOrigin: func(r *http.Request) bool { return true },
		// 处理 Sec-WebSocket-Protocol Header
		Subprotocols: []string{c.GetHeader("Sec-WebSocket-Protocol")},
	}

	conn, err := upGrader.Upgrade(c.Writer, c.Request, c.Writer.Header())
	if err != nil {
		log.Printf("websocket upgrade failed, %s", id)
		return
	}

	sess := &easyws.Session{
		GroupID: JobGroup,
		ID:      id,
		Request: c.Request,
		Conn:    conn,
		Hub:     WsHub,
	}
	sess.Run()
}

func (Cws) UnWsClient(c *gin.Context) {
	id := c.Param("id")
	WsHub.UnRegister(JobGroup, id)
	servers.OK(c, servers.WithData("ws close success"),
		servers.WithMsg("success"))
}

func SendGroup(group string, msg []byte) {
	WsHub.WriteGroup(group, websocket.TextMessage, []byte("{\"code\":200,\"data\":"+string(msg)+"}"))
}
