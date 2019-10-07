package rabbitmq

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"mas/exception/http_err"
)

type RabbitMQ struct {
	channel 	*amqp.Channel
	Name 		string
	exchange 	string
}

func New(s string) (*RabbitMQ, interface{}) {
	// 连接rabbitMQ服务器
	conn, e := amqp.Dial(s)
	if e != nil {
		return nil, http_err.RabbitmqConnectionFail()
	}
	// 打开唯一并发通道
	ch, e := conn.Channel()
	if e != nil {
		return nil, http_err.RabbitmqConnectionFail()
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
		return nil, http_err.RabbitmqConnectionFail()
	}

	mq := new(RabbitMQ)
	mq.channel = ch
	mq.Name = q.Name
	return mq, nil
}

func (q *RabbitMQ) Bind(exchange string) interface{} {
	// 绑定队列
	e := q.channel.QueueBind(
		q.Name,
		"",
		exchange,
		false,
		nil,
		)
	if e != nil {
		return http_err.RabbitmqBindFail()
	}
	q.exchange = exchange
	return nil
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


