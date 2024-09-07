package client

import (
	"encoding/json"
	"errors"
	wss "github.com/gorilla/websocket"
	"github.com/tidwall/gjson"
	"log"
	"main/constant"
	"main/dto"
	"time"
)

type messageChan chan *dto.WSPayload
type closeErrorChan chan error

type SocketClient struct {
	version         int
	conn            *wss.Conn
	messageQueue    messageChan
	session         *dto.Session
	user            *dto.WSUser
	closeChan       closeErrorChan
	heartBeatTicker *time.Ticker // 用于维持定时心跳
}

func NewSocketClient(session dto.Session) *SocketClient {
	return &SocketClient{
		messageQueue:    make(messageChan, constant.DefaultQueueSize),
		session:         &session,
		closeChan:       make(closeErrorChan, 10),
		heartBeatTicker: time.NewTicker(60 * time.Second),
	}
}

func (c *SocketClient) Connect() error {
	var err error
	c.conn, _, err = wss.DefaultDialer.Dial(c.session.URL, nil)
	if err != nil {
		log.Fatal("%s, connect err: %v", c.session, err)
		return err
	}
	log.Printf("%s, url %s, connected", c.session, c.session.URL)

	return nil
}

func (c *SocketClient) Identify() error {
	// 避免传错 intent
	if c.session.Intent == 0 {
		c.session.Intent = constant.IntentGuilds
	}
	payload := &dto.WSPayload{
		Data: &dto.WSIdentityData{
			Token:   c.session.Token.GetString(),
			Intents: c.session.Intent,
		},
	}
	payload.OPCode = constant.Identify
	return c.Write(payload)
}

func (c *SocketClient) Resume() error {
	payload := &dto.WSPayload{
		Data: &dto.WSResumeData{
			Token:     c.session.Token.GetString(),
			SessionID: c.session.ID,
			Seq:       c.session.LastSeq,
		},
	}
	payload.OPCode = constant.Resume // 内嵌结构体字段，单独赋值
	return c.Write(payload)
}

func (c *SocketClient) Write(message *dto.WSPayload) error {
	m, _ := json.Marshal(message)
	log.Printf("%s write %s message, %v", c.session, constant.OPMeans(message.OPCode), string(m))

	if err := c.conn.WriteMessage(wss.TextMessage, m); err != nil {
		log.Fatalf("%s WriteMessage failed, %v", c.session, err)
		return err
	}
	return nil
}

func (c *SocketClient) Listening() error {
	defer func() {
		// 退出监听则关闭连接
		err := c.conn.Close()
		if err != nil {
			log.Fatalf(err.Error())
			return
		}
		// 关闭心跳定时器
		c.heartBeatTicker.Stop()
	}()
	// reading message
	go c.readMessageToQueue()
	// read message from queue and handle,in goroutine to avoid business logic block closeChan and heartBeatTicker
	go c.listenMessageAndHandle()

	for {
		select {
		case err := <-c.closeChan:
			log.Fatalf("%s Listening stop. err is %v", c.session, err)
			return err
		case <-c.heartBeatTicker.C:
			// 监听心跳包
			log.Printf("%s listened heartBeat", c.session)
			heartBeatEvent := &dto.WSPayload{
				OPCode: constant.Heartbeat,
				Data:   c.session.LastSeq,
			}
			// 不处理错误，Write 内部会处理，如果发生发包异常，会通知主协程退出
			_ = c.Write(heartBeatEvent)
		}
	}
	return nil
}

// startHeartBeatTicker 启动定时心跳
func (c *SocketClient) startHeartBeatTicker(message []byte) {
	helloData := &dto.WSHelloData{}
	data := gjson.Get(string(message), "d")
	if err := json.Unmarshal([]byte(data.String()), helloData); err != nil {
		log.Fatalf("%s hello data parse failed, %v, message %v", c.session, err, message)
	}
	// 根据 hello 的回包，重新设置心跳的定时器时间
	c.heartBeatTicker.Reset(time.Duration(helloData.HeartbeatInterval) * time.Millisecond)
}

func (c *SocketClient) isHandleBuildIn(payload *dto.WSPayload) bool {
	switch payload.OPCode {
	case constant.Hello: // 接收到 hello 后需要开始发心跳
		c.startHeartBeatTicker(payload.RawMessage)
	case constant.HeartbeatACK: // 心跳 ack 不需要业务处理
	case constant.Reconnect: // 达到连接时长，需要重新连接，此时可以通过 resume 续传原连接上的事件
		c.closeChan <- errors.New("重新连接")
	case constant.InvalidSession: // 无效的 sessionLog，需要重新鉴权
		c.closeChan <- errors.New("无效session")
	default:
		return false
	}
	return true
}

func (c *SocketClient) readMessageToQueue() {
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			log.Fatalf("%s read message failed, %v, message %s", c.session, err, string(message))
			close(c.messageQueue)
			return
		}
		payload := &dto.WSPayload{}
		if err := json.Unmarshal(message, payload); err != nil {
			log.Fatalf("%s json failed, %v", c.session, err)
			continue
		}
		payload.RawMessage = message
		log.Printf("%s receive %s message, %s", c.session, constant.OPMeans(payload.OPCode), string(message))
		// 处理内置的一些事件，如果处理成功，则这个事件不再投递给业务
		if c.isHandleBuildIn(payload) {
			continue
		}
		c.messageQueue <- payload
	}
}

func (c *SocketClient) listenMessageAndHandle() {
	// 循环读取消息列表，进行处理
	for {
		select {
		case payload := <-c.messageQueue:
			c.saveSeq(payload.Seq)
			// ready 事件需要特殊处理
			if payload.Type == "READY" {
				c.readyHandler(payload)
				continue
			}
			if err := HandlerProcess(payload.OPCode, payload.Type, payload); err != nil {
				log.Fatalf("%s HandlerProcess failed, %v", c.session, err)
			}
		}
	}
	log.Printf("%s message queue is closed", c.session)
}

func (c *SocketClient) saveSeq(seq uint32) {
	if seq > 0 {
		c.session.LastSeq = seq
	}
}

// readyHandler 针对ready返回的处理，需要记录 sessionID 等相关信息
func (c *SocketClient) readyHandler(payload *dto.WSPayload) {
	readyData := &dto.WSReadyData{}
	data := gjson.Get(string(payload.RawMessage), "d")
	if err := json.Unmarshal([]byte(data.String()), readyData); err != nil {
		log.Fatalf("%s parseReadyData failed, %v, message %v", c.session, err, payload.RawMessage)
	}
	c.version = readyData.Version
	// 基于 ready 事件，更新 session 信息
	c.session.ID = readyData.SessionID
	c.user = &dto.WSUser{
		ID:       readyData.User.ID,
		Username: readyData.User.Username,
		Bot:      readyData.User.Bot,
	}
}
