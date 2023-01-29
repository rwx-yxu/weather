package weather

import (
	"fmt"
	"net/http"
	"time"
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

func (c Client) FormatURL(resource string, id int) string {
	switch resource {
	case "sitelist":
		return fmt.Sprintf("%s/public/data/val/wxfcs/all/json/sitelist?key=%s", c.BaseURL, c.APIKey)
	case "regionlist":
		return fmt.Sprintf("%s/public/data/txt/wxfcs/regionalforecast/json/sitelist?key=%s", c.BaseURL, c.APIKey)

	}
	return ""
}
