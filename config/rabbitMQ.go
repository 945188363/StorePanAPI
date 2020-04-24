package config

const (
	AsyncTransferEnable = true

	RabbitMQUrl = "amqp://guest:guest@127.0.0.1:5672"

	TransferExchangerName = "uploadServer.trans"

	TransferQueueName = "uploadServer.queue"

	TransferErrQueueName = "uploadServer.errQueue"

	TransferRoutingKey = "trans"
)
