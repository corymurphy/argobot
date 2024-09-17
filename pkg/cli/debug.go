package cli

import (
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
)

var debug = &cobra.Command{
	Use:   "debug",
	Short: "debugs the server",
	Long:  `debugs the server`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("debug")
		http.HandleFunc("/", HelloServer)
		http.ListenAndServe(":8080", nil)
	},
}

func init() {
	rootCmd.AddCommand(debug)
}

func HelloServer(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
}
