package main

import "testing"

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
