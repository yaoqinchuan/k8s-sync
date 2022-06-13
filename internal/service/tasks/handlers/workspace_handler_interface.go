package tasks

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/os/gtime"
	"k8s-sync/internal/consts"
	"k8s-sync/internal/model"
	"k8s-sync/internal/service/internal/do"
	"k8s-sync/internal/utils"
	"sync"
)

var (
	handlerMap = make(map[string]WorkspaceHandlerInterface, 5)
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
	handlerMap[consts.WS_STARTING] = &WorkspaceStartingHandler{}
	handlerMap[consts.WS_RESTORING] = &WorkspaceStartingHandler{}
	handlerMap[consts.WS_DELETING] = &WorkspaceDeletingHandler{}
	handlerMap[consts.WS_STOPPING] = &WorkspaceStoppingHandler{}
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
			StartTime:          gtime.Now(),
		}
		doWorkspaceRefreshTask(ctx, params, handler)
	}()
	wg.Wait()
}

func BatchDisPatchHandler(ctx context.Context, workspaceModels []*do.WorkspaceDo) {
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
				StartTime:          gtime.Now(),
			}
			doWorkspaceRefreshTask(ctx, params, handler)
		}()
	}
	wg.Wait()
}
