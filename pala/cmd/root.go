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

	"github.com/daiguadaidai/pilipala/pala/config"
	"github.com/daiguadaidai/pilipala/pala/service"
	"github.com/spf13/cobra"
)

var palaStartConfig *config.PalaStartConfig

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "pala",
	Short: "运行命令工具",
	Long: `
    监听并获取执行命令的通知, 让后启动一个任务.
    ./pala \
        --listen-host="0.0.0.0" \
        --listen-port=19529 \
        --command-path="./pala_command" \
        --run-command-log-path="./log" \
        --run-command-paraller=8 \
        --is-log-dir-prefix-date=true \
        --heartbeat-interval=60 \
        --pili-host="localhost" \
        --pili-port=19528
`,
	Run: func(cmd *cobra.Command, args []string) {
		service.Start(palaStartConfig)
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
	palaStartConfig = new(config.PalaStartConfig)

	rootCmd.PersistentFlags().StringVar(&palaStartConfig.ListenHost, "listen-host",
		config.LISTEN_HOST, "启动Http服务监听host")
	rootCmd.PersistentFlags().IntVar(&palaStartConfig.ListenPort, "listen-port",
		config.LISTEN_PORT, "启动Http服务监听port")
	rootCmd.PersistentFlags().StringVar(&palaStartConfig.CommandPath, "command-path",
		config.COMMAND_PATH, "命令存放位置")
	rootCmd.PersistentFlags().StringVar(&palaStartConfig.RunCommandLogPath, "run-command-log-path",
		config.RUN_COMMAND_LOG_PATH, "命令输出信息存放位置")
	rootCmd.PersistentFlags().IntVar(&palaStartConfig.RunConnamdParaller, "run-command-paraller",
		config.RUN_COMMAND_PARALLER, "运行命令的并发数")
	rootCmd.PersistentFlags().BoolVar(&palaStartConfig.IsLogDirPrefixDate, "is-log-dir-prefix-date",
		config.IS_LOG_DIR_PREFIX_DATE, "日志目录是否需要使用日期作为上级目录")
	rootCmd.PersistentFlags().IntVar(&palaStartConfig.HeartbeatInterval, "heartbeat-interval",
		config.HEARTBEAT_INTERVAL, "心跳检测间隔时间")
	rootCmd.PersistentFlags().StringVar(&palaStartConfig.PiliHost, "pili-host",
		config.PILI_HOST, "调度器(pili)服务Host")
	rootCmd.PersistentFlags().IntVar(&palaStartConfig.PiliPort, "pili-port",
		config.PILI_PORT, "调度器(pili)服务的Port")
}
