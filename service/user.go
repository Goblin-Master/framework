package service

import (
	"context"
	"framework/domain"
	"framework/repo"
	"go.uber.org/zap"
)

type IUserService interface {
	GetUser(ctx context.Context, u domain.User) string
}

var _ IUserService = (*UserService)(nil)

type UserService struct {
	zap      *zap.SugaredLogger
	userRepo *repo.UserRepo
}

func NewUserService(zap *zap.SugaredLogger, userRepo *repo.UserRepo) *UserService {
	return &UserService{
		zap:      zap,
		userRepo: userRepo,
	}
}
func (us *UserService) GetUser(ctx context.Context, u domain.User) string {
	//按道理是这里将前端接收道的数据绑定到domain.User
	us.zap.Info(us.userRepo.GetUser(ctx, u))
	return "user service"
}
