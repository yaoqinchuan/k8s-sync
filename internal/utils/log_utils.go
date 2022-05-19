package utils

import (
	"context"
	"github.com/gogf/gf/v2/os/glog"
)

var Logger glog.Logger

func init() {
	Logger = *getLogger()
}

// 使用配置的方式來初始化錯誤日誌
func getLogger() *glog.Logger {

	logger := glog.New()
	configMap := make(map[string]interface{}, 2)
	configMap["rotateSize"] = "100M"
	configMap["rotateExpire"] = "1d"
	configMap["rotateBackupLimit"] = 200
	configMap["RotateBackupExpire"] = "7d"
	configMap["RotateBackupCompress"] = 9

	err := logger.SetConfigWithMap(configMap)
	if err != nil {
		return nil
	}

	loggerPath := "/var/log"
	err = logger.SetPath(loggerPath)
	if err != nil {
		glog.Panic(context.Background(), "init error log path failed, error:"+err.Error())

	}
	logger.SetStdoutPrint(true)
	logger.SetFile("{Y-m-d}.log")

	// 打印栈信息
	logger.SetStack(true)
	err = logger.SetLevelStr("INFO")
	if err != nil {
		glog.Panic(context.Background(), "init error log level failed, error:"+err.Error())
	}
	logger.SetWriterColorEnable(true)

	return logger
}
