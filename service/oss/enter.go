package oss

import (
	"context"
	"mime/multipart"
)

type Service interface {
	UploadFile(ctx context.Context, file *multipart.FileHeader) (string, error)
}
