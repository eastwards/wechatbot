package handlers

import (
	"github.com/eatmoreapple/openwechat"
	"log"
	"main/gtp"
	"strings"
)

var _ MessageHandlerInterface = (*UserMessageHandler)(nil)

// UserMessageHandler 私聊消息处理
type UserMessageHandler struct {
}

// handle 处理消息
func (g *UserMessageHandler) handle(msg *openwechat.Message) error {
	if msg.IsText() {
		return g.ReplyText(msg)
	}
	return nil
}

// NewUserMessageHandler 创建私聊处理器
func NewUserMessageHandler() MessageHandlerInterface {
	return &UserMessageHandler{}
}

// ReplyText 发送文本消息到群
func (g *UserMessageHandler) ReplyText(msg *openwechat.Message) error {
	// 接收私聊消息
	sender, err := msg.Sender()
	log.Printf("Received User %v Text Msg : %v", sender.NickName, msg.Content)

	quoteText := "「" + sender.Self.NickName
	requestText := msg.Content

	if strings.HasPrefix(msg.Content, quoteText) {
		stringSlice := strings.Split(msg.Content, "=>")
		requestText = strings.ReplaceAll(stringSlice[1], "\n- - - - - - - - - - - - - - -\n", " \r\n ")
		requestText = strings.ReplaceAll(requestText, "」", "")
	} else {
		// 向GPT发起请求
		requestText = strings.TrimSpace(requestText)
		requestText = strings.Trim(requestText, "\n")
	}

	reply, err := gtp.Completions(requestText)
	if err != nil {
		log.Printf("gtp request error: %v \n", err)
		msg.ReplyText("机器人神了，我一会发现了就去修。")
		return err
	}
	if reply == "" {
		return nil
	}

	// 回复用户
	reply = strings.TrimSpace(reply)
	reply = strings.Trim(reply, "\n")
	reply = " => " + requestText + " \r\n " + reply
	_, err = msg.ReplyText(reply)
	if err != nil {
		log.Printf("response user error: %v \n", err)
	}
	return err
}
