package dao

import (
	"testing"
	"github.com/cihub/seelog"
	"fmt"
)

func TestPilipalaHostDao_FindHostByCommandId(t *testing.T) {
	defer seelog.Flush()
	InitSeelog()
	InitDBConfig()

	pilipalaCommandHostDao := new(PilipalaCommandHostDao)
	columnStr := `id, pilipala_command_program_id, pilipala_host_id`
	pilipalaHosts, err := pilipalaCommandHostDao.FindByCommandId(1, columnStr)
	if err != nil {
		seelog.Errorf("保存命令数据错误. %v", err)
		return
	}

	seelog.Info("获取成功该")
	fmt.Println(pilipalaHosts)
}
