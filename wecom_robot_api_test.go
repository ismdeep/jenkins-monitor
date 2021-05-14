package main

import (
	"fmt"
	"strings"
	"testing"
)

func TestSendMessage(t *testing.T) {
	type args struct {
		key string
		msg string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "TestSendMessage-001",
			args: args{
				key: "",
				msg: "こうにちわ、わたしはドラえもんです。",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := WeComRobotService{Key: tt.args.key}
			if err := service.SendMessage(tt.args.msg); (err != nil) != tt.wantErr {
				t.Errorf("SendMessage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWeComRobotService_SendMarkdown(t *testing.T) {
	strs := make([]string, 0)
	strs = append(strs, "> feat: 支持配置式权限控制")
	strs = append(strs, "> fix: 修复权限控制漏洞")
	content := strings.Join(strs, "\n")

	type fields struct {
		Key string
	}
	type args struct {
		msg string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "TestWeComRobotService_SendMarkdown()-001",
			fields: fields{
				Key: "016529bf-c3ec-4b5a-9d8d-23f951d90f9b",
			},
			args: args{
				msg: fmt.Sprintf(`<font color="green">%v</font> 正在打包
> CommitID: %v
%v`, "Doraemon", "dbaeb80bf73a750640431f5c6f55b0798109dc44", content),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			receiver := &WeComRobotService{
				Key: tt.fields.Key,
			}
			if err := receiver.SendMarkdown(tt.args.msg); err != nil {
				t.Errorf("SendMarkdown() error = %v", err)
			}
		})
	}
}
