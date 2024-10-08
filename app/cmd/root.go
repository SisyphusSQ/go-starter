package cmd

import (
	"github.com/spf13/cobra"
	"os"
)

var (
	Version = "1.0.0"

	rootCmd = &cobra.Command{
		Use:     "go-starter",
		Version: Version,
		Short:   "go-starter Management CLI",
		Run: func(cmd *cobra.Command, args []string) {
			httpCmd.Run(cmd, args)
		},
	}
)

func Execute() {
	initAll()
	if err := rootCmd.Execute(); err != nil {
		println(err)
		os.Exit(1)
	}
}

func initAll() {
	httpCmd.Flags().StringVarP(&configuare, "config", "c", "./config/config.yml", "config file path")
	//fmt.Println(configuare)
	//config.SetConfigFile(configuare)
	//cobra.OnInitialize(config.InitConfig)

	rootCmd.AddCommand(httpCmd)
	rootCmd.AddCommand(versionCmd)
}
