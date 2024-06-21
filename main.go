package main

import (
	"ToDoList_self/config"
	"ToDoList_self/pkg/log"
	"ToDoList_self/repository/cache"
	"ToDoList_self/repository/db/dao"
	"ToDoList_self/routes"
)

// @title Self_TDoList
// @version 0.10
// @description A sample demo of ToDoList
// @contact.name Little Chen
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @BasePath /todoList
func main() {
	config.LoadeConf()
	dao.InitMysql()
	cache.InitRedis()
	log.InitLog()
	//mq.InitRabbitMQ()
	//mq.InitFollowRabbitMQ()
	routes.NewRoute()
}
