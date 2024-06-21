package test

import (
	"ToDoList_self/repository/mq"
	"fmt"
	"testing"
)

func TestSimplefellowMQPublish(t *testing.T) {
	//config.LoadeConf()
	mq.InitRabbitMQ()
	mq.InitFollowRabbitMQ()
	fellowAddMQ := mq.SimpleFollowAddMQ
	userId := "2"
	videoId := "10"
	for i := 0; i < 10; i++ {
		msg := fmt.Sprintf("%s-%s", userId, videoId)
		fellowAddMQ.PublishSimpleFollow(msg)
	}
}
