package api

import (
	"ToDoList_self/pkg/ctl"
	"ToDoList_self/service"
	"ToDoList_self/types"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func RelationActionHandle() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var followreq types.FellowReq
		userInfo, err := ctl.GetUserInfo(ctx)
		if err := ctx.ShouldBind(&followreq); err == nil {
			fmt.Printf("Register info:%#v\n", followreq)
		} else {
			fmt.Println("请求参数获取失败：", err)
		}
		l := service.GetFollowSrv()
		switch {
		case 1 == followreq.Type:
			go func() {
				_, err := l.FollowAction(ctx, int64(userInfo.Id), int64(followreq.ID))
				if err != nil {
					log.Println(err)
				}
			}()

		}
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, ctl.RespError(err))
			return
		} else {
			ctx.JSON(http.StatusOK, ctl.RespSuccess())
			return
		}

	}
}

func FollowListHandle() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {

	}
}

func FollowerListHandle() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {

	}
}
