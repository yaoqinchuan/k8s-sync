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

func (*WorkspaceComponentDao) GetWorkspaceComponentByWorkspaceId(ctx context.Context, workspaceId int64) (*[]do.WorkspaceComponentDo, error) {
	connect := g.DB("default")
	record, err := connect.GetArray(ctx, fmt.Sprint("select * from workspace_component where workspace_id = ?"), workspaceId)
	if err != nil {
		return nil, err
	}
	var result []do.WorkspaceComponentDo
	for i := 0; i < len(record); i++ {
		var tmpResult = &do.WorkspaceComponentDo{}
		err = record[i].Struct(tmpResult)
		if err != nil {
			return nil, err
		}
		result = append(result, *tmpResult)
	}

	return &result, nil
}

func (*WorkspaceComponentDao) GetWorkspaceComponentByWorkspaceIdAndKind(ctx context.Context, workspaceId int, kind string) (*[]do.WorkspaceComponentDo, error) {
	connect := g.DB("default")
	record, err := connect.GetAll(ctx, fmt.Sprint("select * from workspace_component where workspace_id = ? and kind=?"), workspaceId, kind)
	if err != nil {
		return nil, err
	}
	var result []do.WorkspaceComponentDo
	for i := 0; i < len(record); i++ {
		var tmpResult = &do.WorkspaceComponentDo{}
		err = record[i].Struct(tmpResult)
		if err != nil {
			return nil, err
		}
		result = append(result, *tmpResult)
	}
	return &result, nil
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
func (*WorkspaceComponentDao) UpdateWorkspaceComponentByWorkspaceIdAndKind(ctx context.Context, updateMap *gdb.Map, workspaceId int,
	kind string) (int64, error) {
	connect := g.DB("default")
	result, err := connect.Update(ctx, "workspace_component", updateMap, "workspace_id=? and kind=?", workspaceId, kind)
	if err != nil {
		return -1, err
	}
	rowAffect, _ := result.RowsAffected()
	return rowAffect, nil
}

func (*WorkspaceComponentDao) DeleteWorkspaceComponentsByWorkspaceId(ctx context.Context, workspaceId int) error {
	connect := g.DB("default")
	_, err := connect.Update(ctx, "workspace_component", "deleted=0", "workspace_id=?", workspaceId)
	if err != nil {
		return err
	}
	return nil
}

func (*WorkspaceComponentDao) DeleteWorkspaceComponentById(ctx context.Context, workspaceId int) error {
	connect := g.DB("default")
	_, err := connect.Update(ctx, "workspace_component", "deleted=0", "id=?", workspaceId)
	if err != nil {
		return err
	}
	return nil
}
