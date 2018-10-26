package server

import (
	"github.com/daiguadaidai/pilipala/pala/config"
	"net/http"
	"github.com/cihub/seelog"
	"sync"
)

func StartHttpServer(_wg *sync.WaitGroup) {
	defer _wg.Done()

	// 获取pala启动配置信息
	palaStartConfig := config.GetPalaStartConfig()

	// http.HandleFunc("/ClientParams", ) // 运行一个命令

	seelog.Infof("Pala监听地址为: %v", palaStartConfig.PalaAddress())
	err := http.ListenAndServe(palaStartConfig.PalaAddress(), nil)
	if err != nil {
		seelog.Errorf("pala启动服务出错: %v", err)
	}
}
