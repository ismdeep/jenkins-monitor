package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

// WeComRobotService 微信机器人服务
type WeComRobotService struct {
	Key string
}

// RespData 接口访问返回数据
type RespData struct {
	ErrCode int    `json:"errcode"`
	Msg     string `json:"errmsg"`
}

type MarkdownMessagePostData struct {
	MsgType  string `json:"msgtype"`
	Markdown struct {
		Content string `json:"content"`
	} `json:"markdown"`
}

// TextMessagePostData 企业微信机器人Post请求数据结构体。接口文档： https://work.weixin.qq.com/api/doc/90000/90136/91770
type TextMessagePostData struct {
	MsgType       string   `json:"msgtype"`
	MentionedList []string `json:"mentioned_list"`
	Text          struct {
		Content string `json:"content"`
	} `json:"text"`
}

// SendContent 发送内容
func (receiver *WeComRobotService) SendContent(msgType string, postData interface{}) error {
	httpClient := &http.Client{}

	data, err := json.Marshal(postData)
	if err != nil {
		return err
	}
	resp, err := httpClient.Post(fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=%v", receiver.Key), "application/json", bytes.NewReader(data))
	if err != nil {
		return err
	}

	all, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	respData := &RespData{}

	if err := json.Unmarshal(all, respData); err != nil {
		return err
	}

	if respData.ErrCode != 0 {
		return errors.New(respData.Msg)
	}

	return nil
}

// SendMessage 发送消息
func (receiver *WeComRobotService) SendMessage(msg string) error {
	return receiver.SendContent("text", TextMessagePostData{
		MsgType:       "text",
		MentionedList: make([]string, 0),
		Text: struct {
			Content string `json:"content"`
		}{
			Content: msg,
		},
	})
}

// SendMarkdown 发送MD文本
func (receiver *WeComRobotService) SendMarkdown(msg string) error {
	return receiver.SendContent("markdown", MarkdownMessagePostData{
		MsgType: "markdown",
		Markdown: struct {
			Content string `json:"content"`
		}{
			Content: msg,
		},
	})
}
