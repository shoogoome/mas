package physicalTransmission

import (
	"mas/utils/config"
	"mas/utils/rabbitmq"
	"strconv"
	"sync"
	"time"
)

type dataServerInfo struct {
	serverIp string
	registerTime time.Time
}

var dataServers = make([]dataServerInfo, 1)
var mutex sync.RWMutex

// 监听数据服务
func ListenHearbeat() {
	q := rabbitmq.New(config.SystemConfig.RabbitMQ.Host)
	defer q.Close()

	q.Bind(config.SystemConfig.RabbitMQ.Queue)
	c := q.Consume()
	go removeExpiredDataServer()
	for msg := range c {
		dataServer, e := strconv.Unquote(string(msg.Body))
		if e != nil {
			panic(e)
		}
		mutex.Lock()
		dataServers = append(dataServers, dataServerInfo {
			dataServer, time.Now(),
		})
		mutex.Unlock()
	}
}

// 每5s扫描一次数据服务信息 超过10s没响应则删除
func removeExpiredDataServer() {
	for {
		time.Sleep(5 * time.Second)
		index := make([]int, 1)
		mutex.Lock()
		// 筛选过期服务
		for i, v := range dataServers {
			if v.registerTime.Add(10 * time.Second).Before(time.Now()) {
				index = append(index, i)
			}
		}
		// 移除过期服务
		for _, i := range index {
			dataServers[i] = dataServers[len(dataServers) - 1]
			dataServers = dataServers[:len(dataServers) - 1]
		}
		mutex.Unlock()
	}
}

// 随机获取数据服务



