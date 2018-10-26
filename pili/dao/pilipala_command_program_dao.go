package dao

import (
	"github.com/daiguadaidai/pilipala/pili/model"
	"github.com/jinzhu/gorm"
	"github.com/daiguadaidai/pilipala/pili/gdbc"
)

type PilipalaCommandProgramDao struct{}

// 创建一个命令
func (this *PilipalaCommandProgramDao) Create(_pilipalaCommandProgram *model.PilipalaCommandProgram) error {
	ormInstance := gdbc.GetOrmInstance()

	err := ormInstance.DB.Create(_pilipalaCommandProgram).Error
	if err != nil {
		return err
	}

	return nil
}

// 通过ID获取命令
func (this *PilipalaCommandProgramDao) GetByID(
	_id int,
	_columnStr string,
) (*model.PilipalaCommandProgram, error) {
	ormInstance := gdbc.GetOrmInstance()

	pilipalaCommandProgram := new(model.PilipalaCommandProgram)
	err := ormInstance.DB.Select(_columnStr).
		Where("id = ?", _id).
		First(&pilipalaCommandProgram).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return pilipalaCommandProgram, nil
}

// 获取第一页数据
func (this *PilipalaCommandProgramDao) PaginationFirstFind(
	_offset int,
	_columnStr string,
) ([]model.PilipalaCommandProgram, error) {
	ormInstance := gdbc.GetOrmInstance()

	pilipalaCommandPrograms := []model.PilipalaCommandProgram{}
	err := ormInstance.DB.Select(_columnStr).
		Order("id DESC").
		Limit(_offset).
		Find(&pilipalaCommandPrograms).Error
	if err != nil {
		return pilipalaCommandPrograms, err
	}

	return pilipalaCommandPrograms, nil
}

// 分页获取数据
func (this *PilipalaCommandProgramDao) PaginationFind(
	_min_pk int,
	_offset int,
	_columnStr string,
) ([]model.PilipalaCommandProgram, error) {
	ormInstance := gdbc.GetOrmInstance()

	pilipalaCommandPrograms := []model.PilipalaCommandProgram{}
	err := ormInstance.DB.Select(_columnStr).
		Where("id < ?", _min_pk).
		Order("id DESC").
		Limit(_offset).
		Find(&pilipalaCommandPrograms).Error
	if err != nil {
		return pilipalaCommandPrograms, err
	}

	return pilipalaCommandPrograms, nil
}
