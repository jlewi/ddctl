package cmd

import (
	"github.com/spf13/cobra"
)

func NewLogsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "logs",
	}
	cmd.AddCommand(NewQueryToURL())
	return cmd
}
