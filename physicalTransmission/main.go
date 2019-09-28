package mian

import (
	"google.golang.org/grpc"
	"log"
	pb "mas/models/physicalTransmission"
	"net"
	"time"
)

const (
	port = ":5432"
)

// 物理存储服务层启动入口
func main() {

	// tcp连接
TCP:
	lis, err := net.Listen("tcp", port)
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
}

