package handlers

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/os/gtime"
	"k8s-sync/internal/model"
	"k8s-sync/internal/service"
	"k8s-sync/internal/utils"
	"sync"
)

const (
	WORKSPACE_TASK_PRE_EXECUTE = "PRE_EXECUTE"
	WORKSPACE_TASK_EXECUTING   = "EXECUTING"
	WORKSPACE_TASK_SUCCESS     = "SUCCESS"
	WORKSPACE_TASK_ERROR       = "ERROR"
	WORKSPACE_TASK_TIMEOUT     = "TIMEOUT"
	WORKSPACE_TASK_FINISH      = "FINISH"
)

var (
	handlerMap       = make(map[string]WorkspaceHandlerInterface, 5)
	workspaceService = service.WorkspaceService{}
)

const (
	WS_STARTING_TIMEOUT_PERIOD  = 600
	WS_RESTORING_TIMEOUT_PERIOD = 600
	WS_STOPPING_TIMEOUT_PERIOD  = 600
	WS_DELETING_TIMEOUT_PERIOD  = 600

	WS_PENDING = "PENDING"

	WS_STARTING         = "STARTING"
	WS_RESTORING        = "RESTORING"
	WS_STARTING_TIMEOUT = "STARTING_TIMEOUT"
	WS_RUNNING          = "RUNNING"

	WS_DELETING         = "DELETING"
	WS_DELETED          = "DELETED"
	WS_DELETING_TIMEOUT = "DELETING_TIMEOUT"

	WS_STOPPING         = "STOPPING"
	WS_STOPPED          = "STOPPED"
	WS_STOPPING_TIMEOUT = "STOPPING_TIMEOUT"
	WS_ERROR            = "ERROR"
)

type WorkspaceHandlerInterface interface {
	PreExec(ctx context.Context, workspaceHandlerParams *model.WorkspaceHandlerParams) error
	DoExec(ctx context.Context, workspaceHandlerParams *model.WorkspaceHandlerParams) error
	DoOnSuccess(ctx context.Context, workspaceHandlerParams *model.WorkspaceHandlerParams) error
	DoOnError(ctx context.Context, workspaceHandlerParams *model.WorkspaceHandlerParams) error
	DoOnTimeout(ctx context.Context, workspaceHandlerParams *model.WorkspaceHandlerParams) error
	CheckTimeout(ctx context.Context, workspaceHandlerParams *model.WorkspaceHandlerParams) bool
}

func init() {
	handlerMap[WS_STARTING] = &WorkspaceStartingHandler{}
	handlerMap[WS_RESTORING] = &WorkspaceStartingHandler{}
	handlerMap[WS_DELETING] = &WorkspaceDeletingHandler{}
	handlerMap[WS_STOPPING] = &WorkspaceStoppingHandler{}
}

func doWorkspaceRefreshTask(ctx context.Context, workspaceHandlerParams *model.WorkspaceHandlerParams, workspaceHandlerInterface WorkspaceHandlerInterface) {
	defer func() {
		if err := recover(); err != nil {
			err := workspaceHandlerInterface.DoOnError(ctx, workspaceHandlerParams)
			if err != nil {
				utils.Logger.Error(ctx, fmt.Sprintf("workspace %v  do on error failed, error %v", workspaceHandlerParams.WorkspaceName, err.Error()))
			}
		}
	}()

	if workspaceHandlerInterface.CheckTimeout(ctx, workspaceHandlerParams) {
		if err := workspaceHandlerInterface.DoOnTimeout(ctx, workspaceHandlerParams); err != nil {
			panic(err)
		}
	}
	if err := workspaceHandlerInterface.DoExec(ctx, workspaceHandlerParams); err != nil {
		panic(err)
	}
	if err := workspaceHandlerInterface.DoOnSuccess(ctx, workspaceHandlerParams); err != nil {
		panic(err)
	}

}

func DisPatchHandler(ctx context.Context, workspaceModel *model.WorkspaceModel) {
	if nil == workspaceModel {
		return
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer func() {
			wg.Done()
		}()
		handler := handlerMap[workspaceModel.Status]
		params := &model.WorkspaceHandlerParams{
			WorkspaceId:        workspaceModel.Id,
			WorkspaceName:      workspaceModel.Name,
			WorkspaceOldStatus: workspaceModel.Status,
			TaskStatus:         WORKSPACE_TASK_PRE_EXECUTE,
			StartTime:          gtime.Now(),
		}
		doWorkspaceRefreshTask(ctx, params, handler)
	}()
	wg.Wait()
}

func BatchDisPatchHandler(ctx context.Context, workspaceModels []*model.WorkspaceModel) {
	if nil == workspaceModels || 0 == len(workspaceModels) {
		return
	}
	var wg sync.WaitGroup
	wg.Add(len(workspaceModels))
	for i := 0; i < len(workspaceModels); i++ {
		index := i
		go func() {
			defer func() {
				wg.Done()
			}()
			handler := handlerMap[workspaceModels[index].Status]
			params := &model.WorkspaceHandlerParams{
				WorkspaceId:        workspaceModels[index].Id,
				WorkspaceName:      workspaceModels[index].Name,
				WorkspaceOldStatus: workspaceModels[index].Status,
				TaskStatus:         WORKSPACE_TASK_PRE_EXECUTE,
				StartTime:          gtime.Now(),
			}
			doWorkspaceRefreshTask(ctx, params, handler)
		}()
	}
	wg.Wait()
}
