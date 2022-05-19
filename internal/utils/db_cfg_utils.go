// read config from db
package utils

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcfg"
	"github.com/gogf/gf/v2/os/gtime"
)

const AppConfigDbTableName = "app_config"

var dbConnect gdb.DB
var DbConfigData *gcfg.Config
var adapter gcfg.Adapter

type appConfig struct {
	Id          uint        `json:"id"       `
	ConfigKey   string      `json:"config_key" `
	ConfigValue string      `json:"config_value" `
	Describe    string      `json:"describe" `
	CreateAt    *gtime.Time `json:"create_at" `
	UpdateAt    *gtime.Time `json:"update_at" `
}

type DbConfig struct {
	Adapter gcfg.Adapter
}

func init() {
	dbConnect = g.DB("default")
	DbConfigData, _ = gcfg.New()
	DbConfigData.SetAdapter(adapter)
}

func (*DbConfig) Available(ctx context.Context, resource ...string) (ok bool) {
	_, err := dbConnect.GetOne(ctx, fmt.Sprintf("select * from %v limit 1", AppConfigDbTableName))
	if err != nil {
		Logger.Critical(ctx, "check app config failed.")
		return false
	}
	return true
}

func (*DbConfig) Get(ctx context.Context, pattern string) (value interface{}, err error) {
	config, err := dbConnect.GetOne(ctx, fmt.Sprintf("select * from %v where config_name = ? and deleted = 0", AppConfigDbTableName), pattern)
	if err != nil {
		Logger.Error(ctx, fmt.Sprintf("get config by key %v failed, error is %v", pattern, err))
		return nil, err
	}
	var configStruct *appConfig
	err = config.Struct(configStruct)
	if err != nil {
		Logger.Error(ctx, fmt.Sprintf("convert config content %v to appConfig struct failed, error is %v", config, err))
		return nil, err
	}
	return configStruct.ConfigValue, nil
}

func (*DbConfig) Data(ctx context.Context) (data map[string]interface{}, err error) {
	result := make(map[string]interface{}, 10)
	config, err := dbConnect.GetAll(ctx, fmt.Sprintf("select * from %v where deleted = 0", AppConfigDbTableName))
	if err != nil {
		Logger.Error(ctx, fmt.Sprintf("get all config failed, error is %v", err))
		return nil, err
	}
	for i := 0; i < config.Len(); i++ {
		var configStruct *appConfig
		err = config[i].Struct(configStruct)
		if err != nil {
			Logger.Error(ctx, fmt.Sprintf("convert config content %v to appConfig struct failed, error is %v", config, err))
			return nil, err
		}
		result[configStruct.ConfigKey] = configStruct.ConfigValue
	}
	return result, nil
}
