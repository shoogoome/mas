package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"mas/exception"
	"mas/models"
	"os"
	"strconv"
)

var SystemConfig models.SystemConfig

// 初始化读取配置
func InitSystemConfig() {
	yamlFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		exception.OutputException("fail to load config.yaml", err)
	}
	err = yaml.Unmarshal(yamlFile, &SystemConfig)
	if err != nil {
		exception.OutputException("fail to unmarshal config.yaml", err)
	}
}

// 从环境变量读取配置
func InitEnvConfig() {
	// 系统配置
	SystemConfig.Server.Server = os.Getenv("Server")
	SystemConfig.Server.FileRootPath = "/root"
	SystemConfig.Server.FileTempPath = "/tmp"
	SystemConfig.Server.Token = os.Getenv("Token")
	SystemConfig.Server.Key = os.Getenv("Key")
	SystemConfig.Server.Resend, _ = strconv.Atoi(os.Getenv("Resend"))
	SystemConfig.Server.GrpcPort = os.Getenv("GrpcPort")
	SystemConfig.Server.Resend, _ = strconv.Atoi(os.Getenv("GrpcRetry"))
	size, _ := strconv.Atoi(os.Getenv("ChuckMaxSize"))
	SystemConfig.Server.ChuckMaxSize = int64(size)
	if os.Getenv("Gzip") == "true" {
		SystemConfig.Server.Gzip = true
	} else {
		SystemConfig.Server.Gzip = false
	}
}