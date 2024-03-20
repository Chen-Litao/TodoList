package api

import "github.com/gin-gonic/gin"

type CreateReq struct {
	Title   string `form:"title"`
	Status  string `form:"status"`
	Content string `form:"content"`
}

func CreateTaskHandle() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {

	}
}
