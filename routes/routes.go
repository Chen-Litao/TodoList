package routes

import (
	middleware "ToDoList_self/middleware/jwt"
	"ToDoList_self/service"
	"github.com/gin-gonic/gin"
)

func NewRoute() {
	r := gin.Default()
	r.POST("/register", service.RegisterHandle())
	r.POST("/login", service.LoginHandle())
	authed := r.Group("/")
	authed.Use(middleware.JWT())
	{
		authed.POST("task_create", service.CreateTaskHandle())
	}
	r.Run(":8080")
}
