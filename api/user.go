package api

import (
	"ToDoList_self/pkg/ctl"
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
		code, err := l.Register(ctx.Request.Context(), &register)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, ctl.RespError(err, code))
			return
		} else {
			ctx.JSON(http.StatusOK, ctl.RespSuccess())
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
		token, code, err := l.Login(ctx.Request.Context(), &login)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, ctl.RespError(err, code))
			return
		} else {
			ctx.JSON(http.StatusOK, ctl.RespSuccessWithData(token))
			return
		}
	}
}
