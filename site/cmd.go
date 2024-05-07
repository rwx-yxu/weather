package site

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/rwx-yxu/weather/region"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	BaseCmd.AddCommand(listCmd)
	BaseCmd.AddCommand(setCmd)
	BaseCmd.AddCommand(forecastCmd)
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

var setCmd = &cobra.Command{
	Use:          `set`,
	Short:        `set the default region and location`,
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
		var selectedSite int
		siteNames := make([]string, len(sites))
		for i, s := range sites {
			siteNames[i] = fmt.Sprintf("%s - %s", s.Name, s.Region)
		}
		selectedSiteOptions := make([]huh.Option[int], len(sites))
		for i, s := range sites {
			selectedSiteOptions[i] = huh.NewOption(fmt.Sprintf("%s - %s", s.Name, s.Area), i)
		}

		form := huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[int]().
					Options(selectedSiteOptions...).
					Title("Choose your location").
					Description("Please select your default location from the provided list below").
					Value(&selectedSite),
			),
		)

		err = form.Run()

		if err != nil {
			return err
		}
		var regionID string
		regions, err := region.GetList(APIKey)
		if err != nil {
			return err
		}

		for _, r := range regions {
			if r.Name == sites[selectedSite].Region {
				regionID = r.Id
			}
		}

		viper.Set("LocationID", sites[selectedSite].Id)
		viper.Set("RegionID", regionID)
		err = viper.WriteConfig()
		if err != nil {
			return err
		}
		fmt.Println("Site set")
		return nil
	},
}

var forecastCmd = &cobra.Command{
	Use:          `forecast`,
	Short:        `returns the current days weather forecast`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		APIKey := viper.GetString("apikey")
		if APIKey == "" {
			return errors.New("no api key set in .env file. Please run `weather conf init` to initialise conf file")
		}

		siteID := viper.GetString("locationID")
		if siteID == "" {
			return errors.New("site id is not set. Please use the command 'weather site set' to find a specific location")
		}
		forecast, err := GetTodayForecast(APIKey, siteID)
		if err != nil {
			return err
		}
		fmt.Printf("Day %s°C - %s | Night %s °C - %s\n", forecast.Day.Temp, forecast.Day.Description, forecast.Night.Temp, forecast.Night.Description)
		regionID := viper.GetString(`regionID`)
		if regionID == "" {
			return errors.New("region id is not set. Please use the command 'weather site set LOCATION' first. Or, use the command 'weather site find LOCATION' to find a specific location")
		}

		regionForecast, err := region.GetForecast(APIKey, regionID)
		if err != nil {
			return err
		}
		if len(regionForecast) == 0 {
			return errors.New("no forecasts found for region")
		}

		for _, f := range regionForecast {
			if f.Title != "Headline:" {
				fmt.Printf("%v\n\n", f.Content)
				continue
			}
			fmt.Printf("%v %v\n", f.Title, f.Content)
		}
		return nil
	},
}
