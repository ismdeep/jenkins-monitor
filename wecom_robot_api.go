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
	Msg     string `json:"msg"`
}

// TextMessagePostData 企业微信机器人Post请求数据结构体。接口文档： https://work.weixin.qq.com/api/doc/90000/90136/91770
type TextMessagePostData struct {
	MsgType       string   `json:"msgtype"`
	MentionedList []string `json:"mentioned_list"`
	Text          struct {
		Content string `json:"content"`
	} `json:"text"`
}

// SendMessage 发送信息
func (receiver *WeComRobotService) SendMessage(msg string) error {
	httpClient := &http.Client{}
	postData := TextMessagePostData{
		MsgType:       "text",
		MentionedList: make([]string, 0),
		Text: struct {
			Content string `json:"content"`
		}{
			Content: msg,
		},
	}
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
