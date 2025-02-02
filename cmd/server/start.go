package server

import (
	"bitbucket.org/nordcloud/mfacli/config"
	"bitbucket.org/nordcloud/mfacli/pkg/vault"

	"github.com/spf13/cobra"
)

func CreateStartCmd(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "start-server",
		Short: "Start a credentials cache server in the background",
		RunE: func(cmd *cobra.Command, args []string) error {
			return vault.StartServer(cfg)
		},
	}
}
