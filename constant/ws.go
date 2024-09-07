package constant

// https://bot.q.qq.com/wiki/develop/api-v2/dev-prepare/interface-framework/event-emit.html#websocket-%E6%96%B9%E5%BC%8F

// initent 事件
/*
PUBLIC_GUILD_MESSAGES (1 << 30) // 消息事件，此为公域的消息事件
  - AT_MESSAGE_CREATE       // 当收到@机器人的消息时
  - PUBLIC_MESSAGE_DELETE   // 当频道的消息被删除时
*/
type Intent int

const (
	IntentGuilds         Intent = 1 << iota
	IntentGuildAtMessage Intent = 1 << 30 // 只接收@消息事件
)

// 长连接维护 OpCode
type OPCode int

const (
	Dispatch OPCode = iota
	Heartbeat
	Identify
	_ // Presence Update
	_ // Voice State Update
	_
	Resume
	Reconnect
	_ // Request Guild Members
	InvalidSession
	Hello
	HeartbeatACK
	HTTPCallbackAck
)

var opMeans = map[OPCode]string{
	Dispatch:        "Dispatch",       // 服务端进行消息推送
	Heartbeat:       "Heartbeat",      // 客户端或服务端发送心跳
	Identify:        "Identity",       // 客户端发送鉴权
	Resume:          "Resume",         // 客户端恢复连接
	Reconnect:       "Reconnect",      //  服务端通知客户端重新连接
	InvalidSession:  "InvalidSession", //  当 identify 或 resume 的时候，如果参数有错，服务端会返回该消息
	Hello:           "Hello",          // 	当客户端与网关建立 ws 连接之后，网关下发的第一条消息
	HeartbeatACK:    "HeartbeatAck",   // 	当发送心跳成功之后，就会收到该消息
	HTTPCallbackAck: "HeartbeatAck",   // 仅用于 http 回调模式的回包，代表机器人收到了平台推送的数据
}

func OPMeans(op OPCode) string {
	means, ok := opMeans[op]
	if !ok {
		means = "unknown"
	}
	return means
}

// EventType 事件类型
type EventType string

const (
	EventAtMessageCreate EventType = "AT_MESSAGE_CREATE"
)

const DefaultQueueSize = 10000
