package client

import (
	"log"
	"main/conf"
	"main/constant"
	"main/dto"
)

const retryCnt = 10

type Container struct {
	sessionChan chan dto.Session
}

func New() *Container {
	return &Container{}
}

func (c *Container) Start(webSocketIp string, token *conf.Token, intents constant.Intent) error {
	c.sessionChan = make(chan dto.Session, 1)
	sessionTemp := dto.Session{
		URL:     webSocketIp,
		Token:   *token,
		Intent:  intents,
		LastSeq: 0,
		Cnt:     0,
	}
	c.sessionChan <- sessionTemp
	for session := range c.sessionChan {
		// 异常重连
		if session.Cnt >= retryCnt {
			// 重试次数超过10 则返回空
			log.Fatalf("重试次数超过默认次数次: %d 默认退出！", session.Cnt)
			return nil
		}
		session.Cnt++
		go c.newConnect(session)
	}
	return nil
}

func (c *Container) newConnect(session dto.Session) {
	defer func() {
		// 有异常则将session放回chan中，进行重连
		if err := recover(); err != nil {
			c.sessionChan <- session
		}
	}()
	socketClient := NewSocketClient(session)
	if err := socketClient.Connect(); err != nil {
		log.Fatalf("newConnect.Connect has err:%v", err)
		c.sessionChan <- session // 连接失败，丢回去队列排队重连
		return
	}
	// 连接成功则将次数置为 0
	session.Cnt = 0
	var err error
	// 如果 session id 不为空，则执行的是 resume 操作，如果为空，则执行的是 identify 操作
	if session.ID != "" {
		err = socketClient.Resume()
	} else {
		// 初次鉴权
		err = socketClient.Identify()
	}
	if err != nil {
		log.Fatalf("[ws/session] Identify/Resume err %+v", err)
		return
	}
	if err := socketClient.Listening(); err != nil {
		currentSession := socketClient.session
		// 将 session 放到 session chan 中，用于启动新的连接，当前连接退出
		c.sessionChan <- *currentSession
	}
}
