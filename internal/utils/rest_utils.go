package utils

import (
	"github.com/gogf/gf/v2/net/ghttp"
)

type restResult struct {
	Status       string      `json:"status"`
	ErrorCode    string      `json:"errorCode"`
	ErrorMessage string      `json:"errorMessage"`
	Data         interface{} `json:"data"`
}

func RestSuccess(data interface{}, r *ghttp.Request) {
	result := restResult{
		Status:       "OK",
		ErrorCode:    "200",
		ErrorMessage: "",
		Data:         data,
	}
	err := r.Response.WriteJson(result)
	if err != nil {
		panic(err)
	}
}

func RestFailed(errorMessage string, r *ghttp.Request) {
	result := restResult{
		Status:       "ERROR",
		ErrorCode:    "500",
		ErrorMessage: errorMessage,
		Data:         "",
	}
	err := r.Response.WriteJson(result)
	if err != nil {
		panic(err)
	}
}
func RestFailedWithCode(errorCode string, errorMessage string, r *ghttp.Request) {
	result := restResult{
		Status:       "ERROR",
		ErrorCode:    errorCode,
		ErrorMessage: errorMessage,
		Data:         "",
	}
	err := r.Response.WriteJson(result)
	if err != nil {
		panic(err)
	}
}
