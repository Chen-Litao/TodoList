package service

import (
	"ToDoList_self/pkg/ctl"
	"ToDoList_self/pkg/e"
	"ToDoList_self/pkg/log"
	//"ToDoList_self/repository/cache"
	"ToDoList_self/repository/db/dao"
	"ToDoList_self/repository/db/model"
	"ToDoList_self/types"
	"context"
	"sync"
)

//因为后续操作涉及对数据库的CRUD所以采用单例模式，保证数据的一致性

var TaskSrvIns *TaskSrv
var TaskSrvOnce sync.Once

type TaskSrv struct {
}

func GetTaskSrv() *TaskSrv {
	TaskSrvOnce.Do(func() {
		TaskSrvIns = &TaskSrv{}
	})
	return TaskSrvIns
}

func (s *TaskSrv) CreateTask(ctx context.Context, req *types.CreateReq) (code int, err error) {
	code = e.SUCCESS
	userInfo, err := ctl.GetUserInfo(ctx)
	if err != nil {
		code = e.ErrorGetUserInfo
		log.LoggerObj.Error(err, e.GetMsg(code))
		return
	}
	//1.通过context传来的ID信息找到对用用户
	userdao := dao.NewUserDao(ctx)
	user, err := userdao.FindUserByUserID(userInfo.Id)
	if err != nil {
		code = e.ErrorNotExistUser
		log.LoggerObj.Error(err, e.GetMsg(code))
		return
	}
	//将对应用户信息和TodoList信息传入数据库当中
	tempTsak := &model.Task{
		User:    *user,
		Uid:     user.ID,
		Title:   req.Title,
		Status:  0,
		Content: req.Content,
	}
	taskdao := dao.NewTaskDao(ctx)
	err = taskdao.CreateTask(tempTsak)
	if err != nil {
		code = e.ErrorCreateTask
		log.LoggerObj.Error(err, e.GetMsg(code))
		return
	}
	return
}

func (s *TaskSrv) ListTask(ctx context.Context, req *types.ListTasksReq) (resp interface{}, code int, err error) {
	code = e.SUCCESS
	userInfo, err := ctl.GetUserInfo(ctx)
	if err != nil {
		code = e.ErrorGetUserInfo
		log.LoggerObj.Error(err, e.GetMsg(code))
		return
	}
	tasks, total, err := dao.NewTaskDao(ctx).ListTask(req.Start, req.Limit, userInfo.Id)
	if err != nil {
		code = e.ErrorDatabase
		log.LoggerObj.Error(err, e.GetMsg(code))
		return
	}
	taskRepList := make([]*types.ListTasksTx, 0)
	for _, v := range tasks {
		taskRepList = append(taskRepList, &types.ListTasksTx{
			Uid:       v.ID,
			Title:     v.Title,
			Status:    v.Status,
			Content:   v.Content,
			CreatedAt: v.CreatedAt.Unix(),
		})
	}
	return ctl.RespListWithTotal(taskRepList, total), code, nil
}

func (s *TaskSrv) ShowTask(ctx context.Context, req *types.ShowTasksReq) (resp interface{}, code int, err error) {
	code = e.SUCCESS
	userInfo, err := ctl.GetUserInfo(ctx)
	if err != nil {
		code = e.ErrorGetUserInfo
		log.LoggerObj.Error(err, e.GetMsg(code))
		return
	}
	tasks, err := dao.NewTaskDao(ctx).ShowTask(req.Uid, userInfo.Id)
	if err != nil {
		code = e.ErrorDatabase
		log.LoggerObj.Error(err, e.GetMsg(code))
		return
	}
	taskreq := &types.ShowTasksTx{
		Id:     tasks.ID,
		Title:  tasks.Title,
		Status: tasks.Status,
		//ViewCount: cache.View(tasks),
		Content:   tasks.Content,
		CreatedAt: tasks.CreatedAt.Unix(),
	}
	//cache.AddView(tasks)
	return ctl.RespSuccessWithData(taskreq), code, nil
}

func (s *TaskSrv) UpateTask(ctx context.Context, req *types.UpdateTasksReq) (code int, err error) {
	code = e.SUCCESS
	userInfo, err := ctl.GetUserInfo(ctx)
	if err != nil {
		code = e.ErrorGetUserInfo
		log.LoggerObj.Error(err, e.GetMsg(code))
		return
	}
	err = dao.NewTaskDao(ctx).UpdateTask(req.ID, userInfo.Id, req)
	if err != nil {
		code = e.ErrorDatabase
		log.LoggerObj.Error(err, e.GetMsg(code))
		return
	}
	return code, err
}

func (s *TaskSrv) SearchTask(ctx context.Context, req *types.SearchTasksReq) (resp interface{}, code int, err error) {
	code = e.SUCCESS
	userInfo, err := ctl.GetUserInfo(ctx)
	if err != nil {
		code = e.ErrorGetUserInfo
		log.LoggerObj.Error(err, e.GetMsg(code))
		return
	}
	tasks, total, err := dao.NewTaskDao(ctx).SearchTaskByInfo(userInfo.Id, req.Info)
	if err != nil {
		code = e.ErrorDatabase
		log.LoggerObj.Error(err, e.GetMsg(code))
		return
	}
	taskRepList := make([]*types.ListTasksTx, 0)
	for _, v := range tasks {
		taskRepList = append(taskRepList, &types.ListTasksTx{
			Uid:       v.ID,
			Title:     v.Title,
			Status:    v.Status,
			Content:   v.Content,
			CreatedAt: v.CreatedAt.Unix(),
		})
	}
	return ctl.RespListWithTotal(taskRepList, total), code, nil
}

func (s *TaskSrv) DeleteTask(ctx context.Context, req *types.DeleteTasksReq) (code int, err error) {
	code = e.SUCCESS
	userInfo, err := ctl.GetUserInfo(ctx)
	if err != nil {
		code = e.ErrorGetUserInfo
		log.LoggerObj.Error(err, e.GetMsg(code))
		return
	}
	err = dao.NewTaskDao(ctx).DeleteTaskById(req.ID, userInfo.Id)
	if err != nil {
		code = e.ErrorDatabase
		log.LoggerObj.Error(err, e.GetMsg(code))
		return
	}
	return code, err
}
