package repo

import (
	"context"
	"framework/domain"
	"framework/repo/cache"
	"framework/repo/dao"
	"go.uber.org/zap"
)

type IUserRepo interface {
	GetUser(ctx context.Context, u domain.User) string
}

var _ IUserRepo = (*UserRepo)(nil)

type UserRepo struct {
	zap       *zap.SugaredLogger
	userCache *cache.UserCache
	userDao   *dao.UserDao
}

func NewUserRepo(zap *zap.SugaredLogger, userDao *dao.UserDao, userCache *cache.UserCache) *UserRepo {
	return &UserRepo{
		zap:       zap,
		userDao:   userDao,
		userCache: userCache,
	}
}
func (ur *UserRepo) GetUser(ctx context.Context, u domain.User) string {
	//将domain层数据转换成dao层数据
	ur.zap.Info(ur.userCache.Get(nil, u.Id))
	ur.zap.Info(ur.userDao.Get())
	ur.zap.Info(ur.userCache.Set(nil, domain.User{}))
	return "user repo"
}

// 在仓储做模型转换
func (ur *UserRepo) entityToDomain(u dao.User) domain.User {
	return domain.User{
		Name: u.Name,
		Id:   u.Id,
	}
}

func (ur *UserRepo) domainToEntity(u domain.User) dao.User {
	return dao.User{
		Name: u.Name,
		Id:   u.Id,
	}
}
