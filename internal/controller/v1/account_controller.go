package controller

import (
	"encoding/json"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gtime"
	"k8s-sync/internal/model"
	"k8s-sync/internal/service/manager"
	"k8s-sync/internal/utils"
)

var accountService = manager.AccountService{}

func AccountApiHandlerRegister(group *ghttp.RouterGroup) {
	group.POST("account", func(r *ghttp.Request) { addAccount(r) })
	group.DELETE("account", func(r *ghttp.Request) { deleteAccount(r) })
	group.PUT("account", func(r *ghttp.Request) { updateAccount(r) })
	group.GET("account", func(r *ghttp.Request) { getAccount(r) })

}

func getAccount(r *ghttp.Request) {
	var accountModel *model.AccountModel
	var err error
	if accountName := r.Get("userName"); nil != accountName {
		accountModel, err = accountService.GetByName(r.Context(), accountName.String())
	} else if accountId := r.Get("userId"); nil != accountId {
		accountModel, err = accountService.GetByUserId(r.Context(), accountId.String())
	}
	if err != nil {
		utils.RestFailed(err.Error(), r)
		return
	}
	utils.RestSuccess(accountModel, r)
	return
}

func addAccount(r *ghttp.Request) {
	bodyBytes := r.GetBody()
	if len(bodyBytes) == 0 {
		utils.RestFailed("request body is empty.", r)
		return
	}
	var accountModel = &model.AccountModel{}
	err := json.Unmarshal(bodyBytes, accountModel)
	if err != nil {
		utils.RestFailed(err.Error(), r)
		return
	}
	userInfo, err := utils.GetUserInfoByContext(r.Context())
	if err != nil {
		utils.RestFailed(err.Error(), r)
		return
	}
	accountModel.Modifier = userInfo.UserName
	accountModel.UpdateAt = gtime.Now()
	accountModel.CreateAt = gtime.Now()
	id, err := accountService.AddAccount(r.Context(), accountModel)
	if err != nil {
		utils.RestFailed(err.Error(), r)
		return
	}
	utils.RestSuccess(id, r)
	return
}
func updateAccount(r *ghttp.Request) {
	bodyBytes := r.GetBody()
	if len(bodyBytes) == 0 {
		utils.RestFailed("request body is empty.", r)
		return
	}
	accountId := r.Get("userId")
	if nil == accountId {
		utils.RestFailed("user id is empty.", r)
	}
	var accountModel = &model.AccountModel{}
	err := json.Unmarshal(bodyBytes, accountModel)
	if err != nil {
		utils.RestFailed(err.Error(), r)
		return
	}
	userInfo, err := utils.GetUserInfoByContext(r.Context())
	if err != nil {
		utils.RestFailed(err.Error(), r)
		return
	}

	saveMap := g.Map{}
	if "" != accountModel.UserRole {
		saveMap["userRole"] = accountModel.UserRole
	}
	if "" != accountModel.UserName {
		saveMap["userName"] = accountModel.UserName
	}
	if "" != accountModel.UserId {
		saveMap["userId"] = accountModel.UserId
	}
	if "" != accountModel.Email {
		saveMap["email"] = accountModel.Email
	}
	saveMap["modifier"] = userInfo.UserName
	saveMap["updateAt"] = gtime.Now()
	id, err := accountService.UpdateAccountByUserId(r.Context(), &saveMap, accountId.String())
	if err != nil {
		utils.RestFailed(err.Error(), r)
		return
	}
	utils.RestSuccess(id, r)
	return
}
func deleteAccount(r *ghttp.Request) {
	var err error
	if accountId := r.Get("userId"); nil != accountId {
		_, err = accountService.DeleteByUserId(r.Context(), accountId.String())
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
