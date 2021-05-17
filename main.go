package main

import (
	"encoding/json"
	"github.com/ismdeep/args"
	"github.com/ismdeep/log"
	"io/ioutil"
	"time"
)

func main() {
	if !args.Exists("-c") {
		log.Error("main()", "msg", "run program with    -c config.json")
		return
	}

	config := &Config{}
	data, err := ioutil.ReadFile(args.GetValue("-c"))
	if err != nil {
		log.Error("main()", "msg", "load config.json failed")
		return
	}

	if err := json.Unmarshal(data, config); err != nil {
		log.Error("main()", "msg", "extract config.json to &Config{} failed")
		return
	}

	if config.URL == "" || config.JobName == "" || config.ServiceName == "" || config.WeComRobotKey == "" {
		log.Error("main()", "msg", "config.json error")
		return
	}

	log.Info("main()", "msg", "jenkinsMonitorService Start Successfully.")
	log.Info("main()", "GetTimeNow()", GetTimeNow(TimeZoneShangHai))

	service := &JenkinsMonitorService{}
	service.Config = config
	service.WeComRobot = &WeComRobotService{Key: config.WeComRobotKey}
	service.SleepDuration = 3 * time.Second
	service.ErrSleepDuration = 1 * time.Millisecond
	service.StartMonitor()

	return
}
