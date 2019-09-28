package main

import (
	"mas/physicalTransmission"
	_ "mas/routers"
	"github.com/astaxie/beego"
)


func init() {

	// 监听存储服务信号
	go physicalTransmission.ListenHearbeat()
	// 初始化gRPC连接池
	go physicalTransmission.InitGrpcClientPool()

}

func main() {
	beego.Run()
}

