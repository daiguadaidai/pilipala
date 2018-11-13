package model

import (
	"github.com/daiguadaidai/pilipala/common/types"
)

type PilipalaCommandProgram struct {
	Id           types.NullInt64  `json:"id"            gorm:"primary_key;not null;AUTO_INCREMENT"`                                              // 主键ID
	Title        types.NullString `json:"title"         gorm:"column:title;type:varchar(150);not null"`                                          // 干什么的
	FileName     types.NullString `json:"file_name"     gorm:"column:file_name;type:varchar(150);not null"`                                      // 文件名
	HaveDedicate types.NullInt64  `json:"have_dedicate" gorm:"column:have_dedicate;not null;default:0"`                                          // 主键ID
	Params       types.NullString `json:"params"        gorm:"column:params;type:varchar(500);not null;default:''"`                              // 命令参数
	UpdatedAt    types.NullTime   `json:"updated_at"    gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"` // 更新时间
	CreatedAt    types.NullTime   `json:"created_at"    gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP"`                             // 创建时间
}

func (PilipalaCommandProgram) TableName() string {
	return "pilipala_command_program"
}
