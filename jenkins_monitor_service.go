package main

import (
	"fmt"
	"github.com/ismdeep/log"
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
	WeComRobot       *WeComRobotService
	SleepDuration    time.Duration
	ErrSleepDuration time.Duration
	RetryCount       int
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

	jenkinsDetail, err := GetJenkinsRunDetail(receiver.Config.URL, receiver.Config.JobName, jenkinsRun.ID)
	if err != nil {
		log.Error("SendMessage()", "msg", "获取详情失败。", "err", err)
		return
	}

	strs := make([]string, 0)
	for _, str := range jenkinsDetail.Changes {
		strs = append(strs, fmt.Sprintf("> %v (%v)", str.CommitMsg, str.CommitID))
	}

	markdownContent := fmt.Sprintf(`<font color="green">%v</font> 正在打包
> CommitID: %v
%v`, jenkinsRun.Name, jenkinsDetail.RevisionID, strings.Join(strs, "\n"))

	if err := receiver.WeComRobot.SendMarkdown(markdownContent); err != nil {
		log.Error("SendMessage()", "msg", "发送信息失败。", "err", err)
		return
	}

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

		if jenkinsRun.Status == "SUCCESS" {
			msg := fmt.Sprintf("打包成功：%v", jenkinsRun.Name)
			_ = receiver.WeComRobot.SendMessage(msg)
			log.Info("MonitorFunc()", "msg", msg)
			if receiver.Config.CallbackShell != "" {
				go func() {
					log.Info("MonitorFunc()", "msg", "执行打包成功回调脚本")
					_ = exec.Command(receiver.Config.CallbackShell).Run()
					_ = receiver.WeComRobot.SendMessage(fmt.Sprintf("服务发布成功: %v", jenkinsRun.Name))
					log.Info("MonitorFunc()", "msg", "执行打包成功回调脚本成功")
				}()
			}
			break
		}

		if jenkinsRun.Status == "FAILED" {
			msg := fmt.Sprintf("打包失败：%v", jenkinsRun.Name)
			_ = receiver.WeComRobot.SendMessage(msg)
			log.Info("MonitorFunc()", "msg", msg)
			break
		}
		time.Sleep(receiver.SleepDuration)
	}
}
