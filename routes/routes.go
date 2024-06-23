package routes

import (
	"ToDoList_self/api"
	_ "ToDoList_self/docs"
	middleware "ToDoList_self/middleware/jwt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewRoute() {
	r := gin.Default()
	store := cookie.NewStore([]byte("something-very-secret"))
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler)) // 开启swag
	r.Use(sessions.Sessions("mysession", store))
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

			authed.POST("follow_action", api.RelationActionHandle())
			authed.GET("follow_list", api.FollowListHandle())
			authed.GET("follower_list", api.FollowerListHandle())

		}
	}

	r.Run(":4080")
}
