package main

import "testing"

func TestMillsToHumanText(t *testing.T) {
	type args struct {
		val int64
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "TestMillsToHumanText-001",
			args: args{
				val: 100,
			},
			want: "100毫秒",
		},
		{
			name: "TestMillsToHumanText-002",
			args: args{
				val: 1000,
			},
			want: "1.00秒",
		},
		{
			name: "TestMillsToHumanText-003",
			args: args{
				val: 1000 * 70,
			},
			want: "1.17分钟",
		},
		{
			name: "TestMillsToHumanText-004",
			args: args{
				val: 1000 * 60 * 70,
			},
			want: "1.17小时",
		},
		{
			name: "TestMillsToHumanText-005",
			args: args{
				val: 1000 * 60 * 60 * 50,
			},
			want: "2.08天",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MillsToHumanText(tt.args.val); got != tt.want {
				t.Errorf("MillsToHumanText() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetTimeNow(t *testing.T) {
	type args struct {
		locationStr string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "TestGetTimeNow-001",
			args: args{
				locationStr: TimeZoneShangHai,
			},
		},
		{
			name: "TestGetTimeNow-002",
			args: args{
				locationStr: TimeZoneTokyo,
			},
		},
		{
			name: "TestGetTimeNow-003",
			args: args{
				locationStr: TimeZoneChicago,
			},
		},
		{
			name: "TestGetTimeNow-004",
			args: args{
				locationStr: TimeZoneLondon,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetTimeNow(tt.args.locationStr)
			t.Logf("TestGetTimeNow(), got = %v", got)
		})
	}
}
