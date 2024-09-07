package client

import (
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"log"
	"main/conf"
	"main/constant"
	"time"
)

// 沙箱网关,如果是生产环境，需要更换
var testGetWay = "https://sandbox.api.sgroup.qq.com"

/*
返回的ws url

	{
	  "url": "wss://api.sgroup.qq.com/websocket/"
	}
*/
type WsUrl struct {
	URL string `json:"url"`
}

type HttpClient struct {
	client *resty.Client
}

func (h *HttpClient) NewHttpClient(token *conf.Token) {
	// https://bot.q.qq.com/wiki/develop/api/#%E6%8E%A5%E5%8F%A3%E8%AF%B4%E6%98%8E
	// 使用申请机器人时平台返回的机器人 appID + token 拼接而成。此时，所有的操作都是以机器人身份来完成的
	// example: Authorization: Bot 100000.Cl2FMQZnCjm1XVW7vRze4b7Cq4se7kKWs
	authReqHeader := token.GetString()
	h.client = resty.New().
		SetAuthToken(authReqHeader).
		SetAuthScheme("Bot").
		SetTimeout(10 * time.Second)
	h.client.SetBaseURL(testGetWay)
}

func (h *HttpClient) GetWebSocket(url string) string {
	// https://bot.q.qq.com/wiki/develop/api-v2/openapi/wss/url_get.html
	var wsUrl = ""
	resp, err := h.client.R().Get(url)
	if err != nil {
		log.Fatalf("请求错误: %v", err)
		return ""
	}
	if resp != nil && resp.StatusCode() == constant.SUCCESS {
		var wsUrlStruct WsUrl
		if err := json.Unmarshal(resp.Body(), &wsUrlStruct); err != nil {
			log.Fatalf("GetWebSocket.Unmarshal has err: %v", err)
			return ""
		}
		wsUrl = wsUrlStruct.URL
	}

	log.Printf("GetWebSocket.url:%s", wsUrl)
	return wsUrl
}

func (h *HttpClient) Post(method string, paramKey string, paramValue string, body interface{}) error {
	// 发送请求
	resp, err := h.client.R().
		SetPathParam(paramKey, paramValue).
		SetBody(body).
		Post(method)
	// 检查错误
	if err != nil {
		log.Fatalf("Post has err: %v", err)
		return err
	}
	// 打印响应状态码和内容
	log.Println("Post.响应状态码:", resp.StatusCode())
	log.Println("Post.响应体:", resp.String())

	return nil
}
