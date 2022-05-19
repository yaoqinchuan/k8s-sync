package model

import "github.com/gogf/gf/v2/os/gtime"

type AccountModel struct {
	Id       uint        `json:"id" c:"id"       `
	UserName string      `json:"user_name"  c:"user_name"`
	UserId   string      `json:"user_id" c:"user_id"`
	UserRole int         `json:"user_role" c:"user_role"`
	Email    string      `json:"email" c:"email"`
	CreateAt *gtime.Time `json:"create_at" c:"create_at" `
	UpdateAt *gtime.Time `json:"update_at" c:"update_at"`
}
