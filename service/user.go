package service

import (
	"ToDoList_self/config"
	"ToDoList_self/repository/db/model"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"os"
	"time"
)

type RegisterReq struct {
	User     string `form:"user"`
	Password string `form:"password"`
}
type LoginReq struct {
	User     string `form:"user" json:"user"`
	Password string `form:"password" json:"password"`
}
type Claims struct {
	User     string
	Password string
	jwt.StandardClaims
}

func RegisterHandle() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var register RegisterReq
		//获取前端发来的json数据
		if err := ctx.ShouldBind(&register); err == nil {
			fmt.Printf("Register info:%#v\n", register)
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

var JwtSecret = []byte(os.Getenv("JWT_SECRET"))

func LoginHandle() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var login LoginReq
		//获取前端发来的json数据
		if err := ctx.ShouldBind(&login); err == nil {
			fmt.Printf("Register info:%#v\n", login)
		} else {
			fmt.Println("请求参数获取失败：", err)
		}
		Userreq := model.User{}
		//判断当前按用户是否存在
		err := config.DB.Model(&model.User{}).Where("user_name = ?", login.User).First(&Userreq).Error
		fmt.Println("用户：", Userreq)
		if err == nil {
			//检索到了用户
			//验证密码是否匹配
			if login.Password == Userreq.Password {
				//生成token
				nowTime := time.Now()
				expireTime := nowTime.Add(24 * time.Hour)
				claims := Claims{
					User:     login.User,
					Password: login.Password,
					StandardClaims: jwt.StandardClaims{
						ExpiresAt: expireTime.Unix(),
						Issuer:    "to-do-list",
					},
				}
				tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
				jwtStr, err := tokenClaims.SignedString(JwtSecret)
				if err != nil {
					fmt.Println("生成jwt参数失败：", err)
					return
				}
				ctx.JSON(http.StatusOK, gin.H{
					"Token":   jwtStr,
					"message": "验证成功，生成token",
				})
			} else {
				fmt.Println("密码验证失败：")
				ctx.JSON(http.StatusBadGateway, gin.H{
					"user":    login.User,
					"message": "密码验证失败",
				})
			}
		} else {
			//为检索到相关用户
			fmt.Println("未找到相关用户", err)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "用户未找到，请注册",
			})
		}

	}
}
