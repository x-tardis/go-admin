package ws

import (
	"context"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

// 消息类型
const (
	MsgTypeSingle    = iota // 单体消息
	MsgTypeGroup            // 组播消息
	MsgTypeBroadcast        // 广播消息
)

// Manager 所有 websocket 信息
type Manager struct {
	register   chan *Client    // 注册
	unRegister chan UnRegister // 注销
	message    chan *Message   // 消息

	// 以下需要持锁
	mu          sync.Mutex
	groups      map[string]map[string]*Client
	groupCount  uint
	clientCount uint
}

// Client 单个 websocket 信息
type Client struct {
	Id         string
	Group      string
	Context    context.Context
	CancelFunc context.CancelFunc
	Socket     *websocket.Conn
	Message    chan []byte

	manager *Manager
}

// 取消注册信息
type UnRegister struct {
	group string
	id    string
}

type Message struct {
	Type  int
	Id    string
	Group string
	Data  []byte
}

// 读信息，从 websocket 连接直接读取数据
func (c *Client) Read(ctx context.Context) {
	defer func() {
		c.manager.UnRegisterClient(c.Group, c.Id)
		log.Printf("client [%s] disconnect", c.Id)
		if err := c.Socket.Close(); err != nil {
			log.Printf("client [%s] disconnect err: %s", c.Id, err)
		}
	}()

	for {
		if ctx.Err() != nil {
			break
		}
		messageType, message, err := c.Socket.ReadMessage()
		if err != nil || messageType == websocket.CloseMessage {
			break
		}
		log.Printf("client [%s] receive message: %s", c.Id, string(message))
		c.Message <- message
	}
}

// 写信息，从 channel 变量 SendSingle 中读取数据写入 websocket 连接
func (c *Client) Write(ctx context.Context) {
	defer func() {
		log.Printf("client [%s] disconnect", c.Id)
		if err := c.Socket.Close(); err != nil {
			log.Printf("client [%s] disconnect err: %s", c.Id, err)
		}
	}()

	for {
		if ctx.Err() != nil {
			break
		}
		select {
		case message, ok := <-c.Message:
			if !ok {
				_ = c.Socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			log.Printf("client [%s] write message: %s", c.Id, string(message))
			err := c.Socket.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				log.Printf("client [%s] writemessage err: %s", c.Id, err)
			}
		case <-c.Context.Done():
			break
		}
	}
}

// 启动 websocket 管理器
func (sf *Manager) Start() {
	log.Printf("websocket manage start")
	go sf.SendMessageLoop()

	for {
		select {
		case client := <-sf.register: // 注册
			log.Printf("register: client [%s] to group [%s]", client.Id, client.Group)

			sf.mu.Lock()
			if sf.groups[client.Group] == nil {
				sf.groups[client.Group] = make(map[string]*Client)
				sf.groupCount++
			}
			sf.groups[client.Group][client.Id] = client
			sf.clientCount++
			sf.mu.Unlock()

		case ur := <-sf.unRegister: // 注销
			log.Printf("unregister client [%s] from group [%s]", ur.id, ur.group)

			sf.mu.Lock()
			if group, ok := sf.groups[ur.group]; ok {
				if client, ok := group[ur.id]; ok {
					delete(group, ur.id)
					sf.clientCount--
					if len(group) == 0 {
						delete(sf.groups, ur.group)
						sf.groupCount--
					}
					close(client.Message)
					client.CancelFunc()
				}
			}
			sf.mu.Unlock()
		}
	}
}

func (sf *Manager) SendMessageLoop() {
	for msg := range sf.message {
		switch msg.Type {
		case MsgTypeSingle:
			if groupMap, ok := sf.groups[msg.Group]; ok {
				if conn, ok := groupMap[msg.Id]; ok {
					conn.Message <- msg.Data
				}
			}
		case MsgTypeGroup:
			if groups, ok := sf.groups[msg.Group]; ok {
				for _, conn := range groups {
					conn.Message <- msg.Data
				}
			}
		case MsgTypeBroadcast:
			// TODO: 分发到组会更快??
			for _, v := range sf.groups {
				for _, conn := range v {
					conn.Message <- msg.Data
				}
			}
		}
	}
}

// SendSingle 发送单体消息
func (sf *Manager) SendSingle(id string, group string, message []byte) {
	sf.message <- &Message{
		Type:  MsgTypeSingle,
		Id:    id,
		Group: group,
		Data:  message,
	}
}

// SendGroup 发送组播消息
func (sf *Manager) SendGroup(group string, message []byte) {
	sf.message <- &Message{Type: MsgTypeGroup, Group: group, Data: message}
}

// SendBroadcast 发送广播消息
func (sf *Manager) SendBroadcast(message []byte) {
	sf.message <- &Message{Type: MsgTypeBroadcast, Data: message}
}

// 注册
func (sf *Manager) RegisterClient(client *Client) {
	client.manager = sf
	sf.register <- client
}

// 注销
func (sf *Manager) UnRegisterClient(group, id string) {
	sf.unRegister <- UnRegister{group, id}
}

// 当前组个数
func (sf *Manager) GroupLen() uint {
	return sf.groupCount
}

// 当前连接个数
func (sf *Manager) ClientLen() uint {
	return sf.clientCount
}

// 获取 wsManager 管理器信息
func (sf *Manager) Info() map[string]interface{} {
	return map[string]interface{}{
		"groupLen":          sf.groupCount,
		"clientLen":         sf.clientCount,
		"chanRegisterLen":   len(sf.register),
		"chanUnregisterLen": len(sf.unRegister),
		"chanMessageLen":    len(sf.message),
	}
}
