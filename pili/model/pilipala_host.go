package model

import (
	"github.com/daiguadaidai/pilipala/common/sql_type_util"
)

type PilipalaHost struct {
	Id                  sql_type_util.NullInt64  `gorm:"primary_key;not null;AUTO_INCREMENT"`                                              // 主键ID
	IsValid             sql_type_util.NullInt64  `gorm:"column:is_valid;not null;default:0"`                                               // 是否有效的
	IsDedicate          sql_type_util.NullInt64  `gorm:"column:is_dedicate;not null;default:0"`                                            // 是否是专用的
	Host                sql_type_util.NullString `gorm:"column:host;type:varchar(50);not null"`                                            // host
	RunningTaskCount    sql_type_util.NullInt64  `gorm:"column:running_task_count;not null;default:0"`                                     // 当前运行了多少个任务
	UpdatedAt           sql_type_util.NullTime `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"` // 更新时间
	CreatedAt           sql_type_util.NullTime `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP"`                             // 创建时间
}

func (PilipalaHost) TableName() string {
	return "pilipala_host"
}
