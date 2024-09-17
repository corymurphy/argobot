package cli

import (
	"os"

	"github.com/corymurphy/argobot/pkg/argocd"
	"github.com/corymurphy/argobot/pkg/env"
	"github.com/corymurphy/argobot/pkg/logging"
	"github.com/corymurphy/argobot/pkg/server"
	"github.com/spf13/cobra"
)

const (
	WebhookSecretEnvVar = "ARGOBOT_GH_WEBHOOK_SECRET"
)

var (
	opts server.Options = server.Options{}
)

var run = &cobra.Command{
	Use:   "start",
	Short: "runs the server",
	Long:  `runs the server`,
	Run: func(cmd *cobra.Command, args []string) {

		config, err := env.ReadConfig(opts.Path)
		if err != nil {
			panic(err)
		}

		content, err := os.ReadFile(config.AppConfig.PrivateKeyFilePath)

		if err != nil {
			panic(err)
		}

		config.Github.App.PrivateKey = string(content)
		config.Github.App.WebhookSecret = os.Getenv(WebhookSecretEnvVar)

		server.NewServer(config, logging.NewLogger(logging.Info), argocd.NewCliClient(config.ArgoCliConfig)).Start()
	},
}

func init() {
	run.Flags().StringVar(&opts.Path, "config-path", "config.yml", "path to the config.yml file")

	rootCmd.AddCommand(run)
}
