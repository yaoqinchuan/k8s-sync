package tasks

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/os/gcron"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gtime"
	"k8s-sync/internal/service/internal/dao"
	"k8s-sync/internal/service/internal/do"
	"k8s-sync/internal/service/tasks"
	"k8s-sync/internal/utils"
)

const (
	ASYNC_TASK_INTERVAL   = 3
	WORKSPACE_OPS_TIMEOUT = 30 * 60
)

var (
	asyncTaskService = tasks.AsyncTaskService{}
	asyncDao         = dao.AsyncTaskDao{}
)

func asyncWorkspaceTask(ctx context.Context) {
	result, err := asyncDao.GetReadyAsyncTaskByName(ctx, tasks.TASK_NAME_WORKSPACE_OPS)
	if err != nil {
		utils.Logger.Error(ctx, fmt.Sprintf("query async task %v failed", tasks.TASK_NAME_WORKSPACE_OPS))
		return
	}
	localIP := utils.GetLocalIp()
	if localIP == "" {
		utils.Logger.Error(ctx, "get local ip failed.")
		return
	}
	if result == nil {
		_, err := asyncDao.AddAsyncTask(ctx, &do.AsyncTask{
			TaskName:        tasks.TASK_NAME_WORKSPACE_OPS,
			TaskAttributes:  "",
			Ip:              "",
			Status:          tasks.TASK_STAUS_PENDING,
			ErrorInfo:       "",
			RetryTime:       1,
			TotalRetryTime:  1,
			TaskStartTime:   nil,
			TaskTimeoutTime: WORKSPACE_OPS_TIMEOUT,
			TaskEndTime:     nil,
			CreateAt:        gtime.Now(),
			UpdateAt:        gtime.Now(),
			Deleted:         0,
		})
		if err != nil {
			return
		}
		return
	} else {
		cas, err := asyncDao.UpdateAsyncTaskIpByIdAndStatusCAS(ctx, utils.GetLocalIp(), result.Ip,
			tasks.TASK_STAUS_PENDING, result.Status, gtime.Now(), WORKSPACE_OPS_TIMEOUT, result.Id)
		if err != nil {
			return
		}
		if cas == 0 {
			return
		}
	}
	asyncTaskService.ExecTask(ctx, result)
}

func refreshTimeoutTask(ctx context.Context) {
	err := asyncDao.UpdateTimeoutAsyncTaskToReady(ctx)
	if err != nil {
		utils.Logger.Error(ctx, fmt.Sprintf("update timeout async task %v failed", tasks.TASK_NAME_WORKSPACE_OPS))
		return
	}
}
func init() {
	_, err := gcron.Add(gctx.New(), fmt.Sprintf("*/10 * * * * *"), func(ctx context.Context) {
		refreshTimeoutTask(ctx)
	}, tasks.TASK_NAME_TASK_TIMEOUT_REFRESH)
	if err != nil {
		return
	}

	_, err = gcron.Add(gctx.New(), fmt.Sprintf("*/%v * * * * *", ASYNC_TASK_INTERVAL), func(ctx context.Context) {
		asyncWorkspaceTask(ctx)
	}, tasks.TASK_NAME_WORKSPACE_OPS)
	if err != nil {
		return
	}
}
