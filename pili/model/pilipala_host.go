package model

import (
	"github.com/daiguadaidai/pilipala/common/types"
)

type PilipalaHost struct {
	Id               types.NullInt64  `json:"id"                 gorm:"primary_key;not null;AUTO_INCREMENT"`                                              // 主键ID
	IsValid          types.NullInt64  `json:"is_valid"           gorm:"column:is_valid;not null;default:0"`                                               // 是否有效的
	IsDedicate       types.NullInt64  `json:"is_dedicate"        gorm:"column:is_dedicate;not null;default:0"`                                            // 是否是专用的
	Host             types.NullString `json:"host"               gorm:"column:host;type:varchar(50);not null"`                                            // host
	RunningTaskCount types.NullInt64  `json:"running_task_count" gorm:"column:running_task_count;not null;default:0"`                                     // 当前运行了多少个任务
	UpdatedAt        types.NullTime   `json:"updated_at"         gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"` // 更新时间
	CreatedAt        types.NullTime   `json:"created_at"         gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP"`                             // 创建时间
}

func (PilipalaHost) TableName() string {
	return "pilipala_host"
}
