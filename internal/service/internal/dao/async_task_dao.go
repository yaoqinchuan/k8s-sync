package dao

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"k8s-sync/internal/service/internal/do"
	"k8s-sync/internal/service/tasks"
)

type AsyncTaskDao struct {
}

func (*AsyncTaskDao) GetReadyAsyncTaskByName(ctx context.Context, name string) (*do.AsyncTask, error) {
	connect := g.DB("default")
	record, err := connect.GetOne(ctx, fmt.Sprint("select * from async_task where name = ? and status in (?, ?) and deleted = 0"), name, tasks.TASK_STAUS_PENDING, tasks.TASK_STAUS_PENDING)
	if err != nil {
		return nil, err
	}
	if record == nil {
		return nil, nil
	}
	var result = &do.AsyncTask{}
	err = record.Struct(result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (*AsyncTaskDao) UpdateTimeoutAsyncTaskToReady(ctx context.Context) error {
	connect := g.DB("default")
	_, err := connect.Update(ctx, "async_task", "status=?", "deleted = 0 and task_timeout_time>NOW()", tasks.TASK_STAUS_PENDING)
	if err != nil {
		return err
	}
	return nil
}

func (*AsyncTaskDao) GetAsyncTaskById(ctx context.Context, id int64) (*do.AsyncTask, error) {
	connect := g.DB("default")
	record, err := connect.GetOne(ctx, fmt.Sprint("select * from async_task where id = ?"), id)
	if err != nil {
		return nil, err
	}
	var result = &do.AsyncTask{}
	err = record.Struct(result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (*AsyncTaskDao) AddAsyncTask(ctx context.Context, data *do.AsyncTask) (int64, error) {
	connect := g.DB("default")
	result, err := connect.Save(ctx, "async_task", data)
	if err != nil {
		return -1, err
	}
	line, _ := result.LastInsertId()
	return line, nil
}
func (*AsyncTaskDao) UpdateAsyncTaskById(ctx context.Context, updateMap *gdb.Map, id int64) (int64, error) {
	connect := g.DB("default")
	result, err := connect.Update(ctx, "async_task", updateMap, "id=?", id)
	if err != nil {
		return -1, err
	}
	rowAffect, _ := result.RowsAffected()
	return rowAffect, nil
}

func (*AsyncTaskDao) UpdateAsyncTaskIpByIdAndStatusCAS(ctx context.Context, newIp string, oldIp string,
	newStatus string, oldStatus string, taskStartTime *gtime.Time, taskTimeoutTime int, id int64) (int64, error) {
	connect := g.DB("default")
	result, err := connect.Update(ctx, "async_task", "ip=?,status=?,task_start_time=?,task_timeout_time=?", " ip=? and id=? and status=?",
		newIp, newStatus, taskStartTime, taskTimeoutTime, oldIp, id, oldStatus)
	if err != nil {
		return -1, err
	}
	rowAffect, _ := result.RowsAffected()
	return rowAffect, nil
}

func (*AsyncTaskDao) DeleteByUserId(ctx context.Context, id int64) (int64, error) {
	connect := g.DB("default")
	result, err := connect.Update(ctx, "async_task", "deleted=1", "id=?", id)
	if err != nil {
		return -1, err
	}
	rowAffect, _ := result.RowsAffected()
	return rowAffect, nil
}
