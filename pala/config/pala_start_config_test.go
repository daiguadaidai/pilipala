package config

import (
	"testing"
	"fmt"
)

func TestPalaStartConfig_CheckAndStore(t *testing.T) {
	palaStartConfig := NewPalaStartConfig()

	fmt.Println("ListenHost:", palaStartConfig.ListenHost)
	fmt.Println("ListenPort:", palaStartConfig.ListenPort)
	fmt.Println("PiliHost:", palaStartConfig.PiliHost)
	fmt.Println("PiliPort:", palaStartConfig.PiliPort)
	fmt.Println("CommandPath:", palaStartConfig.CommandPath)
	fmt.Println("RunCommandInfoLogPath:", palaStartConfig.RunCommandInfoLogPath)
	fmt.Println("RunCommandErrorLogPath:", palaStartConfig.RunCommandErrorLogPath)

	palaStartConfig.CheckAndStore()
}
