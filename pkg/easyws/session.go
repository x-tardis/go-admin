package easyws

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
)

// Session 会话
type Session struct {
	GroupID  string
	ID       string
	Conn     *websocket.Conn
	Request  *http.Request
	alive    int32
	lctx     context.Context
	cancel   context.CancelFunc
	outBound chan *message

	Hub *Hub
}

// WriteMessage 写消息 (websocket.TextMessage, websocket.BinaryMessage)
func (sf *Session) WriteMessage(msgType int, data []byte) error {
	if !(msgType == websocket.TextMessage || msgType == websocket.BinaryMessage) {
		return ErrBadmsgType
	}
	sf.outBound <- &message{msgType, data}
	return nil
}

// WriteControl 写控制消息 (websocket.CloseMessage, websocket.PingMessage and websocket.PongMessage.)
func (sf *Session) WriteControl(msgType int, data []byte) error {
	return sf.Conn.WriteControl(msgType, data,
		time.Now().Add(sf.Hub.SessionConfig.WriteTimeout))
}

// Run
func (sf *Session) Run() {
	cfg := sf.Hub.SessionConfig

	sf.outBound = make(chan *message, cfg.MessageBufferSize)
	sf.lctx, sf.cancel = context.WithCancel(sf.Hub.ctx)
	sf.Hub.register(sf)
	sf.Hub.connectHandler(sf)
	defer func() {
		sf.Conn.Close()
		sf.cancel()
		sf.Hub.UnRegister(sf.GroupID, sf.ID)
		sf.Hub.disconnectHandler(sf)
	}()
	go sf.writePump()

	readTimeout := cfg.KeepAlive * time.Duration(cfg.Ratio) / 100 * 4

	// 设置 pong handler
	sf.Conn.SetPongHandler(func(message string) error {
		atomic.StoreInt32(&sf.alive, 0)
		sf.Conn.SetReadDeadline(time.Now().Add(readTimeout))
		sf.Hub.pongHandler(sf, message)
		return nil
	})

	// 设置 ping handler
	sf.Conn.SetPingHandler(func(message string) error {
		atomic.StoreInt32(&sf.alive, 0)
		sf.Conn.SetReadDeadline(time.Now().Add(readTimeout))
		err := sf.Conn.WriteControl(websocket.PongMessage, []byte(message), time.Now().Add(cfg.WriteTimeout))
		if err != nil {
			if e, ok := err.(net.Error); !(ok && e.Temporary() || err == websocket.ErrCloseSent) {
				return err
			}
		}
		sf.Hub.pingHandler(sf, message)
		return nil
	})

	sf.Conn.SetCloseHandler(func(code int, text string) error {
		return sf.Hub.closeHandler(sf, code, text)
	})

	if cfg.MaxMessageSize > 0 {
		sf.Conn.SetReadLimit(cfg.MaxMessageSize)
	}

	for {
		sf.Conn.SetReadDeadline(time.Now().Add(readTimeout))
		msgType, data, err := sf.Conn.ReadMessage()
		if err != nil {
			sf.Hub.errorHandler(sf, fmt.Errorf("read message %w", err))
			return
		}
		atomic.StoreInt32(&sf.alive, 0)
		sf.Hub.receiveHandler(sf, msgType, data)
	}
}

// writePump
func (sf *Session) writePump() {
	var retries int

	cfg := sf.Hub.SessionConfig
	monTick := time.NewTicker(cfg.KeepAlive * time.Duration(cfg.Ratio) / 100)
	defer func() {
		monTick.Stop()
		sf.Conn.Close()
	}()
	for {
		select {
		case <-sf.lctx.Done():
			return
		case msg := <-sf.outBound:
			err := sf.Conn.WriteMessage(msg.msgType, msg.data)
			if err != nil {
				sf.Hub.errorHandler(sf, fmt.Errorf("ws: write message %w", err))
				return
			}
			sf.Hub.sendHandler(sf, msg.msgType, msg.data)
		case <-monTick.C:
			if atomic.AddInt32(&sf.alive, 1) > 1 {
				retries++
				if retries > 3 {
					sf.Hub.errorHandler(sf, errors.New("ws: keep alive failed"))
					return
				}
				err := sf.Conn.WriteControl(websocket.PingMessage, []byte{},
					time.Now().Add(cfg.WriteTimeout))
				if err != nil {
					sf.Hub.errorHandler(sf, fmt.Errorf("ws: write control %w", err))
					return
				}
			} else {
				retries = 0
			}
		}
	}
}
