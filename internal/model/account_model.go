package model

import "github.com/gogf/gf/v2/os/gtime"

type AccountModel struct {
	// p: params gconv/c/json: map convert tag recognise by order  v:validate
	Id       uint        `json:"id" p:"id" c:"id" v:"require"      `
	UserName string      `json:"userName" p:"userName" c:"userName"`
	UserId   string      `json:"userId" p:"userId" c:"userId" v:"require"`
	UserRole string      `json:"userRole" p:"userRole" c:"userRole" v:"require"`
	Email    string      `json:"email" p:"email" c:"email" v:"require"`
	Modifier string      `json:"modifier" p:"modifier" c:"modifier"`
	Deleted  int         `json:"deleted" p:"deleted" c:"deleted"`
	CreateAt *gtime.Time `json:"createAt" p:"createAt" c:"createAt" `
	UpdateAt *gtime.Time `json:"updateAt" p:"updateAt" c:"updateAt"`
}
