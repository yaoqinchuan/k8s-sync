package utils

import (
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/net/ghttp"
)

type restResult struct {
	Status       string      `json:"status"`
	ErrorCode    string      `json:"error_code"`
	ErrorMessage string      `json:"error_message"`
	Data         interface{} `json:"data"`
}

func RestSuccess(body interface{}, r *ghttp.Request) {
	result := restResult{
		Status:       "OK",
		ErrorCode:    "200",
		ErrorMessage: "",
		Data:         gjson.Marshal(body),
	}
	r.Response.Write(result)
}

func RestFailed(errorMessage string, r *ghttp.Request) {
	result := restResult{
		Status:       "ERROR",
		ErrorCode:    "500",
		ErrorMessage: errorMessage,
		Data:         "",
	}
	r.Response.WriteJsonP(result)
}
func RestFailedWithCode(errorCode string, errorMessage string, r *ghttp.Request) {
	result := restResult{
		Status:       "ERROR",
		ErrorCode:    errorCode,
		ErrorMessage: errorMessage,
		Data:         "",
	}
	r.Response.WriteJsonP(result)
}
