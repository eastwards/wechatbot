package handlers

import (
	"github.com/eatmoreapple/openwechat"
	"log"
	"main/config"
	"time"
)

// MessageHandlerInterface 消息处理接口
type MessageHandlerInterface interface {
	handle(*openwechat.Message) error
	ReplyText(*openwechat.Message) error
}

type HandlerType string

const (
	GroupHandler = "group"
	UserHandler  = "user"
)

// handlers 所有消息类型类型的处理器
var handlers map[HandlerType]MessageHandlerInterface

func init() {
	handlers = make(map[HandlerType]MessageHandlerInterface)
	handlers[GroupHandler] = NewGroupMessageHandler()
	handlers[UserHandler] = NewUserMessageHandler()
}

// Handler 全局处理入口
func Handler(msg *openwechat.Message) {
	log.Printf("hadler Received msg : %v", msg.Content)
	log.Printf("msg Status:%v, StatusNotifyCode:%v, CreateTime:%v, SystemNow:%v", msg.Status, msg.StatusNotifyCode, msg.CreateTime, time.Now().Unix())

	// 多线程处理
	go func() {

		// 如果已接受过，不再进行处理
		if msg.IsNotify() {
			log.Printf("msg was Notified. never explain.")
			return
		}

		// 如果消息时间超过10秒，不再处理
		if 10 < time.Now().Unix()-msg.CreateTime {
			log.Printf("msg was expired. never explain.")
			return
		}

		// 处理群消息
		if msg.IsSendByGroup() {
			handlers[GroupHandler].handle(msg)
			return
		}

		// 好友申请
		if msg.IsFriendAdd() {
			if config.LoadConfig().AutoPass {
				_, err := msg.Agree("你好我是基于chatGPT引擎开发的微信机器人，你可以向我提问任何问题。")
				if err != nil {
					log.Fatalf("add friend agree error : %v", err)
					return
				}
			}
		}

		// 私聊
		handlers[UserHandler].handle(msg)
	}()
}
