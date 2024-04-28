package site

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	BaseCmd.AddCommand(listCmd)
}

// sayCmd represents the say command
var BaseCmd = &cobra.Command{
	Use:   "site",
	Short: "top level site command to list find and set met office sites",
}

var listCmd = &cobra.Command{
	Use:          `list`,
	Short:        `prints all met office site locations to standard output`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		APIKey := viper.GetString("APIKey")
		if APIKey == "" {
			return errors.New("no api key set in .env file. Please run `weather conf init` to initialise conf file")
		}
		sites, err := GetList(APIKey)
		if err != nil {
			return err
		}
		fmt.Println("Available sites:")
		for _, s := range sites {
			fmt.Printf("Name: %s, Region: %s, Area: %s\n", s.Name, s.Name, s.Region, s.Area)
		}
		return nil
	},
}
