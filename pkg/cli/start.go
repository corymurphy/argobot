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
	LogLevelEnvVar      = "ARGOBOT_LOG_LEVEL"
	DefaultLogLevel     = logging.Info
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

		content, err := os.ReadFile(config.PrivateKeyFilePath)

		if err != nil {
			panic(err)
		}

		config.Github.App.PrivateKey = string(content)
		config.Github.App.WebhookSecret = os.Getenv(WebhookSecretEnvVar)

		logLevel := DefaultLogLevel
		if serializedLevel, exists := os.LookupEnv(LogLevelEnvVar); exists {
			logLevel, err = logging.GetLogLevel(serializedLevel)
			if err != nil {
				panic(err)
			}
		}

		argoClient := &argocd.ApplicationsClient{
			BaseUrl: config.ArgoCdUrl,
		}

		if apiKey, exists := os.LookupEnv("ARGOBOT_ARGOCD_API_KEY"); exists {
			argoClient.Token = apiKey
		}

		server.NewServer(
			config,
			logging.NewLogger(logLevel),
			argoClient,
		).Start()
	},
}

func init() {
	run.Flags().StringVar(&opts.Path, "config-path", "config.yml", "path to the config.yml file")

	rootCmd.AddCommand(run)
}
