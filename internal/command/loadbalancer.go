package command

import (
	"github.com/onkarbanerjee/roundbalancer/pkg/dispatcher"
	"github.com/spf13/cobra"
)

func configureLoadBalancingServerCommand(command *cobra.Command) {
	loadbalancerCommand := &cobra.Command{
		Use:   "load-balancing-server",
		Short: "load balancing server",
	}
	serverStartCommand := &cobra.Command{
		Use:   "start",
		Short: "start server",
		RunE:  dispatcher.Start,
	}
	command.AddCommand(loadbalancerCommand)
	loadbalancerCommand.AddCommand(serverStartCommand)
}
