package main

import (
	"github.com/astaxie/beego"
	"mas/models"
	"mas/physicalTransmission/run"
	_ "mas/routers"
)

func init() {

	// 监听存储服务信号
	//go physicalTransmission.ListenHearbeat()
	// 初始化rs纠删配置
	go models.InitRsConfig()
	// 启动存储服务端
	go run.Run()

}

func main() {
	beego.Run()
}
