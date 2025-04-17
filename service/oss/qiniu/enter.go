package qiniu

import (
	"context"
	"fmt"
	"framework/service/oss"
	"github.com/qiniu/go-sdk/v7/storagev2/uploader"
	"go.uber.org/zap"
	"mime/multipart"
)

type Service struct {
	client     *uploader.UploadManager
	bucketName *string
	url        *string
	zap        *zap.SugaredLogger
}

var _ oss.Service = (*Service)(nil)

func NewQiNiuOSSService(zap *zap.SugaredLogger, client *uploader.UploadManager, bucketName string, url string) *Service {
	return &Service{
		zap:        zap,
		client:     client,
		bucketName: &bucketName,
		url:        &url,
	}
}

func (s *Service) UploadFile(ctx context.Context, file *multipart.FileHeader) (string, error) {
	key := fmt.Sprintf("%s/%s", *s.bucketName, file)
	reader, err := file.Open()
	if err != nil {
		s.zap.Errorf("open file error:%v", err)
		return "", err
	}
	defer reader.Close()
	err = s.client.UploadReader(ctx, reader, &uploader.ObjectOptions{
		BucketName: *s.bucketName,
		ObjectName: &key,
		FileName:   file.Filename,
	}, nil)
	s.zap.Errorf("upload file error:%v", err)
	return fmt.Sprintf("%s/%s", *s.url, key), err
}
