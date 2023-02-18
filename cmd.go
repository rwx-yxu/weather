package weather

import (
	"errors"
	"fmt"
	"strings"

	Z "github.com/rwxrob/bonzai/z"
	"github.com/rwxrob/help"
	"github.com/rwxrob/vars"
)

func init() {
	Z.Vars.SoftInit()
}

var Cmd = &Z.Cmd{

	Name:      `weather`,
	Summary:   `a command that prints out current day routine schedule`,
	Version:   `v0.4.1`,
	Copyright: `Copyright 2023 Yongle Xu`,
	License:   `Apache-2.0`,
	Site:      `yonglexu.dev`,
	Source:    `git@github.com:rwx-yxu/routine.git`,
	Issues:    `github.com/rwx-yxu/routine/issues`,

	Commands: []*Z.Cmd{
		nowCmd, sitesCmd,
		// standard external branch imports (see rwxrob/{help,conf,vars})
		help.Cmd, vars.Cmd,
	},

	// Add custom BonzaiMark template extensions (or overwrite existing ones).

	Description: `
		{{cmd .Name}} is a tool that queries the Open weather map API for current weather information.
			`,
}

var nowCmd = &Z.Cmd{
	Name:     `now`,
	Summary:  `print current weather conditions to standard output (default)`,
	Commands: []*Z.Cmd{help.Cmd},
	Call: func(x *Z.Cmd, _ ...string) error {
		fmt.Println("Hello world")
		return nil
	},
}

var sitesCmd = &Z.Cmd{
	Name:     `site`,
	Summary:  `site commands that will output met office sites`,
	Commands: []*Z.Cmd{help.Cmd, siteListCmd, siteFindCmd, siteSetCmd, siteForecastCmd},
}

var siteListCmd = &Z.Cmd{
	Name:     `list`,
	Summary:  `prints all met office site locations to standard output (default)`,
	Commands: []*Z.Cmd{help.Cmd},
	Call: func(x *Z.Cmd, _ ...string) error {
		APIKey := Z.Vars.Get(`.apikey`)
		if APIKey == "" {
			return errors.New("API Key not set. Please use the command 'weather var set apikey'")
		}

		sites, err := GetSiteList(APIKey)
		if err != nil {
			return err
		}
		fmt.Println("Available sites:")
		for _, s := range sites {
			fmt.Printf("Name: %s, Region: %s, Area: %s\n", s.Name, s.Region, s.Area)
		}
		return nil
	},
}

var siteFindCmd = &Z.Cmd{
	Name:     `find`,
	Summary:  `filter sites from the sites list`,
	Commands: []*Z.Cmd{help.Cmd},
	Call: func(x *Z.Cmd, args ...string) error {
		APIKey := Z.Vars.Get(`.apikey`)
		if APIKey == "" {
			return errors.New("API Key not set. Please use the command 'weather var set         apikey'")
		}

		sites, err := GetSiteList(APIKey)
		if err != nil {
			return err
		}

		for _, s := range sites {
			out := fmt.Sprintf("Name: %s, Region: %s, Area: %s", s.Name, s.Region, s.Area)
			if strings.Contains(out, args[0]) {
				fmt.Println(out)
			}
		}
		return nil
	},
}

var siteSetCmd = &Z.Cmd{
	Name:     `set`,
	Summary:  `set the met office site location as the default for weather commands`,
	Commands: []*Z.Cmd{help.Cmd},
	Call: func(x *Z.Cmd, args ...string) error {
		APIKey := Z.Vars.Get(`.apikey`)
		if APIKey == "" {
			return errors.New("API Key not set. Please use the command 'weather var set         apikey'")
		}

		sites, err := GetSiteList(APIKey)
		if err != nil {
			return err
		}
		for _, s := range sites {
			if s.Name == args[0] {
				err := Z.Vars.Set("locationID", s.Id)
				if err != nil {
					return err
				}

				regions, err := GetRegionList(APIKey)
				if err != nil {
					return err
				}

				for _, r := range regions {
					if r.Name == s.Region {
						err := Z.Vars.Set("regionID", r.Id)
						if err != nil {
							return err
						}
					}
				}
				return nil
			}
		}
		return errors.New("site is not in list. Please use weather site list to find list all sites")
	},
}

var siteForecastCmd = &Z.Cmd{
	Name:     `forecast`,
	Summary:  `get the met office forecast for today using the set location id.`,
	Commands: []*Z.Cmd{help.Cmd},
	Call: func(x *Z.Cmd, args ...string) error {
		APIKey := Z.Vars.Get(`.apikey`)
		if APIKey == "" {
			return errors.New("API Key not set. Please use the command 'weather var set apikey'")
		}

		siteID := Z.Vars.Get(`locationID`)
		if siteID == "" {
			return errors.New("Site id is not set. Please use the command 'weather site set LOCATION' first. Or, use the command 'weather site find LOCATION' to find a specific location")
		}
		forecast, err := GetTodayForecast(APIKey, siteID)
		if err != nil {
			return err
		}
		fmt.Printf("Day %s°C - %s | Night %s °C - %s\n", forecast.Day.Temp, forecast.Day.Description, forecast.Night.Temp, forecast.Night.Description)

		regionID := Z.Vars.Get(`regionID`)
		if siteID == "" {
			return errors.New("Region id is not set. Please use the command 'weather site set LOCATION' first. Or, use the command 'weather site find LOCATION' to find a specific location")
		}
		regionForecast, err := GetRegionForecast(APIKey, regionID)
		if err != nil {
			return err
		}
		if len(regionForecast) == 0 {
			return errors.New("No forecasts found for region")
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
