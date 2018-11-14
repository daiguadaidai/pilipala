package task

import (
	"bufio"
	"fmt"
	"github.com/cihub/seelog"
	"github.com/daiguadaidai/pilipala/common"
	"github.com/daiguadaidai/pilipala/pala/config"
	"io"
	"os"
	"os/exec"
	"sync"
	"syscall"
)

var TaskChan chan *Task
var RunningTaskMap sync.Map

func init() {
	TaskChan = make(chan *Task, 10000)
}

// 启动任务执行服务
func StartTaskServer(_wg *sync.WaitGroup) {
	defer _wg.Done()

	palaStartConfig := config.GetPalaStartConfig()

	seelog.Info("启动任务服务")
	for i := 0; i < palaStartConfig.RunConnamdParaller; i++ {
		seelog.Infof("启动任务服务 %d 号 启动成功", i)
		_wg.Add(1)
		go startTask(_wg)
	}
}

// 启动task
func startTask(_wg *sync.WaitGroup) {
	defer _wg.Done()

	for task := range TaskChan {
		go task.Start()
	}
}

type Task struct {
	Command  string `form:"command" json:"command" binding:"required"`
	TaskUUID string `form:"task_uuid" json:"task_uuid" binding:"required"`
	Params   string `form:"params" json:"params"`
	Config   *config.PalaStartConfig
	LogPath  string
	LogFile  *os.File
}

// 开始一个任务task
func (t *Task) Start() {
	t.Config = config.GetPalaStartConfig()
	t.InitLogFile()
	defer t.LogFile.Close()

	if err := t.UpdateLogPath(); err != nil {
		seelog.Error(err.Error())
	}

	// 检测命令是否存在
	if err := t.CheckAndDownloadCommand(); err != nil {
		t.TaskRunFail(err)
		return
	}

	// 检测命令是否有执行权限
	if err := t.CheckCommandIsExecutable(); err != nil {
		t.TaskRunFail(err)
		return
	}

	// 运行
	if err := t.Run(); err != nil {
		t.TaskRunFail(err)
		return
	}

	t.TaskRunSuccess()
}

// 运行命令
func (t *Task) Run() error {
	wg := new(sync.WaitGroup) // 再次创建一个并发控制器. 只提供运行命令中使用

	// 获取命令路经
	cmdPath := t.Config.GetCommandFilePath(t.Command)
	args, err := common.GetArgs(t.Params)
	if err != nil {
		return fmt.Errorf("解析命令参数出错. params: %s. %v", t.Params, err)
	}
	// 创建命令执行器
	cmd := exec.Command(cmdPath, args...)
	// 主进程退出子进程也退出
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("获取输出管道错误. task uuid: %s. %s %s. %v",
			t.TaskUUID, cmdPath, t.Params, err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("获取<错误>输出管道错误. task uuid: %s. %s %s. %v",
			t.TaskUUID, cmdPath, t.Params, err)
	}

	// 开始执行命令
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("启动命令出错. task uuid: %s. %s %s. %v",
			t.TaskUUID, cmdPath, t.Params, err)
	}
	seelog.Infof("命令执行中. pid: %v. task uuid: %s. %s %s",
		cmd.Process.Pid, t.TaskUUID, cmdPath, t.Params)

	// 保存执行命令 pid
	RunningTaskMap.Store(t.TaskUUID, cmd.Process.Pid)
	defer func() {
		RunningTaskMap.Delete(t.TaskUUID)
	}()

	// 记录命令的输出
	wg.Add(1)
	go t.LogOutput(wg, stdout)
	wg.Add(1)
	go t.LogOutput(wg, stderr)

	wg.Wait()
	// 等待结束
	if err := cmd.Wait(); err != nil {
		seelog.Errorf("Wait err: %v", err)
	}

	// 执行失败
	if !cmd.ProcessState.Success() {
		return fmt.Errorf("命令执行失败. pid: %v. task uuid: %s. %s %s",
			cmd.Process.Pid, t.TaskUUID, cmdPath, t.Params)
	}

	return nil
}

// 出入info日志
func (t *Task) LogOutput(_wg *sync.WaitGroup, _stdout io.ReadCloser) {
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
		t.Log(string(output))
	}
}

// 检测和下载命令
func (t *Task) CheckAndDownloadCommand() error {
	commandPath := fmt.Sprintf("%s/%s", t.Config.CommandPath, t.Command)
	exists, err := common.PathExists(commandPath)
	if err != nil {
		return err
	}
	if !exists {
		seelog.Warnf("命令不存在: %s", t.Command)
		if err1 := common.DownloadFile(t.Config.PiliDownloadCommandUrl(t.Command), commandPath); err1 != nil {
			return fmt.Errorf("%v. %s", err1, t.Config.PiliDownloadCommandUrl(t.Command))
		}
		seelog.Warnf("命令下载成功: %s", t.Command)
	}

	return nil
}

func (t *Task) CheckCommandIsExecutable() error {
	commandPath := fmt.Sprintf("%s/%s", t.Config.CommandPath, t.Command)
	executable, err := common.FileIsExecutable(commandPath)
	if err != nil {
		return err
	}
	if !executable {
		seelog.Warnf("命令不可执行: %s", t.Command)
		if err1 := common.ChmodFile(commandPath); err1 != nil {
			return err1
		}
		seelog.Warnf("命令可执行权限设置成功: %s", t.Command)
	}

	return nil
}

// 通知任务执行成功
func (t *Task) TaskRunSuccess() {
	if err := common.GetUrl(t.Config.PiliTaskSuccessUrl(t.TaskUUID), ""); err != nil {
		seelog.Errorf("通知失败<任务完成>. UUID: %s, command: %s, params: %s. %v",
			t.TaskUUID, t.Command, t.Params, err)
		return
	}

	seelog.Infof("通知成功<任务完成>. UUID: %s, command: %s, params: %s",
		t.TaskUUID, t.Command, t.Params)
}

// 通知任务执行失败
func (t *Task) TaskRunFail(_err error) {
	seelog.Errorf("%v", _err)

	if err := common.GetUrl(t.Config.PiliTaskFailUrl(t.TaskUUID), ""); err != nil {
		seelog.Errorf("通知失败<任务执行失败>. UUID: %s, command: %s, params: %s. %v",
			t.TaskUUID, t.Command, t.Params, err)
		return
	}

	seelog.Infof("通知成功<任务执行失败>. UUID: %s, command: %s, params: %s",
		t.TaskUUID, t.Command, t.Params)
}

// 初始化日志文件
func (t *Task) InitLogFile() {
	t.LogPath = t.Config.GetLogPath(t.TaskUUID)
	seelog.Infof("任务: %s. 命令: %s. 输出文件: %s", t.TaskUUID, t.Command, t.LogPath)

	var err error
	t.LogFile, err = os.Create(t.LogPath)
	if err != nil {
		seelog.Errorf("创建错误日志文件失败. task uuid: %s. logfile: %s",
			t.TaskUUID, t.LogPath)
		return
	}
}

// 记录输出信息
func (t *Task) Log(info string) {
	if _, err := fmt.Fprintln(t.LogFile, info); err != nil {
		seelog.Errorf("写入自建日志出错. %v", err)
	}
}

func (t *Task) UpdateLogPath() error {
	data := make(map[string]string)
	data["task_uuid"] = t.TaskUUID
	data["log_path"] = t.LogPath

	if err := common.PutUrl(t.Config.PiliTaskUpdateURL(), data); err != nil {
		return fmt.Errorf("更新任务日志地址出错: %v", err)
	}

	return nil
}
