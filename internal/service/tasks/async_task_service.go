package tasks

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/os/gtimer"
	"k8s-sync/internal/service/internal/dao"
	"k8s-sync/internal/service/internal/do"
	"k8s-sync/internal/utils"
	"os"
	"os/exec"
	"time"
)

type AsyncTaskService struct {
	AsyncTaskInterface
}

var (
	asyncTaskDao = dao.AsyncTaskDao{}
)

func (asyncTaskService *AsyncTaskService) ExecTask(ctx context.Context, task *do.AsyncTask) {

	defer func() {
		if err := recover(); err != nil {

			if task.RetryTime <= task.TotalRetryTime {
				asyncTaskDao.UpdateAsyncTaskById(ctx, &g.Map{
					"status":     TASK_STAUS_PENDING,
					"update_at":  gtime.Now(),
					"retry_time": task.RetryTime + 1,
					"error_info": err,
				}, task.Id)
			} else {
				asyncTaskDao.UpdateAsyncTaskById(ctx, &g.Map{
					"task_end_time": gtime.Now(),
					"status":        TASK_STATUS_FAILED,
					"update_at":     gtime.Now(),
					"retry_time":    task.RetryTime + 1,
					"deleted":       1,
					"error_info":    err,
				}, task.Id)
				asyncTaskService.DoOnError(ctx, task)
			}
			utils.Logger.Error(ctx, fmt.Sprintf("execute task %v failed, time is %v, error %v", task.Id, task.RetryTime, err))
		}
	}()
	taskPid := os.Getpid()

	gtimer.AddTimes(ctx, time.Second, 10, func(ctx context.Context) {
		if time.Duration(time.Now().UnixNano()-task.TaskStartTime.UnixNano()).Seconds() > task.TaskTimeoutTime {
			cmd := exec.Command("/bin/bash", "-c", fmt.Sprintf("kill -9 %v", taskPid))
			asyncTaskDao.UpdateAsyncTaskById(ctx, &g.Map{
				"task_start_time": gtime.Now(),
				"status":          TASK_STATUS_TIMEOUT,
				"update_at":       gtime.Now(),
				"retry_time":      task.RetryTime + 1,
			}, task.Id)
			err := cmd.Run()
			if err != nil {
				utils.Logger.Error(ctx, fmt.Sprintf("stop task for timeout failed"))
				gtimer.Exit()
			}
		}
	})

	asyncTaskService.PreExec(ctx, task)
	asyncTaskDao.UpdateAsyncTaskById(ctx, &g.Map{
		"status":    TASK_STATUS_RUNNING,
		"update_at": gtime.Now(),
	}, task.Id)

	asyncTaskService.DoExec(ctx, task)

	asyncTaskService.DoOnSuccess(ctx, task)

	asyncTaskDao.UpdateAsyncTaskById(ctx, &g.Map{
		"task_end_time": gtime.Now(),
		"status":        TASK_STATUS_SUCCESS,
		"update_at":     gtime.Now(),
		"retry_time":    task.RetryTime + 1,
		"deleted":       1,
	}, task.Id)
}
