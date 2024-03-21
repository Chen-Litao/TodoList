package main

import (
	"ToDoList_self/config"
	"ToDoList_self/pkg/log"
	"ToDoList_self/routes"
)

func main() {
	config.Loade()
	log.InitLog()
	routes.NewRoute()
}
