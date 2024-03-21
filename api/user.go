package api

import (
	"ToDoList_self/service"
	"ToDoList_self/types"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func RegisterHandle() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var register types.RegisterReq
		//获取前端发来的json数据
		if err := ctx.ShouldBind(&register); err == nil {
			fmt.Printf("Register info:%#v\n", register)
		} else {
			fmt.Println("请求参数获取失败：", err)
		}
		l := service.GetUserSrv()
		err := l.Register(ctx.Request.Context(), &register)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"ERROR":   err,
				"message": "用户注册失败",
			})
			return
		} else {
			ctx.JSON(http.StatusOK, gin.H{
				"message": "用户注册成功",
			})
			return
		}
	}
}

func LoginHandle() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var login types.LoginReq
		//获取前端发来的json数据
		if err := ctx.ShouldBind(&login); err == nil {
			fmt.Printf("Register info:%#v\n", login)
		} else {
			fmt.Println("请求参数获取失败：", err)
		}
		//判断当前按用户是否存在
		l := service.GetUserSrv()
		token, err := l.Login(ctx.Request.Context(), &login)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"ERROR":   err,
				"message": "用户注册失败",
			})
			return
		} else {
			ctx.JSON(http.StatusOK, gin.H{
				"token":   token,
				"message": "用户注册成功",
			})
			return
		}
	}
}
