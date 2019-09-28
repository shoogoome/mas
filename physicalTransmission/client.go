package physicalTransmission

import (
	"apimachinery/pkg/util/rand"
	"google.golang.org/grpc"
	"log"
	pb "mas/models/physicalTransmission"
	"mas/utils/config"

	"time"
)

// gRPC连接池实体
var clientPool = make(chan pb.PhysicalTransmissionClient, config.SystemConfig.Server.GrpcClientNumber)

// 单连接接口
type PhysicalTransmission interface {
	Close()
}

// 单链接实体
type physicalTransmission struct {
	pb.PhysicalTransmissionClient
}

// 关闭连接则放回连接池
func (client physicalTransmission) Close() {
	clientPool <- client.PhysicalTransmissionClient
}

// 获取连接（连接池为空时等待）
func NewPhysicalTransmission() PhysicalTransmission {
	return physicalTransmission{<-clientPool}
}

// 新建gRPC连接
func newGrpcConnection() pb.PhysicalTransmissionClient {

	rand.Seed(time.Now().Unix())
	index := rand.Intn(config.SystemConfig.Server.ServerNum)
gRPC:
	conn, err := grpc.Dial(config.SystemConfig.Server.ServerIp[index], grpc.WithInsecure())
	if err != nil {
		time.Sleep(time.Second * 3)
		log.Printf("[!] 存储服务gRPC连接失败 [%s]: %v\n", config.SystemConfig.Server.ServerIp[index], err)
		goto gRPC
	}
	c := pb.NewPhysicalTransmissionClient(conn)
	return c
}

// 初始化连接池 (定义连接池最大数量可能很大，所以直接启动goroutine运行)
func InitGrpcClientPool() {
	// 填充最大
	go func() {
		for i := 0; i < config.SystemConfig.Server.GrpcClientNumber; i++ {
			clientPool <- newGrpcConnection()
		}
	} ()
}


