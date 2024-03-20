package routes

import (
	"ToDoList_self/api"
	middleware "ToDoList_self/middleware/jwt"
	"github.com/gin-gonic/gin"
)

func NewRoute() {
	r := gin.Default()
	r.POST("/register", api.RegisterHandle())
	r.POST("/login", api.LoginHandle())
	authed := r.Group("/")
	authed.Use(middleware.JWT())
	{
		authed.GET("ping", func(c *gin.Context) {
			c.JSON(200, "success")
		})
		authed.POST("task_create", api.CreateTaskHandle())
	}
	r.Run(":8080")
}
