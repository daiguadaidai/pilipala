package dao

import (
	"github.com/daiguadaidai/pilipala/pili/model"
	"github.com/jinzhu/gorm"
	"github.com/daiguadaidai/pilipala/pili/gdbc"
)

type PilipalaHostDao struct{}

// 创建一个命令
func (this *PilipalaHostDao) Create(_pilipalaHost *model.PilipalaHost) error {
	ormInstance := gdbc.GetOrmInstance()

	err := ormInstance.DB.Create(_pilipalaHost).Error
	if err != nil {
		return err
	}

	return nil
}

// 通过ID获取命令
func (this *PilipalaHostDao) GetByID(
	_id int,
	_columnStr string,
) (*model.PilipalaHost, error) {
	ormInstance := gdbc.GetOrmInstance()

	pilipalaHost := new(model.PilipalaHost)
	err := ormInstance.DB.Select(_columnStr).
		Where("id = ?", _id).
		First(&pilipalaHost).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return pilipalaHost, nil
}

// 获取所有的IP
func (this *PilipalaHostDao) FindAll(
	_columnStr string,
	_isDadicate int,
) ([]model.PilipalaHost, error) {
	ormInstance := gdbc.GetOrmInstance()

	pilipalaHosts := []model.PilipalaHost{}
	err := ormInstance.DB.Select(_columnStr).
		Where("is_dedicate = ?", _isDadicate).
		Find(&pilipalaHosts).Error
	if err != nil {
		return pilipalaHosts, err
	}

	return pilipalaHosts, nil
}

// 获取第一页数据
func (this *PilipalaHostDao) PaginationFirstFind(
	_offset int,
	_columnStr string,
) ([]model.PilipalaHost, error) {
	ormInstance := gdbc.GetOrmInstance()

	pilipalaHosts := []model.PilipalaHost{}
	err := ormInstance.DB.Select(_columnStr).
		Order("id DESC").
		Limit(_offset).
		Find(&pilipalaHosts).Error
	if err != nil {
		return pilipalaHosts, err
	}

	return pilipalaHosts, nil
}

// 分页获取数据
func (this *PilipalaHostDao) PaginationFind(
	_min_pk int,
	_offset int,
	_columnStr string,
) ([]model.PilipalaHost, error) {
	ormInstance := gdbc.GetOrmInstance()

	pilipalaHosts := []model.PilipalaHost{}
	err := ormInstance.DB.Select(_columnStr).
		Where("id < ?", _min_pk).
		Order("id DESC").
		Limit(_offset).
		Find(&pilipalaHosts).Error
	if err != nil {
		return pilipalaHosts, err
	}

	return pilipalaHosts, nil
}
