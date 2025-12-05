package tools

import (
	"errors"

	"github.com/spf13/cobra"
)

// cdCmd 是 cd 命令
var cdCmd = &cobra.Command{
	Use:   "cd [path]",
	Short: "change current directory",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runCd(args)
	},
}

// runCd 是 cd 命令实际核心执行逻辑
func runCd(args []string) error {
	if len(args) < 1 {
		return errors.New("arg not found")
	}
	if len(args) > 1 {
		return errors.New("need only one arg")
	}
	// TODO cd 实现
	return nil
}
