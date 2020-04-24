package mq

import (
	"StorePanAPI/config"
	"github.com/streadway/amqp"
	"log"
)

var conn *amqp.Connection
var channel *amqp.Channel

// 初始化channel
func initChannel() bool {
	if channel != nil {
		return true
	}
	var err error
	// 创建连接
	conn, err = amqp.Dial(config.RabbitMQUrl)
	if err != nil {
		log.Fatal(err)
		return false
	}
	// 初始化channel
	channel, err = conn.Channel()
	if err != nil {
		log.Fatal(err)
		return false
	}
	return true
}

// 发布消息方法
func Publish(exchange, routingKey string, msg []byte) bool {
	// 初始化channel
	if !initChannel() {
		return false
	}
	// 发布消息
	err := channel.Publish(
		exchange,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        msg,
		},
	)
	if err != nil {
		log.Fatal(err)
		return false
	}
	return true
}
