package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// JenkinsRun JenkinsRun结构体
type JenkinsRun struct {
	Links struct {
		Self struct {
			Href string `json:"href"`
		} `json:"self"`
		ChangeSets struct {
			Href string `json:"href"`
		} `json:"changesets"`
	} `json:"_links"`
	ID                  string `json:"id"`                  // ID
	Name                string `json:"name"`                // 名称
	Status              string `json:"status"`              // 状态 IN_PROGRESS
	StartTimeMillis     int64  `json:"startTimeMillis"`     // 开始时间
	EndTimeMillis       int64  `json:"endTimeMillis"`       // 结束时间
	DurationMillis      int64  `json:"durationMillis"`      // 耗时
	QueueDurationMillis int64  `json:"queueDurationMillis"` // 队列耗时
	PauseDurationMillis int64  `json:"pauseDurationMillis"` // 暂停耗时
	Stages              []struct {
		Links struct {
			Self struct {
				Href string `json:"href"`
			} `json:"self"`
		} `json:"_links"`
		ID                  string `json:"id"`
		Name                string `json:"name"`
		ExecNode            string `json:"execNode"`
		Status              string `json:"status"`
		StartTimeMillis     int64  `json:"startTimeMillis"`
		DurationMillis      int64  `json:"durationMillis"`
		PauseDurationMillis int64  `json:"pauseDurationMillis"`
	} `json:"stages"`
}

// GetJenkinsRunList 获取JenkinsRun列表
func GetJenkinsRunList(url string, jobName string) ([]*JenkinsRun, error) {
	httpClient := &http.Client{}
	resp, err := httpClient.Get(fmt.Sprintf("%v/job/%v/wfapi/runs", url, jobName))
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	jenkinsRunList := make([]*JenkinsRun, 0)
	if err := json.Unmarshal(data, &jenkinsRunList); err != nil {
		return nil, err
	}

	return jenkinsRunList, nil
}

// GetJenkinsRunDetail 获取JenkinsRun详情
func GetJenkinsRunDetail(url string, jobName string, jenkinsRunID string) (*JenkinsRun, error) {
	httpClient := &http.Client{}
	resp, err := httpClient.Get(fmt.Sprintf("%v/job/%v/%v/wfapi/describe", url, jobName, jenkinsRunID))
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	jenkinsRun := &JenkinsRun{}
	if err := json.Unmarshal(data, jenkinsRun); err != nil {
		return nil, err
	}
	return jenkinsRun, nil
}
