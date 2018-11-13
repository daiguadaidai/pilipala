package dao

import (
	"github.com/daiguadaidai/pilipala/pili/model"
	"github.com/jinzhu/gorm"
	"github.com/daiguadaidai/pilipala/pili/gdbc"
)

type PilipalaCommandHostDao struct{}


// 通过command_id获取commmand host 相关信息
func (this *PilipalaCommandHostDao) FindByCommandId(
	_commandId int64,
	_columnStr string,
) ([]model.PilipalaCommandHost, error) {
	ormInstance := gdbc.GetOrmInstance()

	pilipalaCommandHosts := []model.PilipalaCommandHost{}
	if err := ormInstance.DB.Select(_columnStr).
		Where("pilipala_command_program_id = ?", _commandId).
		Find(&pilipalaCommandHosts).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return pilipalaCommandHosts, nil
		}
		return nil, err
	}

	return pilipalaCommandHosts, nil
}

// 通过command_id获取commmand host 相关信息
func (this *PilipalaCommandHostDao) FindIdsByCommandId(
	_commandId int64,
) ([]int64, error) {
	ormInstance := gdbc.GetOrmInstance()

	pilipalaCommandHosts := []model.PilipalaCommandHost{}
	if err := ormInstance.DB.Select("pilipala_host_id").
		Where("pilipala_command_program_id = ?", _commandId).
		Find(&pilipalaCommandHosts).Error; err != nil {
		return nil, err
	}

	ids := make([]int64, 0, 1)
	for _, pilipalaCommandHost := range pilipalaCommandHosts {
		ids = append(ids, pilipalaCommandHost.PilipalaHostId.Int64)
	}

	return ids, nil
}
