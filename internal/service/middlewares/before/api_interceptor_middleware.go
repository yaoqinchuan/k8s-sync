package middleware

import (
	"fmt"
	"github.com/gogf/gf/v2/net/ghttp"
	"k8s-sync/internal/utils"
)

func ApiRequestRecord(r *ghttp.Request) {
	utils.Logger.Info(r.Context(), fmt.Sprintf("request %v received", r.Request.URL))
	r.Middleware.Next()
}
