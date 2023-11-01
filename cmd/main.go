package cmd

import(
	"os"
	"fmt"

	"github.com/eFishery/CiBi/cmd/servicemap"
	"github.com/eFishery/CiBi/cmd/notionlog"

	cobra "github.com/spf13/cobra"
	viper "github.com/spf13/viper"
)

var (
    buildTime string
    version   string
)

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(serviceMapCmd)
	rootCmd.AddCommand(logCmd)
	logCmd.AddCommand(logDeployCmd)
	completion.Hidden = true
	rootCmd.AddCommand(completion)

	rootCmd.PersistentFlags().StringP("file", "f", "cibi.yaml", "Define the file path of different file name for cibi.yml")

	viper.BindPFlag("file", rootCmd.PersistentFlags().Lookup("file"))
}

// @see https://cobra.dev/#getting-started
var rootCmd = &cobra.Command{
	Use:     	"cibi",
	Short:   	"ea ea ea ea",
	Version: 	version,
}

var completion = &cobra.Command{
	Use:		"completion",
	Short:		"idk, I just want to hide this",
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Hugo",
	Long:  `All software has versions. This is Hugo's`,
	Run: func(cmd *cobra.Command, args []string) {
	  fmt.Println("SUGOI BANG")
	},
  }

var serviceMapCmd = &cobra.Command{
	Use:     	"servicemap",
	Short:   	"set a list of service map",
	Version: 	version,
	Args:		cobra.MaximumNArgs(1),
	Run:     	servicemap.RunServiceMap,
}

var logCmd = &cobra.Command{
	Use:     	"log",
	Short:   	"Insert Deployment Log to notion database",
	Version: 	version,
}

var logDeployCmd = &cobra.Command{
	Use:     	"deployed",
	Short:   	"Mark as Deployed",
	Version: 	version,
	Run:     	notionlog.RunDeployedLog,
}

func Execute(){
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}