package cli

import (
	"github.com/corymurphy/argobot/pkg/server"
	"github.com/spf13/cobra"
)

var (
	opts server.Options = server.Options{}
)

var run = &cobra.Command{
	Use:   "start",
	Short: "runs the server",
	Long:  `runs the server`,
	Run: func(cmd *cobra.Command, args []string) {
		server.Start(opts)
	},
}

func init() {
	run.Flags().StringVar(&opts.Path, "config-path", "config.yml", "path to the config.yml file")

	rootCmd.AddCommand(run)
}
