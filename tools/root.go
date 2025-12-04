package tools

import "github.com/spf13/cobra"

var RootCmd = &cobra.Command{
	Use:   "sysgo",
	Short: "sysgo --- a Linux command implement with go, command prefix",
}

func init() {
	RootCmd.AddCommand(lsCmd)
}
