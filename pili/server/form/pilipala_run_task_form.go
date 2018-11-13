package form

import (
	"fmt"
	"github.com/Unknwon/com"
	"github.com/daiguadaidai/pilipala/common/types"
	"github.com/daiguadaidai/pilipala/pili/model"
	"github.com/gin-gonic/gin"
)

type StartTaskForm struct {
	Command        string `form:"command" json:"command" binding:"required"`
	Params         string `form:"params" json:"params"`
	ParentId       int64  `form:"parent_id" json:"parent_id"`
	IsNeedTaskUUID bool   `form:"is_need_task_uuid" json:"is_need_task_uuid"`
}

// 获取访问pala的参数
func (this *StartTaskForm) GetPalaRunParam(_uuid string) map[string]string {
	params := make(map[string]string)
	params["command"] = this.Command
	params["params"] = this.Params
	params["task_uuid"] = _uuid

	if this.IsNeedTaskUUID {
		params["params"] = fmt.Sprintf("%s --task-uuid=%s", params["params"], _uuid)
	}

	return params
}

type UpdateTaskForm struct {
	TaskUUID   types.NullString `form:"task_uuid" json:"task_uuid" binding:"required"`
	Host       types.NullString `form:"host" json:"host"`
	FileName   types.NullString `form:"file_name" json:"file_name"`
	Params     types.NullString `form:"params" json:"params"`
	Pid        types.NullInt64  `form:"pid" json:"pid"`
	LogPath    types.NullString `form:"log_path" json:"log_path"`
	NotifyInfo types.NullString `form:"notify_info" json:"notify_info"`
	RealInfo   types.NullString `form:"real_info" json:"real_info"`
	Status     types.NullInt64  `form:"status" json:"status"`
}

func (this *UpdateTaskForm) NewTaskModel() *model.PilipalaRunTask {
	return &model.PilipalaRunTask{
		TaskUUID:   this.TaskUUID,
		Host:       this.Host,
		FileName:   this.FileName,
		Params:     this.Params,
		Pid:        this.Pid,
		LogPath:    this.LogPath,
		NotifyInfo: this.NotifyInfo,
		RealInfo:   this.RealInfo,
		Status:     this.Status,
	}
}

const (
	DEFAULT_TAIL_SIZE_STR = "20480"
	DEFAULT_TAIL_SIZE     = 20480
)

func GetTailSize(_c *gin.Context) int64 {
	sizeStr := _c.DefaultQuery("size", DEFAULT_TAIL_SIZE_STR)
	size, err := com.StrTo(sizeStr).Int64()
	if err != nil {
		size = DEFAULT_TAIL_SIZE
	}
	return size
}
