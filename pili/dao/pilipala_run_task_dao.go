package dao

import (
	"github.com/daiguadaidai/pilipala/pili/model"
	"github.com/jinzhu/gorm"
	"github.com/daiguadaidai/pilipala/pili/gdbc"
)

type PilipalaRunTaskDao struct{}

// 通过ID获取任务
func (this *PilipalaRunTaskDao) GetByID(
	_id int,
	_columnStr string,
) (*model.PilipalaRunTask, error) {
	ormInstance := gdbc.GetOrmInstance()

	pilipalaRunTask := new(model.PilipalaRunTask)
	err := ormInstance.DB.Select(_columnStr).
		Where("id = ?", _id).
		First(&pilipalaRunTask).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return pilipalaRunTask, nil
}

// 通过task_uuid获取任务
func (this *PilipalaRunTaskDao) GetByTaskUUID(
	_taskUUID string,
	_columnStr string,
) (*model.PilipalaRunTask, error) {
	ormInstance := gdbc.GetOrmInstance()

	pilipalaRunTask := new(model.PilipalaRunTask)
	err := ormInstance.DB.Select(_columnStr).
		Where("task_uuid = ?", _taskUUID).
		First(&pilipalaRunTask).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return pilipalaRunTask, nil
}

// 创建一个任务
func (this *PilipalaRunTaskDao) Create(_pilipalaRunTask *model.PilipalaRunTask) error {
	ormInstance := gdbc.GetOrmInstance()

	err := ormInstance.DB.Create(_pilipalaRunTask).Error
	if err != nil {
		return err
	}

	return nil
}
