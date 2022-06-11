package model

import "github.com/gogf/gf/v2/os/gtime"

type WorkspaceHandlerParams struct {
	WorkspaceId        int64       `json:"workspace_id"`
	WorkspaceName      string      `json:"workspace_name"`
	WorkspaceOldStatus string      `json:"workspace_old_name"`
	TaskStatus         string      `json:"task_status"`
	StartTime          *gtime.Time `json:"start_time"`
	EndTime            *gtime.Time `json:"end_time"`
}
