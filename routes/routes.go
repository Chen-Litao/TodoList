package routes

import (
	"ToDoList_self/api"
	_ "ToDoList_self/docs"
	middleware "ToDoList_self/middleware/jwt"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewRoute() {
	r := gin.Default()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler)) // 开启swag
	todoList := r.Group("todoList")
	{
		todoList.POST("/register", api.RegisterHandle())
		todoList.POST("/login", api.LoginHandle())
		authed := todoList.Group("/")
		authed.Use(middleware.JWT())
		{
			authed.GET("ping", func(c *gin.Context) {
				c.JSON(200, "success")
			})
			authed.POST("task_create", api.CreateTaskHandle())
		}
	}

	r.Run(":8080")
}
