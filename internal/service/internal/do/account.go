package do

import (
	"github.com/gogf/gf/v2/os/gtime"
)

type Account struct {
	Id       uint        `json:"id"       `
	UserName string      `json:"config_key" `
	UserId   string      `json:"config_value" `
	UserRole int         `json:"describe" `
	Email    string      `json:"email"`
	CreateAt *gtime.Time `json:"create_at" `
	UpdateAt *gtime.Time `json:"update_at" `
}
