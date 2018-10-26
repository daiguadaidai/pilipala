package model

import (
	"github.com/daiguadaidai/pilipala/common/sql_type_util"
)

type PilipalaCommandHost struct {
	Id                        sql_type_util.NullInt64 `gorm:"primary_key;not null;AUTO_INCREMENT"`                                              // 主键ID
	PilipalaCommandProgramId  sql_type_util.NullInt64 `gorm:"column:pilipala_command_program_id;type:bigint;not null"`                          // 命令的主键, 执行了是那个命令
	UpdatedAt                 sql_type_util.NullTime  `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"` // 更新时间
	CreatedAt                 sql_type_util.NullTime  `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP"`                             // 创建时间
}

func (PilipalaCommandHost) TableName() string {
	return "pilipala_command_host"
}
