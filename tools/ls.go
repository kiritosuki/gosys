package tools

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// lsOptions 是选项集合结构体
type lsOptions struct {
	lsAll bool
}

// opts 全局选项集合变量
var opts = &lsOptions{}

// lsCmd 是 ls 命令
var lsCmd = &cobra.Command{
	Use:   "ls [paths...]",
	Short: "list files in directory",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runLs(args)
	},
}

// init 用来初始化 lsCmd 配置
func init() {
	lsCmd.Flags().BoolVarP(&opts.lsAll, "all", "a", false, "list all files in directory including hidden files")
}

// runLs 是 ls 命令实际核心执行逻辑
func runLs(args []string) error {
	if len(args) == 0 {
		return singleLs("")
	} else if len(args) == 1 {
		return singleLs(args[0])
	} else {
		for _, arg := range args {
			absPath, err := filepath.Abs(arg)
			if err != nil {
				return err
			}
			fmt.Println(absPath + ": ")
			err = singleLs(arg)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// singleLs 是 ls 后面只跟 一个或零个参数的实现
func singleLs(arg string) error {
	path := "."
	if arg != "" {
		path = arg
	}
	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		name := entry.Name()
		if !opts.lsAll && name[0] == '.' {
			continue
		}
		if entry.IsDir() {
			name += "/"
		}
		fmt.Print(name + "  ")
	}
	fmt.Println()
	return nil
}
