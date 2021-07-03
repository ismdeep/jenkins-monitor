package main

import (
	"fmt"
	"github.com/ismdeep/log"
	"github.com/ismdeep/wecombot"
	"os/exec"
	"strings"
	"sync"
	"time"
)

// JenkinsMonitorService Jenkins监控服务
type JenkinsMonitorService struct {
	Config           *Config
	JenkinsRunIDs    []string
	mutex            sync.Mutex
	WeComBot         *wecombot.Bot
	SleepDuration    time.Duration
	ErrSleepDuration time.Duration
	RetryCount       int
}

// GetJenkinsRunResultMarkdown 获取Jenkins打包结果内容
func (receiver *JenkinsMonitorService) GetJenkinsRunResultMarkdown(jenkinsRun *JenkinsRun) (string, error) {
	jenkinsDetail, err := GetJenkinsRunDetail(receiver.Config.URL, receiver.Config.JobName, jenkinsRun.ID)
	if err != nil {
		log.Error("SendMessage()", "msg", "获取详情失败。", "err", err)
		return "", err
	}

	strList := make([]string, 0)
	for _, str := range jenkinsDetail.Changes {
		strList = append(strList, fmt.Sprintf("> %v (%v)", str.CommitMsg, str.CommitID))
	}

	statusText := "正在打包"
	statusClass := "comment"
	if jenkinsRun.Status == "SUCCESS" {
		statusText = "打包成功"
		statusClass = "info"
	}

	if jenkinsRun.Status == "FAILED" {
		statusText = "打包失败"
		statusClass = "warning"
	}

	markdownContent := fmt.Sprintf(`<font color="%v">%v</font> [%v] %v
%v
> 打包详情：[点击查看](%v/view/web/job/%v/%v/)
> 打包耗时：%v
> 打包时间：%v`, statusClass, jenkinsRun.Name, receiver.Config.Branch, statusText,
		strings.Join(strList, "\n"),
		receiver.Config.URL, receiver.Config.JobName, jenkinsRun.ID,
		MillsToHumanText(jenkinsRun.DurationMillis),
		GetTimeNow(TimeZoneShangHai))

	return markdownContent, nil
}

// StartMonitor 启动监控
func (receiver *JenkinsMonitorService) StartMonitor() {
	for {
		jenkinsRunList, err := GetJenkinsRunList(receiver.Config.URL, receiver.Config.JobName)
		if err != nil {
			log.Error("StartMonitor()", "msg", "GetJenkinsRunList() failed", "err", err)
			time.Sleep(receiver.ErrSleepDuration)
			continue
		}

		for _, jenkinsRun := range jenkinsRunList {
			if jenkinsRun.Name == receiver.Config.ServiceName && jenkinsRun.Status == "IN_PROGRESS" {
				jenkinsDetail, err := GetJenkinsRunDetail(receiver.Config.URL, receiver.Config.JobName, jenkinsRun.ID)
				if err != nil {
					log.Error("StartMonitor()", "msg", "GetJenkinsRunDetail() failed", "err", err)
					continue
				}
				if jenkinsDetail.Branch == receiver.Config.Branch {
					receiver.MonitorFunc(jenkinsRun)
				}
			}
		}
		time.Sleep(receiver.SleepDuration)
	}
}

// MonitorFunc 监控函数
func (receiver *JenkinsMonitorService) MonitorFunc(jenkinsRun *JenkinsRun) {

	// ---- 检测JenkinsRunID是否已经存在 start
	receiver.mutex.Lock()
	for _, tmpID := range receiver.JenkinsRunIDs {
		if tmpID == jenkinsRun.ID {
			receiver.mutex.Unlock()
			return
		}
	}
	receiver.JenkinsRunIDs = append(receiver.JenkinsRunIDs, jenkinsRun.ID)
	receiver.mutex.Unlock()
	// ---- 检测JenkinsRunID是否已经存在 end

	log.Info("MonitorFunc()", "msg", "正在打包", "jenkinsRun", jenkinsRun)

	retryCount := receiver.RetryCount
	for {
		jenkinsRun, err := GetJenkinsRun(receiver.Config.URL, receiver.Config.JobName, jenkinsRun.ID)
		if err != nil {
			// 重试次数用完了
			if retryCount <= 0 {
				log.Error("重试次数用完了，监控退出。")
				break
			}
			retryCount--
			time.Sleep(receiver.ErrSleepDuration)
			continue
		}

		// 打包成功
		if jenkinsRun.Status == "SUCCESS" {
			msg, _ := receiver.GetJenkinsRunResultMarkdown(jenkinsRun)
			_ = receiver.WeComBot.SendMarkdown(msg)
			log.Info("MonitorFunc()", "msg", "打包成功", "jenkinsRun", jenkinsRun)
			if receiver.Config.CallbackShell != "" {
				go func() {
					log.Info("MonitorFunc()", "msg", "执行打包成功回调脚本")
					startTime := time.Now().UnixNano()
					_ = exec.Command(receiver.Config.CallbackShell).Run()
					endTime := time.Now().UnixNano()
					msg := fmt.Sprintf(`<font color="info">%v</font> [%v] 服务发布成功
> 服务地址：[%v](%v)
> 发布耗时：%v
> 发布时间：%v`,
						jenkinsRun.Name, receiver.Config.Branch,
						receiver.Config.PublishURL, receiver.Config.PublishURL,
						MillsToHumanText((endTime-startTime)/1000000),
						GetTimeNow(TimeZoneShangHai))
					_ = receiver.WeComBot.SendMarkdown(msg)
					log.Info("MonitorFunc()", "msg", "执行打包成功回调脚本成功")
				}()
			}
			break
		}

		// 打包失败
		if jenkinsRun.Status == "FAILED" {
			msg, _ := receiver.GetJenkinsRunResultMarkdown(jenkinsRun)
			_ = receiver.WeComBot.SendMarkdown(msg)
			log.Info("MonitorFunc()", "msg", "打包失败", "jenkinsRun", jenkinsRun)
			break
		}
		time.Sleep(receiver.SleepDuration)
	}
}
