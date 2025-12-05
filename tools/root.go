package tools

import "github.com/spf13/cobra"

// RootCmd 是根命令 即sysgo
var RootCmd = &cobra.Command{
	Use:   "sysgo",
	Short: "sysgo is the root command",
}

// init 用来初始化 RootCmd 配置
func init() {
	RootCmd.AddCommand(lsCmd)
	RootCmd.AddCommand(cdCmd)
}
