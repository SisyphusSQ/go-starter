package dto

import larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"

type LarkContent struct {
	Tag    string `json:"tag,omitempty"`
	Text   string `json:"text,omitempty"`
	Href   string `json:"href,omitempty"`
	UserID string `json:"user_id,omitempty"`
}

type LarkPost struct {
	ZhCN struct {
		Title   string          `json:"title"`
		Content [][]LarkContent `json:"content"`
	} `json:"zh_cn"`
}

type BotMessage struct {
	MsgType string `json:"msg_type"`
	Content struct {
		Post LarkPost `json:"post"`
	} `json:"content"`
}

func NewBotMessage(title string, contents []LarkContent) BotMessage {
	var message BotMessage
	message.MsgType = "post"
	message.Content.Post.ZhCN.Title = title
	message.Content.Post.ZhCN.Content = [][]LarkContent{contents}
	return message
}

type BotMessageResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type LarkMessageRequest struct {
	Contacts []string `json:"contacts"`
	Message  string   `json:"message"`
}

type LarkMessageResponse struct {
	Responses []*larkim.CreateMessageResp `json:"responses"`
}
