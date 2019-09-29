package physicalTransmission

import (
	"google.golang.org/grpc"
	"log"
	"mas/exception/http_err"
	"mas/models/physicalTransmission"
	"mas/utils/config"
	"mas/utils/rs"
	"math/rand"
	"sync"
	"time"
)

/**
原本是打算做客户端连接池
但是忽然意识到在本系统中连接是有状态的，
这样每个客户端都需要维护 单连接x服务数 的连接池
这样的成本还不如在每次传输时new连接
 */
/*
// gRPC连接池实体
var clientPool = make(chan pb.PhysicalTransmissionClient, config.SystemConfig.Server.GrpcSingleClientNumber)

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




// 初始化连接池 (定义连接池最大数量可能很大，所以直接启动goroutine运行)
func InitGrpcClientPool() {
	// 填充最大
	for i := 0; i < config.SystemConfig.Server.GrpcSingleClientNumber; i++ {
		clientPool <- newGrpcConnection()
	}
}
*/

var phmutex sync.RWMutex

// 创建指定服务端gRPC连接
func NewAppointGrpcConnection(server []string) (chan physicalTransmission.PhysicalTransmissionClient, chan string, interface{}) {
	return newGrpcClientConnection(server)
}

// 创建随机服务端gRPC连接
func NewRandomGrpcConnection() (chan physicalTransmission.PhysicalTransmissionClient, chan string, interface{}){
	return newGrpcClientConnection(getRandomServerIp())
}

// 新建gRPC连接
func newGrpcConnection(ip string, ch chan physicalTransmission.PhysicalTransmissionClient, realIps chan string, lock chan bool) {

	rand.Seed(time.Now().Unix())
	index := rand.Intn(config.SystemConfig.Server.ServerNum)

	retry := config.SystemConfig.Server.GrpcRetry
GRPC:
	conn, err := grpc.Dial(ip, grpc.WithInsecure())
	if err != nil {
		log.Printf("[!] 存储服务gRPC连接失败 [%s]: %v\n", config.SystemConfig.Server.ServerIp[index], err)
		// 重试
		if retry > 0 {
			retry -= 1
			goto GRPC
		} else {
			lock <- false
			return
		}
	}
	mutex.Lock()
	ch <- physicalTransmission.NewPhysicalTransmissionClient(conn)
	realIps <- ip
	phmutex.Unlock()
	lock <- true
}

// 创建服务端gRPC连接
func newGrpcClientConnection(server []string) (chan physicalTransmission.PhysicalTransmissionClient, chan string, interface{}) {

	dataLenght := len(server)
	// 存储服务数量小于总数据分片数则直接报错
	tolerant := dataLenght - rs.RsConfig.AllShards
	if tolerant < 0 {
		return nil, nil, http_err.StorageServerInsufficient()
	}

	conn := make(chan physicalTransmission.PhysicalTransmissionClient, dataLenght)
	realIps := make(chan string, dataLenght)
	lock := make(chan bool)
	// 随机填充gRPC服务连接
	for _, ip := range server {
		go newGrpcConnection(ip + config.SystemConfig.Server.GrpcPort, conn, realIps, lock)
	}
	// 监控连接服务
	// 连接失败次数大于最大容错数直接报错
	for i := 0; i < dataLenght; i++{
		success := <- lock
		if !success {
			if tolerant > 0 {
				tolerant -= 1
			} else {
				return nil, nil, http_err.StorageServerInsufficient()
			}
		}
	}
	return conn, realIps, nil
}




