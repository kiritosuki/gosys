package tools

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

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
// TODO 子进程不能修改父进程
func runCd(args []string) error {
	if len(args) < 1 {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		err = os.Chdir(homeDir)
		if err != nil {
			return err
		}
		return nil
	}
	if len(args) > 1 {
		return errors.New("need only one arg")
	}
	// cd 实现
	path := args[0]
	var targetDir string
	if path == "~" {
		var err error
		targetDir, err = os.UserHomeDir()
		if err != nil {
			return err
		}
	} else if strings.HasPrefix(path, "~/") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		targetDir = strings.Replace(path, "~/", homeDir+"/", 1)
	} else {
		targetDir = path
	}
	absDir, err := filepath.Abs(filepath.Clean(targetDir))
	if err != nil {
		return err
	}
	err = os.Chdir(absDir)
	if err != nil {
		return err
	}
	return nil
}
