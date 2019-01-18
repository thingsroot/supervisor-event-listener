package config

import (
	"flag"
	"fmt"
	"github.com/ouqiang/supervisor-event-listener/utils"
	"gopkg.in/ini.v1"
	"os"
	"strings"
)

type Config struct {
	NotifyType string
	WebHook    WebHook
	MailServer MailServer
	MailUser   MailUser
	Slack      Slack
	WechatCrop WechatCrop
}

type WechatCrop struct {
	CropID string
	AgentID int
	Secret string
}

type WebHook struct {
	Url string
}

type Slack struct {
	WebHookUrl string
	Channel    string
}

// 邮件服务器
type MailServer struct {
	User     string
	Password string
	Host     string
	Port     int
}

// 接收邮件的用户
type MailUser struct {
	Email []string
}

func ParseConfig() *Config {
	var configFile string
	flag.StringVar(&configFile, "c", "/etc/supervisor-event-listener.ini", "config file")
	flag.Parse()
	configFile = strings.TrimSpace(configFile)
	if configFile == "" {
		Exit("请指定配置文件路径")
	}
	file, err := ini.Load(configFile)
	if err != nil {
		Exit("读取配置文件失败#" + err.Error())
	}
	section := file.Section("default")
	notifyType := section.Key("notify_type").String()
	notifyType = strings.TrimSpace(notifyType)
	if !utils.InStringSlice([]string{"mail", "slack", "webhook", "wechat_crop"}, notifyType) {
		Exit("不支持的通知类型-" + notifyType)
	}

	config := &Config{}
	config.NotifyType = notifyType
	switch notifyType {
	case "mail":
		config.MailServer = parseMailServer(section)
		config.MailUser = parseMailUser(section)
	case "slack":
		config.Slack = parseSlack(section)
	case "webhook":
		config.WebHook = parseWebHook(section)
	case "wechat_crop":
		config.WechatCrop = parseWechatCrop(section)
	}

	return config
}

func parseMailServer(section *ini.Section) MailServer {
	user := section.Key("mail.server.user").String()
	user = strings.TrimSpace(user)
	password := section.Key("mail.server.password").String()
	password = strings.TrimSpace(password)
	host := section.Key("mail.server.host").String()
	host = strings.TrimSpace(host)
	port, portErr := section.Key("mail.server.port").Int()
	if user == "" || password == "" || host == "" || portErr != nil {
		Exit("邮件服务器配置错误")
	}

	mailServer := MailServer{}
	mailServer.User = user
	mailServer.Password = password
	mailServer.Host = host
	mailServer.Port = port

	return mailServer
}

func parseMailUser(section *ini.Section) MailUser {
	user := section.Key("mail.user").String()
	user = strings.TrimSpace(user)
	if user == "" {
		Exit("邮件收件人配置错误")
	}
	mailUser := MailUser{}
	mailUser.Email = strings.Split(user, ",")

	return mailUser
}

func parseSlack(section *ini.Section) Slack {
	webHookUrl := section.Key("slack.webhook_url").String()
	webHookUrl = strings.TrimSpace(webHookUrl)
	channel := section.Key("slack.channel").String()
	channel = strings.TrimSpace(channel)
	if webHookUrl == "" || channel == "" {
		Exit("Slack配置错误")
	}

	slack := Slack{}
	slack.WebHookUrl = webHookUrl
	slack.Channel = channel

	return slack
}

func parseWebHook(section *ini.Section) WebHook {
	url := section.Key("webhook_url").String()
	url = strings.TrimSpace(url)
	if url == "" {
		Exit("WebHookUrl配置错误")
	}
	webHook := WebHook{}
	webHook.Url = url

	return webHook
}

func parseWechatCrop(section *ini.Section) WechatCrop {
	crop_id := section.Key("wechat_crop.crop_id").String()
	crop_id = strings.TrimSpace(crop_id)
	agent_id, agent_idErr := section.Key("wechat_crop.agent_id").Int()

	secret := section.Key("wechat_crop.secret").String()
	secret = strings.TrimSpace(secret)

	if crop_id == "" || secret == "" || agent_idErr != nil {
		Exit("微信企业号配置错误")
	}

	wechatCrop := WechatCrop{}
	wechatCrop.CropID = crop_id
	wechatCrop.AgentID = agent_id
	wechatCrop.Secret = secret

	return wechatCrop
}

func Exit(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(1)
}
