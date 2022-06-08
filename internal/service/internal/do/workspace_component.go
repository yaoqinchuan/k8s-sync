package do

import "github.com/gogf/gf/v2/os/gtime"

type WorkspaceComponentDo struct {
	Id          int         `json:"id"`
	WorkspaceId int         `json:"workspace_id"`
	Kind        string      `json:"kind"`
	Name        string      `json:"name"`
	Des         string      `json:"desc"`
	Deleted     int         `json:"deleted"`
	CreateAt    *gtime.Time `json:"create_at"`
	UpdateAt    *gtime.Time `json:"update_at"`
}
