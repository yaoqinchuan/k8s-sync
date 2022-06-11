package tasks

import (
	"context"
	"k8s-sync/internal/service/internal/dao"
	"k8s-sync/internal/service/internal/do"
	"k8s-sync/internal/service/tasks/handlers"
)

type WorkspaceAsyncTaskImpl struct {
	AsyncTaskInterface
}

var (
	workspaceDao     = dao.WorkspaceDao{}
	workspaceHandler = handlers.WorkspaceHandler{}
)

func (*WorkspaceAsyncTaskImpl) PreExec(ctx context.Context, task *do.AsyncTask) error {
	return nil
}
func (*WorkspaceAsyncTaskImpl) DoExec(ctx context.Context, task *do.AsyncTask) error {
	workspaces, err := workspaceDao.GetMiddleStatusWorkspace(ctx)
	if err != nil {
		panic(err)
	}
	workspaceHandler.BatchDisPatchHandler(ctx, workspaces)

}
func (*WorkspaceAsyncTaskImpl) DoOnSuccess(ctx context.Context, task *do.AsyncTask) error {
	return nil
}
func (*WorkspaceAsyncTaskImpl) DoOnError(ctx context.Context, task *do.AsyncTask) error {
	return nil
}
func (*WorkspaceAsyncTaskImpl) DoOnTimeout(ctx context.Context, task *do.AsyncTask) error {
	return nil
}
