package tools

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

type LsOptions struct {
	LsAll bool // -a
}

var lsCmd = &cobra.Command{
	Use:   "ls [paths...]",
	Short: "list directory contents",
	RunE: func(cmd *cobra.Command, args []string) error {
		lsAll, _ := cmd.Flags().GetBool("all")
		options := LsOptions{
			LsAll: lsAll,
		}
		return runLs(args, options)
	},
}

func init() {
	lsCmd.Flags().BoolP("all", "a", false, "show all files includes hidden")
}

func runLs(args []string, options LsOptions) error {
	path := "."
	if len(args) > 0 {
		path = args[0]
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		name := entry.Name()
		// 默认不显示.开头的隐藏文件
		if !options.LsAll && name[0] == '.' {
			continue
		}

		// 目录后面加“/”
		if entry.IsDir() {
			name += "/"
		}
		fmt.Print(name + "  ")
	}
	fmt.Println()
	return nil
}
