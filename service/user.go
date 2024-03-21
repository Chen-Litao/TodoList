package service

import (
	"ToDoList_self/pkg/e"
	"ToDoList_self/pkg/util"
	"ToDoList_self/repository/db/dao"
	"ToDoList_self/repository/db/model"
	"ToDoList_self/types"
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"sync"
)

// 因为后续操作涉及对数据库的CRUD所以采用单例模式，保证数据的一致性
var UserSrvIns *UserSrv
var UserSrvOnce sync.Once

type UserSrv struct {
}

func GetUserSrv() *UserSrv {
	UserSrvOnce.Do(func() {
		UserSrvIns = &UserSrv{}
	})
	return UserSrvIns
}
func (s *UserSrv) Register(ctx context.Context, req *types.RegisterReq) (code int, err error) {
	userdao := dao.NewUserDao(ctx)
	_, err = userdao.FindUserByUserName(req.User)
	code = e.SUCCESS
	//判断当前按用户是否存在
	switch err {
	case gorm.ErrRecordNotFound:
		req.Password, err = model.SetPassword(req.Password)
		if err != nil {
			fmt.Println("密码加密错误")
			code = e.ErrorFailEncryption
			return
		}
		regUser := model.User{UserName: req.User, Password: req.Password}
		err = userdao.CreateUser(&regUser)
		if err != nil {
			fmt.Println("创建用户出错")
			code = e.ErrorCreateUser
			return
		}
	case nil:
		err = errors.New("注册用户失败，用户已存在")
		code = e.ErrorExistUser
		return
	default:
		return
	}
	return
}

func (s *UserSrv) Login(ctx context.Context, req *types.LoginReq) (Token string, code int, err error) {
	userdao := dao.NewUserDao(ctx)
	user, err := userdao.FindUserByUserName(req.User)
	code = e.SUCCESS
	if err != nil {
		fmt.Println("未找到相关用户", err)
		code = e.ErrorNotExistUser
		return
	}
	//已经找到用户开始匹配密码
	if model.CheckPassword(user.Password, req.Password) {
		Token, err = util.CreateToken(req.User, req.Password)
		if err != nil {
			fmt.Println("token生成失败", err)
			code = e.ErrorAuthToken
			return
		}
		return
	} else {
		err = errors.New("密码匹配错误")
		code = e.ErrorNotCompare
		return
	}
}
