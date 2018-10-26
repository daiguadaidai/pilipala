package common

import (
	"os"
	"path/filepath"
	"fmt"
	"github.com/cihub/seelog"
	"time"
)

// 文件/目录 是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}

// 获取绝对路径
func AbsPath(_path string) string {
	absPath, err := filepath.Abs(_path)
	if err != nil {
		return _path
	}

	return absPath
}

// 创建一个目录
func CreateDir(_path string) error {
	err := os.MkdirAll(_path, os.ModePerm)

	return err
}

/* 检测和创建路径
Params:
    _path: 路径
    _pathUsage: 是做什么用的
 */
func CheckAndCreatePath(_path string, _pathUsage string) error {
	if exists, err := PathExists(_path); err != nil {
		return fmt.Errorf("%v, 检测失败. %v", _pathUsage, err)
	} else {
		if exists {
			seelog.Infof("%v, 已经存在: %s(%s)",
				_pathUsage, _path, AbsPath(_path))
		} else { // 不存在需要创建目录
			seelog.Warnf("%v, 不存在: %s(%s)",
				_pathUsage, _path, AbsPath(_path))

			err1 := CreateDir(_path)
			if err1 != nil {
				return fmt.Errorf("%v失败: %s(%s). %v",
					_pathUsage, _path, AbsPath(_path), err1)
			}
			seelog.Warnf("创建%v成功: %s(%s)",
				_pathUsage, _path, AbsPath(_path))
		}
	}

	return nil
}

// 生成一个全局唯一的文件名
func CreateUUIDFileName() string {
	t := time.Now()
	fileName := fmt.Sprintf("%v", t.Format("20060102150405123456"))

	return fileName
}
