package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/gogf/gf/v2/database/gdb"
	"k8s-sync/internal/consts"
	"k8s-sync/internal/model"
	"k8s-sync/internal/service/internal/dao"
	"k8s-sync/internal/service/internal/do"
	"k8s-sync/internal/utils"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var clientSet *kubernetes.Clientset

func init() {
	configPath, err := utils.ConfigData.Get(context.TODO(), "k8sConfigFile")
	if nil != err {
		panic(err)
	}
	if configPath.IsEmpty() {
		panic("k8s config cert file is empty.")
	}
	config, err := clientcmd.BuildConfigFromFlags("", configPath.String())
	if err != nil {
		panic(err)
	}
	clientSet, err = kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}
}

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

func (*WorkspaceService) GetWorkspaceById(ctx context.Context, workspaceId int) (*model.WorkspaceModel, error) {
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
	input, err := GenerateWorkspace(workspaceModel.Input)
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
func (*WorkspaceService) UpdateWorkspaceById(ctx context.Context, updateMap *gdb.Map, id int) (int64, error) {
	return workspaceDao.UpdateWorkspaceById(ctx, updateMap, id)
}

func (*WorkspaceService) DeleteById(ctx context.Context, id int) (int64, error) {
	return workspaceDao.DeleteById(ctx, id)
}
func (workspaceService *WorkspaceService) CreateWorkspace(ctx context.Context, workspaceModel *model.WorkspaceModel) error {
	workspaceModel.Status = consts.WS_PENDING
	_, err := workspaceService.AddWorkspace(ctx, workspaceModel)
	if err != nil {
		return err
	}
	return nil
}

func (workspaceService *WorkspaceService) StartingWorkspace(ctx context.Context, workspaceModel *model.WorkspaceModel) error {
	workspaceModel.Status = consts.WS_STARTING
	err := StartWorkspace(ctx, clientSet, workspaceModel.Input)
	if err != nil {
		return err
	}
	return nil
}

func (workspaceService *WorkspaceService) StoppingWorkspace(ctx context.Context, workspaceModel *model.WorkspaceModel) error {
	workspaceModel.Status = consts.WS_STOPPING
	err := StopWorkspace(ctx, clientSet, workspaceModel.Input)
	if err != nil {
		return err
	}
	return nil
}

func (workspaceService *WorkspaceService) DeletingWorkspace(ctx context.Context, workspaceModel *model.WorkspaceModel) error {
	workspaceModel.Status = consts.WS_DELETING
	err := DeleteWorkspace(ctx, clientSet, workspaceModel.Input)
	if err != nil {
		return err
	}
	return nil
}

func (workspaceService *WorkspaceService) AddAndStartWorkspace(ctx context.Context, workspaceModel *model.WorkspaceModel) error {
	_, err := workspaceService.AddWorkspace(ctx, workspaceModel)
	if err != nil {
		return err
	}
	err = StartWorkspace(ctx, clientSet, workspaceModel.Input)
	if err != nil {
		return err
	}
	return nil
}
