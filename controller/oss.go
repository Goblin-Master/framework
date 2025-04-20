package controller

import (
	"framework/service/oss"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type OSSController struct {
	zap        *zap.SugaredLogger
	OSSService oss.Service
}

func NewCommonController(zap *zap.SugaredLogger, ossService oss.Service) *OSSController {
	return &OSSController{
		zap:        zap,
		OSSService: ossService,
	}
}
func (cc *OSSController) RegisterRouter(r *gin.Engine) {
	rg := r.Group("oss")
	rg.POST("file", cc.UploadFile)
}

func (cc *OSSController) UploadFile(c *gin.Context) {
	cc.zap.Info("upload file")
	file, err := c.FormFile("file")
	if err != nil {
		cc.zap.Errorf("get file error:%v", err)
		c.JSON(400, gin.H{
			"message": "上传文件失败",
		})
		return
	}
	filePath, err := cc.OSSService.UploadFile(c, file)
	if err != nil {
		cc.zap.Errorf("upload file error:%v", err)
		c.JSON(400, gin.H{
			"message": "上传文件失败",
		})
		return
	}
	c.JSON(200, gin.H{
		"message":  "上传文件成功",
		"filePath": filePath,
	})
}
