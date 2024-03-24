package ctl

import (
	"context"
	"errors"
)

type key string

const userOptKey key = "user"

// 仅在context的value上传输用户ID
type UserInfo struct {
	Id uint `json:"id"`
}

func NewUserOptContext(ctx context.Context, u *UserInfo) context.Context {
	return context.WithValue(ctx, userOptKey, u)
}

func GetUserInfo(ctx context.Context) (*UserInfo, error) {
	user, ok := ctx.Value(userOptKey).(*UserInfo)
	if !ok {
		return nil, errors.New("获取用户信息错误")
	}
	return user, nil
}
