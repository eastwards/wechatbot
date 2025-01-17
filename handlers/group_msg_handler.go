package handlers

import (
	"github.com/eatmoreapple/openwechat"
	"log"
	"main/gtp"
	"strings"
)

var _ MessageHandlerInterface = (*GroupMessageHandler)(nil)

// GroupMessageHandler 群消息处理
type GroupMessageHandler struct {
}

// handle 处理消息
func (g *GroupMessageHandler) handle(msg *openwechat.Message) error {
	if msg.IsText() {
		return g.ReplyText(msg)
	}
	return nil
}

// NewGroupMessageHandler 创建群消息处理器
func NewGroupMessageHandler() MessageHandlerInterface {
	return &GroupMessageHandler{}
}

// ReplyText 发送文本消息到群
func (g *GroupMessageHandler) ReplyText(msg *openwechat.Message) error {
	// 接收群消息
	sender, err := msg.Sender()
	group := openwechat.Group{sender}
	log.Printf("Received Group %v Text Msg : %v", group.NickName, msg.Content)

	quoteText := "「" + sender.Self.NickName

	// 不是@的不处理或者不是引用机器人说的话
	if !msg.IsAt() && !strings.HasPrefix(msg.Content, quoteText) {
		return nil
	}

	requestText := ""

	if msg.IsAt() {
		// 替换掉@文本，然后向GPT发起请求
		replaceText := "@" + sender.Self.NickName
		requestText = strings.TrimSpace(strings.ReplaceAll(msg.Content, replaceText, ""))
	} else {
		stringSlice := strings.Split(msg.Content, "=>")
		requestText = strings.ReplaceAll(stringSlice[1], "\n- - - - - - - - - - - - - - -\n", " \r\n ")
		requestText = strings.ReplaceAll(requestText, "」", "")
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

	// 获取@我的用户
	groupSender, err := msg.SenderInGroup()
	if err != nil {
		log.Printf("get sender in group error :%v \n", err)
		return err
	}

	// 回复@我的用户
	reply = strings.TrimSpace(reply)
	reply = strings.Trim(reply, "\n")
	atText := "@" + groupSender.NickName
	replyText := atText + " => " + requestText + " \r\n " + reply
	_, err = msg.ReplyText(replyText)
	if err != nil {
		log.Printf("response group error: %v \n", err)
	}
	return err
}
