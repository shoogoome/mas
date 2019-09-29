package main

import (
	"google.golang.org/grpc"
	"log"
	pb "mas/models/physicalTransmission"
	"mas/utils/config"
	"mas/utils/rabbitmq"
	"net"
	"time"
)

// 定期发送心跳信号
func StartHeartbeat() {
	q := rabbitmq.New(config.SystemConfig.RabbitMQ.Host)
	defer q.Close()
	for {
		q.Publish(config.SystemConfig.RabbitMQ.Queue, config.SystemConfig.Server.Server)
		time.Sleep(5 * time.Second)
	}
}

// 物理存储服务层启动入口
func main() {

	// tcp连接
TCP:
	lis, err := net.Listen("tcp", config.SystemConfig.Server.GrpcPort)
	if err != nil {
		time.Sleep(time.Second * 3)
		log.Println("[!] tcp连接错误: ", err)
		goto TCP
	}
	s := grpc.NewServer()
	// gRPC注册
REGISTER:
	pb.RegisterPhysicalTransmissionServer(s, &server{})
	if err = s.Serve(lis); err != nil {
		time.Sleep(time.Second * 3)
		log.Println("[!] gRPC注册失败: ", err)
		goto REGISTER
	}
	// 启动活跃心跳信号
	go StartHeartbeat()
}

