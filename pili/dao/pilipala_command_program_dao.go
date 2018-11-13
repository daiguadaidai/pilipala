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

// 判断文件是否存在
func (this *PilipalaCommandProgramDao) FileExists(_fileName string) (bool, error) {
	ormInstance := gdbc.GetOrmInstance()

	var count int
	err := ormInstance.DB.Model(model.PilipalaCommandProgram{}).
		Where("file_name = ?", _fileName).
		Count(&count).Error
	if err != nil {
		return false, err
	}

	if count > 0 { // 已经存在
		return true, nil
	}

	return false, nil
}

// 通过文件名获取记录
func (this *PilipalaCommandProgramDao) GetByFileName(
	_columnStr string,
	_fileName string,
) (*model.PilipalaCommandProgram, error) {
	ormInstance := gdbc.GetOrmInstance()

	pilipalaCommandProgram := new(model.PilipalaCommandProgram)
	err := ormInstance.DB.Select(_columnStr).
		Where("file_name = ?", _fileName).
		First(&pilipalaCommandProgram).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return pilipalaCommandProgram, nil
}

// 通过文件名获取记录
func (this *PilipalaCommandProgramDao) FindByFileName(
	_columnStr string,
	_fileName string,
) ([]model.PilipalaCommandProgram, error) {
	ormInstance := gdbc.GetOrmInstance()

	pilipalaCommandPrograms := []model.PilipalaCommandProgram{}
	err := ormInstance.DB.Select(_columnStr).
		Where("file_name = ?", _fileName).
		Find(&pilipalaCommandPrograms).Error
	if err != nil {
		return pilipalaCommandPrograms, err
	}

	return pilipalaCommandPrograms, nil
}

// 比较命令文件是否存在和其他命令相互冲突
func (this *PilipalaCommandProgramDao) FileIsConflict(_id int, _name string) (bool, error) {
	pilipalaCommandPrograms, err := this.FindByFileName("id", _name)
	if err != nil {
		return false, err
	}
	for _, pilipalaCommandProgram := range pilipalaCommandPrograms {
		if pilipalaCommandProgram.Id.Int64 != int64(_id) {
			return true, nil
		}
	}

	return false, nil
}
