package handler

import (
	"encoding/json"
	"fmt"
	"github.com/cihub/seelog"
	"github.com/daiguadaidai/pilipala/common"
	"github.com/daiguadaidai/pilipala/common/types"
	"github.com/daiguadaidai/pilipala/pili/config"
	"github.com/daiguadaidai/pilipala/pili/dao"
	"github.com/daiguadaidai/pilipala/pili/model"
	"github.com/daiguadaidai/pilipala/pili/server/form"
	"github.com/daiguadaidai/pilipala/pili/server/message"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
)

func init() {
	handler := new(PilipalaRunTaskHandler)
	AddHandler(handler)
}

// 注册route
func (this *PilipalaRunTaskHandler) Register(_router *gin.Engine) {
	_router.GET("/pili/tasks/start", this.Start)
	_router.POST("/pili/tasks/start", this.Start)
	_router.GET("/pili/tasks/success/:uuid", this.TaskSuccess)
	_router.GET("/pili/tasks/fail/:uuid", this.TaskFail)
	_router.GET("/pili/tasks/running/:uuid", this.TaskRunning)
	_router.PUT("/pili/tasks", this.Update)
	_router.GET("/pili/tasks/tail/:uuid", this.TaskTail)
}

type PilipalaRunTaskHandler struct{}

// 执行一个任务
func (this *PilipalaRunTaskHandler) Start(_c *gin.Context) {
	res := message.NewResponseMessage()

	// 获取前端传来参数
	startTaskForm := new(form.StartTaskForm)
	if err := _c.ShouldBind(startTaskForm); err != nil {
		res.Code = 30000
		res.Message = fmt.Sprintf("获取参数错误. %v", err)
		seelog.Error(res.Message)
		_c.JSON(http.StatusOK, res)
		return
	}
	seelog.Info(startTaskForm)

	res = this.run(startTaskForm)
	_c.JSON(http.StatusOK, res)
}

// 标记任务成功
func (this *PilipalaRunTaskHandler) TaskSuccess(_c *gin.Context) {
	res := message.NewResponseMessage()

	uuid := _c.Param("uuid")
	taskDao := new(dao.PilipalaRunTaskDao)
	if err := taskDao.UpdateTaskStatus(uuid, model.TASK_STATUS_SUCCESS); err != nil {
		res.Code = 30000
		res.Message = fmt.Sprintf("设置任务<执行成功>失败. %v", err)
		seelog.Error(res.Message)
		_c.JSON(http.StatusOK, res)
		return
	}

	// 获取task
	task, err := taskDao.GetByTaskUUID(uuid, "host")
	if err != nil {
		seelog.Warnf("获取指定任务出错 %v. %v", uuid, err)
		_c.JSON(http.StatusOK, res)
		return
	}

	// 该host上的任务数 -1
	hostDao := new(dao.PilipalaHostDao)
	if err := hostDao.DecrTaskByHost(task.Host.String); err != nil {
		seelog.Warnf("任务启动成功. 添加当前host(%v)任务数失败", task.Host.String)
	}

	_c.JSON(http.StatusOK, res)
}

// 标记任务失败
func (this *PilipalaRunTaskHandler) TaskFail(_c *gin.Context) {
	res := message.NewResponseMessage()

	uuid := _c.Param("uuid")
	taskDao := new(dao.PilipalaRunTaskDao)
	if err := taskDao.UpdateTaskStatus(uuid, model.TASK_STATUS_FAIL); err != nil {
		res.Code = 30000
		res.Message = fmt.Sprintf("设置任务<执行失败>失败. %v", err)
		seelog.Error(res.Message)
		_c.JSON(http.StatusOK, res)
		return
	}

	// 获取task
	task, err := taskDao.GetByTaskUUID(uuid, "host")
	if err != nil {
		seelog.Warnf("获取指定任务出错 %v. %v", uuid, err)
		_c.JSON(http.StatusOK, res)
		return
	}

	// 该host上的任务数 -1
	hostDao := new(dao.PilipalaHostDao)
	if err := hostDao.DecrTaskByHost(task.Host.String); err != nil {
		seelog.Warnf("任务启动成功. 减少当前host(%v)任务数失败", task.Host.String)
	}

	_c.JSON(http.StatusOK, res)
}

