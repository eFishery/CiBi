package cmd

import(
	"os"
	"fmt"

	"github.com/eFishery/CiBi/cmd/servicemap"

	cobra "github.com/spf13/cobra"
)

var (
    buildTime string
    version   string
)

// @see https://cobra.dev/#getting-started
var rootCmd = &cobra.Command{
	Use:     	"servicemap [Cibi.yaml]",
	Short:   	"running service map by default will run cibi.yaml",
	Version: 	version,
	Args:		cobra.MaximumNArgs(1),
	Run:     	servicemap.RunServiceMap,
}

func Execute(){
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}