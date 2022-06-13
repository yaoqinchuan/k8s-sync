package manager

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/os/gtime"
	"k8s-sync/internal/consts"
	"k8s-sync/internal/model"
	"k8s-sync/internal/service/internal/dao"
	"k8s-sync/internal/service/internal/do"
	"k8s-sync/internal/service/internal/k8s"
	"k8s-sync/internal/service/tasks/handlers"
)

type WorkspaceService struct {
}

var workspaceDao = dao.WorkspaceDao{}
var workspaceComponentDao = dao.WorkspaceComponentDao{}

func (*WorkspaceService) GetWorkspaceByName(ctx context.Context, name string) (*model.WorkspaceModel, error) {
	workspaceDo, err := workspaceDao.GetWorkspaceByName(ctx, name)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if nil == workspaceDo || err == sql.ErrNoRows {
		return nil, nil
	}
	input := &model.WorkspaceSpecModel{}
	if err := json.Unmarshal([]byte(workspaceDo.Spec), input); nil != err {
		return nil, err
	}
	accountModel := &model.WorkspaceModel{
		Id:         workspaceDo.Id,
		Name:       workspaceDo.Name,
		Attributes: workspaceDo.Attributes,
		Spec:       workspaceDo.Spec,
		Input:      input,
		Runtime:    workspaceDo.Runtime,
		Status:     workspaceDo.Status,
		CreateAt:   workspaceDo.CreateAt,
		UpdateAt:   workspaceDo.UpdateAt,
		Modifier:   workspaceDo.Modifier,
	}
	return accountModel, nil
}

func (*WorkspaceService) GetWorkspaceById(ctx context.Context, workspaceId int64) (*model.WorkspaceModel, error) {
	workspaceDo, err := workspaceDao.GetWorkspaceById(ctx, workspaceId)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if nil == workspaceDo || err == sql.ErrNoRows {
		return nil, nil
	}
	input := &model.WorkspaceSpecModel{}
	if err := json.Unmarshal([]byte(workspaceDo.Spec), input); nil != err {
		return nil, err
	}

	accountModel := &model.WorkspaceModel{
		Id:         workspaceDo.Id,
		Name:       workspaceDo.Name,
		Attributes: workspaceDo.Attributes,
		Spec:       workspaceDo.Spec,
		Input:      input,
		Runtime:    workspaceDo.Runtime,
		Status:     workspaceDo.Status,
		CreateAt:   workspaceDo.CreateAt,
		UpdateAt:   workspaceDo.UpdateAt,
		Modifier:   workspaceDo.Modifier,
	}
	return accountModel, nil
}

func (*WorkspaceService) AddWorkspace(ctx context.Context, workspaceModel *model.WorkspaceModel) (int64, error) {
	spec, err := json.Marshal(workspaceModel.Spec)
	if err != nil {
		return 0, err
	}
	input, err := k8s.GenerateWorkspace(workspaceModel.Input)
	if err != nil {
		return 0, err
	}
	inputBytes, err := json.Marshal(input)
	if err != nil {
		return 0, err
	}
	workspaceDo := &do.WorkspaceDo{
		Id:         workspaceModel.Id,
		Name:       workspaceModel.Name,
		Attributes: workspaceModel.Attributes,
		Spec:       string(spec),
		Input:      string(inputBytes),
		Runtime:    workspaceModel.Runtime,
		Status:     workspaceModel.Status,
		CreateAt:   workspaceModel.CreateAt,
		UpdateAt:   workspaceModel.UpdateAt,
		Modifier:   workspaceModel.Modifier,
	}
	return workspaceDao.AddWorkspace(ctx, workspaceDo)
}
func (*WorkspaceService) UpdateWorkspaceById(ctx context.Context, updateMap *gdb.Map, id int64) (int64, error) {
	return workspaceDao.UpdateWorkspaceById(ctx, updateMap, id)
}

func (*WorkspaceService) UpdateWorkspaceStatusById(ctx context.Context, status string, id int64) (int64, error) {
	return workspaceDao.UpdateWorkspaceStatusById(ctx, status, id)
}
func (*WorkspaceService) UpdateWorkspaceStatusWithCASById(ctx context.Context, newStatus string, oldStatus string, id int64) (int64, error) {
	return workspaceDao.UpdateWorkspaceStatusWithCASById(ctx, newStatus, oldStatus, id)
}
func (*WorkspaceService) DeleteById(ctx context.Context, id int64) (int64, error) {
	return workspaceDao.DeleteById(ctx, id)
}
func (workspaceService *WorkspaceService) CreateWorkspace(ctx context.Context, workspaceModel *model.WorkspaceModel) (int64, error) {
	workspaceModel.Status = handlers.WS_PENDING
	workspaceId, err := workspaceService.AddWorkspace(ctx, workspaceModel)
	if err != nil {
		return 0, err
	}
	service := &do.WorkspaceComponentDo{
		Kind:        consts.SERVICE,
		Name:        fmt.Sprintf("svc-%v", workspaceModel.Name),
		WorkspaceId: workspaceId,
		CreateAt:    gtime.Now(),
		UpdateAt:    gtime.Now(),
		Deleted:     0,
	}
	_, err = workspaceComponentDao.AddWorkspaceComponent(ctx, service)
	if err != nil {
		return workspaceId, err
	}

	if 0 != workspaceModel.Temporary {
		pvc := &do.WorkspaceComponentDo{
			Kind:        consts.PersistentVolumeClaim,
			Name:        fmt.Sprintf("pvc-%v", workspaceModel.Name),
			WorkspaceId: workspaceId,
			CreateAt:    gtime.Now(),
			UpdateAt:    gtime.Now(),
			Deleted:     0,
		}
		_, err = workspaceComponentDao.AddWorkspaceComponent(ctx, pvc)
		if err != nil {
			return workspaceId, err
		}
	}
	return workspaceId, nil
}

