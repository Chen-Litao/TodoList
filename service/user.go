package service

import (
	"ToDoList_self/config"
	"ToDoList_self/repository/db/model"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

type RegisterReq struct {
	User     string `form:"user" json:"user"`
	Password string `form:"password" json:"password"`
}

func RegisterHandle() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var register RegisterReq
		//获取前端发来的json数据
		if err := ctx.ShouldBind(&register); err == nil {
			fmt.Printf("Register info:%#v\n", register)
			//ctx.JSON(http.StatusOK, gin.H{
			//	"user":     register.User,
			//	"password": register.Password,
			//})
		} else {
			fmt.Println("请求参数获取失败：", err)
		}
		var User model.User
		//判断当前按用户是否存在
		err := config.DB.Model(&model.User{}).Where("user_name = ?", register.User).First(&User).Error
		if err == nil {
			fmt.Println("用户身份已创建：")
			return
		} else if errors.Is(err, gorm.ErrRecordNotFound) {
			// 用户不存在，尝试创建新用户
			//TODO  现在密码是明文存入数据库，后续需要更改为密文输入数据库
			regUser := model.User{UserName: register.User, Password: register.Password}
			err := config.DB.Model(&model.User{}).Create(&regUser).Error
			if err != nil {
				fmt.Println("写入数据库失败", err)
				// 这里可以返回错误信息或者错误响应
			} else {
				// 用户创建成功，返回成功消息
				ctx.JSON(http.StatusOK, gin.H{
					"user":     register.User,
					"password": register.Password,
					"message":  "用户注册成功",
				})
			}
		} else {
			fmt.Println("查询用户时发生错误:", err)
		}
	}
}
