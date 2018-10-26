package config

import (
	"github.com/daiguadaidai/pilipala/common"
	"fmt"
)

const (
	LISTEN_HOST = "0.0.0.0"
	LISTEN_PORT = 19528

	COMMAND_PATH = "./pili_commands"
	UPLOAD_COMMAND_PATH = "./pili_upload_commands"
)

var piliStartConfig *PiliStartConfig

type PiliStartConfig struct {
	ListenHost string // 启动服务绑定的IP
	ListenPort int    // 启动服务绑定的端口

	CommandPath string // 命令存放的路径
	UploadCommandPath string // 上传命令临时使用目录
}

func NewPalaStartConfig() *PiliStartConfig {
	piliStartConfig := new(PiliStartConfig)

	piliStartConfig.ListenHost = LISTEN_HOST
	piliStartConfig.ListenPort = LISTEN_PORT

	piliStartConfig.CommandPath = COMMAND_PATH
	piliStartConfig.UploadCommandPath = UPLOAD_COMMAND_PATH

	return piliStartConfig
}

// 设置 piliStartconfig
func SetPiliStartConfig(_piliStartConfig *PiliStartConfig) {
	piliStartConfig = _piliStartConfig
}

func GetPiliStartConfig() *PiliStartConfig {
	return piliStartConfig
}

// 检测配置信息, 初始化一些需要的东西
func (this *PiliStartConfig) CheckAndStore() error {

	// 检测和创建命令存放目录(临时)
	if err := common.CheckAndCreatePath(this.UploadCommandPath,
		"(临时)命令存放目录"); err != nil {
		return err
	}

	// 检测和创建命令存放目录(最终)
	if err := common.CheckAndCreatePath(this.CommandPath,
		"(最终)命令存放目录"); err != nil {
		return err
	}

	return nil
}

// 获取pala监听地址
func (this *PiliStartConfig) PiliAddress() string {
	return fmt.Sprintf("%v:%v", this.ListenHost, this.ListenPort)
}

// 上传文件临时存放路径
func (this *PiliStartConfig) UploadCommandFilePath(_fileName string) string {
	return fmt.Sprintf("%v/%v", this.UploadCommandPath, _fileName)
}
