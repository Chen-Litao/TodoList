package types

type RegisterReq struct {
	User     string `form:"user"`
	Password string `form:"password"`
}
type LoginReq struct {
	User     string `form:"user" json:"user"`
	Password string `form:"password" json:"password"`
}
