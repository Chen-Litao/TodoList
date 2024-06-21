package dao

import (
	"ToDoList_self/repository/db/model"
	"context"
	"gorm.io/gorm"
	"log"
)

type FollowDao struct {
	*gorm.DB
}

// 创建一个可被追踪链路的上下文
func NewFollowDao(ctx context.Context) *FollowDao {
	if ctx == nil {
		ctx = context.Background()
	}
	return &FollowDao{NewDBClient(ctx)}
}
func (dao *FollowDao) FindEverFollowing(userId int64, targetId int64) (*model.Follow, error) {
	follow := model.Follow{}
	err := dao.Model(&model.Follow{}).
		Where("user_id = ? AND following_id = ?", userId, targetId).
		Take(&follow).Error
	if nil != err {
		// 当没查到记录报错时，不当做错误处理。
		if "record not found" == err.Error() {
			return nil, nil
		}
		log.Println(err.Error())
		return nil, err
	}
	// 正常情况，返回取到的关系和空err.
	return &follow, nil
}

func (dao *FollowDao) InsertFollowRelation(userId int64, targetId int64) error {
	// 生成需要插入的关系结构体。
	follow := model.Follow{
		UserId:      userId,
		FollowingId: targetId,
		Followed:    1,
	}
	err := dao.Model(&model.Follow{}).Create(&follow).Error
	if nil != err {
		log.Println(err.Error())
		return err
	}
	return nil
}

// UpdateFollowRelation 给定用户和目标用户的id，更新他们的关系为取消关注或再次关注。
func (dao *FollowDao) UpdateFollowRelation(userId int64, targetId int64, followed int8) error {
	// 更新用户与目标用户的关注记录（正在关注或者取消关注）
	err := dao.Model(&model.Follow{}).
		Where("user_id = ?", userId).
		Where("following_id = ?", targetId).
		Update("followed", followed).Error
	// 更新失败，返回错误。
	if nil != err {
		// 更新失败，打印错误日志。
		log.Println(err.Error())
		return err
	}
	// 更新成功。
	return nil
}

// GetFollowingsInfo 返回当前用户正在关注的用户信息列表，包括当前用户正在关注的用户ID列表和正在关注的用户总数
func (dao *FollowDao) GetFollowingsInfo(userId int64) ([]int64, int64, error) {

	var followingCnt int64
	var followingId []int64

	// user_id -> following_id
	result := dao.Model(&model.Follow{}).Where("user_id = ? AND followed = ?", userId, 1).Find(&followingId)
	followingCnt = result.RowsAffected

	if nil != result.Error {
		log.Println(result.Error.Error())
		return nil, 0, result.Error
	}

	return followingId, followingCnt, nil

}
func (dao *FollowDao) GetFollowersInfo(userId int64) ([]int64, int64, error) {

	var followerCnt int64
	var followerId []int64

	// following_id -> user_id
	result := dao.Model(&model.Follow{}).Where("following_id = ?", userId).Where("followed = ?", 1).Pluck("user_id", &followerId)
	followerCnt = result.RowsAffected

	if nil != result.Error {
		log.Println(result.Error.Error())
		return nil, 0, result.Error
	}

	return followerId, followerCnt, nil
}

// GetUserName 在user表中根据id查询用户姓名，放在followDao文件中并不妥当，后续可能废弃
func (dao *FollowDao) GetUserName(userId int64) (string, error) {
	var name string

	err := dao.Table("users").Where("id = ?", userId).Pluck("user_name", &name).Error

	if nil != err {
		log.Println(err.Error())
		return "", err
	}

	return name, nil
}

// GetFollowerCnt 给定当前用户id，查询relation表中该用户的粉丝数。
func (dao *FollowDao) GetFollowerCnt(userId int64) (int64, error) {
	// 用于存储当前用户粉丝数的变量
	var cnt int64
	// 当查询出现错误的情况，日志打印err msg，并返回err.
	if err := dao.
		Model(&model.Follow{}).
		Where("following_id = ?", userId).
		Where("followed = ?", 1).
		Count(&cnt).Error; nil != err {
		log.Println(err.Error())
		return 0, err
	}
	// 正常情况，返回取到的粉丝数。
	return cnt, nil
}

// GetFollowingCnt 给定当前用户id，查询relation表中该用户关注了多少人。
func (dao *FollowDao) GetFollowingCnt(userId int64) (int64, error) {
	// 用于存储当前用户关注了多少人。
	var cnt int64
	// 查询出错，日志打印err msg，并return err
	if err := dao.Model(&model.Follow{}).
		Where("user_id = ?", userId).
		Where("followed = ?", 1).
		Count(&cnt).Error; nil != err {
		log.Println(err.Error())
		return 0, err
	}
	// 查询成功，返回人数。
	return cnt, nil
}
