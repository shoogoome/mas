package rabbitmq

import (
	"encoding/json"
	"github.com/streadway/amqp"
)

type RabbitMQ struct {
	channel 	*amqp.Channel
	Name 		string
	exchange 	string
}

func New(s string) *RabbitMQ {
	// 连接rabbitMQ服务器
	conn, e := amqp.Dial(s)
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
		"",
		false,
		true,
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
	return mq
}

func (q *RabbitMQ) Bind(exchange string) {
	// 绑定队列
	e := q.channel.QueueBind(
		q.Name,
		"",
		exchange,
		false,
		nil,
		)
	if e != nil {
		panic(e)
	}
	q.exchange = exchange
}


func (q *RabbitMQ) Send(queue string, body interface{}) {
	str, e := json.Marshal(body)
	if e != nil {
		panic(e)
	}
	// 发送信息
	e = q.channel.Publish(
		"",
			queue,
			false,
			false,
			amqp.Publishing{
				ReplyTo: q.Name,
				Body: []byte(str),
			})
	if e != nil {
		panic(e)
	}
}


func (q *RabbitMQ) Publish(exchange string, body interface{}) {
	str, e := json.Marshal(body)
	if e != nil {
		panic(e)
	}

	e = q.channel.Publish(
		exchange,
		"",
		false,
		false,
		amqp.Publishing{
			ReplyTo: q.Name,
			Body: []byte(str),
		})
	if e != nil {
		panic(e)
	}
}

func (q *RabbitMQ) Consume() <- chan amqp.Delivery {
	// 获取队列信息
	c, e := q.channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil)
	if e != nil {
		panic(e)
	}
	return c
}

func (q *RabbitMQ) Close() {
	_ = q.channel.Close()
}


