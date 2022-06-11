package tasks

import (
	"context"
	"k8s-sync/internal/service/internal/do"
)

const (
	TASK_STAUS_PENDING  = "Pending"
	TASK_STATUS_RUNNING = "Running"
	TASK_STATUS_FAILED  = "Failed"
	TASK_STATUS_TIMEOUT = "Timeout"
	TASK_STATUS_SUCCESS = "Success"

	TASK_NAME_WORKSPACE_OPS = "TASK_NAME_WORKSPACE_OPS"
)

type AsyncTaskInterface interface {
	PreExec(ctx context.Context, task *do.AsyncTask)
	DoExec(ctx context.Context, task *do.AsyncTask)
	DoOnSuccess(ctx context.Context, task *do.AsyncTask)
	DoOnError(ctx context.Context, task *do.AsyncTask)
	DoOnTimeout(ctx context.Context, task *do.AsyncTask)
}
