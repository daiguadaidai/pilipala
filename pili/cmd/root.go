// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/daiguadaidai/pilipala/pili/config"
	"github.com/daiguadaidai/pilipala/pili/service"
)

var piliStartConfig *config.PiliStartConfig
var dbConfig *config.DBConfig

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "pili",
	Short: "任务管理的接口管理工具",
	Long: `主要用于任务的管理, 分发, 通知等等
./pili \
    --listen-host=0.0.0.0 \
    --listen-port=19528 \
    --command-path=./pili_commands \
    --upload-command-path=./pili_upload_commands \
    --db-host=127.0.0.1 \
    --db-port=3306 \
    --db-schema=boom \
    --db-username=root \
    --db-password=root \
    --db-auto-commit=true \
    --db-charset=utf8mb4 \
    --db-timeout=10 \
    --db-max-idel-conns=100 \
    --db-max-open-conns=100
`,
	Run: func(cmd *cobra.Command, args []string) {
		service.Start(piliStartConfig, dbConfig)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	piliStartConfig = new(config.PiliStartConfig)
	dbConfig = new(config.DBConfig)

	// http 配置
	rootCmd.PersistentFlags().StringVar(&piliStartConfig.ListenHost, "listen-host",
		config.LISTEN_HOST, "Http服务监听的IP")
	rootCmd.PersistentFlags().IntVar(&piliStartConfig.ListenPort, "listen-port",
		config.LISTEN_PORT, "Http服务监听的端口")
	rootCmd.PersistentFlags().StringVar(&piliStartConfig.CommandPath, "command-path",
		config.COMMAND_PATH, "上传命令(最终)存放的路径")
	rootCmd.PersistentFlags().StringVar(&piliStartConfig.UploadCommandPath, "upload-command-path",
		config.UPLOAD_COMMAND_PATH, "上传命令(临时)存放的路径")

	// 链接的数据库配置
	rootCmd.PersistentFlags().StringVar(&dbConfig.Host, "db-host",
		config.DB_HOST, "数据库host")
	rootCmd.PersistentFlags().IntVar(&dbConfig.Port, "db-port",
		config.DB_PORT, "数据库port")
	rootCmd.PersistentFlags().StringVar(&dbConfig.Username, "db-username",
		config.DB_USERNAME, "数据库用户名")
	rootCmd.PersistentFlags().StringVar(&dbConfig.Password, "db-password",
		config.DB_PASSWORD, "数据库密码")
	rootCmd.PersistentFlags().StringVar(&dbConfig.Database, "db-schema",
		config.DB_SCHEMA, "数据库名称")
	rootCmd.PersistentFlags().StringVar(&dbConfig.CharSet, "db-charset",
		config.DB_CHARSET, "数据库字符集")
	rootCmd.PersistentFlags().IntVar(&dbConfig.Timeout, "db-timeout",
		config.DB_TIMEOUT, "数据库timeout")
	rootCmd.PersistentFlags().IntVar(&dbConfig.MaxIdelConns, "db-max-idel-conns",
		config.DB_MAX_IDEL_CONNS, "数据库最大空闲连接数")
	rootCmd.PersistentFlags().IntVar(&dbConfig.MaxOpenConns, "db-max-open-conns",
		config.DB_MAX_OPEN_CONNS, "数据库最大连接数")
	rootCmd.PersistentFlags().BoolVar(&dbConfig.AutoCommit, "db-auto-commit",
		config.DB_AUTO_COMMIT, fmt.Sprintf("数据库自动提交 (default:%v)", config.DB_AUTO_COMMIT))
}
