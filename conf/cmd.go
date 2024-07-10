package conf

import (
	"errors"
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	SoftInit()
	err := Load()
	if err != nil {
		log.Fatalf("Unable to load conf file %s", err)
	}
	BaseCmd.AddCommand(initCmd)
	BaseCmd.AddCommand(dataCmd)
}

// sayCmd represents the say command
var BaseCmd = &cobra.Command{
	Use:   "conf",
	Short: "top level conf command to initialise the .env configuration file for the weather application",
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Reinitialise the configuration .env file in the user cache directory",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Attempting to create env file...")
		if err := Init(); err != nil && !errors.Is(err, ErrorConfigFileExists) {
			return err
		}
		fmt.Printf(".env initialised at: %s\n", filePath())
		return nil
	},
}

var dataCmd = &cobra.Command{
	Use:   "data",
	Short: fmt.Sprintf("Prints out all of the stored env variables at: %s", filePath()),
	RunE: func(cmd *cobra.Command, args []string) error {
		keys, err := keys()
		if err != nil {
			return err
		}
		for _, key := range keys {
			fmt.Printf("%s=%s\n", key, viper.Get(key))
		}
		return nil
	},
}
