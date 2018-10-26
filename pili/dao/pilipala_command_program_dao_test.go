package dao

import (
	"testing"
	"github.com/cihub/seelog"
	"github.com/liudng/godump"
	"github.com/daiguadaidai/pilipala/pili/config"
	"github.com/daiguadaidai/pilipala/pili/model"
	"database/sql"
	"github.com/daiguadaidai/pilipala/common/sql_type_util"
)

func InitDBConfig() {
	dbConfig := &config.DBConfig{
		Host: "10.10.10.21",
		Port: 3307,
		Username: "HH",
		Password: "oracle12",
		Database: "boom",
		CharSet: "utf8mb4",
		AutoCommit: true,
		MaxOpenConns: 100,
		MaxIdelConns: 100,
		Timeout: 10,
	}

	config.SetDBConfig(dbConfig)
}

func InitSeelog() {
	logger, _ := seelog.LoggerFromConfigAsBytes([]byte(config.LogDefautConfig()))
	seelog.ReplaceLogger(logger)
}

func TestPilipalaCommandProgramDao_Create(t *testing.T) {
	defer seelog.Flush()
	InitSeelog()
	InitDBConfig()

	pilipalaCommandProgram := &model.PilipalaCommandProgram {
		Title: sql_type_util.NullString{sql.NullString{"测试1", true}},
		FileName: sql_type_util.NullString{sql.NullString{"test1.py", true}},
		Params: sql_type_util.NullString{sql.NullString{"--host=0.0.0.0 --port=3306", true}},
		HaveDedicate: sql_type_util.NullInt64{sql.NullInt64{1, true}},
	}

	pilipalaCommandProgramDao := new(PilipalaCommandProgramDao)
	err := pilipalaCommandProgramDao.Create(pilipalaCommandProgram)
	if err != nil {
		seelog.Errorf("保存命令数据错误. %v", err)
		return
	}

	seelog.Info("创建成功")
}

func TestPilipalaCommandProgramDao_GetByID(t *testing.T) {
	defer seelog.Flush()
	InitSeelog()
	InitDBConfig()

	pilipalaCommandProgramDao := new(PilipalaCommandProgramDao)

	id := 1
	columnStr := "*"

	pilipalaCommandProgram, err := pilipalaCommandProgramDao.GetByID(id, columnStr)
	if err != nil {
		seelog.Errorf("从出具库中获取命令错误, %v", err)
		return
	}

	if pilipalaCommandProgram == nil {
		seelog.Warnf("没有获取到数据")
	} else {
		godump.Dump(pilipalaCommandProgram)
	}
}

func TestPilipalaCommandProgramDao_PaginationFind(t *testing.T) {
	defer seelog.Flush()
	InitSeelog()
	InitDBConfig()

	pilipalaCommandProgramDao := new(PilipalaCommandProgramDao)

	min_pk := 3
	offset := 50
	columnStr := "*"

	datas, err := pilipalaCommandProgramDao.PaginationFind(min_pk, offset, columnStr)
	if err != nil {
		seelog.Errorf("分页获取命令出错, %v", err)
		return
	}

	for _, pilipalaCommandProgram := range datas {
		seelog.Info(pilipalaCommandProgram)
	}
}


