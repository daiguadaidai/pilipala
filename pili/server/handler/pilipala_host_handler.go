package handler

import (
	"fmt"
	"github.com/cihub/seelog"
	"github.com/daiguadaidai/pilipala/pili/dao"
	"github.com/daiguadaidai/pilipala/pili/server/message"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (this *PilipalaHostHandler) Register(_router *gin.Engine) {
	_router.GET("/pili/hosts", this.List)
	_router.GET("/pili/hosts/heartbeat/:host", this.Heartbeat)
}

func init() {
	handler := new(PilipalaHostHandler)
	AddHandler(handler)
}

type PilipalaHostHandler struct{}

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

// 接收心跳
func (this *PilipalaHostHandler) Heartbeat(_c *gin.Context) {
	res := message.NewResponseMessage()

	host := _c.Param("host")
	if host == "" {
		res.Code = 30000
		res.Message = "请输入正确的host"
		_c.JSON(http.StatusOK, res)
		return
	}

	d := new(dao.PilipalaHostDao)
	if err := d.UpdateIsValidByHost(host, 1); err != nil {
		res.Code = 30000
		res.Message = fmt.Sprintf("更新%s心跳失败. %v", host, err)
		_c.JSON(http.StatusOK, res)
		return
	}

	_c.JSON(http.StatusOK, res)
}
