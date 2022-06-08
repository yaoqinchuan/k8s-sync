package dao

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"k8s-sync/internal/service/internal/do"
)

type WorkspaceComponentDao struct {
}

func (*WorkspaceComponentDao) GetWorkspaceComponentByWorkspaceId(ctx context.Context, workspaceId int) (*do.WorkspaceComponentDo, error) {
	connect := g.DB("default")
	record, err := connect.GetArray(ctx, fmt.Sprint("select * from workspace_component where workspace_id = ?"), workspaceId)
	if err != nil {
		return nil, err
	}
	result := &[]do.WorkspaceComponentDo{}
	for i := 0; i < len(record); i++ {
		var tmpResult = &do.WorkspaceComponentDo{}
		err = record[i].Struct(tmpResult)
		if err != nil {
			return nil, err
		}
		result = append(result, tmpResult)
	}

	return result, nil
}

func (*WorkspaceComponentDao) GetWorkspaceComponentByWorkspaceIdAndKind(ctx context.Context, workspaceId int, kind string) (*do.WorkspaceComponentDo, error) {
	connect := g.DB("default")
	record, err := connect.GetOne(ctx, fmt.Sprint("select * from workspace_component where workspace_id = ? and kind=?"), workspaceId, kind)
	if err != nil {
		return nil, err
	}
	var result = &do.WorkspaceComponentDo{}
	err = record.Struct(result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (*WorkspaceComponentDao) AddWorkspaceComponent(ctx context.Context, data *do.WorkspaceComponentDo) (int64, error) {
	connect := g.DB("default")
	result, err := connect.Save(ctx, "workspace_component", data)
	if err != nil {
		return -1, err
	}
	line, _ := result.LastInsertId()
	return line, nil
}
func (*WorkspaceDao) UpdateWorkspaceComponentByWorkspaceIdAndKind(ctx context.Context, updateMap *gdb.Map, workspaceId int,
	kind string) (int64, error) {
	connect := g.DB("default")
	result, err := connect.Update(ctx, "workspace_component", updateMap, "workspace_id=? and kind=?", workspaceId, kind)
	if err != nil {
		return -1, err
	}
	rowAffect, _ := result.RowsAffected()
	return rowAffect, nil
}

func (*WorkspaceDao) DeleteById(ctx context.Context, workspaceId int) (int64, error) {
	connect := g.DB("default")
	result, err := connect.Delete(ctx, "workspace", "id=?", workspaceId)
	if err != nil {
		return -1, err
	}
	rowAffect, _ := result.RowsAffected()
	return rowAffect, nil
}
