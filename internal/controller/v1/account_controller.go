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

var accountDao = service.AccountService{}

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
		accountModel, err = accountDao.GetByName(r.Context(), accountName.String())
	} else if accountId := r.Get("userId"); nil != accountId {
		accountModel, err = accountDao.GetByUserId(r.Context(), accountId.String())
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
	id, err := accountDao.AddAccount(r.Context(), accountModel)
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
	accountId := r.Get("account_id")
	if nil == accountId {
		utils.RestFailed("account id is empty.", r)
	}
	var accountModel *model.AccountModel
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
	saveMap := gconv.Map(*accountModel)
	id, err := accountDao.UpdateAccountByUserId(r.Context(), &saveMap, accountId.String())
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
		_, err = accountDao.DeleteByUserId(r.Context(), accountId.String())
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
