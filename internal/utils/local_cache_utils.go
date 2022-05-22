package utils

import (
	"context"
	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/os/gcache"
	"time"
)

type LocalCacheUtils struct {
}

var cache *gcache.Cache

func init() {
	cache = gcache.New()
}

func (*LocalCacheUtils) GetCache(ctx context.Context, key interface{}) (interface{}, error) {
	varCache, err := cache.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	return varCache.Val(), nil
}

func (*LocalCacheUtils) SetCache(ctx context.Context, key interface{}, value interface{}, expireTime time.Duration) error {
	err := cache.Set(ctx, key, value, expireTime)
	if err != nil {
		return err
	}
	return nil
}

func (*LocalCacheUtils) UpdateCache(ctx context.Context, key interface{}, value interface{}) (*gvar.Var, bool, error) {
	return cache.Update(ctx, key, value)
}

func (*LocalCacheUtils) GetKeys(ctx context.Context) (keys []interface{}, err error) {
	return cache.Keys(ctx)
}
func (*LocalCacheUtils) UpdateExpire(ctx context.Context, key interface{}, duration time.Duration) (oldDuration time.Duration, err error) {
	return cache.UpdateExpire(ctx, key, duration)
}
func (*LocalCacheUtils) GetExpire(ctx context.Context, key interface{}) (time.Duration, error) {
	return cache.GetExpire(ctx, key)
}

func (*LocalCacheUtils) Clear(ctx context.Context) error {
	return cache.Clear(ctx)
}

func (*LocalCacheUtils) Contains(ctx context.Context, key interface{}) (bool, error) {
	return cache.Contains(ctx, key)
}

func (*LocalCacheUtils) Size(ctx context.Context) (size int, err error) {
	return cache.Size(ctx)
}
