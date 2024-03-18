package main

import (
	"ToDoList_self/config"
	"ToDoList_self/routes"
)

func main() {
	config.Loade()
	routes.NewRoute()
}
