package types

type CreateReq struct {
	Title   string `form:"title"`
	Status  int    `form:"status"`
	Content string `form:"content"`
}

type ListTasksReq struct {
	Limit int `form:"limit"`
	Start int `form:"start"`
}

type ShowTasksReq struct {
	Uid uint `form:"id"`
}

type UpdateTasksReq struct {
	ID      uint   `form:"id"`
	Title   string `form:"title"`
	Status  int    `form:"status"`
	Content string `form:"content"`
}
type SearchTasksReq struct {
	Info string `form:"info"`
}
type DeleteTasksReq struct {
	ID uint `form:"id"`
}

type FellowReq struct {
	ID   uint `form:"to_user_id"`
	Type uint `form:"action_type"`
}
type ListTasksTx struct {
	Uid       uint
	Title     string
	Status    int
	Content   string
	StartTime int64
	CreatedAt int64
	EndTime   int64
}
type ShowTasksTx struct {
	Id        uint
	Title     string
	Status    int
	ViewCount uint64
	Content   string
	StartTime int64
	CreatedAt int64
	EndTime   int64
}
