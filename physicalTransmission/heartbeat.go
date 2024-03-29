package physicalTransmission

import (
	"fmt"
	"mas/utils/config"
	"mas/utils/rabbitmq"
	"math/rand"
	"strconv"
	"sync"
	"time"
)



var dataServers map[string]time.Time
var mutex sync.RWMutex

// 监听数据服务
func ListenHearbeat() {

	dataServers = make(map[string]time.Time)
Connection:
	q, e := rabbitmq.New(config.SystemConfig.RabbitMQ.Host)
	if e != nil {
		time.Sleep(time.Second * 3)
		fmt.Println("[!] rabbitmq connection fail, try to reconnect")
		goto Connection
	}
	defer q.Close()
Bind:
	err := q.Bind(config.SystemConfig.RabbitMQ.Queue)
	if err != nil {
		time.Sleep(time.Second * 3)
		fmt.Println("[!] rabbitmq bind fail, try to rebind")
		goto Bind
	}
	c := q.Consume()
	go removeExpiredDataServer()
	for msg := range c {
		dataServer, e := strconv.Unquote(string(msg.Body))
		if e == nil {
			mutex.Lock()
			dataServers[dataServer] = time.Now()
			mutex.Unlock()
		} else {
			fmt.Println(fmt.Sprintf("[!] hearbeat Wrongful: %v", e))
		}
	}
}

// 每5s扫描一次数据服务信息 超过10s没响应则删除
func removeExpiredDataServer() {
	for {
		time.Sleep(5 * time.Second)
		mutex.Lock()
		for s, v := range dataServers {
			if v.Add(10 * time.Second).Before(time.Now()) {
				delete(dataServers, s)
			}
		}
		mutex.Unlock()
	}
}

// 获取随机IP
func GetRandomServerIp() (server []string){

	rand.Seed(time.Now().Unix())

	serverIP := make([]string, 0)
	mutex.Lock()
	// 获取全部服务ip
	for k := range dataServers {
		serverIP = append(serverIP, k)
	}
	// 随机填充
	for _, k := range rand.Perm(len(dataServers)) {
		server = append(server, serverIP[k])
	}
	mutex.Unlock()
	return server
}