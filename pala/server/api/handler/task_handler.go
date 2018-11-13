package handler

import (
	"fmt"
	"github.com/cihub/seelog"
	"github.com/daiguadaidai/pilipala/common"
	"github.com/daiguadaidai/pilipala/pala/config"
	"github.com/daiguadaidai/pilipala/pala/server/api/handler/form"
	"github.com/daiguadaidai/pilipala/pala/server/api/message"
	"github.com/daiguadaidai/pilipala/pala/server/task"
	"github.com/gin-gonic/gin"
	"github.com/hpcloud/tail"
	"net/http"
	"os"
	"strings"
)

func (this *TaskHandler) Register(_router *gin.Engine) {
	_router.POST("/pala/tasks/start", this.Start)
	_router.GET("/pala/tasks/remove/:command", this.RemoveCommand)
	_router.GET("/pala/tasks/kill/:task_uuid", this.KillCommand)
	_router.GET("/pala/tasks/tail", this.Tail)
}

func init() {
	handler := new(TaskHandler)
	AddHandler(handler)
}

type TaskHandler struct{}

func (this *TaskHandler) Start(_c *gin.Context) {
	res := message.NewResponseMessage()

	// 获取前端传来参数
	t := new(task.Task)
	if err := _c.ShouldBind(t); err != nil {
		res.Code = 30000
		res.Message = fmt.Sprintf("获取参数错误. %v", err)
		seelog.Error(res.Message)
		_c.JSON(http.StatusOK, res)
		return
	}
	seelog.Info(t)

	task.TaskChan <- t

	_c.JSON(http.StatusOK, res)
}

// 删除命令文件
func (this *TaskHandler) RemoveCommand(_c *gin.Context) {
	res := message.NewResponseMessage()
	command := _c.Param("command")
	if command == "" {
		res.Code = 30000
		res.Message = fmt.Sprintf("没有输入需要删除的命令")
		seelog.Error(res.Message)
		_c.JSON(http.StatusOK, res)
		return
	}
	cf := config.GetPalaStartConfig()
	if err := os.Remove(cf.GetCommandFilePath(command)); err != nil {
		res.Code = 30000
		res.Message = fmt.Sprintf("删除命令失败. command: %s. %v", command, err)
		seelog.Error(res.Message)
		_c.JSON(http.StatusOK, res)
		return
	}

	_c.JSON(http.StatusOK, res)
}

// kill掉一个命令
func (this *TaskHandler) KillCommand(_c *gin.Context) {
	res := message.NewResponseMessage()
	taskUUID := _c.Param("task_uuid")
	if taskUUID == "" {
		res.Code = 30000
		res.Message = fmt.Sprintf("没有输入需要强制停止的任务")
		seelog.Error(res.Message)
		_c.JSON(http.StatusOK, res)
		return
	}
	pid, ok := task.RunningTaskMap.Load(taskUUID)
	if !ok {
		res.Code = 30000
		res.Message = fmt.Sprintf("该任务已经不存在了")
		seelog.Error(res.Message)
		_c.JSON(http.StatusOK, res)
		return
	}
	if err := common.KillProcess(pid.(int)); err != nil {
		res.Code = 30000
		res.Message = fmt.Sprintf("停止任务失败: pid: %d. %v", pid, err)
		seelog.Error(res.Message)
		_c.JSON(http.StatusOK, res)
		return
	}
	task.RunningTaskMap.Delete(taskUUID)

	_c.JSON(http.StatusOK, res)
}

func (this *TaskHandler) Tail(_c *gin.Context) {
	res := message.NewResponseMessage()

	tff, err := form.NewTailFileForm(_c)
	if err != nil {
		res.Code = 30000
		res.Message = err.Error()
		seelog.Error(res.Message)
		_c.JSON(http.StatusOK, res)
		return
	}

	cf := tail.Config{
		MustExist: true,
		Location: &tail.SeekInfo{
			Offset: tff.Size,
			Whence: 2,
		},
	}
	tt, err := tail.TailFile(tff.LogPath, cf)
	if err != nil {
		res.Code = 30000
		res.Message = err.Error()
		seelog.Error(res.Message)
		_c.JSON(http.StatusOK, res)
		return
	}
	lines := make([]string, 1)
	for line := range tt.Lines {
		lines = append(lines, line.Text)
	}
	res.Data = strings.Join(lines, "\n")
	_c.JSON(http.StatusOK, res)
}
