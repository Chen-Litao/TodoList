package dao

import (
	"ToDoList_self/config"
	"ToDoList_self/repository/db/model"
	"ToDoList_self/types"
	"context"
	"gorm.io/gorm"
)

type TaskDao struct {
	*gorm.DB
}

// 创建一个可被追踪链路的上下文
func NewTaskDao(ctx context.Context) *TaskDao {
	if ctx == nil {
		ctx = context.Background()
	}
	return &TaskDao{config.NewDBClient(ctx)}
}

func (dao *TaskDao) CreateTask(userInfo *model.Task) (err error) {
	err = dao.Model(&model.Task{}).Create(userInfo).Error
	return
}

func (s *TaskDao) ListTask(start, limit int, userId uint) (r []*model.Task, total int64, err error) {
	//Preload:外键预加载
	err = s.Model(&model.Task{}).Preload("User").Where("uid = ?", userId).
		Count(&total).
		Limit(limit).Offset((start - 1) * limit).
		Find(&r).Error

	return
}

func (s *TaskDao) ShowTask(id, userId uint) (r *model.Task, err error) {
	err = s.Model(&model.Task{}).Preload("User").
		Where("id = ? AND uid = ?", id, userId).
		First(&r).Error
	return
}

func (s *TaskDao) UpdateTask(id, userId uint, req *types.UpdateTasksReq) error {
	r := new(model.Task)
	err := s.Model(&model.Task{}).Preload("User").
		Where("id = ? AND uid = ?", id, userId).
		First(&r).Error
	if err != nil {
		return err
	}
	if req.Status != 0 {
		r.Status = req.Status
	}

	if req.Title != "" {
		r.Title = req.Title
	}

	if req.Content != "" {
		r.Content = req.Content
	}
	return s.Save(r).Error
}

func (s *TaskDao) SearchTaskByInfo(userId uint, info string) (r []*model.Task, total int64, err error) {
	err = s.Model(&model.Task{}).
		Where("title LIKE ? OR content LIKE ? AND uid = ?", "%"+info+"%", "%"+info+"%", userId).
		Count(&total).
		Find(&r).Error
	return
}

func (s *TaskDao) DeleteTaskById(id, userId uint) error {
	r, err := s.ShowTask(id, userId)
	if err != nil {
		return err
	}
	return s.Delete(&r).Error
}
