package ws

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/x-tardis/go-admin/pkg/servers"
)

// 初始化 wsManager 管理器
var WebsocketManager = Manager{
	groups:      make(map[string]map[string]*Client),
	register:    make(chan *Client, 128),
	unRegister:  make(chan UnRegister, 128),
	message:     make(chan *Message, 512),
	groupCount:  0,
	clientCount: 0,
}

// gin 处理 websocket handler
func (sf *Manager) WsClient(c *gin.Context) {
	log.Println("-------------> ws register")
	upGrader := websocket.Upgrader{
		// cross origin domain
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		// 处理 Sec-WebSocket-Protocol Header
		Subprotocols: []string{c.GetHeader("Sec-WebSocket-Protocol")},
	}

	conn, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("websocket connect error: %s", c.Param("channel"))
		return
	}

	fmt.Println("token: ", c.Query("token"))

	id, channel := c.Param("id"), c.Param("channel")
	ctx, cancel := context.WithCancel(context.Background())
	client := &Client{
		Id:         id,
		Group:      channel,
		Context:    ctx,
		CancelFunc: cancel,
		Socket:     conn,
		Message:    make(chan []byte, 1024),
	}

	sf.RegisterClient(client)
	go client.Read(ctx)
	go client.Write(ctx)
	time.Sleep(time.Second * 15)

	FileMonitoringById(ctx, "temp/job.log", id, channel, SendSingle)
}

func (*Manager) UnWsClient(c *gin.Context) {
	id := c.Param("id")
	group := c.Param("channel")
	WebsocketManager.UnRegisterClient(group, id)
	fmt.Println(WebsocketManager.Info())
	servers.OK(c,
		servers.WithData("ws close success"),
		servers.WithMsg("success"))
}

func SendGroup(msg []byte) {
	WebsocketManager.SendGroup("leffss", []byte("{\"code\":200,\"data\":"+string(msg)+"}"))
	fmt.Println(WebsocketManager.Info())
}

func SendBroadcast(msg []byte) {
	WebsocketManager.SendBroadcast([]byte("{\"code\":200,\"data\":" + string(msg) + "}"))
	fmt.Println(WebsocketManager.Info())
}

func SendSingle(ctx context.Context, id string, group string, msg []byte) {
	WebsocketManager.SendSingle(id, group, []byte("{\"code\":200,\"data\":"+string(msg)+"}"))
	fmt.Println(WebsocketManager.Info())
}
