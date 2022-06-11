package model

import "github.com/gogf/gf/v2/os/gtime"

type WorkspaceModel struct {
	Id         int64               `json:"id" v:"required" c:"id" p:"id"`
	Name       string              `json:"name" v:"required" c:"name" p:"name"`
	Attributes string              `json:"attributes" v:"required" c:"attributes" p:"attributes"`
	Spec       string              `json:"spec" v:"required" c:"spec" p:"spec"`
	Input      *WorkspaceSpecModel `json:"input" v:"required" c:"spec" p:"spec"`
	Runtime    string              `json:"runtime" v:"required" c:"runtime" p:"runtime"`
	Status     string              `json:"status" v:"required" c:"status" p:"status"`
	Temporary  int                 `json:"temporary" v:"required" c:"temporary" p:"temporary"`
	Creator    string              `json:"creator" v:"required" c:"creator" p:"creator"`
	Modifier   string              `json:"modifier" v:"required" c:"modifier" p:"modifier"`
	Deleted    int                 `json:"deleted" v:"required" c:"deleted" p:"deleted"`
	CreateAt   *gtime.Time         `json:"create_at" v:"required" c:"createAt" p:"createAt"`
	UpdateAt   *gtime.Time         `json:"update_at" v:"required" c:"updateAt" p:"updateAt"`
}
