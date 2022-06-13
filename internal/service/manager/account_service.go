package manager

import (
	"context"
	"database/sql"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/os/gtime"
	"k8s-sync/internal/model"
	"k8s-sync/internal/service/internal/dao"
	"k8s-sync/internal/service/internal/do"
)

type AccountService struct {
}

var accountDao = dao.AccountDao{}

func (*AccountService) GetByName(ctx context.Context, name string) (*model.AccountModel, error) {
	accountDo, err := accountDao.GetByName(ctx, name)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if nil == accountDo || err == sql.ErrNoRows {
		return nil, nil
	}

	accountModel := &model.AccountModel{
		Id:       accountDo.Id,
		UserName: accountDo.UserName,
		UserId:   accountDo.UserId,
		UserRole: accountDo.UserRole,
		Email:    accountDo.Email,
		CreateAt: accountDo.CreateAt,
		UpdateAt: accountDo.UpdateAt,
		Modifier: accountDo.Modifier,
	}
	return accountModel, nil
}

func (*AccountService) GetByUserId(ctx context.Context, userId string) (*model.AccountModel, error) {
	accountDo, err := accountDao.GetByUserId(ctx, userId)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if nil == accountDo || err == sql.ErrNoRows {
		return nil, nil
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
		CreateAt: gtime.Now(),
		UpdateAt: gtime.Now(),
		Modifier: data.Modifier,
		Deleted:  0,
	}
	return accountDao.AddAccount(ctx, accountDo)
}
func (*AccountService) UpdateAccountByUserId(ctx context.Context, updateMap *gdb.Map, userId string) (int64, error) {
	return accountDao.UpdateAccountByUserId(ctx, updateMap, userId)
}

func (*AccountService) DeleteByUserId(ctx context.Context, userId string) (int64, error) {
	return accountDao.DeleteByUserId(ctx, userId)
}
