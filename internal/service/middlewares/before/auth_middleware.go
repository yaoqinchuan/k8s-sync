package middleware

import (
	"github.com/gogf/gf/v2/encoding/gbase64"
	"github.com/gogf/gf/v2/net/ghttp"
	"k8s-sync/internal/service"
	"k8s-sync/internal/utils"
)

func AuthMiddleware(r *ghttp.Request) {
	userIdEncode := r.GetHeader("userId")
	if "" == userIdEncode {
		utils.RestFailed("userId is needed", r)
		return
	}
	userId, err := gbase64.DecodeString(userIdEncode)
	if err != nil {
		utils.RestFailed(err.Error(), r)
		return
	}
	account := &service.AccountService{}
	userInfo, err := account.GetByUserId(r.Context(), string(userId))
	if err != nil {
		utils.RestFailed(err.Error(), r)
		return
	}
	if userInfo == nil {
		utils.RestFailed("you are forbidden to request k8s-sync service, please contact admin.", r)
		return
	}
	r.Middleware.Next()
}
