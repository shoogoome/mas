package test

import (
	"github.com/streadway/amqp"
	"testing"
)
type RabbitMQ struct {
	channel 	*amqp.Channel
	Name 		string
	exchange 	string
}

func TestRabbitmq(t *testing.T) {


	conn, e := amqp.Dial("amqp://root:12345678@localhost:5672")
	if e != nil {
		panic(e)
	}
	// 打开唯一并发通道
	ch, e := conn.Channel()
	if e != nil {
		panic(e)
	}
	// 创建队列
	q, e := ch.QueueDeclare(
		"testtest",
		true,
		false,
		false,
		false,
		nil,
	)
	if e != nil {
		panic(e)
	}

	mq := new(RabbitMQ)
	mq.channel = ch
	mq.Name = q.Name
	//return mq

}
