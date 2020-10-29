package easyws

import (
	"context"
	"errors"
	"sync"

	"github.com/gorilla/websocket"
)

// 错误返回
var (
	ErrHubClosed       = errors.New("ws: hub is closed")
	ErrSessionNotFound = errors.New("ws: session not found")
	ErrBadmsgType      = errors.New("ws: bad write tMessage type")
)

// 广播/组播 消息体
type tMessage struct {
	isBroadcast bool // 广播/组播
	groupID     string
	msgType     int
	data        []byte
}

// message session 消息体
type message struct {
	msgType int
	data    []byte
}

// Hub 管理中心
type Hub struct {
	tMessage chan *tMessage
	ctx      context.Context
	cancel   context.CancelFunc
	config

	// 以下需要持锁
	mu       sync.Mutex
	sessions map[string]map[string]*Session
	groupCnt int
	sessCnt  int
}

// New 创建管理中心
func New(opts ...Option) *Hub {
	ctx, cancel := context.WithCancel(context.Background())
	hub := &Hub{
		sessions: make(map[string]map[string]*Session),
		config:   defaultConfig(),
		ctx:      ctx,
		cancel:   cancel,
	}

	for _, opt := range opts {
		opt(hub)
	}

	hub.tMessage = make(chan *tMessage, hub.MessageBufferSize<<2)
	return hub
}

func (sf *Hub) register(sess *Session) {
	sf.mu.Lock()
	if sf.sessions[sess.GroupID] == nil {
		sf.sessions[sess.GroupID] = make(map[string]*Session)
		sf.groupCnt++
	}
	if oldSess, ok := sf.sessions[sess.GroupID][sess.ID]; ok {
		oldSess.cancel()
	} else {
		sf.sessCnt++
	}
	sf.sessions[sess.GroupID][sess.ID] = sess
	sf.mu.Unlock()
}

func (sf *Hub) UnRegister(groupID, id string) {
	sf.mu.Lock()
	if group, ok := sf.sessions[groupID]; ok {
		if sess, ok := group[id]; ok {
			delete(group, id)
			sf.sessCnt--
			if len(group) == 0 {
				delete(sf.sessions, groupID)
				sf.groupCnt--
			}
			sess.cancel()
		}
	}
	sf.mu.Unlock()
}

// Run 运行管理中心
func (sf *Hub) Run(ctx context.Context) {
	defer func() {
		sf.mu.Lock()
		for _, group := range sf.sessions {
			for _, sess := range group {
				sess.cancel()
			}
		}
		sf.sessions = make(map[string]map[string]*Session)
		sf.mu.Unlock()
	}()

	for {
		select {
		case msg := <-sf.tMessage:
			sf.mu.Lock()
			if msg.isBroadcast {
				group, ok := sf.sessions[msg.groupID]
				if ok {
					for _, sess := range group {
						sess.WriteMessage(msg.msgType, msg.data) // nolint: errcheck
					}
				}
			} else {
				for _, group := range sf.sessions {
					for _, sess := range group {
						sess.WriteMessage(msg.msgType, msg.data) // nolint: errcheck
					}
				}
			}
			sf.mu.Unlock()
		case <-ctx.Done():
			sf.cancel() // local cancel mark it closed
			return
		case <-sf.ctx.Done():
			return
		}
	}
}

// WriteBroadcast 广播消息 (websocket.TextMessage, websocket.BinaryMessage)
func (sf *Hub) WriteBroadcast(msgType int, data []byte) error {
	if !(msgType == websocket.TextMessage || msgType == websocket.BinaryMessage) {
		return ErrBadmsgType
	}
	select {
	case sf.tMessage <- &tMessage{true, "", msgType, data}:
		return nil
	case <-sf.ctx.Done():
		return ErrHubClosed
	}
}

// WriteGroup 组播 (websocket.TextMessage, websocket.BinaryMessage)
func (sf *Hub) WriteGroup(groupID string, msgType int, data []byte) error {
	if !(msgType == websocket.TextMessage || msgType == websocket.BinaryMessage) {
		return ErrBadmsgType
	}
	select {
	case sf.tMessage <- &tMessage{false, groupID, msgType, data}:
		return nil
	case <-sf.ctx.Done():
		return ErrHubClosed
	}
}

// WriteMessage 单播 (websocket.TextMessage, websocket.BinaryMessage)
func (sf *Hub) WriteMessage(groupID, id string, msgType int, data []byte) error {
	if !(msgType == websocket.TextMessage || msgType == websocket.BinaryMessage) {
		return ErrBadmsgType
	}
	select {
	case <-sf.ctx.Done():
		return ErrHubClosed
	default:
	}
	sf.mu.Lock()
	sess, err := sf.findSession(groupID, id)
	sf.mu.Unlock()
	if err != nil {
		return err
	}
	return sess.WriteMessage(msgType, data)
}

// WriteControl 单播控制信息
func (sf *Hub) WriteControl(groupID, id string, msgType int, data []byte) error {
	sf.mu.Lock()
	sess, err := sf.findSession(groupID, id)
	sf.mu.Unlock()
	if err != nil {
		return err
	}
	return sess.WriteControl(msgType, data)
}

// SessionLen 客户端会话的数量
func (sf *Hub) SessionLen() (count int) {
	sf.mu.Lock()
	count = sf.sessCnt
	sf.mu.Unlock()
	return
}

// GroupLen 组个数
func (sf *Hub) GroupLen() (count int) {
	sf.mu.Lock()
	count = sf.groupCnt
	sf.mu.Unlock()
	return
}

// Close 关闭
func (sf *Hub) Close() {
	sf.cancel()
}

// IsClosed 判断是否关闭
func (sf *Hub) IsClosed() bool {
	return sf.ctx.Err() != nil
}

func (sf *Hub) findSession(groupID, id string) (*Session, error) {
	if group, ok := sf.sessions[groupID]; ok {
		if sess, ok := group[id]; ok {
			return sess, nil
		}
	}
	return nil, ErrSessionNotFound
}
