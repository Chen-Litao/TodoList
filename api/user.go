package api

import (
	"ToDoList_self/pkg/ctl"
	"ToDoList_self/service"
	"ToDoList_self/types"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// RegisterHandle

// RegisterHandle @Tags USER
// @Summary 用户注册接口接口
// @Accept json
// @Produce json
// @Param data body types.RegisterReq true "用户名, 密码"
// @Success 200 {object} ctl.Response "{"status":200,"data":{},"msg":"ok"}"
// @Failure 500  {object} ctl.Response "{"status":500,"data":{},"Msg":{},"Error":"error"}"
// @Router /register [post]
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

// LoginHandle @Tags USER
// @Summary 用户登录接口接口
// @Accept json
// @Produce json
// @Param data body types.LoginReq true "用户名, 密码"
// @Success 200 {object} ctl.Response "{"status":200,"data":{},"msg":"ok"}"
// @Failure 500  {object} ctl.Response "{"status":500,"data":{},"Msg":{},"Error":"error"}"
// @Router /login [post]
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
