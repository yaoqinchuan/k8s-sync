package tasks

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gtime"
	"k8s-sync/internal/consts"
	"k8s-sync/internal/model"
	"k8s-sync/internal/service/internal/k8s"
	"k8s-sync/internal/service/manager"
	"k8s-sync/internal/utils"
)

type WorkspaceRestoringHandler struct {
	WorkspaceHandlerInterface
}

var (
	workspaceService = manager.WorkspaceService{}
)

func (workspaceRestoringHandler *WorkspaceRestoringHandler) PreExec(ctx context.Context, workspaceHandlerParams *model.WorkspaceHandlerParams) error {
	workspaceHandlerParams.TaskStatus = consts.WORKSPACE_TASK_PRE_EXECUTE
	return nil
}
func (workspaceRestoringHandler *WorkspaceRestoringHandler) DoExec(ctx context.Context, workspaceHandlerParams *model.WorkspaceHandlerParams) error {
	workspaceHandlerParams.TaskStatus = consts.WORKSPACE_TASK_EXECUTING
	workspace, err := workspaceService.GetWorkspaceById(ctx, workspaceHandlerParams.WorkspaceId)
	if err != nil {
		return gerror.New(fmt.Sprintf(" get workspace %v info failed, error %v", workspaceHandlerParams.WorkspaceName, err.Error()))
	}
	running, err := k8s.CheckWorkspaceRunning(ctx, k8s.ClientSet, workspace.Input)
	if err != nil {
		return gerror.New(fmt.Sprintf(" check workspace %v status failed, error %v", workspaceHandlerParams.WorkspaceName, err.Error()))
	}
	if running {
		workspaceHandlerParams.TaskStatus = consts.WORKSPACE_TASK_SUCCESS
		workspaceHandlerParams.EndTime = gtime.Now()
	}
	return nil
}
func (workspaceRestoringHandler *WorkspaceRestoringHandler) DoOnSuccess(ctx context.Context, workspaceHandlerParams *model.WorkspaceHandlerParams) error {
	effectLine, err := workspaceService.UpdateWorkspaceStatusWithCASById(ctx, consts.WS_RUNNING, consts.WS_RESTORING, workspaceHandlerParams.WorkspaceId)
	if err != nil {
		return gerror.New(fmt.Sprintf(" update workspace %v to running status failed, error %v", workspaceHandlerParams.WorkspaceName, err.Error()))
	}
	if effectLine == 0 {
		utils.Logger.Warning(ctx, fmt.Sprintf("update workspace %v to running status by cas failed", workspaceHandlerParams.WorkspaceName))
	}
	return nil
}
func (workspaceRestoringHandler *WorkspaceRestoringHandler) DoOnError(ctx context.Context, workspaceHandlerParams *model.WorkspaceHandlerParams) error {
	workspaceHandlerParams.TaskStatus = consts.WORKSPACE_TASK_ERROR
	workspaceHandlerParams.EndTime = gtime.Now()
	effectLine, err := workspaceService.UpdateWorkspaceStatusWithCASById(ctx, consts.WS_ERROR, consts.WS_RESTORING, workspaceHandlerParams.WorkspaceId)
	if err != nil {
		return gerror.New(fmt.Sprintf("update workspace %v to stopped status failed, error %v", workspaceHandlerParams.WorkspaceName, err.Error()))
	}
	if effectLine == 0 {
		utils.Logger.Warning(ctx, fmt.Sprintf("update workspace %v to error status by cas failed", workspaceHandlerParams.WorkspaceName))
	}
	return nil
}
func (workspaceRestoringHandler *WorkspaceRestoringHandler) CheckTimeout(ctx context.Context, workspaceHandlerParams *model.WorkspaceHandlerParams) bool {
	now := gtime.Now().Timestamp()
	if now-workspaceHandlerParams.StartTime.Timestamp() >= consts.WS_RESTORING_TIMEOUT_PERIOD {
		return true
	}
	return false
}

func (workspaceRestoringHandler *WorkspaceRestoringHandler) DoOnTimeout(ctx context.Context, workspaceHandlerParams *model.WorkspaceHandlerParams) error {
	effectLine, err := workspaceService.UpdateWorkspaceStatusWithCASById(ctx, consts.WS_STARTING_TIMEOUT, consts.WS_RESTORING, workspaceHandlerParams.WorkspaceId)
	workspaceHandlerParams.WorkspaceOldStatus = consts.WORKSPACE_TASK_TIMEOUT
	workspaceHandlerParams.EndTime = gtime.Now()
	if err != nil {
		return err
	}
	if effectLine == 0 {
		utils.Logger.Warning(ctx, fmt.Sprintf("update workspace %v to stop timeout status by cas failed", workspaceHandlerParams.WorkspaceName))
	}
	return nil
}
