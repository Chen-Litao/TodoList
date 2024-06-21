package mq

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

var BaseRmq *RabbitMQ

type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	// 连接信息
	MqUrl string
}

func InitRabbitMQ() {
	//connString := strings.Join([]string{config.ConfigVal.RabbitMQ.RabbitMQ, "://", config.ConfigVal.RabbitMQ.RabbitMQUser, ":", config.ConfigVal.RabbitMQ.RabbitMQPassWord, "@", config.ConfigVal.RabbitMQ.RabbitMQHost, ":", config.ConfigVal.RabbitMQ.RabbitMQPort, "/"}, "")
	connString := "amqp://guest:guest@10.151.1.122:5672/"
	conn, err := amqp.Dial(connString)
	if err != nil {
		panic(err)
	}
	BaseRmq = &RabbitMQ{
		MqUrl: connString,
	}
	BaseRmq.conn = conn
	BaseRmq.channel, err = conn.Channel()
	BaseRmq.failOnError(err, "Failed to get channel")
}

func (r *RabbitMQ) failOnError(err error, msg string) {
	if err != nil {
		log.Printf("%s: %s", msg, err)
		log.Panicf("%s: %s", msg, err)
	}
}
