package dao

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"k8s-sync/internal/service/internal/do"
)

type AsyncTaskDao struct {
}

func (*AsyncTaskDao) GetAsyncTaskByTaskNameAndStatus(ctx context.Context, name, status string) (*[]do.AsyncTask, error) {
	connect := g.DB("default")
	record, err := connect.GetArray(ctx, fmt.Sprint("select * from async_task where task_name = ? and status = ? order by id asc limit 100"), name, status)
	if err != nil {
		return nil, err
	}
	var result []do.AsyncTask
	for i := 0; i < len(record); i++ {
		var tmpResult = &do.AsyncTask{}
		err = record[i].Struct(tmpResult)
		if err != nil {
			return nil, err
		}
		result = append(result, *tmpResult)
	}
	return &result, nil
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

func (*AsyncTaskDao) UpdateAsyncTaskIpByIdCAS(ctx context.Context, newIp string, oldIp string, id int64) (int64, error) {
	connect := g.DB("default")
	result, err := connect.Update(ctx, "async_task", "ip = ?", " ip=? and id=?", newIp, oldIp, id)
	if err != nil {
		return -1, err
	}
	rowAffect, _ := result.RowsAffected()
	return rowAffect, nil
}

func (*AsyncTaskDao) DeleteByUserId(ctx context.Context, id int64) (int64, error) {
	connect := g.DB("default")
	result, err := connect.Delete(ctx, "async_task", "id=?", id)
	if err != nil {
		return -1, err
	}
	rowAffect, _ := result.RowsAffected()
	return rowAffect, nil
}
