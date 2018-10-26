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
	"github.com/daiguadaidai/pilipala/pala/config"
	"github.com/daiguadaidai/pilipala/pala/service"
)

var palaStartConfig *config.PalaStartConfig

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "pala",
	Short: "运行命令工具",
	Long: `
    监听并获取执行命令的通知, 让后启动一个任务.
    ./pala \\
        listen-host="0.0.0.0" \\
        listen-port=19529 \\
        command-path="./pala_command"
        run-command-info-log-path="./log" \\
        run-command-error-log-path="./log" \\
        pili-host="localhost" \\
        pili-port=19528
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
	rootCmd.PersistentFlags().StringVar(&palaStartConfig.RunCommandInfoLogPath, "run-command-info-log-path",
		config.RUN_COMMAND_INFO_LOG_PATH, "运行命令时, 命令的<正常>输出存放位置")
	rootCmd.PersistentFlags().StringVar(&palaStartConfig.RunCommandInfoLogPath, "run-command-error-log-path",
		config.RUN_COMMAND_ERROR_LOG_PATH, "运行命令时, 命令的<错误>输出存放位")
}
