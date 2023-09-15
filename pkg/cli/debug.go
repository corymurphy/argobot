package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/corymurphy/argobot/pkg/argocd"
)

var debug = &cobra.Command{
	Use:   "debug",
	Short: "debug",
	Long:  `debug`,
	Run: func(cmd *cobra.Command, args []string) {
		// runner = cmd.Command./bin/argocd-linux-amd64 --server https://argocd.local app diff --server-side-generate helloworld{}

		// ./bin/argocd-linux-amd64 --server https://argocd.local app diff --server-side-generate helloworld

		// result, _ := runner.NewShellCommand().Run(
		// 	"./bin/argocd-linux-amd64",
		// 	"--server",
		// 	"https://argocd.local",
		// 	"app",
		// 	"diff",
		// 	"--server-side-generate",
		// 	"helloworld",
		// )

		// fmt.Println(string(result))

		// result, err := runner.CliCommand{}.Run()
		client := argocd.NewCliClient()
		diff, err := client.Diff("helloworld", "")

		fmt.Println(err)
		fmt.Println(diff)
	},
}

func init() {
	rootCmd.AddCommand(debug)
}
