package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/daiguadaidai/pilipala/pili/dao"
	"github.com/daiguadaidai/pilipala/pili/server/message"
	"fmt"
	"github.com/cihub/seelog"
	"strconv"
	"net/http"
)

type PilipalaHostHandler struct {}

func (this *PilipalaHostHandler) Create(_c *gin.Context) {
}

func (this *PilipalaHostHandler) Delete(_c *gin.Context) {
}

func (this *PilipalaHostHandler) Update(_c *gin.Context) {
}

func (this *PilipalaHostHandler) PartialUpdate(_c *gin.Context) {
}

func (this *PilipalaHostHandler) List(_c *gin.Context) {
	minPK := _c.DefaultQuery("min_pk", "")
	offset := _c.DefaultQuery("offset", "50")
	columnStr := "*"

	res := message.NewResponseMessage()

	pilipalaHostDao := new(dao.PilipalaHostDao)

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
		items, err := pilipalaHostDao.PaginationFirstFind(offsetInt, columnStr)
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
	items, err := pilipalaHostDao.PaginationFind(minPKInt, offsetInt, columnStr)
	if err != nil {
		res.Code = 30000
		res.Message = fmt.Sprintf("获取数据失败. %v", err)
		seelog.Error(res.Message)
	}
	res.Data["items"] = items

	_c.JSON(http.StatusOK, res)
}

func (this *PilipalaHostHandler) ListAll(_c *gin.Context) {
	columnStr := _c.DefaultQuery("columnStr", "*")
	isDedicate := _c.DefaultQuery("isDedicate", "0")
	// 是否是专用机器
	isDedicateInt, err := strconv.Atoi(isDedicate)
	if err != nil {
		isDedicateInt = 0
	}

	res := message.NewResponseMessage()

	pilipalaHostDao := new(dao.PilipalaHostDao)

	items, err := pilipalaHostDao.FindAll(columnStr, isDedicateInt)
	if err != nil {
		res.Code = 30000
		res.Message = fmt.Sprintf("获取数据失败. %v", err)
		seelog.Error(res.Message)
	}
	res.Data["items"] = items

	_c.JSON(http.StatusOK, res)
}


func (this *PilipalaHostHandler) Retrieve(_c *gin.Context) {
}

func (this *PilipalaHostHandler) Register(_router *gin.Engine) {
	_router.POST("/pilipala_host", this.Create)
	_router.DELETE("/pilipala_host/:pk", this.Delete)
	_router.PUT("/pilipala_host/:pk", this.Update)
	_router.PATCH("/pilipala_host/:pk", this.PartialUpdate)
	_router.GET("/pilipala_host", this.List)
	_router.GET("/pilipala_host/:pk", this.Retrieve)
	_router.GET("/pilipala_host_all", this.ListAll)
}

func init() {
	handler := new(PilipalaHostHandler)
	AddHandler(handler)
}
