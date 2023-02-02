package weather

import (
	"fmt"
	"net/http"
	"time"
)

type Resource string

const (
	SiteList            Resource = "sitelist"
	RegionList          Resource = "regionlist"
	ForecastTodaySite   Resource = "forecasttodaysite"
	ForecastTodayRegion Resource = "forecasttodayregion"
)

type Client struct {
	APIKey     string
	BaseURL    string
	HTTPClient *http.Client
}

func NewClient(apiKey string) *Client {
	return &Client{
		APIKey:  apiKey,
		BaseURL: "http://datapoint.metoffice.gov.uk",
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c Client) FormatURL(res Resource, id int) string {
	switch res {
	case SiteList:
		return fmt.Sprintf("%s/public/data/val/wxfcs/all/json/sitelist?key=%s", c.BaseURL, c.APIKey)
	case RegionList:
		return fmt.Sprintf("%s/public/data/txt/wxfcs/regionalforecast/json/sitelist?key=%s", c.BaseURL, c.APIKey)
	case ForecastTodaySite:
		return fmt.Sprintf("%s/public/data/txt/wxfcs/all/json/%v?key=%s&res=daily", c.BaseURL, id, c.APIKey)
	case ForecastTodayRegion:
		return fmt.Sprintf("%s/public/data/txt/wxfcs/regionalforecast/json/%v?key=%s", c.BaseURL, id, c.APIKey)
	}
	return ""
}
