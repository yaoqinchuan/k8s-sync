package tasks

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcron"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/os/gtimer"
	"k8s-sync/internal/service/internal/dao"
	"k8s-sync/internal/service/internal/do"
	"k8s-sync/internal/service/tasks"
	"k8s-sync/internal/utils"
	"os"
	"os/exec"
	"time"
)

const (
	ASYNC_TASK_INTERVAL       = 3
	WORKSPACE_OPS_TIMEOUT     = 30 * 60
	IP_HEART_BEAT_LIFE_PERIOD = 60
)

var (
	asyncTaskService = tasks.AsyncTaskService{}
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
	taskPid := os.Getpid()
	gtimer.AddTimes(ctx, time.Second, 10, func(ctx context.Context) {
		if time.Duration(time.Now().UnixNano()-result.TaskStartTime.UnixNano()).Seconds() > result.TaskTimeoutTime {
			cmd := exec.Command("/bin/bash", "-c", fmt.Sprintf("kill -9 %v", taskPid))
			asyncDao.UpdateAsyncTaskById(ctx, &g.Map{
				"status":        WORKSPACE_OPS_TIMEOUT,
				"update_at":     gtime.Now(),
				"task_end_time": gtime.Now(),
				"error_info":    fmt.Sprintf("task %v timeout", result.TaskName),
			}, result.Id)
			err := cmd.Run()
			if err != nil {
				utils.Logger.Error(ctx, fmt.Sprintf("stop task for timeout failed"))
				gtimer.Exit()
			}
		}
	})
	asyncTaskService.ExecTask(ctx, result)
}

func init() {
	gcron.Add(gctx.New(), fmt.Sprintf("*/%v * * * * *", ASYNC_TASK_INTERVAL), func(ctx context.Context) {

	})
}
