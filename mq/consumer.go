package mq

import (
	"log"
)

var done chan bool

func StartConsume(qName, cName string, callback func(msg []byte) bool) {
	msgs, err := channel.Consume(
		qName,
		cName,
		true,
		false,
		false,
		false,
		nil)
	if err != nil {
		log.Fatal(err)
		return
	}
	done = make(chan bool)
	// 异步逐条消费信息
	go func() {
		for msg := range msgs {
			isSuc := callback(msg.Body)
			if !isSuc {
				log.Fatal("consume message fail")
				// TODO 写入 ErrQueue队列中进行补偿
			}
		}
	}()

	// 阻塞等待消息消费完毕
	<-done

	// 关闭channel
	channel.Close()
}
