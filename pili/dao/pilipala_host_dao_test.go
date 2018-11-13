package dao

import (
	"fmt"
	"github.com/daiguadaidai/pilipala/common/types"
	"github.com/daiguadaidai/pilipala/pili/model"
	"github.com/liudng/godump"
	"testing"
)

func TestPilipalaHostDao_FindAll(t *testing.T) {
	InitDBConfig()

	host := &model.PilipalaHost{
		IsValid:    types.GetNullInt64(1),
		IsDedicate: types.GetNullInt64(0),
	}

	hostDao := new(PilipalaHostDao)
	newHosts, err := hostDao.FindAll("id, host", host)
	if err != nil {
		t.Errorf("获取host出错, %v", err)
		return
	}

	godump.Dump(newHosts)
}

// 获取最优的host
func TestPilipalaHostDao_GetOptimalHost(t *testing.T) {
	InitDBConfig()

	host := &model.PilipalaHost{
		IsValid:    types.GetNullInt64(1),
		IsDedicate: types.GetNullInt64(0),
	}

	hostDao := new(PilipalaHostDao)
	newHost, err := hostDao.GetOptimalHost("id, host", host, []int64{})
	if err != nil {
		t.Errorf("获取host出错, %v", err)
		return
	}

	fmt.Println(newHost)
}

// 获取最优的host
func TestPilipalaHostDao_GetOptimalHost_Map(t *testing.T) {
	InitDBConfig()

	where := make(map[string]interface{})
	where["is_dedicate"] = 0
	where["is_valid"] = 3
	ids := []int64{1, 2, 3}

	hostDao := new(PilipalaHostDao)
	newHost, err := hostDao.GetOptimalHost("id, host", where, ids)
	if err != nil {
		t.Errorf("获取host出错, %v", err)
		return
	}

	fmt.Println(newHost)
}

// host任务数 自增
func TestPilipalaHostDao_IncrTaskByHost(t *testing.T) {
	InitDBConfig()

	hostDao := new(PilipalaHostDao)
	err := hostDao.IncrTaskByHost("10.10.10.21")
	if err != nil {
		t.Errorf("host 运行任务书自增失败, %v", err)
		return
	}

	fmt.Println("自增成功")
}

// host任务数 自减
func TestPilipalaHostDao_DecrTaskByHost(t *testing.T) {
	InitDBConfig()

	hostDao := new(PilipalaHostDao)
	err := hostDao.DecrTaskByHost("10.10.10.21")
	if err != nil {
		t.Errorf("host 运行任务书自减失败, %v", err)
		return
	}

	fmt.Println("自减成功")
}

// host任务数 自减
func TestPilipalaHostDao_UpdateIsValidByHost(t *testing.T) {
	InitDBConfig()

	hostDao := new(PilipalaHostDao)
	err := hostDao.UpdateIsValidByHost("10.10.10.21", 1)
	if err != nil {
		t.Errorf("host 更新心跳失败, %v", err)
		return
	}

	fmt.Println("host 更新心跳成功")
}
