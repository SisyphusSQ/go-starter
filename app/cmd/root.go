package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/SisyphusSQ/go-starter/v2/vars"
)

var (
	rootCmd = &cobra.Command{
		Use:     "go-starter",
		Version: vars.AppVersion,
		Short:   "go-starter Management CLI",
		RunE: func(cmd *cobra.Command, args []string) error {
			return httpCmd.RunE(cmd, args)
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
	httpCmd.Flags().StringVarP(&configure, "config", "c", "./config/config.yml", "config file path")
	rootCmd.AddCommand(httpCmd)
	rootCmd.AddCommand(versionCmd)
}
