package dao

import (
	"github.com/daiguadaidai/pilipala/pili/gdbc"
	"github.com/daiguadaidai/pilipala/pili/model"
	"github.com/jinzhu/gorm"
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
	_where interface{},
) ([]model.PilipalaHost, error) {
	ormInstance := gdbc.GetOrmInstance()

	pilipalaHosts := []model.PilipalaHost{}
	err := ormInstance.DB.Select(_columnStr).
		Where(_where).
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

// 获取优的host
func (this *PilipalaHostDao) GetOptimalHost(
	_columnStr,
	_where interface{},
	_ids []int64,
) (*model.PilipalaHost, error) {
	ormInstance := gdbc.GetOrmInstance()

	pilipalaHost := new(model.PilipalaHost)
	err := ormInstance.DB.Select(_columnStr).
		Where(_ids).Where(_where).
		Order("running_task_count ASC").
		First(pilipalaHost).Error
	if err != nil {
		return pilipalaHost, err
	}

	return pilipalaHost, nil
}

// 任务数自增1
func (this *PilipalaHostDao) IncrTaskByHost(_host string) error {
	ormInstance := gdbc.GetOrmInstance()

	if err := ormInstance.DB.Model(&model.PilipalaHost{}).Where("host = ?", _host).
		Update("running_task_count", gorm.Expr("running_task_count + ?", 1)).
		Error; err != nil {
		return err
	}

	return nil
}

// 任务数自减1
func (this *PilipalaHostDao) DecrTaskByHost(_host string) error {
	ormInstance := gdbc.GetOrmInstance()

	if err := ormInstance.DB.Model(&model.PilipalaHost{}).Where("host = ?", _host).
		Update("running_task_count", gorm.Expr("running_task_count - ?", 1)).
		Error; err != nil {
		return err
	}

	return nil
}

// 任务数自减1
func (this *PilipalaHostDao) UpdateIsValidByHost(_host string, _isValid int64) error {
	ormInstance := gdbc.GetOrmInstance()

	host := new(model.PilipalaHost)
	ormInstance.DB.Model(&model.PilipalaHost{}).Where("host = ?", _host).First(host)

	if err := ormInstance.DB.Model(&model.PilipalaHost{}).Where("host = ?", _host).
		Update("is_valid", _isValid).Error; err != nil {
		return err
	}

	return nil
}