// todo add compenennt create
func (workspaceService *WorkspaceService) StartWorkspace(ctx context.Context, workspaceId int64) error {
	workspaceModel, err := workspaceService.GetWorkspaceById(ctx, workspaceId)
	if err != nil {
		workspaceService.UpdateWorkspaceStatusById(ctx, handlers.WS_ERROR, workspaceModel.Id)
		return err
	}

	workspaceService.UpdateWorkspaceStatusById(ctx, handlers.WS_STARTING, workspaceModel.Id)

	_, err = workspaceService.UpdateWorkspaceStatusById(ctx, handlers.WS_STARTING, workspaceModel.Id)
	if err != nil {
		return err
	}
	err = k8s.DoCreateWorkspace(ctx, k8s.ClientSet, workspaceModel.Input)
	if err != nil {
		workspaceService.UpdateWorkspaceStatusById(ctx, handlers.WS_ERROR, workspaceModel.Id)
		return err
	}

	err = k8s.DoCreatePVC(ctx, k8s.ClientSet, workspaceModel.Input)
	if err != nil {
		workspaceService.UpdateWorkspaceStatusById(ctx, handlers.WS_ERROR, workspaceModel.Id)
		return err
	}

	err = k8s.DoCreateSVC(ctx, k8s.ClientSet, workspaceModel.Input)
	if err != nil {
		workspaceService.UpdateWorkspaceStatusById(ctx, handlers.WS_ERROR, workspaceModel.Id)
		return err
	}
	return nil
}

func (workspaceService *WorkspaceService) RestoringWorkspace(ctx context.Context, workspaceId int64) (bool, error) {
	workspaceModel, err := workspaceService.GetWorkspaceById(ctx, workspaceId)
	if err != nil {
		workspaceService.UpdateWorkspaceStatusById(ctx, handlers.WS_ERROR, workspaceModel.Id)
		return false, err
	}
	workspaceService.UpdateWorkspaceStatusById(ctx, handlers.WS_STARTING, workspaceModel.Id)
	_, err = workspaceService.UpdateWorkspaceStatusById(ctx, handlers.WS_STARTING, workspaceModel.Id)
	if err != nil {
		return false, err
	}
	err = k8s.DoRestoreWorkspace(ctx, k8s.ClientSet, workspaceModel.Input)
	if err != nil {
		workspaceService.UpdateWorkspaceStatusById(ctx, handlers.WS_ERROR, workspaceModel.Id)
		return false, err
	}
	return true, nil
}

func (workspaceService *WorkspaceService) StoppingWorkspace(ctx context.Context, workspaceId int64) error {
	workspaceModel, err := workspaceService.GetWorkspaceById(ctx, workspaceId)
	if err != nil {
		workspaceService.UpdateWorkspaceStatusById(ctx, handlers.WS_ERROR, workspaceModel.Id)
		return err
	}
	workspaceModel.Status = handlers.WS_STOPPING
	_, err = workspaceService.UpdateWorkspaceStatusById(ctx, handlers.WS_STOPPING, workspaceModel.Id)
	if err != nil {
		return err
	}
	err = k8s.DoStopWorkspace(ctx, k8s.ClientSet, workspaceModel.Input)
	if err != nil {
		return err
	}

	return nil
}

// delete workspace must in Stopped status
func (workspaceService *WorkspaceService) DeletingWorkspace(ctx context.Context, workspaceId int64) error {
	workspaceModel, err := workspaceService.GetWorkspaceById(ctx, workspaceId)
	if err != nil {
		workspaceService.UpdateWorkspaceStatusById(ctx, handlers.WS_ERROR, workspaceModel.Id)
		return err
	}
	workspaceModel.Status = handlers.WS_DELETING
	workspaceComponents, err := workspaceComponentDao.GetWorkspaceComponentByWorkspaceId(ctx, workspaceModel.Id)
	if err != nil {
		return err
	}
	for i := 0; i < len(*workspaceComponents); i++ {
		workspaceComponent := (*workspaceComponents)[i]
		if workspaceComponent.Kind == consts.SERVICE {
			err := k8s.DoDeleteService(ctx, k8s.ClientSet, workspaceComponent.Name, workspaceModel.Input.NameSpace)
			if err != nil {
				return err
			}
		}
		if workspaceComponent.Kind == consts.PersistentVolumeClaim {
			err := k8s.DoDeletePVC(ctx, k8s.ClientSet, workspaceComponent.Name, workspaceModel.Input.NameSpace)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (workspaceService *WorkspaceService) AddAndStartWorkspace(ctx context.Context, workspaceModel *model.WorkspaceModel) error {
	workspaceId, err := workspaceService.CreateWorkspace(ctx, workspaceModel)
	if err != nil {
		return err
	}
	err = workspaceService.StartWorkspace(ctx, workspaceId)
	return err
}
