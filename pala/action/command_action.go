package action

import (
	"github.com/daiguadaidai/pilipala/pala/config"
	"fmt"
	"sync"
	"bufio"
	"io"
)

type CommandAction struct {
	CommandName string
	ParamStr string
	TaskUUID string
}

// 运行命令
func (this *CommandAction) Run() {
	wg := new(sync.WaitGroup)


	wg.Wait()


}

/* 输出错误日志到文件
Params:
    _wg: bin
 */
func (this *CommandAction) OutputError(_wg *sync.WaitGroup, _stdout io.ReadCloser) {
	defer _wg.Done()

	outputBuf := bufio.NewReader(_stdout)

	for {
		output, _, err := outputBuf.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				fmt.Printf("Error: 获取输出错误. %v\n", err)
				return
			}
		}

		fmt.Printf("info ------------- %s\n", output)
	}
}

func (this *CommandAction) CommandPath() string {
	palaStartConfig := config.GetPalaStartConfig()

	return fmt.Sprintf("%s/%s", palaStartConfig.CommandPath, this.CommandName)
}
