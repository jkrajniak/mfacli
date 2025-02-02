package cmd

import (
	"bitbucket.org/nordcloud/mfacli/cmd/add"
	"bitbucket.org/nordcloud/mfacli/cmd/doc"
	"bitbucket.org/nordcloud/mfacli/cmd/generate"
	"bitbucket.org/nordcloud/mfacli/cmd/list"
	"bitbucket.org/nordcloud/mfacli/cmd/remove"
	"bitbucket.org/nordcloud/mfacli/cmd/server"
	"bitbucket.org/nordcloud/mfacli/config"

	"fmt"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	globalCfg   config.Config
	versionFlag bool
	logFile     string
	logLevel    string
)

func initDataDir() (string, error) {
	dir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	dir = filepath.Join(dir, config.DataDirName)
	err = os.MkdirAll(dir, os.FileMode(0755))
	if err != nil {
		return "", err
	}

	return dir, nil
}

func initLogger() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stderr)
	log.SetLevel(log.DebugLevel)
}

func createRootCmd() *cobra.Command {
	return &cobra.Command{
		Use:          "mfacli",
		Short:        "A tool to generate MFA TOTP codes",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if !versionFlag {
				return fmt.Errorf("A subcommand should be provided")
			}
			fmt.Println(config.Version)
			return nil
		},
	}
}

func addFlags(rootCmd *cobra.Command, defaultSocket, defaultVault string) {
	rootCmd.Flags().BoolVarP(&versionFlag, "version", "v", false, "Show version")

	rootCmd.PersistentFlags().StringVarP(&globalCfg.SocketPath, "socket", "S", defaultSocket, "custom Unix socket path to bind (if server) or to connect (if client)")
	rootCmd.PersistentFlags().StringVarP(&globalCfg.VaultPath, "vault", "V", defaultVault, "custom encrypted vault file")
	rootCmd.PersistentFlags().StringVar(&globalCfg.ServerLogFile, config.FlagServerLogFile, "", "Server log file")
	rootCmd.PersistentFlags().BoolVar(&globalCfg.NoCache, "no-cache", false, "don't use vault cache server")
	rootCmd.PersistentFlags().Var(&globalCfg.Password, "password", "vault password in a format accepted by openssl (env:*, file:* or pass:*)")
}

func addSubcommands(rootCmd *cobra.Command) {
	rootCmd.AddCommand(generate.CreatePrintCmd(&globalCfg))
	rootCmd.AddCommand(generate.CreateClipboardCmd(&globalCfg))
	rootCmd.AddCommand(generate.CreateTypeCmd(&globalCfg))
	rootCmd.AddCommand(add.Create(&globalCfg))
	rootCmd.AddCommand(list.Create(&globalCfg))
	rootCmd.AddCommand(remove.Create(&globalCfg))
	rootCmd.AddCommand(server.CreateRunCmd(&globalCfg))
	rootCmd.AddCommand(server.CreateStartCmd(&globalCfg))
	rootCmd.AddCommand(server.CreateStopCmd(&globalCfg))
	doc.Bind(rootCmd)
}

func Execute() error {
	initLogger()

	dataDir, err := initDataDir()
	if err != nil {
		return err
	}

	rootCmd := createRootCmd()

	socketPath := filepath.Join(dataDir, config.DefaultSocketName)
	vaultPath := filepath.Join(dataDir, config.DefaultVaultName)
	addFlags(rootCmd, socketPath, vaultPath)

	addSubcommands(rootCmd)

	return rootCmd.Execute()
}
