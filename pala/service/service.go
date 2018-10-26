package service

import (
	"github.com/daiguadaidai/pilipala/pala/config"
	"github.com/cihub/seelog"
	"sync"
	"github.com/daiguadaidai/pilipala/pala/server"
)

func Start(_palaStartConfig *config.PalaStartConfig) {
	defer seelog.Flush()
	logger, _ := seelog.LoggerFromConfigAsBytes([]byte(config.LogDefautConfig()))
	seelog.ReplaceLogger(logger)

	// 检测和创建指定和需要的目录
	err := _palaStartConfig.CheckAndStore()
	if err != nil {
		seelog.Errorf("检测启动配置文件错误: %v", err)
		return
	}

	config.SetPalaStartConfig(_palaStartConfig)

	wg := new(sync.WaitGroup)

	// 启动palaserver
	wg.Add(1)
	go server.StartHttpServer(wg)

	wg.Wait()
}
