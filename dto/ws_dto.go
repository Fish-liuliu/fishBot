package dto

import (
	"main/conf"
	"main/constant"
)

type Session struct {
	ID      string
	URL     string
	Token   conf.Token
	Intent  constant.Intent
	LastSeq uint32
	Cnt     int
}

// WSUser 当前连接的用户信息
type WSUser struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Bot      bool   `json:"bot"`
}

/*
鉴权 example
{
  "op": 0,
  "s": 1,
  "t": "READY",
  "d": {
    "version": 1,
    "session_id": "082ee18c-0be3-491b-9d8b-fbd95c51673a",
    "user": {
      "id": "6158788878435714165",
      "username": "群pro测试机器人",
      "bot": true
    },
    "shard": [0, 0]
  }
}

*/

// WSPayload websocket 消息结构
type WSPayload struct {
	OPCode     constant.OPCode    `json:"op"`
	Seq        uint32             `json:"s,omitempty"`
	Type       constant.EventType `json:"t,omitempty"`
	Data       interface{}        `json:"d,omitempty"`
	RawMessage []byte             `json:"-"` // 原始的 message 数据
	// shard
}

// WSIdentityData 鉴权数据
type WSIdentityData struct {
	Token      string          `json:"token"`
	Intents    constant.Intent `json:"intents"`
	Shard      []uint32        `json:"shard"` // array of two integers (shard_id, num_shards)
	Properties struct {
		Os      string `json:"$os,omitempty"`
		Browser string `json:"$browser,omitempty"`
		Device  string `json:"$device,omitempty"`
	} `json:"properties,omitempty"`
}

// WSResumeData 重连数据
type WSResumeData struct {
	Token     string `json:"token"`
	SessionID string `json:"session_id"`
	Seq       uint32 `json:"seq"`
}

// WSHelloData hello 返回
type WSHelloData struct {
	HeartbeatInterval int `json:"heartbeat_interval"`
}

// WSReadyData ready，鉴权后返回
type WSReadyData struct {
	Version   int    `json:"version"`
	SessionID string `json:"session_id"`
	User      struct {
		ID       string `json:"id"`
		Username string `json:"username"`
		Bot      bool   `json:"bot"`
	} `json:"user"`
	Shard []uint32 `json:"shard"`
}
