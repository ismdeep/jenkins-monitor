package main

import "testing"

func TestJenkinsMonitorService_Add(t *testing.T) {
	type fields struct {
		JenkinsRunIDs []string
	}
	type args struct {
		jenkinsRun *JenkinsRun
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "TestJenkinsMonitorService_Add-001",
			fields: fields{
				JenkinsRunIDs: make([]string, 0),
			},
			args: args{
				jenkinsRun: &JenkinsRun{
					Links: struct {
						Self struct {
							Href string `json:"href"`
						} `json:"self"`
						ChangeSets struct {
							Href string `json:"href"`
						} `json:"changesets"`
					}{},
					ID:                  "",
					Name:                "",
					Status:              "",
					StartTimeMillis:     0,
					EndTimeMillis:       0,
					DurationMillis:      0,
					QueueDurationMillis: 0,
					PauseDurationMillis: 0,
					Stages:              nil,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			receiver := &JenkinsMonitorService{
				JenkinsRunIDs: tt.fields.JenkinsRunIDs,
			}
			receiver.MonitorFunc(tt.args.jenkinsRun)
			t.Logf("Add() receiver = %v", receiver.JenkinsRunIDs)
		})
	}
}
