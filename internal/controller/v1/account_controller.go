package controller

import (
	"encoding/json"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/util/gconv"
	"k8s-sync/internal/model"
	"k8s-sync/internal/service"
	"k8s-sync/internal/utils"
)

var accountDao = service.AccountService{}

func AccountApiHandlerRegister(s *ghttp.Server) {

	s.BindHandler("GET://api/v1/account", func(r *ghttp.Request) {
		GetAccount(r)
	})
	s.BindHandler("POST://api/v1/account", func(r *ghttp.Request) {
		AddAccount(r)
	})
	s.BindHandler("DELETE://api/v1/account", func(r *ghttp.Request) {
		DeleteAccount(r)
	})
	s.BindHandler("PUT://api/v1/account", func(r *ghttp.Request) {
		UpdateAccount(r)
	})
}

func GetAccount(r *ghttp.Request) {
	var accountModel *model.AccountModel
	var err error
	if accountName := r.Get("account_name"); nil != accountName {
		accountModel, err = accountDao.GetByName(r.Context(), accountName.String())
	} else if accountId := r.Get("account_id"); nil != accountId {
		accountModel, err = accountDao.GetByUserId(r.Context(), accountId.Int())
	}
	if err != nil {
		utils.RestFailed(err.Error(), r)
		return
	}
	utils.RestSuccess(accountModel, r)
	return
}

func AddAccount(r *ghttp.Request) {
	bodyBytes := r.GetBody()
	if len(bodyBytes) == 0 {
		utils.RestFailed("request body is empty.", r)
		return
	}
	var accountModel *model.AccountModel
	err := json.Unmarshal(bodyBytes, accountModel)
	if err != nil {
		utils.RestFailed(err.Error(), r)
		return
	}
	id, err := accountDao.AddAccount(r.Context(), accountModel)
	if err != nil {
		utils.RestFailed(err.Error(), r)
		return
	}
	utils.RestSuccess(id, r)
	return
}
func UpdateAccount(r *ghttp.Request) {
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
	saveMap := gconv.Map(*accountModel)
	id, err := accountDao.UpdateAccountByUserId(r.Context(), &saveMap, accountId.Int64())
	if err != nil {
		utils.RestFailed(err.Error(), r)
		return
	}
	utils.RestSuccess(id, r)
	return
}
func DeleteAccount(r *ghttp.Request) {
	var err error
	if accountId := r.Get("account_id"); nil != accountId {
		_, err = accountDao.DeleteByUserId(r.Context(), accountId.Int64())
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
