package service

import (
	"context"
	"github.com/gogf/gf/v2/database/gdb"
	"k8s-sync/internal/model"
	"k8s-sync/internal/service/internal/dao"
	"k8s-sync/internal/service/internal/do"
)

type AccountService struct {
}

var accountDao = dao.AccountDao{}

func (*AccountService) GetByName(ctx context.Context, name string) (*model.AccountModel, error) {
	accountDo, err := accountDao.GetByName(ctx, name)
	if err != nil {
		return nil, err
	}
	accountModel := &model.AccountModel{
		Id:       accountDo.Id,
		UserName: accountDo.UserName,
		UserId:   accountDo.UserId,
		UserRole: accountDo.UserRole,
		Email:    accountDo.Email,
		CreateAt: accountDo.CreateAt,
		UpdateAt: accountDo.UpdateAt,
	}
	return accountModel, nil
}

func (*AccountService) GetByUserId(ctx context.Context, userId string) (*model.AccountModel, error) {
	accountDo, err := accountDao.GetByUserId(ctx, userId)
	if err != nil {
		return nil, err
	}
	accountModel := &model.AccountModel{
		Id:       accountDo.Id,
		UserName: accountDo.UserName,
		UserId:   accountDo.UserId,
		UserRole: accountDo.UserRole,
		Email:    accountDo.Email,
		CreateAt: accountDo.CreateAt,
		UpdateAt: accountDo.UpdateAt,
	}
	return accountModel, nil
}

func (*AccountService) AddAccount(ctx context.Context, data *model.AccountModel) (int64, error) {

	accountDo := &do.Account{
		Id:       data.Id,
		UserName: data.UserName,
		UserId:   data.UserId,
		UserRole: data.UserRole,
		Email:    data.Email,
		CreateAt: data.CreateAt,
		UpdateAt: data.UpdateAt,
	}
	return accountDao.AddAccount(ctx, accountDo)
}
func (*AccountService) UpdateAccountByUserId(ctx context.Context, updateMap *gdb.Map, userId string) (int64, error) {
	return accountDao.UpdateAccountByUserId(ctx, updateMap, userId)
}

func (*AccountService) DeleteByUserId(ctx context.Context, userId string) (int64, error) {
	return accountDao.DeleteByUserId(ctx, userId)
}
