package main

import (
	"os"

	"github.com/rwx-yxu/weather/conf"
	"github.com/rwx-yxu/weather/site"
	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "weather",
		Short: "Weather cli",
	}
	rootCmd.AddCommand(site.BaseCmd)
	rootCmd.AddCommand(conf.BaseCmd)
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
