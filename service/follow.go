package service

import (
	"ToDoList_self/repository/cache"
	"ToDoList_self/repository/db/dao"
	"ToDoList_self/repository/mq"
	"context"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

var FollowSrvIns *FollowSrv
var FollowSrvOnce sync.Once

type FollowSrv struct {
}

func GetFollowSrv() *FollowSrv {
	TaskSrvOnce.Do(func() {
		//FollowSrvIns = &FollowSrv{}
	})
	return FollowSrvIns
}
func CacheTimeGenerator() time.Duration {
	// 先设置随机数 - 这里比较重要
	rand.Seed(time.Now().Unix())
	// 再设置缓存时间
	// 10 + [0~20) 分钟的随机时间
	return time.Duration((10 + rand.Int63n(20)) * int64(time.Minute))
}

// FollowAction 关注操作的业务
func (followService *FollowSrv) FollowAction(ctx context.Context, userId int64, targetId int64) (bool, error) {

	followDao := dao.NewFollowDao(ctx)
	follow, err := followDao.FindEverFollowing(userId, targetId)
	// 寻找SQL 出错。
	if nil != err {
		return false, err
	}
	// 获取关注的消息队列
	followAddMQ := mq.SimpleFollowAddMQ
	// 曾经关注过，只需要update一下followed即可。
	if nil != follow {
		//发送消息队列
		err := followAddMQ.PublishSimpleFollow(fmt.Sprintf("%d-%d-%s", userId, targetId, "update"))
		if err != nil {
			return false, err
		}
		//更新Redis
		followService.AddToRDBWhenFollow(ctx, userId, targetId)
		return true, nil

	}
	//发送消息队列
	err = followAddMQ.PublishSimpleFollow(fmt.Sprintf("%d-%d-%s", userId, targetId, "insert"))
	if err != nil {
		return false, err
	}
	followService.AddToRDBWhenFollow(ctx, userId, targetId)
	return true, nil
}

func (followService *FollowSrv) AddToRDBWhenFollow(ctx context.Context, userId int64, targetId int64) {
	followDao := dao.NewFollowDao(ctx)
	// 尝试给following数据库追加user关注target的记录
	stringID := strconv.FormatInt(userId, 10)
	keyCnt1 := cache.UserFollowings.Exists(context.Background(), stringID)

	if keyCnt1.Err() != nil {
		log.Println(keyCnt1.Err().Error())
	}

	// 只判定键是否不存在，若不存在即从数据库导入
	if keyCnt1.Val() <= 0 {
		userFollowingsId, _, err := followDao.GetFollowingsInfo(userId)
		if err != nil {
			log.Println(err.Error())
			return
		}
		ImportToRDBFollowing(ctx, userId, userFollowingsId)
	}
	// 数据库导入到redis结束后追加记录
	cache.UserFollowings.SAdd(ctx, strconv.FormatInt(userId, 10), targetId)

	// 尝试给follower数据库追加target的粉丝有user的记录
	keyCnt2 := cache.UserFollowings.Exists(context.Background(), strconv.FormatInt(targetId, 10))

	if keyCnt2.Err() != nil {
		log.Println(keyCnt2.Err().Error())
	}
	//
	if keyCnt2.Val() <= 0 {
		//获取target的粉丝，直接刷新，关注时刷新target的粉丝
		userFollowersId, _, err := followDao.GetFollowersInfo(targetId)
		if err != nil {
			log.Println(err.Error())
			return
		}
		ImportToRDBFollower(ctx, targetId, userFollowersId)
	}

	cache.UserFollowers.SAdd(ctx, strconv.FormatInt(targetId, 10), userId)
}

// ImportToRDBFollowing 将登陆用户的关注id列表导入到following数据库中
func ImportToRDBFollowing(ctx context.Context, userId int64, ids []int64) {
	// 将传入的userId及其关注用户id更新至redis中
	for _, id := range ids {
		cache.UserFollowings.SAdd(ctx, strconv.FormatInt(userId, 10), int(id))
	}

	cache.UserFollowings.Expire(ctx, strconv.FormatInt(userId, 10), CacheTimeGenerator())
}

// ImportToRDBFollower 将登陆用户的关注id列表导入到follower数据库中
func ImportToRDBFollower(ctx context.Context, userId int64, ids []int64) {
	// 将传入的userId及其粉丝id更新至redis中
	for _, id := range ids {
		cache.UserFollowers.SAdd(ctx, strconv.FormatInt(userId, 10), int(id))
	}

	cache.UserFollowers.Expire(ctx, strconv.FormatInt(userId, 10), CacheTimeGenerator())
}

// CancelFollowAction 取关操作的业务
func (followService *FollowSrv) CancelFollowAction(ctx context.Context, userId int64, targetId int64) (bool, error) {

	// 获取取关的消息队列
	followDelMQ := mq.SimpleFollowDelMQ
	followDao := dao.NewFollowDao(ctx)
	follow, err := followDao.FindEverFollowing(userId, targetId)
	// 寻找 SQL 出错。
	if nil != err {
		return false, err
	}
	// 曾经关注过，只需要update一下cancel即可。
	if nil != follow {
		err := followDelMQ.PublishSimpleFollow(fmt.Sprintf("%d-%d-%s", userId, targetId, "update"))
		if err != nil {
			return false, err
		}
		DelToRDBWhenCancelFollow(userId, targetId)
		return true, nil

	}
	// 没有关注关系
	return false, nil
}
func DelToRDBWhenCancelFollow(userId int64, targetId int64) {
	// 当a取关b时，redis的三个关注数据库会有以下操作
	cache.UserFollowings.SRem(cache.Ctx, strconv.FormatInt(userId, 10), targetId)

	cache.UserFollowers.SRem(cache.Ctx, strconv.FormatInt(targetId, 10), userId)

}
