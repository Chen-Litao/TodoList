package model

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UserName string `gorm:"unique"`
	Password string
}

func SetPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// 第一个参数是数据库存储的密码，第二个是消息头传入的需要验证的密码
func CheckPassword(password, checkpassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(password), []byte(checkpassword))
	return err == nil
}
