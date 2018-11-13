package dao

import (
	"testing"
	"github.com/daiguadaidai/pilipala/pili/model"
	"github.com/daiguadaidai/pilipala/common"
	"github.com/daiguadaidai/pilipala/common/types"
	"fmt"
)

func TestPilipalaRunTaskDao_Create(t *testing.T) {
	InitDBConfig()

	taskDao := new(PilipalaRunTaskDao)
	task := &model.PilipalaRunTask{
		TaskUUID: types.GetNullString(common.GetUUID()),
		PilipalaCommandProgramId: types.GetNullInt64(1),
		Host: types.GetNullString("10.10.10.55"),
		FileName: types.GetNullString("test.py"),
		Params: types.GetNullString("--task-uuid="),
		Pid: types.GetNullInt64(0),
	}
	if err := taskDao.Create(task); err!= nil {
		t.Errorf("创建任务失败, %v", err)
	}

	fmt.Println("创建成功")
}

func TestPilipalaRunTaskDao_UpdateTaskStatus(t *testing.T) {
	InitDBConfig()

	uuid := "201811051149441151149446"
	taskDao := new(PilipalaRunTaskDao)
	err := taskDao.UpdateTaskStatus(uuid, model.TASK_STATUS_FAIL)
	if err != nil {
		t.Fatalf("更新状态失败. %v", err)
	}

	fmt.Println("更新成功")
}
