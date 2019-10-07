package main

import (
	"github.com/astaxie/beego"
	"mas/models"
	"mas/physicalTransmission"
	physicalTransmissionRun "mas/physicalTransmission/run"
	_ "mas/routers"
	"mas/utils/config"
	"mas/utils/mongo"
)

func init() {

	/* 初始化配置不宜使用goroutine
	   应等待所有配置初始化完毕
	*/
	// 初始化系统配置
	config.InitSystemConfig()
	// 初始化环境变量
	config.InitEnvConfig()
	// 初始化rs纠删配置
	models.InitRsConfig()
	// 初始化mongo连接
	go mongo.InitMongoClient()
	// 监听存储服务信号
	go physicalTransmission.ListenHearbeat()
	// 启动存储服务端
	go physicalTransmissionRun.Run()

}

func main() {
	beego.Run()
}
