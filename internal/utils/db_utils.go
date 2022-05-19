package utils

import (
	"fmt"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

var connectMap map[string]gdb.DB = make(map[string]gdb.DB, 10)

func GetConnect(configNames ...string) (gdb.DB, error) {
	configName := "default"
	if len(configNames) != 1 {
		return nil, fmt.Errorf("get connet params number should not lager than 1, input is %v", configNames)
	} else if len(configNames) == 1 {
		configName = configNames[0]
	}
	connect, ok := connectMap[configName]
	if ok {
		return connect, nil
	}
	dbConnect = g.DB(configName)
	connectMap[configName] = dbConnect
	return dbConnect, nil
}