// 标记任务运行中
func (this *PilipalaRunTaskHandler) TaskRunning(_c *gin.Context) {
	res := message.NewResponseMessage()

	uuid := _c.Param("uuid")
	taskDao := new(dao.PilipalaRunTaskDao)
	if err := taskDao.UpdateTaskStatus(uuid, model.TASK_STATUS_RUNNING); err != nil {
		res.Code = 30000
		res.Message = fmt.Sprintf("设置任务<执行中>失败. %v", err)
		seelog.Error(res.Message)
		_c.JSON(http.StatusOK, res)
		return
	}

	// 获取task
	task, err := taskDao.GetByTaskUUID(uuid, "host")
	if err != nil {
		seelog.Warnf("获取指定任务出错 %v. %v", uuid, err)
		_c.JSON(http.StatusOK, res)
		return
	}

	// 该host上的任务数 +1
	hostDao := new(dao.PilipalaHostDao)
	if err := hostDao.IncrTaskByHost(task.Host.String); err != nil {
		seelog.Warnf("任务启动成功. 减少当前host(%v)任务数失败", task.Host.String)
	}

	_c.JSON(http.StatusOK, res)
}

func (this *PilipalaRunTaskHandler) Update(_c *gin.Context) {
	res := message.NewResponseMessage()

	upForm := new(form.UpdateTaskForm)
	if err := _c.ShouldBind(upForm); err != nil {
		res.Code = 30000
		res.Message = fmt.Sprintf("解析更新任务参数失败. %v", err)
		seelog.Error(res.Message)
		_c.JSON(http.StatusOK, res)
		return
	}

	taskDao := new(dao.PilipalaRunTaskDao)
	if err := taskDao.UpdateByUUID(upForm.NewTaskModel()); err != nil {
		res.Code = 30000
		res.Message = fmt.Sprintf("更新任务失败. %v", err)
		seelog.Error(res.Message)
		_c.JSON(http.StatusOK, res)
		return
	}

	_c.JSON(http.StatusOK, res)
}

func (this *PilipalaRunTaskHandler) TaskTail(_c *gin.Context) {
	res := message.NewResponseMessage()

	uuid := _c.Param("uuid")
	if uuid == "" {
		res.Code = 30000
		res.Message = fmt.Sprintf("必须输入任务uuid")
		seelog.Error(res.Message)
		_c.JSON(http.StatusOK, res)
		return
	}
	size := form.GetTailSize(_c)

	taskDao := new(dao.PilipalaRunTaskDao)
	t, err := taskDao.GetByTaskUUID(uuid, "*")
	if err != nil {
		res.Code = 30000
		res.Message = fmt.Sprintf("获取日志路径失败. %v", err)
		seelog.Error(res.Message)
		_c.JSON(http.StatusOK, res)
		return
	}
	// 获取url
	tailURL := fmt.Sprintf(config.PALA_TASK_TAIL_LOG_URL, t.Host.String)
	queryMap := make(map[string]interface{})
	queryMap["size"] = size
	queryMap["log_path"] = t.LogPath.String
	query := common.GetURLQuery(queryMap)

	tailLogMSG := new(message.TailLogMessage)
	body, err := common.GetUrlRaw(tailURL, query)
	err = json.Unmarshal(body, tailLogMSG)
	if err != nil {
		res.Code = 30000
		res.Message = fmt.Sprintf("解析日志数据失败. %v", err)
		seelog.Error(res.Message)
		_c.JSON(http.StatusOK, res)
		return
	}

	if tailLogMSG.Code != 20000 {
		res.Code = 30000
		_c.JSON(http.StatusOK, tailLogMSG)
		return
	}

	_c.JSON(http.StatusOK, tailLogMSG)
}

