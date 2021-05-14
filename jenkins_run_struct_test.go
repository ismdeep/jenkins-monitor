package main

import (
	"github.com/ismdeep/log"
	"testing"
)

func TestFetchJenkinsRuns(t *testing.T) {
	type args struct {
		url     string
		jobName string
	}
	tests := []struct {
		name    string
		args    args
		want    []JenkinsRun
		wantErr bool
	}{
		{
			name: "TestFetchJenkinsRuns-001",
			args: args{
				url:     "https://jenkinswh.uniontech.com",
				jobName: "gitlab-flow-builder",
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetJenkinsRunList(tt.args.url, tt.args.jobName)
			if err != nil {
				t.Errorf("FetchJenkinsRuns() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			log.Info("FetchJenkinsRuns()", "got", got)
		})
	}
}

func TestGetJenkinsRun(t *testing.T) {
	type args struct {
		url          string
		jobName      string
		jenkinsRunID string
	}
	tests := []struct {
		name    string
		args    args
		want    *JenkinsRun
		wantErr bool
	}{
		{
			name: "TestGetJenkinsRunDetail-001",
			args: args{
				url:          "https://jenkinswh.uniontech.com",
				jobName:      "gitlab-flow-builder",
				jenkinsRunID: "4083",
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetJenkinsRun(tt.args.url, tt.args.jobName, tt.args.jenkinsRunID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetJenkinsRunDetail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("GetJenkinsRunDetail() got = %v", got)
		})
	}
}

func TestGetJenkinsRunDetail(t *testing.T) {
	type args struct {
		url          string
		jobName      string
		jenkinsRunID string
	}
	tests := []struct {
		name    string
		args    args
		want    *JenkinsRunDetail
		wantErr bool
	}{
		{
			name: "TestGetJenkinsRunDetail-001",
			args: args{
				url:          "https://jenkinswh.uniontech.com",
				jobName:      "gitlab-flow-builder",
				jenkinsRunID: "4079",
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetJenkinsRunDetail(tt.args.url, tt.args.jobName, tt.args.jenkinsRunID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetJenkinsRunDetail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			log.Info("TestGetJenkinsRunDetail()", "got", got)
		})
	}
}
