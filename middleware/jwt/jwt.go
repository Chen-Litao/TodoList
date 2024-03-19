package middleware

import (
	"ToDoList_self/service"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
)

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusOK, gin.H{
				"code": 2003,
				"msg":  "请求头中auth为空",
			})
			c.Abort()
			return
		}
		token, err := jwt.ParseWithClaims(authHeader, &service.Claims{}, func(token *jwt.Token) (interface{}, error) {
			return service.JwtSecret, nil
		})
		if err != nil {
			fmt.Println("token解析失败", err)
			c.Abort()
			return
		}
		if !token.Valid {
			fmt.Println("Invalid token")
			c.Abort()
			return
		}
		c.Next()
		//https://andblog.cn/2941
	}
}
