package dao

import (
	"ToDoList_self/config"
	"ToDoList_self/repository/db/model"
	"context"
	"gorm.io/gorm"
)

type UserDao struct {
	*gorm.DB
}

// 创建一个可被追踪链路的上下文
func NewUserDao(ctx context.Context) *UserDao {
	if ctx == nil {
		ctx = context.Background()
	}
	return &UserDao{config.NewDBClient(ctx)}
}

// 通过用户名找对象
func (dao *UserDao) FindUserByUserName(userName string) (user *model.User, err error) {
	err = dao.Model(&model.User{}).Where("user_name=?", userName).First(&user).Error
	return
}

// 通过id找对象
func (dao *UserDao) FindUserByUserID(userName string) (user *model.User, err error) {
	err = dao.Model(&model.User{}).Where("id=?", userName).First(&user).Error
	return
}

// 创建user
func (dao *UserDao) CreateUser(user *model.User) (err error) {
	err = dao.Model(&model.User{}).Create(user).Error
	return
}
