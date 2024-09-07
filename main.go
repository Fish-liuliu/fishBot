package main

import (
	"encoding/json"
	"github.com/tidwall/gjson"
	"log"
	"main/client"
	"main/conf"
	"main/constant"
	"main/dto"
	"main/service"
)

var RobotBusiness service.RobotBussiness

func main() {
	// 获取token
	token, err := conf.GetToken()
	if err != nil {
		log.Fatal(err)
		return
	}

	// 创建http连接
	httpClient := &client.HttpClient{}
	httpClient.NewHttpClient(token)

	// 获取ws连接
	webSocket := httpClient.GetWebSocket(constant.GetWebSocketURI)

	// 获取数据库连接
	db, err := conf.GetGormConnect()
	if err != nil {
		panic(err)
	}
	cache, err := conf.GetRedisConnect()
	if err != nil {
		panic(err)
	}
	bussiness := service.NewRobotBussiness(httpClient, db, cache)

	// 订阅事件
	intent := constant.IntentGuildAtMessage
	// 注册 at消息 处理器
	client.RegisterHandler(constant.Dispatch, constant.EventAtMessageCreate, AtMessageHandler(bussiness))
	// 开启websocket
	client.New().Start(webSocket, token, intent)
}

// AtMessageHandler 艾特消息处理器
func AtMessageHandler(bussiness *service.RobotBussiness) client.EventHandler {
	return func(event *dto.WSPayload, message []byte) error {
		data := &dto.Message{}
		if err := ParseData(message, data); err != nil {
			return err
		}
		if err := bussiness.AtRoBot(data); err != nil {
			return err
		}
		return nil
	}
}

func ParseData(message []byte, target interface{}) error {
	data := gjson.Get(string(message), "d")
	return json.Unmarshal([]byte(data.String()), target)
}
