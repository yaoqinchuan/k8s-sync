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
	"time"
)

type AsyncTaskService struct {
	AsyncTaskInterface
}

var (
	asyncTaskDao = dao.AsyncTaskDao{}
)

const (
	TASK_STAUS_PENDING  = "Pending"
	TASK_STATUS_RUNNING = "Running"
	TASK_STATUS_FAILED  = "Failed"
	TASK_STATUS_TIMEOUT = "Timeout"
	TASK_STATUS_SUCCESS = "Success"
)

func (asyncTaskService *AsyncTaskService) ExecTask(ctx context.Context, task *do.AsyncTask) {

	defer func() {
		if err := recover(); err != nil {
			utils.Logger.Error(ctx, fmt.Sprintf("execute task %v failed, error %v", task.Id, err))
		}
	}()

	gtimer.AddTimes(ctx, time.Second, 10, func(ctx context.Context) {
		if time.Duration(time.Now().UnixNano()-task.TaskStartTime.UnixNano()).Seconds() > task.TaskTimeoutTime {

		}
	})

	err := asyncTaskService.PreExec(ctx, task)
	if err != nil {
		if task.RetryTime <= task.TotalRetryTime {
			asyncTaskDao.UpdateAsyncTaskById(ctx, &g.Map{
				"task_start_time": gtime.Now(),
				"status":          TASK_STAUS_PENDING,
				"update_at":       gtime.Now(),
				"retry_time":      task.RetryTime + 1,
			}, task.Id)
		} else {
			asyncTaskDao.UpdateAsyncTaskById(ctx, &g.Map{
				"error_info": err,
				"status":     TASK_STATUS_FAILED,
				"retry_time": task.RetryTime + 1,
				"update_at":  gtime.Now(),
			}, task.Id)
		}
		return
	}

	asyncTaskDao.UpdateAsyncTaskById(ctx, &g.Map{
		"status":    TASK_STATUS_RUNNING,
		"update_at": gtime.Now(),
	}, task.Id)
	err = asyncTaskService.DoExec(ctx, task)
	if err != nil {
		if task.RetryTime <= task.TotalRetryTime {
			asyncTaskDao.UpdateAsyncTaskById(ctx, &g.Map{
				"status":     TASK_STAUS_PENDING,
				"update_at":  gtime.Now(),
				"retry_time": task.RetryTime + 1,
			}, task.Id)
		} else {
			asyncTaskDao.UpdateAsyncTaskById(ctx, &g.Map{
				"error_info": err,
				"status":     TASK_STATUS_FAILED,
				"retry_time": task.RetryTime + 1,
				"update_at":  gtime.Now(),
			}, task.Id)
			asyncTaskService.DoOnError(ctx, task)
		}
		return
	}
	err = asyncTaskService.DoOnSuccess(ctx, task)
	if err != nil {
		if task.RetryTime <= task.TotalRetryTime {
			asyncTaskDao.UpdateAsyncTaskById(ctx, &g.Map{
				"status":     TASK_STAUS_PENDING,
				"update_at":  gtime.Now(),
				"retry_time": task.RetryTime + 1,
			}, task.Id)
		} else {
			asyncTaskDao.UpdateAsyncTaskById(ctx, &g.Map{
				"error_info": err,
				"status":     TASK_STATUS_FAILED,
				"retry_time": task.RetryTime + 1,
				"update_at":  gtime.Now(),
			}, task.Id)
		}
		return
	}

	asyncTaskDao.UpdateAsyncTaskById(ctx, &g.Map{
		"task_end_time": gtime.Now(),
		"status":        TASK_STATUS_SUCCESS,
		"update_at":     gtime.Now(),
		"retry_time":    task.RetryTime + 1,
	}, task.Id)
}
