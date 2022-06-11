package handlers

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gtime"
	"k8s-sync/internal/model"
	"k8s-sync/internal/service/internal/k8s"
	"k8s-sync/internal/utils"
)

type WorkspaceDeletingHandler struct {
	WorkspaceHandlerInterface
}

func (workspaceDeletingHandler *WorkspaceDeletingHandler) PreExec(ctx context.Context, workspaceHandlerParams *model.WorkspaceHandlerParams) error {
	workspaceHandlerParams.TaskStatus = WORKSPACE_TASK_PRE_EXECUTE
	return nil
}
func (workspaceDeletingHandler *WorkspaceDeletingHandler) DoExec(ctx context.Context, workspaceHandlerParams *model.WorkspaceHandlerParams) error {
	workspaceHandlerParams.TaskStatus = WORKSPACE_TASK_EXECUTING
	workspace, err := workspaceService.GetWorkspaceById(ctx, workspaceHandlerParams.WorkspaceId)
	if err != nil {
		return gerror.New(fmt.Sprintf(" get workspace %v info failed, error %v", workspaceHandlerParams.WorkspaceName, err.Error()))
	}

	stopped, err := k8s.CheckWorkspaceDeleted(ctx, k8s.ClientSet, workspace.Input)
	if err != nil {
		return gerror.New(fmt.Sprintf(" check workspace %v status failed, error %v", workspaceHandlerParams.WorkspaceName, err.Error()))
	}
	if stopped {
		workspaceHandlerParams.TaskStatus = WORKSPACE_TASK_SUCCESS
		workspaceHandlerParams.EndTime = gtime.Now()
	}
	return nil
}
func (workspaceDeletingHandler *WorkspaceDeletingHandler) DoOnSuccess(ctx context.Context, workspaceHandlerParams *model.WorkspaceHandlerParams) error {
	effectLine, err := workspaceService.UpdateWorkspaceStatusWithCASById(ctx, WS_DELETED, WS_DELETING, workspaceHandlerParams.WorkspaceId)
	if err != nil {
		return gerror.New(fmt.Sprintf("update workspace %v to deleted status failed, error %v", workspaceHandlerParams.WorkspaceName, err.Error()))
	}
	if effectLine == 0 {
		utils.Logger.Warning(ctx, fmt.Sprintf("update workspace %v to deleted status by cas failed", workspaceHandlerParams.WorkspaceName))
	}
	return nil
}
func (workspaceDeletingHandler *WorkspaceDeletingHandler) DoOnError(ctx context.Context, workspaceHandlerParams *model.WorkspaceHandlerParams) error {
	workspaceHandlerParams.TaskStatus = WORKSPACE_TASK_ERROR
	workspaceHandlerParams.EndTime = gtime.Now()
	effectLine, err := workspaceService.UpdateWorkspaceStatusWithCASById(ctx, WS_ERROR, WS_DELETING, workspaceHandlerParams.WorkspaceId)
	if err != nil {
		return gerror.New(fmt.Sprintf("update workspace %v to deleted status failed, error %v", workspaceHandlerParams.WorkspaceName, err.Error()))
	}
	if effectLine == 0 {
		utils.Logger.Warning(ctx, fmt.Sprintf("update workspace %v to error status by cas failed", workspaceHandlerParams.WorkspaceName))
	}
	return nil
}
func (workspaceDeletingHandler *WorkspaceDeletingHandler) CheckTimeout(ctx context.Context, workspaceHandlerParams *model.WorkspaceHandlerParams) bool {
	now := gtime.Now().Timestamp()
	if now-workspaceHandlerParams.StartTime.Timestamp() >= WS_DELETING_TIMEOUT_PERIOD {
		return true
	}
	return false
}

func (workspaceDeletingHandler *WorkspaceDeletingHandler) DoOnTimeout(ctx context.Context, workspaceHandlerParams *model.WorkspaceHandlerParams) error {
	effectLine, err := workspaceService.UpdateWorkspaceStatusWithCASById(ctx, WS_DELETING_TIMEOUT, WS_DELETING, workspaceHandlerParams.WorkspaceId)
	workspaceHandlerParams.WorkspaceOldStatus = WORKSPACE_TASK_TIMEOUT
	workspaceHandlerParams.EndTime = gtime.Now()
	if err != nil {
		return err
	}
	if effectLine == 0 {
		utils.Logger.Warning(ctx, fmt.Sprintf("update workspace %v to delete timeout status by cas failed", workspaceHandlerParams.WorkspaceName))
	}
	return nil
}
