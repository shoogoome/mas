package run

import (
	"fmt"
	"google.golang.org/grpc"
	pb "mas/models/physicalTransmission"
	"mas/utils/config"
	"mas/utils/rabbitmq"
	"net"
	"time"
)

// 定期发送心跳信号
func StartHeartbeat() {
	fmt.Println("[*] send Heartbeat signal...")
Connection:
	q, e := rabbitmq.New(config.SystemConfig.RabbitMQ.Host)
	if e != nil {
		time.Sleep(time.Second * 3)
		fmt.Println("[!] rabbitmq connection fail, try to reconnect")
		goto Connection
	}
	defer q.Close()
	for {
		q.Publish(config.SystemConfig.RabbitMQ.Queue, config.SystemConfig.Server.Server)
		time.Sleep(5 * time.Second)
	}
}

// 物理存储服务层启动入口
func Run() {

	// tcp连接
TCP:
	lis, err := net.Listen("tcp", config.SystemConfig.Server.GrpcPort)
	if err != nil {
		time.Sleep(time.Second * 3)
		fmt.Println("[!] tcp connection fail...: ", err)
		goto TCP
	}
	s := grpc.NewServer()
	fmt.Println("[*] tcp connection success...")
	// 启动活跃心跳信号
	go StartHeartbeat()
	// gRPC注册
REGISTER:
	pb.RegisterPhysicalTransmissionServer(s, &server{})
	if err = s.Serve(lis); err != nil {
		time.Sleep(time.Second * 3)
		fmt.Println("[!] grpc register fail...: ", err)
		goto REGISTER
	}
	fmt.Println("[*] grpc register success...")
}

