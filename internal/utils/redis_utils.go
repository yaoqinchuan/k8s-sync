package utils

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
)

type RedisUtils struct {
}

func (*RedisUtils) Set(ctx context.Context, key string, value string) error {
	connect, err := g.Redis().Conn(ctx)
	if err != nil {
		return err
	}
	_, err = connect.Do(ctx, "SET", key, value)
	if err != nil {
		return err
	}
	return nil
}
func (*RedisUtils) Expire(ctx context.Context, key string, expireTime int) error {
	connect, err := g.Redis().Conn(ctx)
	if err != nil {
		return err
	}
	_, err = connect.Do(ctx, "EXPIRE", key, expireTime)
	if err != nil {
		return err
	}
	return nil
}
func (*RedisUtils) SetNx(ctx context.Context, key string, value string) error {
	connect, err := g.Redis().Conn(ctx)
	if err != nil {
		return err
	}
	_, err = connect.Do(ctx, "SETNX", key, value)
	if err != nil {
		return err
	}
	return nil
}
func (*RedisUtils) Exists(ctx context.Context, key string) (int64, error) {
	connect, err := g.Redis().Conn(ctx)
	if err != nil {
		return -1, err
	}
	result, err := connect.Do(ctx, "EXISTS", key)
	if err != nil {
		return -1, err
	}
	return result.Int64(), nil
}
func (*RedisUtils) Delete(ctx context.Context, key string) error {
	connect, err := g.Redis().Conn(ctx)
	if err != nil {
		return err
	}
	_, err = connect.Do(ctx, "DELETE", key)
	if err != nil {
		return err
	}
	return nil
}

func (*RedisUtils) GetByKey(ctx context.Context, key string) (string, error) {
	connect, err := g.Redis().Conn(ctx)
	if err != nil {
		return "", err
	}
	result, err := connect.Do(ctx, "GET", key)
	if err != nil {
		return "", err
	}
	return result.String(), nil
}
