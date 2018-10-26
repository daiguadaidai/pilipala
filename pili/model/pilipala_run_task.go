package model

import (
	"github.com/daiguadaidai/pilipala/common/sql_type_util"
)

type PilipalaRunTask struct {
	Id                        sql_type_util.NullInt64  `gorm:"primary_key;not null;AUTO_INCREMENT"`                                              // 主键ID
	PilipalaCommandProgramId  sql_type_util.NullInt64  `gorm:"column:pilipala_command_program_id;type:bigint;not null"`                          // 命令的主键, 执行了是那个命令
	PilipalaHostId            sql_type_util.NullInt64  `gorm:"column:pilipala_host_id;type:bigint;not null"`                          // 命令的主键, 执行了是那个命令
	TaskUUID                  sql_type_util.NullString `gorm:"column:task_uuid;type:varchar(20);not null"`                                       // 任务UUID
	Host                      sql_type_util.NullString `gorm:"column:host;type:varchar(50);not null"`                                            // host
	FileName                  sql_type_util.NullString `gorm:"column:file_name;type:varchar(150);not null"`                                      // 命令文件名
	Params                    sql_type_util.NullString `gorm:"column:params;type:varchar(500);not null;default:''"`                              // 执行任务的参数
	Pid                       sql_type_util.NullInt64  `gorm:"column:pid;type:bigint;not null;default:0"`                                        // 父id, 任务也有层级结构
	NotifyInfo                sql_type_util.NullString `gorm:"column:notify_info;type:varchar(500);not null;default:''"`                         // 会实时读该字段的信息, 一般外部其他程序可以通过修改这个来和任务进行通讯.
	RealInfo                  sql_type_util.NullString `gorm:"column:real_info;type:varchar(500);not null;default:''"`                           // 保存了一些实时需要持久化的信息
	UpdatedAt                 sql_type_util.NullTime   `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"` // 更新时间
	CreatedAt                 sql_type_util.NullTime   `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP"`                             // 创建时间
}

func (PilipalaRunTask) TableName() string {
	return "pilipala_run_task"
}
