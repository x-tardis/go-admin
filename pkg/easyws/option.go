package easyws

import (
	"time"
)

// SessionConfig 会话配置
type SessionConfig struct {
	WriteTimeout      time.Duration // 写超时时间
	KeepAlive         time.Duration // 保活时间
	Ratio             int           // 监控比例, 需大于100,默认系统是110 即比例值1.1
	MaxMessageSize    int64         // 消息最大字节数, 如果为0,使用系统默认设置
	MessageBufferSize int           // 消息缓存数
}

// config 配置
type config struct {
	SessionConfig
	connectHandler    func(sess *Session)
	disconnectHandler func(sess *Session)
	pingHandler       func(sess *Session, msg string)
	pongHandler       func(sess *Session, msg string)
	sendHandler       func(sess *Session, msgType int, data []byte)
	receiveHandler    func(sess *Session, msgType int, data []byte)
	closeHandler      func(sess *Session, code int, text string) error
	errorHandler      func(sess *Session, err error)
}

// defaultConfig 创建默认选项
func defaultConfig() config {
	return config{
		SessionConfig{
			WriteTimeout:      1 * time.Second,
			KeepAlive:         60 * time.Second,
			Ratio:             110,
			MaxMessageSize:    0,
			MessageBufferSize: 128,
		},

		func(sess *Session) {},
		func(sess *Session) {},
		func(sess *Session, msg string) {},
		func(sess *Session, msg string) {},
		func(sess *Session, msgType int, data []byte) {},
		func(sess *Session, msgType int, data []byte) {},
		func(sess *Session, code int, text string) error { return nil },
		func(sess *Session, err error) {},
	}
}

type Option func(hub *Hub)

// WithSessionConfig 设置会话配置
func WithSessionConfig(cfg *SessionConfig) Option {
	return func(hub *Hub) {
		hub.SessionConfig = *cfg
	}
}

// WithConnectHandler 设置连接回调
func WithConnectHandler(f func(sess *Session)) Option {
	return func(hub *Hub) {
		hub.connectHandler = f
	}
}

// WithDisconnectHandler 设置断开连接回调
func WithDisconnectHandler(f func(sess *Session)) Option {
	return func(hub *Hub) {
		hub.disconnectHandler = f
	}
}

// WithPingHandler 设置收到Ping回调
func WithPingHandler(f func(sess *Session, str string)) Option {
	return func(hub *Hub) {
		hub.pingHandler = f
	}
}

// WithPongHandler 设置收到Pong回调
func WithPongHandler(f func(sess *Session, str string)) Option {
	return func(hub *Hub) {
		hub.pongHandler = f
	}
}

// WithSendHandler 设置发送回调
func WithSendHandler(f func(sess *Session, msgType int, data []byte)) Option {
	return func(hub *Hub) {
		hub.sendHandler = f
	}
}

// WithReceiveHandler 设置接收回调
func WithReceiveHandler(f func(sess *Session, msgType int, data []byte)) Option {
	return func(hub *Hub) {
		hub.receiveHandler = f
	}
}

// WithCloseHandler 设置发送回调
func WithCloseHandler(f func(sess *Session, code int, text string) error) Option {
	return func(hub *Hub) {
		hub.closeHandler = f
	}
}

// SetReceiveHandler 设置接收回调
func WithErrorHandler(f func(sess *Session, err error)) Option {
	return func(hub *Hub) {
		hub.errorHandler = f
	}
}
