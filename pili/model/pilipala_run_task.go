package model

import (
	"github.com/daiguadaidai/pilipala/common/types"
)

const (
	_ = iota
	TASK_STATUS_SUCCESS
	TASK_STATUS_RUNNING
	TASK_STATUS_FAIL
)

type PilipalaRunTask struct {
	Id                       types.NullInt64  `gorm:"primary_key;not null;AUTO_INCREMENT"`                     // 主键ID
	PilipalaCommandProgramId types.NullInt64  `gorm:"column:pilipala_command_program_id;type:bigint;not null"` // 命令的主键, 执行了是那个命令
	TaskUUID                 types.NullString `gorm:"column:task_uuid;type:varchar(30);not null"`              // 任务UUID
	Host                     types.NullString `gorm:"column:host;type:varchar(50);not null"`                   // host
	FileName                 types.NullString `gorm:"column:file_name;type:varchar(150);not null"`             // 命令文件名
	Params                   types.NullString `gorm:"column:params;type:varchar(500);not null;default:''"`     // 执行任务的参数
	Pid                      types.NullInt64  `gorm:"column:pid;type:bigint;not null;default:0"`               // 父id, 任务也有层级结构
	LogPath                  types.NullString `gorm:"column:log_path;type:varchar(255);not null;default:''"`
	NotifyInfo               types.NullString `gorm:"column:notify_info;type:varchar(500);not null;default:''"`                         // 会实时读该字段的信息, 一般外部其他程序可以通过修改这个来和任务进行通讯.
	RealInfo                 types.NullString `gorm:"column:real_info;type:varchar(500);not null;default:''"`                           // 保存了一些实时需要持久化的信息
	UpdatedAt                types.NullTime   `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"` // 更新时间
	CreatedAt                types.NullTime   `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP"`                             // 创建时间
	Status                   types.NullInt64  `gorm:"column:status"`
}

func (PilipalaRunTask) TableName() string {
	return "pilipala_run_task"
}
