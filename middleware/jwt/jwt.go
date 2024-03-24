package middleware

import (
	"ToDoList_self/pkg/ctl"
	"ToDoList_self/pkg/util"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": 2003,
				"msg":  "请求头中auth为空",
			})
		}
		tokenclaim, err := util.ParseToken(authHeader)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": "token无效",
			})
			c.Abort()
			return
		} else if tokenclaim.ExpiresAt < time.Now().Unix() {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": "token无效,超过有效期",
			})
			c.Abort()
			return
		}
		c.Request = c.Request.WithContext(ctl.NewUserOptContext(c.Request.Context(), &ctl.UserInfo{Id: tokenclaim.Id}))

	}
}
