package config

import (
	"github.com/daiguadaidai/pilipala/common"
	"fmt"
)

const (
	LISTEN_HOST = "0.0.0.0"
	LISTEN_PORT = 19529

	COMMAND_PATH = "./pala_commands"
	RUN_COMMAND_INFO_LOG_PATH = "./log"
	RUN_COMMAND_ERROR_LOG_PATH = "./log"

	PILI_HOST = "localhost"
	PILI_PORT = 19528
)

var palaStartConfig *PalaStartConfig

type PalaStartConfig struct {
	ListenHost string // 启动服务绑定的IP
	ListenPort int    // 启动服务绑定的端口

	CommandPath string // 命令存放的路径
	RunCommandInfoLogPath string // 运行命令接收日志的输出位置
	RunCommandErrorLogPath string // 运行命令接收错误日志的输出位置

	PiliHost string // 需要访问pili的host
	PiliPort int    // 需要访问pili的port
}

func NewPalaStartConfig() *PalaStartConfig {
	palaStartConfig := new(PalaStartConfig)

	palaStartConfig.ListenHost = LISTEN_HOST
	palaStartConfig.ListenPort = LISTEN_PORT

	palaStartConfig.CommandPath = COMMAND_PATH
	palaStartConfig.RunCommandInfoLogPath = RUN_COMMAND_INFO_LOG_PATH
	palaStartConfig.RunCommandErrorLogPath = RUN_COMMAND_ERROR_LOG_PATH

	palaStartConfig.PiliHost = PILI_HOST
	palaStartConfig.PiliPort = PILI_PORT

	return palaStartConfig
}

// 设置 palaStartconfig
func SetPalaStartConfig(_palaStartConfig *PalaStartConfig) {
	palaStartConfig = _palaStartConfig
}

func GetPalaStartConfig() *PalaStartConfig {
	return palaStartConfig
}

// 检测配置信息, 初始化一些需要的东西
func (this *PalaStartConfig) CheckAndStore() error {

	// 检测和创建命令存放目录
	if err := common.CheckAndCreatePath(this.CommandPath,
		"命令存放目录"); err != nil {
		return err
	}

	// 检测和创建执行命令时的输出目录
	if err := common.CheckAndCreatePath(this.RunCommandInfoLogPath,
		"被执行命令的(正常)输出目录"); err != nil {
		return err
	}

	if err := common.CheckAndCreatePath(this.RunCommandErrorLogPath,
		"被执行命令的(错误)输出目录"); err != nil {
		return err
	}

	return nil
}

// 获取pili监听地址
func (this *PalaStartConfig) PiliAddress() string {
	return fmt.Sprintf("%v:%v", this.PiliHost, this.PiliPort)
}

// 获取pala监听地址
func (this *PalaStartConfig) PalaAddress() string {
	return fmt.Sprintf("%v:%v", this.ListenHost, this.ListenPort)
}
