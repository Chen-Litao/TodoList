package api

import (
	"ToDoList_self/pkg/ctl"
	"ToDoList_self/service"
	"ToDoList_self/types"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CreateTaskHandle() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var createreq types.CreateReq
		//获取前端发来的json数据
		if err := ctx.ShouldBind(&createreq); err == nil {
			fmt.Printf("Register info:%#v\n", createreq)
		} else {
			fmt.Println("请求参数获取失败：", err)
		}
		l := service.GetTaskSrv()
		code, err := l.CreateTask(ctx.Request.Context(), &createreq)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, ctl.RespError(err, code))
			return
		} else {
			ctx.JSON(http.StatusOK, ctl.RespSuccess())
			return
		}
	}
}

func ListTaskHandle() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var listTaskreq types.ListTasksReq
		//获取前端发来的json数据
		if err := ctx.ShouldBind(&listTaskreq); err == nil {
			fmt.Printf("Register info:%#v\n", listTaskreq)
		} else {
			fmt.Println("请求参数获取失败：", err)
		}
		l := service.GetTaskSrv()
		dateInfo, code, err := l.ListTask(ctx.Request.Context(), &listTaskreq)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, ctl.RespError(err, code))
			return
		} else {
			ctx.JSON(http.StatusOK, dateInfo)
			return
		}
	}
}

func ShowTaskHandle() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var showTaskreq types.ShowTasksReq
		if err := ctx.ShouldBind(&showTaskreq); err == nil {
			fmt.Printf("Register info:%#v\n", showTaskreq)
		} else {
			fmt.Println("请求参数获取失败：", err)
		}
		l := service.GetTaskSrv()
		dateInfo, code, err := l.ShowTask(ctx.Request.Context(), &showTaskreq)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, ctl.RespError(err, code))
			return
		} else {
			ctx.JSON(http.StatusOK, dateInfo)
			return
		}
	}
}

func UpdateTaskHandle() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var showTaskreq types.UpdateTasksReq
		if err := ctx.ShouldBind(&showTaskreq); err == nil {
			fmt.Printf("Register info:%#v\n", showTaskreq)
		} else {
			fmt.Println("请求参数获取失败：", err)
		}
		l := service.GetTaskSrv()
		code, err := l.UpateTask(ctx.Request.Context(), &showTaskreq)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, ctl.RespError(err, code))
			return
		} else {
			ctx.JSON(http.StatusOK, ctl.RespSuccess())
			return
		}
	}
}

func SearchTaskHandle() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var searchTaskreq types.SearchTasksReq
		if err := ctx.ShouldBind(&searchTaskreq); err == nil {
			fmt.Printf("Register info:%#v\n", searchTaskreq)
		} else {
			fmt.Println("请求参数获取失败：", err)
		}
		l := service.GetTaskSrv()
		dateInfo, code, err := l.SearchTask(ctx.Request.Context(), &searchTaskreq)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, ctl.RespError(err, code))
			return
		} else {
			ctx.JSON(http.StatusOK, dateInfo)
			return
		}
	}
}

func DeleteTaskHandle() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var deleteTaskreq types.DeleteTasksReq
		if err := ctx.ShouldBind(&deleteTaskreq); err == nil {
			fmt.Printf("Register info:%#v\n", deleteTaskreq)
		} else {
			fmt.Println("请求参数获取失败：", err)
		}
		l := service.GetTaskSrv()
		code, err := l.DeleteTask(ctx.Request.Context(), &deleteTaskreq)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, ctl.RespError(err, code))
			return
		} else {
			ctx.JSON(http.StatusOK, ctl.RespSuccess())
			return
		}
	}
}
