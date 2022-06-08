package do

import "github.com/gogf/gf/v2/os/gtime"

type WorkspaceDo struct {
	Id         int         `json:"id"`
	Name       string      `json:"name"`
	Attributes string      `json:"attributes"`
	Input      string      `json:"input"`
	Spec       string      `json:"spec"`
	Runtime    string      `json:"runtime"`
	Status     string      `json:"status"`
	Temporary  int         `json:"temporary"`
	Creator    string      `json:"creator"`
	Modifier   string      `json:"modifier"`
	Deleted    int         `json:"deleted"`
	CreateAt   *gtime.Time `json:"create_at"`
	UpdateAt   *gtime.Time `json:"update_at"`
}
