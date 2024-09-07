package dto

// ReplyMessage 发送消息结构体定义
type ReplyMessage struct {
	Content          string            `json:"content,omitempty"`
	MsgID            string            `json:"msg_id,omitempty"` // 要回复的消息id
	MessageReference *MessageReference `json:"message_reference,omitempty"`
}
