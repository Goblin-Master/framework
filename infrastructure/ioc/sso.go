package ioc

import (
	"framework/infrastructure/config"
	"framework/service/oss"
	"framework/service/oss/qiniu"
	"github.com/qiniu/go-sdk/v7/storagev2/credentials"
	"github.com/qiniu/go-sdk/v7/storagev2/http_client"
	"github.com/qiniu/go-sdk/v7/storagev2/uploader"
	"go.uber.org/zap"
)

func InitSSOService(zap *zap.SugaredLogger) oss.Service {
	return InitQiNiuOSS(zap)
}

func InitQiNiuOSS(zap *zap.SugaredLogger) *qiniu.Service {
	mac := credentials.NewCredentials(config.Conf.QiNiu.AccessKey, config.Conf.QiNiu.SecretKey)
	uploadManager := uploader.NewUploadManager(&uploader.UploadManagerOptions{
		Options: http_client.Options{
			Credentials: mac,
		},
	})
	return qiniu.NewQiNiuOSSService(zap, uploadManager, config.Conf.QiNiu.Bucket, config.Conf.QiNiu.Region)
}
