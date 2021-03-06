package service

import (
	"github.com/daiguadaidai/pilipala/pili/config"
	"github.com/cihub/seelog"
	"sync"
	"github.com/daiguadaidai/pilipala/pili/server"
)

func Start(_piliStartConfig *config.PiliStartConfig, _dbConfig *config.DBConfig) {
	defer seelog.Flush()
	logger, _ := seelog.LoggerFromConfigAsBytes([]byte(config.LogDefautConfig()))
	seelog.ReplaceLogger(logger)

	// 检测和创建指定和需要的目录
	err := _piliStartConfig.CheckAndStore()
	if err != nil {
		seelog.Errorf("检测启动配置文件错误: %v", err)
		return
	}
	err = _dbConfig.Check()
	if err != nil {
		seelog.Errorf("检测链接数据库配置错误: %v", err)
		return
	}

	config.SetPiliStartConfig(_piliStartConfig) // 设置全局的http配置文件
	config.SetDBConfig(_dbConfig) // 设置全局的数据库配置文件

	wg := new(sync.WaitGroup)

	// 启动palaserver
	wg.Add(1)
	go server.StartHttpServer(wg)

	wg.Wait()
}
