package dao

import (
	"context"
	"fmt"
	"mygogf/internal/do"
	"mygogf/internal/utils"
)

func GetByName(ctx context.Context, name string) (do.Account, error) {
	connect, err := utils.GetConnect("default")
	if err != nil {
		return nil, err
	}
	record, err := connect.GetOne(ctx, fmt.Sprintf("select * from acount where name = ? limit 1"), name)
	if err != nil {
		return nil, err
	}
	var result *Account
	record.Struct(result)
	return result, nil
}
