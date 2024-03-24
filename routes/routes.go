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
		todoList.GET("/ping", func(c *gin.Context) {
			c.JSON(200, "success")
		})
		todoList.POST("/register", api.RegisterHandle())
		todoList.POST("/login", api.LoginHandle())
		authed := todoList.Group("/")
		authed.Use(middleware.JWT())
		{
			authed.POST("task_create", api.CreateTaskHandle()) //增
			authed.POST("task_list", api.ListTaskHandle())     //查
			authed.POST("task_show", api.ShowTaskHandle())     //查
			authed.POST("task_update", api.UpdateTaskHandle()) //改
			authed.POST("task_search", api.SearchTaskHandle()) //查
			authed.POST("task_delete", api.DeleteTaskHandle()) //删

		}
	}

	r.Run(":8080")
}
