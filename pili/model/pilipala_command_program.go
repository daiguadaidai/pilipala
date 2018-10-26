package model

import (
	"github.com/daiguadaidai/pilipala/common/sql_type_util"
)

type PilipalaCommandProgram struct {
	Id            sql_type_util.NullInt64  `gorm:"primary_key;not null;AUTO_INCREMENT"`                                              // 主键ID
	Title         sql_type_util.NullString `gorm:"column:title;type:varchar(150);not null"`                                          // 干什么的
	FileName      sql_type_util.NullString `gorm:"column:file_name;type:varchar(150);not null"`                                      // 文件名
	HaveDedicate  sql_type_util.NullInt64  `gorm:"column:have_dedicate;not null;default:0"`                                              // 主键ID
	Params        sql_type_util.NullString `gorm:"column:params;type:varchar(500);not null;default:''"`                              // 命令参数
	UpdatedAt     sql_type_util.NullTime   `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"` // 更新时间
	CreatedAt     sql_type_util.NullTime   `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP"`                             // 创建时间
}

func (PilipalaCommandProgram) TableName() string {
	return "pilipala_command_program"
}
