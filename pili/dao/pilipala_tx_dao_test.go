package dao

import (
	"testing"
	"github.com/daiguadaidai/pilipala/pili/model"
	"github.com/daiguadaidai/pilipala/common/types"
)

func TestPilipalaTXDao_CreateCommandAndHost(t *testing.T) {
	InitDBConfig()

	pilipalaCommandProgram := &model.PilipalaCommandProgram{
		Title:        types.GetNullString("事物测试2"),
		FileName:     types.GetNullString("tx02.py"),
		HaveDedicate: types.GetNullInt64(1),
	}

	pilipalaCommandHosts := make([]*model.PilipalaCommandHost, 0, 1)
	pilipalaCommandHost1 := &model.PilipalaCommandHost{
		PilipalaCommandProgramId: types.GetNullInt64(1),
		PilipalaHostId:           types.GetNullInt64(1),
	}
	pilipalaCommandHost2 := &model.PilipalaCommandHost{
		PilipalaCommandProgramId: types.GetNullInt64(1),
		PilipalaHostId:           types.GetNullInt64(2),
	}
	pilipalaCommandHosts = append(pilipalaCommandHosts, pilipalaCommandHost1, pilipalaCommandHost2)

	pilipalaTXDao := new(PilipalaTXDao)
	err := pilipalaTXDao.CreateCommandAndHost(pilipalaCommandProgram, pilipalaCommandHosts)
	if err != nil {
		t.Error(err)
	}

	t.Log("创建成功")
}


func TestPilipalaTXDao_UpdateCommandAndHost(t *testing.T) {
	InitDBConfig()

	// 需要修改的command
	pilipalaCommandProgram := &model.PilipalaCommandProgram{
		Id:           types.GetNullInt64(1),
		Title:        types.GetNullString("update test"),
		FileName:     types.GetNullString("update_tx02.py"),
		HaveDedicate: types.GetNullInt64(1),
	}

	// 需要添加的 command host
	needAddHosts := make([]*model.PilipalaCommandHost, 0, 1)
	needAddHost := &model.PilipalaCommandHost{
		PilipalaCommandProgramId: types.GetNullInt64(1),
		PilipalaHostId:           types.GetNullInt64(3),
	}
	needAddHosts = append(needAddHosts, needAddHost)

	// 需要删除的 command host
	needDeleteHostIds := []int64{1, 2}

	pilipalaTXDao := new(PilipalaTXDao)
	err := pilipalaTXDao.UpdateCommandAndHost(pilipalaCommandProgram, needAddHosts, needDeleteHostIds)
	if err != nil {
		t.Error(err)
		return
	}

	t.Log("更新成功")
}
