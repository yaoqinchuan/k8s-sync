package dao

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"k8s-sync/internal/service/internal/do"
)

type WorkspaceDao struct {
}

func (*WorkspaceDao) GetWorkspaceByName(ctx context.Context, name string) (*do.WorkspaceDo, error) {
	connect := g.DB("default")
	record, err := connect.GetOne(ctx, fmt.Sprint("select * from workspace where name = ?"), name)
	if err != nil {
		return nil, err
	}
	var result = &do.WorkspaceDo{}
	err = record.Struct(result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (*WorkspaceDao) GetMiddleStatusWorkspace(ctx context.Context) ([]*do.WorkspaceDo, error) {
	connect := g.DB("default")
	record, err := connect.GetArray(ctx, fmt.Sprint("select * from workspace where status in ('STARTING STOPPING DELETING RESTORING') limit 100"))
	if err != nil {
		return nil, err
	}
	var result []*do.WorkspaceDo
	for i := 0; i < len(record); i++ {
		var tmpResult = &do.WorkspaceDo{}
		err = record[i].Struct(tmpResult)
		if err != nil {
			return nil, err
		}
		result = append(result, tmpResult)
	}

	return result, nil
}

func (*WorkspaceDao) GetWorkspaceById(ctx context.Context, workspaceId int64) (*do.WorkspaceDo, error) {
	connect := g.DB("default")
	record, err := connect.GetOne(ctx, fmt.Sprint("select * from workspace where id = ?"), workspaceId)
	if err != nil {
		return nil, err
	}
	var result = &do.WorkspaceDo{}
	err = record.Struct(result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
func (*WorkspaceDao) AddWorkspace(ctx context.Context, data *do.WorkspaceDo) (int64, error) {
	connect := g.DB("default")
	result, err := connect.Save(ctx, "workspace", data)
	if err != nil {
		return -1, err
	}
	line, _ := result.LastInsertId()
	return line, nil
}
func (*WorkspaceDao) UpdateWorkspaceById(ctx context.Context, updateMap *gdb.Map, workspaceId int64) (int64, error) {
	connect := g.DB("default")
	result, err := connect.Update(ctx, "workspace", updateMap, "id=?", workspaceId)
	if err != nil {
		return -1, err
	}
	rowAffect, _ := result.RowsAffected()
	return rowAffect, nil
}
func (*WorkspaceDao) UpdateWorkspaceByName(ctx context.Context, updateMap *gdb.Map, name string) (int64, error) {
	connect := g.DB("default")
	result, err := connect.Update(ctx, "workspace", updateMap, "name=?", name)
	if err != nil {
		return -1, err
	}
	rowAffect, _ := result.RowsAffected()
	return rowAffect, nil
}
func (*WorkspaceDao) UpdateWorkspaceStatusById(ctx context.Context, workspaceStatus string, workspaceId int64) (int64, error) {
	connect := g.DB("default")
	result, err := connect.Update(ctx, "workspace", g.Map{
		"status": workspaceStatus,
	}, "id=?", workspaceId)
	if err != nil {
		return -1, err
	}
	rowAffect, _ := result.RowsAffected()
	return rowAffect, nil
}

func (*WorkspaceDao) UpdateWorkspaceStatusWithCASById(ctx context.Context, newStatus string, oldStatus string, workspaceId int64) (int64, error) {
	connect := g.DB("default")
	result, err := connect.Update(ctx, "workspace", g.Map{
		"status": newStatus,
	}, "id=? and status=?", workspaceId, oldStatus)
	if err != nil {
		return -1, err
	}
	rowAffect, _ := result.RowsAffected()
	return rowAffect, nil
}
func (*WorkspaceDao) DeleteById(ctx context.Context, workspaceId int64) (int64, error) {
	connect := g.DB("default")
	result, err := connect.Delete(ctx, "workspace", "id=?", workspaceId)
	if err != nil {
		return -1, err
	}
	rowAffect, _ := result.RowsAffected()
	return rowAffect, nil
}
