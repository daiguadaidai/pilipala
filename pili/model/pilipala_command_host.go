package model

import (
	"github.com/daiguadaidai/pilipala/common/types"
)

type PilipalaCommandHost struct {
	Id                       types.NullInt64 `json:"id"                          gorm:"primary_key;not null;AUTO_INCREMENT"`                                              // 主键ID
	PilipalaCommandProgramId types.NullInt64 `json:"pilipala_command_program_id" gorm:"column:pilipala_command_program_id;type:bigint;not null"`                          // 命令的主键, 执行了是那个命令
	PilipalaHostId           types.NullInt64 `json:"pilipala_host_id"            gorm:"column:pilipala_host_id;type:bigint;not null"`                                     // host id
	UpdatedAt                types.NullTime  `json:"updated_at"                  gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"` // 更新时间
	CreatedAt                types.NullTime  `json:"created_at"                  gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP"`                             // 创建时间
}

func (PilipalaCommandHost) TableName() string {
	return "pilipala_command_host"
}
