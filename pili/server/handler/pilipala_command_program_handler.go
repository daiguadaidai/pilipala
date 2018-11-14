package handler

import (
	"fmt"
	"github.com/cihub/seelog"
	"github.com/daiguadaidai/pilipala/common"
	"github.com/daiguadaidai/pilipala/common/types"
	"github.com/daiguadaidai/pilipala/pili/config"
	"github.com/daiguadaidai/pilipala/pili/dao"
	"github.com/daiguadaidai/pilipala/pili/model"
	"github.com/daiguadaidai/pilipala/pili/server/form"
	"github.com/daiguadaidai/pilipala/pili/server/message"
	"github.com/deckarep/golang-set"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func init() {
	handler := new(PilipalaCommandProgramHandler)
	AddHandler(handler)
}

// 注册route
func (this *PilipalaCommandProgramHandler) Register(_router *gin.Engine) {
	_router.GET("/pili/command_programs", this.List)
	_router.POST("/pili/command_programs", this.Create)
	_router.PUT("/pili/command_programs", this.Update)
	_router.POST("/pili/command_programs/upload_create_command", this.UploadCreateCommand)
	_router.POST("/pili/command_programs/upload_edit_command", this.UploadEditCommand)
	_router.GET("/pili/command_programs/download/:command", this.Download)
}

type PilipalaCommandProgramHandler struct{}

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

func (this *PilipalaCommandProgramHandler) Create(_c *gin.Context) {
	res := message.NewResponseMessage()

	// 获取前端传来参数
	createForm := new(form.CreateCommandForm)
	if err := _c.ShouldBindJSON(createForm); err != nil {
		res.Code = 30000
		res.Message = fmt.Sprintf("获取参数错误. %v", err)
		seelog.Error(res.Message)
		_c.JSON(http.StatusOK, res)
		return
	}

	// 移动文件
	config := config.GetPiliStartConfig()
	oldFilePath := config.UploadCommandFilePath(createForm.TmpFileName)
	newFilePath := config.CommandFilePath(createForm.FileName)
	if err := os.Rename(oldFilePath, newFilePath); err != nil {
		res.Code = 30000
		res.Message = fmt.Sprintf("将命令文件移动到指定目录出错. %v", err)
		seelog.Error(res.Message)
		_c.JSON(http.StatusOK, res)
		return
	}

	// 保存数据库
	// 生成 命令数据
	pilipalaCommandProgram := &model.PilipalaCommandProgram{
		Title:        types.GetNullString(createForm.Title),
		FileName:     types.GetNullString(createForm.FileName),
		Params:       types.GetNullString(createForm.Params),
		HaveDedicate: types.GetNullInt64(createForm.HaveDedicate),
	}

	// 生成命令装用机器数据
	pilipalaCommandHosts := make([]*model.PilipalaCommandHost, 0, 1)
	if createForm.HaveDedicate == 1 {
		for _, hostId := range createForm.DedicateHosts {
			pilipalaCommandHost := new(model.PilipalaCommandHost)
			pilipalaCommandHost.PilipalaHostId = types.GetNullInt64(hostId)
			pilipalaCommandHosts = append(pilipalaCommandHosts, pilipalaCommandHost)
		}
	}

	// 将数据保存到数据库
	txDao := new(dao.PilipalaTXDao)
	err := txDao.CreateCommandAndHost(pilipalaCommandProgram, pilipalaCommandHosts)
	if err != nil {
		res.Code = 30000
		res.Message = fmt.Sprintf("保存命令出错. %v", err)
		seelog.Error(res.Message)
		_c.JSON(http.StatusOK, res)
		return
	}

	_c.JSON(http.StatusOK, res)
}

func (this *PilipalaCommandProgramHandler) Update(_c *gin.Context) {
	res := message.NewResponseMessage()

	// 获取前端传来参数
	updateForm := new(form.UpdateCommandForm)
	if err := _c.ShouldBindJSON(updateForm); err != nil {
		res.Code = 30000
		res.Message = fmt.Sprintf("获取参数错误. %v", err)
		seelog.Error(res.Message)
		_c.JSON(http.StatusOK, res)
		return
	}

	// 保存数据库
	// 生成 命令数据
	pilipalaCommandProgram := &model.PilipalaCommandProgram{
		Id:           types.GetNullInt64(updateForm.Id),
		Title:        types.GetNullString(updateForm.Title),
		FileName:     types.GetNullString(updateForm.FileName),
		Params:       types.GetNullString(updateForm.Params),
		HaveDedicate: types.GetNullInt64(updateForm.HaveDedicate),
	}

	// 获取所有的host
	pilipalaCommandHostDao := new(dao.PilipalaCommandHostDao)
	pilipalaCommandHosts, err := pilipalaCommandHostDao.FindByCommandId(
		pilipalaCommandProgram.Id.Int64, "pilipala_host_id")
	if err != nil {
		res.Code = 30000
		res.Message = fmt.Sprintf("获取是该命令的专用机器错误. %v", err)
		seelog.Error(res.Message)
		_c.JSON(http.StatusOK, res)
		return
	}

	needAddHostIds, needDeleteHostIds := this.GetNeedAddDeleteHostIds(pilipalaCommandHosts, updateForm.DedicateHosts)
	// 生成命令装用机器数据
	needAddHosts := make([]*model.PilipalaCommandHost, 0, 1)
	if updateForm.HaveDedicate == 1 {
		for _, hostId := range needAddHostIds {
			needAddHost := new(model.PilipalaCommandHost)
			needAddHost.PilipalaCommandProgramId = types.GetNullInt64(updateForm.Id)
			needAddHost.PilipalaHostId = types.GetNullInt64(hostId)
			needAddHosts = append(needAddHosts, needAddHost)
		}
	}

	// 将数据保存到数据库
	txDao := new(dao.PilipalaTXDao)
	err = txDao.UpdateCommandAndHost(pilipalaCommandProgram, needAddHosts, needDeleteHostIds)
	if err != nil {
		res.Code = 30000
		res.Message = fmt.Sprintf("更新命令出错. %v", err)
		seelog.Error(res.Message)
		_c.JSON(http.StatusOK, res)
		return
	}

	// 有修改文件的话, 需要移动文件
	if updateForm.TmpFileName != "" {
		config := config.GetPiliStartConfig()
		oldFilePath := config.UploadCommandFilePath(updateForm.TmpFileName)
		newFilePath := config.CommandFilePath(updateForm.FileName)
		if err := os.Rename(oldFilePath, newFilePath); err != nil {
			res.Code = 30000
			res.Message = fmt.Sprintf("将命令文件移动到指定目录出错. %v", err)
			seelog.Error(res.Message)
			_c.JSON(http.StatusOK, res)
			return
		}
	}

	_c.JSON(http.StatusOK, res)
}

// 上传命令创建命令
func (this *PilipalaCommandProgramHandler) UploadCreateCommand(_c *gin.Context) {
	res := message.NewResponseMessage()
	file, err := _c.FormFile("file")
	if err != nil {
		res.Code = 30000
		res.Message = fmt.Sprintf("后台获取文件失败. %v", err)
		seelog.Error(res.Message)
		_c.JSON(http.StatusOK, res)
		return
	}

	// 判断文件是否存在
	pilipalaCommandProgramDao := new(dao.PilipalaCommandProgramDao)
	fileExists, err := pilipalaCommandProgramDao.FileExists(file.Filename)
	if err != nil {
		res.Code = 30000
		res.Message = fmt.Sprintf("判断文件名是否存在出错. %v", err)
		seelog.Error(res.Message)
		_c.JSON(http.StatusOK, res)
		return
	}
	if fileExists {
		res.Code = 30000
		res.Message = fmt.Sprintf("文件 %v 已经存在.", file.Filename)
		seelog.Error(res.Message)
		_c.JSON(http.StatusOK, res)
		return
	}

	piliStartConfig := config.GetPiliStartConfig()
	tmpFileName := common.CreateUUIDFileName()
	filePath := piliStartConfig.UploadCommandFilePath(tmpFileName)
	seelog.Warnf("Uplaod File(create) Name: %v, FilePath: %v", file.Filename, filePath)
	if err := _c.SaveUploadedFile(file, filePath); err != nil {
		res.Code = 30000
		res.Message = fmt.Sprintf("后台保存文件失败. %v", err)
		seelog.Error(res.Message)
		_c.JSON(http.StatusOK, res)
		return
	}

	res.Data["file_name"] = file.Filename
	res.Data["tmp_file_name"] = tmpFileName
	_c.JSON(http.StatusOK, res)
}

// 上传命令编辑命令
func (this *PilipalaCommandProgramHandler) UploadEditCommand(_c *gin.Context) {
	res := message.NewResponseMessage()
	id := _c.PostForm("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		res.Code = 30000
		res.Message = fmt.Sprintf("将传入id转化int出错. %v", err)
		seelog.Error(res.Message)
		_c.JSON(http.StatusOK, res)
		return
	}

	file, err := _c.FormFile("file")
	if err != nil {
		res.Code = 30000
		res.Message = fmt.Sprintf("后台获取文件失败. %v", err)
		seelog.Error(res.Message)
		_c.JSON(http.StatusOK, res)
		return
	}

	// 判断文件是否和其他命令文件有冲突
	pilipalaCommandProgramDao := new(dao.PilipalaCommandProgramDao)
	isConflict, err := pilipalaCommandProgramDao.FileIsConflict(idInt, file.Filename)
	if err != nil {
		res.Code = 30000
		res.Message = fmt.Sprintf("判断需要替换的文件是否有冲突出错. %v", err)
		seelog.Error(res.Message)
		_c.JSON(http.StatusOK, res)
		return
	}
	if isConflict {
		res.Code = 30000
		res.Message = fmt.Sprintf("文件 %v 和其他命令文件有冲突.", file.Filename)
		seelog.Error(res.Message)
		_c.JSON(http.StatusOK, res)
		return
	}

	// 将上传的文件重命名并且移动到上传目录中
	piliStartConfig := config.GetPiliStartConfig()
	tmpFileName := common.CreateUUIDFileName()
	filePath := piliStartConfig.UploadCommandFilePath(tmpFileName)
	seelog.Warnf("Uplaod File(edit) Name: %v, FilePath: %v", file.Filename, filePath)
	if err := _c.SaveUploadedFile(file, filePath); err != nil {
		res.Code = 30000
		res.Message = fmt.Sprintf("后台保存文件失败. %v", err)
		seelog.Error(res.Message)
		_c.JSON(http.StatusOK, res)
		return
	}

	res.Data["file_name"] = file.Filename
	res.Data["tmp_file_name"] = tmpFileName
	_c.JSON(http.StatusOK, res)
}

// 下载命令
func (this *PilipalaCommandProgramHandler) Download(_c *gin.Context) {
	fileName := _c.Param("command")
	if strings.TrimSpace(fileName) == "" {
		errMSG := "下载命令(没有指定命令)"
		seelog.Error(errMSG)
		_c.String(http.StatusForbidden, errMSG)
		return
	}

	// 判断命令是否存在
	pilipalaConfig := config.GetPiliStartConfig()
	fileNamePath := pilipalaConfig.CommandFilePath(fileName)
	if exists, _ := common.PathExists(fileNamePath); !exists {
		errMSG := fmt.Sprintf("命令文件不存在 %v", fileNamePath)
		seelog.Error(errMSG)
		_c.String(http.StatusForbidden, errMSG)
		return
	}

	_c.Header("Content-Description", "File Transfer")
	_c.Header("Content-Transfer-Encoding", "binary")
	_c.Header("Content-Disposition", "attachment; filename="+fileName)
	_c.Header("Content-Type", "application/octet-stream")
	_c.File(fileNamePath)
}

/* 获取需要删除和添加的hostid
Return:
	1. 需要添加的host
    2. 需要删除的host
*/
func (this *PilipalaCommandProgramHandler) GetNeedAddDeleteHostIds(
	_pilipalaCommandHosts []model.PilipalaCommandHost,
	_hostIds []int64,
) ([]int64, []int64) {
	storedHostIds := mapset.NewSet()    // 已经存在的id
	uncertainHostIds := mapset.NewSet() // 前端传入的id

	for _, pilipalaCommandHost := range _pilipalaCommandHosts {
		storedHostIds.Add(pilipalaCommandHost.PilipalaHostId.Int64)
	}

	for _, hostId := range _hostIds {
		uncertainHostIds.Add(hostId)
	}

	// 获取还需要添加的id
	needAddHostIds := make([]int64, 0, 1)
	for value := range uncertainHostIds.Difference(storedHostIds).Iter() {
		needAddHostIds = append(needAddHostIds, value.(int64))
	}

	// 获取还需要删除的id
	needDeleteHostIds := make([]int64, 0, 1)
	for value := range storedHostIds.Difference(uncertainHostIds).Iter() {
		needDeleteHostIds = append(needDeleteHostIds, value.(int64))
	}

	return needAddHostIds, needDeleteHostIds
}
