package tasks

import (
	"context"
	"k8s-sync/internal/service/internal/dao"
	"k8s-sync/internal/service/internal/do"
	"k8s-sync/internal/service/tasks"
)

type WorkspaceAsyncTaskImpl struct {
	AsyncTaskInterface
}

var (
	workspaceDao = dao.WorkspaceDao{}
)

func (*WorkspaceAsyncTaskImpl) PreExec(ctx context.Context, task *do.AsyncTask) {
}
func (*WorkspaceAsyncTaskImpl) DoExec(ctx context.Context, task *do.AsyncTask) {
	workspaces, err := workspaceDao.GetMiddleStatusWorkspace(ctx)
	if err != nil {
		panic(err)
	}
	tasks.BatchDisPatchHandler(ctx, workspaces)
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
