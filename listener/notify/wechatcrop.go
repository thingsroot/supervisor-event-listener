package notify

import (
	"regexp"
	"errors"
	"fmt"
	"os"
	"github.com/yanjunhui/chat/crop"
	"github.com/ouqiang/supervisor-event-listener/event"
	"strings"
)

type WechatCrop struct{
	client *crop.Client
}

func (this *WechatCrop) Send(message event.Message) error {
	// fmt.Fprint(os.Stdout, message.String() + "\n")

	if this.client == nil {
		// fmt.Fprint(os.Stdout, "Init Wechat Crop\n")
		// fmt.Fprint(os.Stdout, Conf.WechatCrop.Secret + "\n")
		this.client = crop.New(Conf.WechatCrop.CropID, Conf.WechatCrop.AgentID, Conf.WechatCrop.Secret)
	}
	content := message.String()

	r := regexp.MustCompile(`(\[(.*?)])`)
	result := r.FindAllStringSubmatch(content, -1)

	text := ""
	if result != nil {
		contentList := []string{}
		for _, v := range result {
			if len(v) == 3 && v[2] != "" {
				contentList = append(contentList, v[2])
			}
		}
		text = strings.Join(contentList, "\n")
	} else {
		text = content
	}

	msg := crop.Message{}
	msg.ToUser = "@all"
	msg.MsgType = "text"
	msg.Text = crop.Content{Content: text}

	err := this.client.Send(msg)
	if err != nil {
		errorMessage := fmt.Sprintf("发送Wechat企业号消息失败-%s",
			err.Error())

		return errors.New(errorMessage)
	}
	return nil
}

