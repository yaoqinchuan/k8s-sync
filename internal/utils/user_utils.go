package utils

import (
	"context"
	"github.com/gogf/gf/v2/util/gconv"
	"k8s-sync/internal/model"
)

func GetUserInfoByContext(ctx context.Context) (*model.AccountModel, error) {

	userInfo := &model.AccountModel{}
	err := gconv.Struct(ctx.Value("userInfo"), userInfo)
	if err != nil {
		return nil, err
	}
	return userInfo, nil
}
