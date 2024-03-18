package routes

import (
	"ToDoList_self/service"
	"github.com/gin-gonic/gin"
)

func NewRoute() {
	r := gin.Default()
	r.POST("/register", service.RegisterHandle())
	r.POST("/login", service.LoginHandle())
	r.Run(":8080")
}
