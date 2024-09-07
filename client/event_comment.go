package client

import (
	"log"
	"main/constant"
	"main/dto"
)

var eventHandlerMap = map[constant.OPCode]map[constant.EventType]EventHandler{}

type EventHandler func(event *dto.WSPayload, message []byte) error

func RegisterHandler(code constant.OPCode, eventT constant.EventType, handler EventHandler) {
	// 添加到map中
	if _, exists := eventHandlerMap[code]; !exists {
		eventHandlerMap[code] = make(map[constant.EventType]EventHandler)
	}
	eventHandlerMap[code][eventT] = handler
}

func HandlerProcess(code constant.OPCode, eventT constant.EventType, payload *dto.WSPayload) error {
	if tempMap, ok := eventHandlerMap[code]; ok {
		if handler, ok1 := tempMap[eventT]; ok1 {
			// 调用 eventHandler
			handler(payload, payload.RawMessage)
		} else {
			log.Printf("没有添加该事件的处理器,opCode:%d eventType:%s", code, eventT)
		}
	}
	return nil
}
