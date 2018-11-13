package dao

import (
	"github.com/daiguadaidai/pilipala/pili/model"
	"github.com/daiguadaidai/pilipala/pili/gdbc"
)

type PilipalaTXDao struct {}

// 创建命令和其专属的机器
func (this *PilipalaTXDao) CreateCommandAndHost(
	_pilipalaCommandProgram *model.PilipalaCommandProgram,
	_pilipalaCommandHosts []*model.PilipalaCommandHost,
) error {
	ormInstance := gdbc.GetOrmInstance()
	tx := ormInstance.DB.Begin()

	if err := tx.Create(_pilipalaCommandProgram).Error; err != nil {
		tx.Rollback()
		return err
	}

	if len(_pilipalaCommandHosts) > 0 {
		for _, pilipalaCommandHost := range _pilipalaCommandHosts {
			pilipalaCommandHost.PilipalaCommandProgramId = _pilipalaCommandProgram.Id
			if err := tx.Create(pilipalaCommandHost).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	tx.Commit()
	return nil
}

// 更新命令已经host相关信息
func (this *PilipalaTXDao) UpdateCommandAndHost(
	_pilipalaCommandProgram *model.PilipalaCommandProgram,
	_needAddHosts []*model.PilipalaCommandHost,
	_needDeleteHostIds []int64,
) error {
	ormInstance := gdbc.GetOrmInstance()
	tx := ormInstance.DB.Begin()

	if err := tx.Model(_pilipalaCommandProgram).
		Updates(_pilipalaCommandProgram).Error; err != nil {
		tx.Rollback()
		return err
	}

	if _pilipalaCommandProgram.HaveDedicate.Int64 == 1 {
		// 添加新增的command host
		if len(_needAddHosts) > 0 {
			for _, needAddHost := range _needAddHosts {
				if err := tx.Create(needAddHost).Error; err != nil {
					tx.Rollback()
					return err
				}
			}
		}

		// 删除不需要的command host
		if len(_needDeleteHostIds) > 0 {
			err := tx.Where("pilipala_command_program_id = ? and pilipala_host_id in(?)",
				_pilipalaCommandProgram.Id.Int64, _needDeleteHostIds).
				Delete(model.PilipalaCommandHost{}).Error
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	} else {
		if err := tx.Where("pilipala_command_program_id = ?", _pilipalaCommandProgram.Id.Int64).
			Delete(model.PilipalaCommandHost{}).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
	return nil
}


