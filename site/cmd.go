package site

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

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
		APIKey := viper.GetString("apikey")
		if APIKey == "" {
			return errors.New("no api key set in .env file. Please run `weather conf init` to initialise conf file")
		}
		sites, err := GetList(APIKey)
		if err != nil {
			return err
		}
		var out strings.Builder
		out.WriteString("Available sites:\n")
		for _, s := range sites {
			out.WriteString(fmt.Sprintf("Name: %s, Region: %s, Area: %s\n", s.Name, s.Region, s.Area))
		}
		lessCmd := exec.Command("less")
		lessCmd.Stdin = strings.NewReader(out.String())
		lessCmd.Stdout = os.Stdout
		lessCmd.Stderr = os.Stderr
		err = lessCmd.Run()
		if err != nil {
			return err
		}

		return nil
	},
}
