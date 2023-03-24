package site

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/rwx-yxu/weather/app"
)

type Site struct {
	Id     string
	Name   string
	Region string
	Area   string
}
type Forecast struct {
	Day
	Night
}

type Day struct {
	Temp        string
	Description string
}

type Night struct {
	Temp        string
	Description string
}

type ForecastResp struct {
	SiteRep struct {
		DV struct {
			Location struct {
				Name   string
				Period []struct {
					Rep []struct {
						DayTemp   string `json:"FDm"`
						NightTemp string `json:"FNm"`
						Period    string `json:"$"`
						Weather   string `json:"W"`
					}
				}
			}
		} `json:"DV"`
	} `json:"SiteRep"`
}

type ListResp struct {
	Locations struct {
		Location []struct {
			Id              string
			Name            string
			UnitaryAuthArea string
			Region          string
		}
	}
}

func ParseResponse(data []byte) ([]Site, error) {
	var resp ListResp
	err := json.Unmarshal(data, &resp)
	if err != nil {
		return []Site{}, fmt.Errorf("invalid API response %q: %w", data, err)
	}
	if len(resp.Locations.Location) < 1 {
		return []Site{}, fmt.Errorf("invalid API response %q: want at least one Weather       element", data)
	}
	var sites = []Site{}
	for _, loc := range resp.Locations.Location {
		site := Site{
			Id:     loc.Id,
			Name:   loc.Name,
			Region: loc.Region,
			Area:   loc.UnitaryAuthArea,
		}

		sites = append(sites, site)
	}

	return sites, nil
}

func ParseForecastResponse(data []byte) (Forecast, error) {
	var resp ForecastResp
	err := json.Unmarshal(data, &resp)
	if err != nil {
		return Forecast{}, fmt.Errorf("invalid API response %q: %w", data, err)
	}
	if len(resp.SiteRep.DV.Location.Period) < 1 {
		return Forecast{}, fmt.Errorf("invalid API response %s: want at least one forecast period", data)
	}
	var forecast = Forecast{}
	for _, f := range resp.SiteRep.DV.Location.Period[0].Rep {
		if f.Period == "Day" {
			forecast.Day.Temp = f.DayTemp
			forecast.Day.Description = ValueOfWeather(f.Weather)
		} else {
			forecast.Night.Temp = f.NightTemp
			forecast.Night.Description = ValueOfWeather(f.Weather)
		}
	}

	return forecast, nil
}

func ValueOfWeather(val string) string {
	switch val {
	case "N/A":
		return "Not available"
	case "0":
		return "Clear night"
	case "1":
		return "Sunny day"
	case "2":
		fallthrough
	case "3":
		return "Partly cloudy"
	case "4":
		return "Not used"
	case "5":
		return "Mist"
	case "6":
		return "Fog"
	case "7":
		return "Cloudy"
	case "8":
		return "Overcast"
	case "9":
		fallthrough
	case "10":
		return "Light rain shower"
	case "11":
		return "Drizzle"
	case "12":
		return "Light rain"
	case "13":
		fallthrough
	case "14":
		return "Heavy rain shower"
	case "15":
		return "Heavy rain"
	case "16":
		fallthrough
	case "17":
		return "Sleet shower"
	case "18":
		return "Sleet"
	case "19":
		fallthrough
	case "20":
		return "Hail shower"
	case "21":
		return "Hail"
	case "22":
		fallthrough
	case "23":
		return "Light snow shower"
	case "24":
		return "Light snow"
	case "25":
		fallthrough
	case "26":
		return "Heavy snow shower"
	case "27":
		return "Heavy snow"
	case "28":
		fallthrough
	case "29":
		return "Thunder snow"
	case "30":
		return "Thunder"
	default:
		return "N/A"
	}
}

func GetList(key string) ([]Site, error) {
	c := app.NewClient(key)

	sites, err := SiteReq(c)
	if err != nil {
		return []Site{}, err
	}
	return sites, err
}

func SiteReq(c *app.Client) ([]Site, error) {
	URL := c.FormatURL(app.SiteList, "")
	resp, err := c.HTTPClient.Get(URL)
	if err != nil {
		return []Site{}, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return []Site{}, fmt.Errorf("unexpected response status %q", resp.Status)
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return []Site{}, err
	}
	sites, err := ParseResponse(data)
	if err != nil {
		return []Site{}, err
	}
	return sites, nil
}

func GetTodayForecast(key, siteID string) (Forecast, error) {
	c := app.NewClient(key)
	forecast, err := ForecastReq(siteID, c)
	if err != nil {
		return Forecast{}, err
	}
	return forecast, err
}

func ForecastReq(id string, c *app.Client) (Forecast, error) {
	URL := c.FormatURL(app.ForecastTodaySite, id)
	resp, error := c.HTTPClient.Get(URL)
	if error != nil {
		return Forecast{}, fmt.Errorf("unable to make forecast request: %w", error)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Forecast{}, fmt.Errorf("unexpected response status %q", resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return Forecast{}, err
	}

	forecast, err := ParseForecastResponse(data)
	if err != nil {
		return Forecast{}, err
	}

	return forecast, nil
}
