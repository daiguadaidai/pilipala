package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/daiguadaidai/pilipala/pili/server/message"
	"net/http"
	"strconv"
	"fmt"
	"github.com/cihub/seelog"
	"github.com/daiguadaidai/pilipala/pili/dao"
)

func init() {
	handler := new(PilipalaCommandHostHandler)
	AddHandler(handler)
}

type PilipalaCommandHostHandler struct {}

func (this *PilipalaCommandHostHandler) Register(_router *gin.Engine) {
	_router.GET("/pili/command_host/:command_id", this.ListByCommandId)
}

func (this *PilipalaCommandHostHandler) ListByCommandId(_c *gin.Context) {
	res := message.NewResponseMessage()

	columnStr := "id, pilipala_command_program_id, pilipala_host_id"
	commandId := _c.Param("command_id")
	commandIdInt, err := strconv.Atoi(commandId)
	if err != nil {
		res.Code = 30000
		res.Message = fmt.Sprintf("传入的命令id有误. %v", err)
		seelog.Error(res.Message)
		_c.JSON(http.StatusOK, res)
		return
	}

	pilipalaCommandDao := new(dao.PilipalaCommandHostDao)
    pipalaHosts, err := pilipalaCommandDao.FindByCommandId(int64(commandIdInt), columnStr)
    if err != nil {
		res.Code = 30000
		res.Message = fmt.Sprintf("获取命令专用机器出错. %v", err)
		seelog.Error(res.Message)
		_c.JSON(http.StatusOK, res)
		return
	}

	res.Data["items"] = pipalaHosts
	_c.JSON(http.StatusOK, res)
}


