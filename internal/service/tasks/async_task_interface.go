package tasks

import (
	"context"
	"k8s-sync/internal/service/internal/do"
)

type AsyncTaskInterface interface {
	PreExec(ctx context.Context, task *do.AsyncTask) error
	DoExec(ctx context.Context, task *do.AsyncTask) error
	DoOnSuccess(ctx context.Context, task *do.AsyncTask) error
	DoOnError(ctx context.Context, task *do.AsyncTask) error
	DoOnTimeout(ctx context.Context, task *do.AsyncTask) error
}
