package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/daiguadaidai/pilipala/pili/dao"
	"github.com/daiguadaidai/pilipala/pili/server/message"
	"fmt"
	"github.com/cihub/seelog"
	"strconv"
	"net/http"
	"github.com/daiguadaidai/pilipala/pili/config"
	"github.com/daiguadaidai/pilipala/common"
)

type PilipalaCommandProgramHandler struct {}

func (this *PilipalaCommandProgramHandler) Create(_c *gin.Context) {
}

func (this *PilipalaCommandProgramHandler) Delete(_c *gin.Context) {
}

func (this *PilipalaCommandProgramHandler) Update(_c *gin.Context) {
}

func (this *PilipalaCommandProgramHandler) PartialUpdate(_c *gin.Context) {
}

func (this *PilipalaCommandProgramHandler) List(_c *gin.Context) {
	minPK := _c.DefaultQuery("min_pk", "")
	offset := _c.DefaultQuery("offset", "50")
	columnStr := "*"

	res := message.NewResponseMessage()

	pilipalaCommandProgramDao := new(dao.PilipalaCommandProgramDao)

	// 转化每一页的偏移量
	offsetInt, err := strconv.Atoi(offset)
	if err != nil {
		res.Code = 30000
		res.Message = fmt.Sprintf("不能获取到正确的分页offset. %v", err)
		seelog.Error(res.Message)
		_c.JSON(http.StatusOK, res)
		return
	}

	// 不能获取到每页中的最小主键, 刚刚进入页面, 代表是第一页
	if minPK == "" {
		items, err := pilipalaCommandProgramDao.PaginationFirstFind(offsetInt, columnStr)
		if err != nil {
			res.Code = 30000
			res.Message = fmt.Sprintf("获取第一页数据失败. %v", err)
			seelog.Error(res.Message)
		}
		res.Data["items"] = items
		_c.JSON(http.StatusOK, res)
		return
	}

	// 创换每一页中最小的主键
	minPKInt, err := strconv.Atoi(offset)
	if err != nil {
		res.Code = 30000
		res.Message = fmt.Sprintf("不能获取到正确的一页中最小的ID. %v", err)
		seelog.Error(res.Message)
	}

	// 分页获取数据
	items, err := pilipalaCommandProgramDao.PaginationFind(minPKInt, offsetInt, columnStr)
	if err != nil {
		res.Code = 30000
		res.Message = fmt.Sprintf("获取数据失败. %v", err)
		seelog.Error(res.Message)
	}
	res.Data["items"] = items

	_c.JSON(http.StatusOK, res)
}

func (this *PilipalaCommandProgramHandler) Retrieve(_c *gin.Context) {
}

// 上传命令
func (this *PilipalaCommandProgramHandler) UploadCommand(_c *gin.Context) {
	res := message.NewResponseMessage()
	file, err := _c.FormFile("file")
	if err != nil {
		res.Code = 30000
		res.Message = fmt.Sprintf("后台获取文件失败. %v", err)
		seelog.Error(res.Message)
		_c.JSON(http.StatusOK, res)
		return
	}

	piliStartConfig := config.GetPiliStartConfig()
	tmpFileName := common.CreateUUIDFileName()
	filePath := piliStartConfig.UploadCommandFilePath(tmpFileName)
	seelog.Warnf("Uplaod File Name: %v, FilePath: %v", file.Filename, filePath)
	if err := _c.SaveUploadedFile(file, filePath); err != nil {
		res.Code = 30000
		res.Message = fmt.Sprintf("后台保存文件失败. %v", err)
		seelog.Error(res.Message)
		_c.JSON(http.StatusOK, res)
		return
	}

	res.Data["FileName"] = file.Filename
	res.Data["TmpFileName"] = tmpFileName
	_c.JSON(http.StatusOK, res)
}

func (this *PilipalaCommandProgramHandler) Register(_router *gin.Engine) {
	_router.POST("/pilipala_command_program", this.Create)
	_router.DELETE("/pilipala_command_program/:pk", this.Delete)
	_router.PUT("/pilipala_command_program/:pk", this.Update)
	_router.PATCH("/pilipala_command_program/:pk", this.PartialUpdate)
	_router.GET("/pilipala_command_program", this.List)
	_router.GET("/pilipala_command_program/:pk", this.Retrieve)
	_router.POST("/pilipala_command_program/upload_command", this.UploadCommand)
}

func init() {
	handler := new(PilipalaCommandProgramHandler)
	AddHandler(handler)
}
