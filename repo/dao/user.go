package dao

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type User struct {
	Id       int64
	Name     string
	Province string
	City     string
	District string
	//创建时间
	Ctime int64
	//更新时间
	Utime int64
}

type IUserDao interface {
	Get() User
}

var _ IUserDao = (*UserDao)(nil)

type UserDao struct {
	zap *zap.SugaredLogger
	db  *gorm.DB
}

func NewUserDao(zap *zap.SugaredLogger, db *gorm.DB) *UserDao {
	return &UserDao{
		db:  db,
		zap: zap,
	}
}
func (ud *UserDao) Get() User {
	ud.zap.Infof("user dao")
	return User{}
}
