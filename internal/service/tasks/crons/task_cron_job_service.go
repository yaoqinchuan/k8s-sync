package tasks

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/os/gcron"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/os/gtimer"
	"k8s-sync/internal/service/internal/dao"
	"k8s-sync/internal/service/internal/do"
	"k8s-sync/internal/service/tasks"
	"k8s-sync/internal/utils"
	"time"
)

const (
	ASYNC_TASK_INTERVAL       = 3
	WORKSPACE_OPS_TIMEOUT     = 30 * 60
	IP_HEART_BEAT_LIFE_PERIOD = 60
)

func asyncWorkspaceTask(ctx context.Context) {
	asyncDao := dao.AsyncTaskDao{}
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
	}

	gtimer.AddTimes(ctx, time.Second, 10, func(ctx context.Context) {

	})
}

func init() {
	gcron.Add(gctx.New(), fmt.Sprintf("*/%v * * * * *", ASYNC_TASK_INTERVAL), func(ctx context.Context) {

	})
}
