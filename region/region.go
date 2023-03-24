package region

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/rwx-yxu/weather/app"
)

type Region struct {
	Id   string
	Name string
}

type ListResp struct {
	Locations struct {
		Location []struct {
			Id   string `json:"@id"`
			Name string `json:"@name"`
		} `json:"Location"`
	} `json:"Locations"`
}

type RegionalFcst struct {
	CreatedOn   string      `json:"createdOn"`
	IssuedAt    string      `json:"issuedAt"`
	RegionId    string      `json:"regionId"`
	FcstPeriods FcstPeriods `json:"FcstPeriods"`
}

type FcstPeriods struct {
	Period []Period `json:"Period"`
}

type Period struct {
	Id        string      `json:"id"`
	Paragraph interface{} `json:"Paragraph"`
}

type Paragraph struct {
	Title string `json:"title"`
	Text  string `json:"$"`
}

type RegionForecast struct {
	Title   string
	Content string
}

func ParseResponse(data []byte) ([]Region, error) {
	var resp ListResp
	err := json.Unmarshal(data, &resp)
	if err != nil {
		return []Region{}, fmt.Errorf("invalid API response %q: %w", data, err)
	}
	if len(resp.Locations.Location) < 1 {
		return []Region{}, fmt.Errorf("invalid API response %q: want at least one Region element", resp)
	}
	var regions = []Region{}
	for _, reg := range resp.Locations.Location {
		region := Region{
			Id:   reg.Id,
			Name: reg.Name,
		}

		regions = append(regions, region)
	}

	return regions, nil
}

func ParseRegionForecastResponse(data []byte) ([]RegionForecast, error) {
	var forecasts = []RegionForecast{}
	var result map[string]RegionalFcst
	err := json.Unmarshal(data, &result)
	if err != nil {
		return []RegionForecast{}, err
	}
	periods := result["RegionalFcst"].FcstPeriods.Period
	for _, period := range periods {
		if paragraphs, ok := period.Paragraph.([]any); ok {
			for _, p := range paragraphs {
				paragraph := p.(map[string]any)
				f := RegionForecast{
					Title:   paragraph["title"].(string),
					Content: paragraph["$"].(string),
				}
				forecasts = append(forecasts, f)
			}
		}
	}
	return forecasts, nil
}

func GetList(key string) ([]Region, error) {
	c := app.NewClient(key)

	regions, err := RegionReq(c)
	if err != nil {
		return []Region{}, err
	}
	return regions, err
}

func RegionReq(c *app.Client) ([]Region, error) {
	URL := c.FormatURL(app.RegionList, "")
	resp, err := c.HTTPClient.Get(URL)
	if err != nil {
		return []Region{}, fmt.Errorf("region GET URL error:%v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return []Region{}, fmt.Errorf("unexpected response status %q", resp.Status)
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return []Region{}, err
	}
	regions, err := ParseResponse(data)
	if err != nil {
		return []Region{}, err
	}
	return regions, nil
}

func GetForecast(key, id string) ([]RegionForecast, error) {
	c := app.NewClient(key)

	rf, err := ForecastReq(id, c)
	if err != nil {
		return []RegionForecast{}, err
	}
	return rf, err
}

func ForecastReq(id string, c *app.Client) ([]RegionForecast, error) {
	URL := c.FormatURL(app.ForecastTodayRegion, id)
	resp, err := c.HTTPClient.Get(URL)
	if err != nil {
		return []RegionForecast{}, fmt.Errorf("region GET URL error:%v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return []RegionForecast{}, fmt.Errorf("unexpected response status %q", resp.Status)
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return []RegionForecast{}, err
	}
	forecasts, err := ParseRegionForecastResponse(data)
	if err != nil {
		return []RegionForecast{}, err
	}
	return forecasts, nil
}
