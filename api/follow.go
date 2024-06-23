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
		//没有成功获取到用户ID
		userInfo, err := ctl.GetUserInfo(ctx.Request.Context())
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
		case 2 == followreq.Type:
			go func() {
				_, err := l.CancelFollowAction(ctx, int64(userInfo.Id), int64(followreq.ID))
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
		var followlist types.FollowListReq
		fmt.Println(followlist.ID)
		if err := ctx.ShouldBind(&followlist); err == nil {
			fmt.Printf("Register info:%#v\n", followlist)
		} else {
			fmt.Println("请求参数获取失败：", err)
		}
		l := service.GetFollowSrv()
		followings, err1 := l.GetFollowings(int64(followlist.ID))
		if err1 != nil {
			fmt.Printf("fail")
			ctx.JSON(http.StatusInternalServerError, ctl.RespError(err1))
			return
		}

		ctx.JSON(http.StatusOK, followings)
	}
}

func FollowerListHandle() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		userInfo, err := ctl.GetUserInfo(ctx.Request.Context())
		if err == nil {
			fmt.Printf("Register info:%#v\n", userInfo)
		} else {
			fmt.Println("请求参数获取失败：", err)
		}
		l := service.GetFollowSrv()
		followers, err1 := l.GetFollower(ctx, int64(userInfo.Id))
		if err1 != nil {
			fmt.Printf("fail")
			ctx.JSON(http.StatusInternalServerError, ctl.RespError(err1))
			return
		}

		ctx.JSON(http.StatusOK, followers)
	}
}
