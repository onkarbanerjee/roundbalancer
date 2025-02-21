package command

import (
	"github.com/onkarbanerjee/roundbalancer/pkg/echo"
	"github.com/spf13/cobra"
)

func configureApplicationServerCommand(command *cobra.Command) {
	applicationServerCommand := &cobra.Command{
		Use:   "application-server",
		Short: "backends application server",
	}
	serverStartCommand := &cobra.Command{
		Use:   "start",
		Short: "start server",
		RunE:  echo.Start,
	}

	serverStartCommand.Flags().IntP("port", "p", 8080, "server port")
	serverStartCommand.Flags().String("id", "", "server id")

	command.AddCommand(applicationServerCommand)
	applicationServerCommand.AddCommand(serverStartCommand)
}
