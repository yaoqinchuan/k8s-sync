package service

import (
	"context"
	"database/sql"
	"github.com/gogf/gf/v2/database/gdb"
	"k8s-sync/internal/model"
	"k8s-sync/internal/service/internal/dao"
	"k8s-sync/internal/service/internal/do"
)

type WorkspaceService struct {
}

var workspaceDao = dao.WorkspaceDao{}

func (*WorkspaceService) GetWorkspaceByName(ctx context.Context, name string) (*model.WorkspaceModel, error) {
	workspaceDo, err := workspaceDao.GetWorkspaceByName(ctx, name)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if nil == workspaceDo || err == sql.ErrNoRows {
		return nil, nil
	}

	accountModel := &model.WorkspaceModel{
		Id:         workspaceDo.Id,
		Name:       workspaceDo.Name,
		Attributes: workspaceDo.Attributes,
		Spec:       workspaceDo.Spec,
		Runtime:    workspaceDo.Runtime,
		Status:     workspaceDo.Status,
		CreateAt:   workspaceDo.CreateAt,
		UpdateAt:   workspaceDo.UpdateAt,
		Modifier:   workspaceDo.Modifier,
	}
	return accountModel, nil
}

func (*WorkspaceService) GetWorkspaceById(ctx context.Context, workspaceId int) (*model.WorkspaceModel, error) {
	workspaceDo, err := workspaceDao.GetWorkspaceById(ctx, workspaceId)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if nil == workspaceDo || err == sql.ErrNoRows {
		return nil, nil
	}
	accountModel := &model.WorkspaceModel{
		Id:         workspaceDo.Id,
		Name:       workspaceDo.Name,
		Attributes: workspaceDo.Attributes,
		Spec:       workspaceDo.Spec,
		Runtime:    workspaceDo.Runtime,
		Status:     workspaceDo.Status,
		CreateAt:   workspaceDo.CreateAt,
		UpdateAt:   workspaceDo.UpdateAt,
		Modifier:   workspaceDo.Modifier,
	}
	return accountModel, nil
}

func (*WorkspaceService) AddWorkspace(ctx context.Context, workspaceModel *model.WorkspaceModel) (int64, error) {
	workspaceDo := &do.WorkspaceDo{
		Id:         workspaceModel.Id,
		Name:       workspaceModel.Name,
		Attributes: workspaceModel.Attributes,
		Spec:       workspaceModel.Spec,
		Runtime:    workspaceModel.Runtime,
		Status:     workspaceModel.Status,
		CreateAt:   workspaceModel.CreateAt,
		UpdateAt:   workspaceModel.UpdateAt,
		Modifier:   workspaceModel.Modifier,
	}
	return workspaceDao.AddWorkspace(ctx, workspaceDo)
}
func (*WorkspaceService) UpdateWorkspaceById(ctx context.Context, updateMap *gdb.Map, id int) (int64, error) {
	return workspaceDao.UpdateWorkspaceById(ctx, updateMap, id)
}

func (*WorkspaceService) DeleteById(ctx context.Context, id int) (int64, error) {
	return workspaceDao.DeleteById(ctx, id)
}
