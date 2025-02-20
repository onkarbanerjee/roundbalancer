package command

import (
	"github.com/spf13/cobra"
)

func Execute() {
	cmd := &cobra.Command{
		Use:   "cli",
		Short: "CLI for running tasks",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	configureApplicationServerCommand(cmd)
	configureLoadBalancingServerCommand(cmd)

	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}
