package model

import "gorm.io/gorm"

type Follow struct {
	gorm.Model
	UserId      int64
	FollowingId int64
	Followed    int8
}
