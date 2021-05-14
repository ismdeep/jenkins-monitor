package main

// Config 配置结构体
type Config struct {
	URL           string `json:"url"`
	JobName       string `json:"job_name"`
	ServiceName   string `json:"service_name"`
	WeComRobotKey string `json:"wecom_robot_key"`
	CallbackShell string `json:"callback_shell"`
}
