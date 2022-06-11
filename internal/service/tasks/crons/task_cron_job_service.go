package tasks

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/os/gcron"
	"github.com/gogf/gf/v2/os/gctx"
	"k8s-sync/internal/service/internal/dao"
	"k8s-sync/internal/service/tasks"
)

const (
	ASYNC_TASK_INTERVAL  = 3
	ASYNC_TASK_BATCHSIZE = 100
)

func asyncTask(ctx context.Context) {
	asyncDao := dao.AsyncTaskDao{}
	result, err := asyncDao.GetAsyncTaskByTaskNameAndStatus(ctx, ASYNC_TASK_BATCHSIZE, tasks.TASK_NAME_WORKSPACE_OPS, tasks.TASK_STAUS_PENDING)
}

func init() {
	gcron.Add(gctx.New(), fmt.Sprintf("*/%v * * * * *", ASYNC_TASK_INTERVAL), func(ctx context.Context) {

	})
}
