package controller

import (
	"encoding/json"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
	"k8s-sync/internal/model"
	"k8s-sync/internal/service"
	"k8s-sync/internal/utils"
)

var workspaceService = service.WorkspaceService{}

func WorkspaceApiHandlerRegister(group *ghttp.RouterGroup) {
	group.POST("workspace", func(r *ghttp.Request) { addWorkspace(r) })
	group.DELETE("workspace", func(r *ghttp.Request) { deleteWorkspace(r) })
	group.PUT("workspace", func(r *ghttp.Request) { updateWorkspace(r) })
	group.GET("workspace", func(r *ghttp.Request) { getWorkspace(r) })

}

func getWorkspace(r *ghttp.Request) {
	var workspaceModel *model.WorkspaceModel
	var err error
	if workspaceName := r.Get("name"); nil != workspaceName {
		workspaceModel, err = workspaceService.GetWorkspaceByName(r.Context(), workspaceName.String())
	} else if workspaceId := r.Get("id"); nil != workspaceId {
		workspaceModel, err = workspaceService.GetWorkspaceById(r.Context(), workspaceId.Int())
	}
	if err != nil {
		utils.RestFailed(err.Error(), r)
		return
	}
	utils.RestSuccess(workspaceModel, r)
	return
}

func startWorkspaceByNameOrId(r *ghttp.Request) {
	var workspaceModel *model.WorkspaceModel
	var err error
	if workspaceName := r.Get("name"); nil != workspaceName {
		workspaceModel, err = workspaceService.GetWorkspaceByName(r.Context(), workspaceName.String())
	} else if workspaceId := r.Get("id"); nil != workspaceId {
		workspaceModel, err = workspaceService.GetWorkspaceById(r.Context(), workspaceId.Int())
	}
	if err != nil {
		utils.RestFailed(err.Error(), r)
		return
	}
	utils.RestSuccess(workspaceModel, r)
	return
}
func addWorkspace(r *ghttp.Request) {
	bodyBytes := r.GetBody()
	if len(bodyBytes) == 0 {
		utils.RestFailed("request body is empty.", r)
		return
	}
	var workspaceModel = &model.WorkspaceModel{}
	err := json.Unmarshal(bodyBytes, workspaceModel)
	if err != nil {
		utils.RestFailed(err.Error(), r)
		return
	}
	userInfo, err := utils.GetUserInfoByContext(r.Context())
	if err != nil {
		utils.RestFailed(err.Error(), r)
		return
	}
	workspaceModel.Modifier = userInfo.UserName
	workspaceModel.UpdateAt = gtime.Now()
	workspaceModel.CreateAt = gtime.Now()
	id, err := workspaceService.AddWorkspace(r.Context(), workspaceModel)
	if err != nil {
		utils.RestFailed(err.Error(), r)
		return
	}
	utils.RestSuccess(id, r)
	return
}

func updateWorkspace(r *ghttp.Request) {
	bodyBytes := r.GetBody()
	if len(bodyBytes) == 0 {
		utils.RestFailed("request body is empty.", r)
		return
	}
	workspaceId := r.Get("id")
	if nil == workspaceId {
		utils.RestFailed("workspace id is empty.", r)
	}
	var workspaceModel *model.WorkspaceModel
	err := json.Unmarshal(bodyBytes, workspaceModel)
	if err != nil {
		utils.RestFailed(err.Error(), r)
		return
	}
	userInfo, err := utils.GetUserInfoByContext(r.Context())
	if err != nil {
		utils.RestFailed(err.Error(), r)
		return
	}
	workspaceModel.Modifier = userInfo.UserName
	workspaceModel.UpdateAt = gtime.Now()
	saveMap := gconv.Map(*workspaceModel)
	id, err := workspaceService.UpdateWorkspaceById(r.Context(), &saveMap, workspaceId.Int())
	if err != nil {
		utils.RestFailed(err.Error(), r)
		return
	}
	utils.RestSuccess(id, r)
	return
}

func deleteWorkspace(r *ghttp.Request) {
	var err error
	if accountId := r.Get("id"); nil != accountId {
		_, err = workspaceService.DeleteById(r.Context(), accountId.Int())
	}
	if err != nil {
		utils.RestFailed(err.Error(), r)
		return
	}
	result := make(map[string]string, 10)
	result["status"] = "success"
	utils.RestSuccess(result, r)
	return
}
