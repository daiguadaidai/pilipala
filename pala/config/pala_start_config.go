package config

import (
	"fmt"
	"github.com/cihub/seelog"
	"github.com/daiguadaidai/pilipala/common"
)

const (
	LISTEN_HOST = "0.0.0.0"
	LISTEN_PORT = 19529

	COMMAND_PATH           = "./pala_commands"
	RUN_COMMAND_LOG_PATH   = "./log"
	RUN_COMMAND_PARALLER   = 8
	IS_LOG_DIR_PREFIX_DATE = true
	HEARTBEAT_INTERVAL     = 60

	PILI_HOST = "localhost"
	PILI_PORT = 19528

	PILI_DOWNLOAD_COMMAND_URL = "http://%s:%d/pili/command_programs/download/%s"
	PILI_TASK_SUCCESS_URL     = "http://%s:%d/pili/tasks/success/%s"
	PILI_TASK_FAIL_URL        = "http://%s:%d/pili/tasks/fail/%s"
	PILI_HEARTBEAT_URL        = "http://%s:%d/pili/hosts/heartbeat/%s"
	PILI_TASK_UPDATE_URL      = "http://%s:%d/pili/tasks"
)

var palaStartConfig *PalaStartConfig

type PalaStartConfig struct {
	ListenHost string // 启动服务绑定的IP
	ListenPort int    // 启动服务绑定的端口

	CommandPath        string // 命令存放的路径
	RunCommandLogPath  string // 运行命令接收日志的输出位置
	RunConnamdParaller int    // 运行命令并发数
	IsLogDirPrefixDate bool   // 日志的目录是否需要使用日期切割
	HeartbeatInterval  int    // 心跳检测间隔时间

	PiliHost string // 需要访问pili的host
	PiliPort int    // 需要访问pili的port
}

func NewPalaStartConfig() *PalaStartConfig {
	palaStartConfig := new(PalaStartConfig)

	palaStartConfig.ListenHost = LISTEN_HOST
	palaStartConfig.ListenPort = LISTEN_PORT

	palaStartConfig.CommandPath = COMMAND_PATH
	palaStartConfig.RunCommandLogPath = RUN_COMMAND_LOG_PATH
	palaStartConfig.RunConnamdParaller = RUN_COMMAND_PARALLER
	palaStartConfig.IsLogDirPrefixDate = IS_LOG_DIR_PREFIX_DATE
	palaStartConfig.HeartbeatInterval = HEARTBEAT_INTERVAL

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

	if err := common.CheckAndCreatePath(this.RunCommandLogPath,
		"被执行命令的(错误)日志目录"); err != nil {
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

// 获取命令相对路径
func (this *PalaStartConfig) GetCommandFilePath(_fileName string) string {
	return fmt.Sprintf("%s/%s", this.CommandPath, _fileName)
}

func (this *PalaStartConfig) GetLogDir() string {
	if !this.IsLogDirPrefixDate {
		return this.RunCommandLogPath
	}

	dir := fmt.Sprintf("%v/%v", this.RunCommandLogPath, common.GetDateStr())
	if err := common.CheckAndCreatePath(dir, "被执行命令的(正常)日志目录"); err != nil {
		seelog.Warnf("创建不了输出日志目录. dir: %s. 使用默认的目录: %s. %v",
			dir, this.RunCommandLogPath, err)
		return this.RunCommandLogPath
	}
	return dir
}

// 获取日志路径
func (this *PalaStartConfig) GetLogPath(_taskUUID string) string {
	return fmt.Sprintf("%s/%s.log", this.GetLogDir(), _taskUUID)
}

// 获取pili下载命令url
func (this *PalaStartConfig) PiliDownloadCommandUrl(_command string) string {
	return fmt.Sprintf(PILI_DOWNLOAD_COMMAND_URL, this.PiliHost, this.PiliPort, _command)
}

// 获取pili任务成功url
func (this *PalaStartConfig) PiliTaskSuccessUrl(_taskUUID string) string {
	return fmt.Sprintf(PILI_TASK_SUCCESS_URL, this.PiliHost, this.PiliPort, _taskUUID)
}

// 获取pili任务失败url
func (this *PalaStartConfig) PiliTaskFailUrl(_taskUUID string) string {
	return fmt.Sprintf(PILI_TASK_FAIL_URL, this.PiliHost, this.PiliPort, _taskUUID)
}

func (this *PalaStartConfig) PiliHeartbeatURL(_host string) string {
	return fmt.Sprintf(PILI_HEARTBEAT_URL, this.PiliHost, this.PiliPort, _host)
}

func (this *PalaStartConfig) PiliTaskUpdateURL() string {
	return fmt.Sprintf(PILI_TASK_UPDATE_URL, this.PiliHost, this.PiliPort)
}
