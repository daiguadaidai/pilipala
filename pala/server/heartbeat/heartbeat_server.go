package heartbeat

import (
	"fmt"
	"github.com/cihub/seelog"
	"github.com/daiguadaidai/pilipala/common"
	"github.com/daiguadaidai/pilipala/pala/config"
	"sync"
	"time"
)

func Start(_wg *sync.WaitGroup) {
	defer _wg.Done()

	cf := config.GetPalaStartConfig()
	ticker := time.NewTicker(time.Second * time.Duration(cf.HeartbeatInterval))

	for {
		select {
		case <-ticker.C:
			if err := HeartBeat(cf); err != nil {
				seelog.Errorf("上报心跳失败. %v", err)
			}
		}
	}
}

// 执行heartbeat
func HeartBeat(cf *config.PalaStartConfig) error {
	ip, err := common.GetIntranetIp()
	if err != nil {
		return err
	}

	if err = common.GetUrl(cf.PiliHeartbeatURL(ip), ""); err != nil {
		return fmt.Errorf("%v %s", err, cf.PiliHeartbeatURL(ip))
	}

	return nil
}
