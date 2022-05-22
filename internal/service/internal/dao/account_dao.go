package dao

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"k8s-sync/internal/service/internal/do"
)

type AccountDao struct {
}

func (*AccountDao) GetByName(ctx context.Context, name string) (*do.Account, error) {
	connect := g.DB("default")
	record, err := connect.GetOne(ctx, fmt.Sprint("select * from account where name = ? limit 1"), name)
	if err != nil {
		return nil, err
	}
	var result *do.Account
	err = record.Struct(result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (*AccountDao) GetByUserId(ctx context.Context, userId string) (*do.Account, error) {
	connect := g.DB("default")
	record, err := connect.GetOne(ctx, fmt.Sprint("select * from account where user_id = ? limit 1"), userId)
	if err != nil {
		return nil, err
	}
	var result *do.Account
	err = record.Struct(result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (*AccountDao) AddAccount(ctx context.Context, data *do.Account) (int64, error) {
	connect := g.DB("default")
	result, err := connect.Save(ctx, "account", data)
	if err != nil {
		return -1, err
	}
	line, _ := result.LastInsertId()
	return line, nil
}
func (*AccountDao) UpdateAccountByUserId(ctx context.Context, updateMap *gdb.Map, userId string) (int64, error) {
	connect := g.DB("default")
	result, err := connect.Update(ctx, "account", updateMap, "user_id=?", userId)
	if err != nil {
		return -1, err
	}
	rowAffect, _ := result.RowsAffected()
	return rowAffect, nil
}

func (*AccountDao) DeleteByUserId(ctx context.Context, userId string) (int64, error) {
	connect := g.DB("default")
	result, err := connect.Delete(ctx, "account", "user_id=?", userId)
	if err != nil {
		return -1, err
	}
	rowAffect, _ := result.RowsAffected()
	return rowAffect, nil
}