func (this *PilipalaRunTaskHandler) run(_form *form.StartTaskForm) *message.ResponseMessage {
	res := message.NewResponseMessage()

	commandProgramDao := new(dao.PilipalaCommandProgramDao)
	commandProgram, err := commandProgramDao.GetByFileName("id, have_dedicate", _form.Command)
	if err != nil {
		res.Code = 30000
		res.Message = fmt.Sprintf("没有指定的命令. %v", err)
		seelog.Error(res.Message)
		return res
	}

	// 获取命令自行的机器
	host, err := this.getRunHost(_form.Command, commandProgram)
	if err != nil {
		res.Code = 30000
		res.Message = fmt.Sprintf("获取命令执行机器失败. %v", err)
		seelog.Error(res.Message)
		return res
	}

	// 获取pala执行命令url
	runTaskUrl := fmt.Sprintf(config.PALA_TASK_RUN_URL, host)
	taskUUID := common.GetUUID()
	data := _form.GetPalaRunParam(taskUUID)

	// 创建一个任务记录
	taskDao := new(dao.PilipalaRunTaskDao)
	task := &model.PilipalaRunTask{
		TaskUUID:                 types.GetNullString(taskUUID),
		PilipalaCommandProgramId: commandProgram.Id,
		Host:     types.GetNullString(host),
		FileName: types.GetNullString(data["command"]),
		Params:   types.GetNullString(data["params"]),
		Pid:      types.GetNullInt64(_form.ParentId),
		Status:   types.GetNullInt64(model.TASK_STATUS_RUNNING),
	}
	if err = taskDao.Create(task); err != nil {
		res.Code = 30000
		res.Message = fmt.Sprintf("创建任务失败. %v", err)
		seelog.Error(res.Message)
		return res
	}

	// 远程启动命令
	if err := common.PostUrl(runTaskUrl, data); err != nil {
		switch err.(type) {
		case *url.Error:
			hostDao := new(dao.PilipalaHostDao)
			seelog.Errorf("链接host %s 有问题, 将标记该host不可用", host)
			if err1 := hostDao.UpdateIsValidByHost(host, 0); err1 != nil {
				seelog.Errorf("标记host %s 失败. %v", host, err1)
			}
		}
		res.Code = 30000
		res.Message = fmt.Sprintf("执行任务失败. %v", err)
		seelog.Error(res.Message)
		if err1 := taskDao.UpdateTaskStatus(task.TaskUUID.String, model.TASK_STATUS_FAIL); err1 != nil {
			res.Message = fmt.Sprintf("%v. 更新任务状态错误. %v", res.Message, err)
		}
		return res
	}

	// 该host上的任务数 +1
	hostDao := new(dao.PilipalaHostDao)
	if err := hostDao.IncrTaskByHost(host); err != nil {
		seelog.Warnf("任务启动成功. 添加当前host(%v)任务数失败", host)
	}

	return res
}

// 获取任务执行的host
func (this *PilipalaRunTaskHandler) getRunHost(
	_fileName string,
	_commandProgram *model.PilipalaCommandProgram,
) (string, error) {
	where := new(model.PilipalaHost)
	where.IsValid = types.GetNullInt64(1)

	hostDao := new(dao.PilipalaHostDao)
	if _commandProgram.HaveDedicate.Int64 == 0 { // 没有专用机器, 使用共用机器
		where.IsDedicate = types.GetNullInt64(0)
		host, err := hostDao.GetOptimalHost("host", where, []int64{})
		if err != nil {
			return "", fmt.Errorf("找不到共用机器: %v", err)
		}
		return host.Host.String, nil
	}

	// 使用专用机器
	// 通过命令id获取host id
	commandHostDao := new(dao.PilipalaCommandHostDao)
	ids, err := commandHostDao.FindIdsByCommandId(_commandProgram.Id.Int64)
	if err != nil {
		return "", fmt.Errorf("该命令还没有添加装用机器, %v", err)
	}

	// 获取host
	host, err := hostDao.GetOptimalHost("host", where, ids)
	if err != nil {
		return "", fmt.Errorf("找不到装用机器: %v", err)
	}
	return host.Host.String, nil
}
