package controller

import (
	"github.com/gogf/gf/v2/encoding/gbase64"
	"github.com/gogf/gf/v2/net/ghttp"
	"k8s-sync/internal/utils"
)

func MaintenanceApiHandlerRegister(group *ghttp.RouterGroup) {

	group.GET("encode/base64", func(r *ghttp.Request) { getEncodeStringByBase64(r) })

	group.GET("decode/base64", func(r *ghttp.Request) { getDecodeStringByBase64(r) })
}

func getEncodeStringByBase64(r *ghttp.Request) {
	plaintext := r.GetHeader("plaintext")
	if plaintext == "" {
		return
	}
	plaintextEncode := gbase64.EncodeString(plaintext)
	utils.RestSuccess(plaintextEncode, r)
	return
}

func getDecodeStringByBase64(r *ghttp.Request) {
	encodeText := r.GetHeader("text")
	if encodeText == "" {
		utils.RestFailed("text is needed", r)
		return
	}
	plaintext, err := gbase64.Decode([]byte(encodeText))
	if err != nil {
		utils.RestFailed(err.Error(), r)
		return
	}
	utils.RestSuccess(string(plaintext), r)
	return
}
