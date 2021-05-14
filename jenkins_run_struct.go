package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/antchfx/htmlquery"
	"github.com/ismdeep/log"
	"golang.org/x/net/html"
	"io/ioutil"
	"net/http"
	"strings"
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

// JenkinsRunDetail JenkinsRun详细信息
type JenkinsRunDetail struct {
	Branch  string
	Changes []struct {
		CommitMsg string
		CommitID  string
	}
	RevisionID string
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

// GetJenkinsRun 获取JenkinsRun
func GetJenkinsRun(url string, jobName string, jenkinsRunID string) (*JenkinsRun, error) {
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

// GetJenkinsRunDetail 获取JenkinsRun详情信息
func GetJenkinsRunDetail(url string, jobName string, jenkinsRunID string) (*JenkinsRunDetail, error) {
	detail := &JenkinsRunDetail{}
	var doc *html.Node
	var err error
	doc, err = htmlquery.LoadURL(fmt.Sprintf("%v/view/web/job/%v/%v/", url, jobName, jenkinsRunID))
	if err != nil {
		log.Error("GetJenkinsRunDetail()", "err", err)
		return nil, err
	}

	nodes := htmlquery.Find(doc, `//div[@id="description"]//a[@rel="nofollow"]`)
	if len(nodes) <= 0 {
		return nil, errors.New("fail to extract data")
	}
	detail.Branch = nodes[0].FirstChild.Data

	changeNodes := htmlquery.Find(doc, `//div[@id='main-panel']//ol/li`)
	for _, changeNode := range changeNodes {
		commitMsg := changeNode.FirstChild.Data
		commitMsg = strings.Split(commitMsg, "\n")[0]
		changeID := htmlquery.Find(changeNode, ".//a//@href")[0].FirstChild.Data
		changeID = strings.ReplaceAll(changeID, "changes#", "")
		detail.Changes = append(detail.Changes, struct {
			CommitMsg string
			CommitID  string
		}{CommitMsg: commitMsg, CommitID: changeID})
	}

	return detail, nil
}
