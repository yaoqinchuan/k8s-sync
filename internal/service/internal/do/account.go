package do

import (
	"github.com/gogf/gf/v2/os/gtime"
)

type Account struct {
	Id       uint        `json:"id"       `
	UserName string      `json:"user_name" `
	UserId   string      `json:"user_id" `
	UserRole string      `json:"user_role" `
	Email    string      `json:"email"`
	CreateAt *gtime.Time `json:"create_at" `
	UpdateAt *gtime.Time `json:"update_at" `
	Modifier string      `json:"modifier"`
	Deleted  int         `json:"deleted"`
}
