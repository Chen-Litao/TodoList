package service

import (
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
func (s *UserSrv) Register(ctx context.Context, req *types.RegisterReq) (err error) {
	userdao := dao.NewUserDao(ctx)
	_, err = userdao.FindUserByUserName(req.User)
	//判断当前按用户是否存在
	switch err {
	case gorm.ErrRecordNotFound:
		regUser := model.User{UserName: req.User, Password: req.Password}
		err = userdao.CreateUser(&regUser)
		if err != nil {
			fmt.Println("创建用户出错")
			return
		}
	case nil:
		err = errors.New("用户已存在")
		return
	default:
		return
	}
	return
}

func (s *UserSrv) Login(ctx context.Context, req *types.LoginReq) (Token string, err error) {
	userdao := dao.NewUserDao(ctx)
	user, err := userdao.FindUserByUserName(req.User)
	if err != nil {
		fmt.Println("未找到相关用户", err)
		return
	}
	//已经找到用户开始匹配密码
	if req.Password == user.Password {
		Token, err = util.CreateToken(req.User, req.Password)
		if err != nil {
			fmt.Println("token生成失败", err)
			return
		}
		return
	} else {
		err = errors.New("密码匹配错误")
		return
	}
}
