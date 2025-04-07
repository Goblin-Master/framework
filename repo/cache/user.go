package cache

import (
	"context"
	"ddd/domain"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

type IUserCache interface {
	Get(ctx context.Context, id int64) (u domain.User, err error)
	Set(ctx context.Context, u domain.User) error
}

var _ IUserCache = (*UserCache)(nil)

type UserCache struct {
	cmd redis.Cmdable
	zap *zap.SugaredLogger
}

func NewUserCache(zap *zap.SugaredLogger, cmd redis.Cmdable) *UserCache {
	return &UserCache{
		zap: zap,
		cmd: cmd,
	}
}
func (uc *UserCache) Get(ctx context.Context, id int64) (u domain.User, err error) {
	uc.zap.Info("cache get")
	return
}
func (uc *UserCache) Set(ctx context.Context, u domain.User) error {
	uc.zap.Info("cache set")
	return nil
}
