package weather

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	APIKey     string
	BaseURL    string
	HTTPClient *http.Client
}

type Site struct {
	Id     string
	Name   string
	Region string
	Area   string
}

type SiteListResp struct {
	Locations struct {
		Location []struct {
			Id              string
			Name            string
			UnitaryAuthArea string
			Region          string
		}
	}
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
	}
	return ""
}

func ParseResponse(data []byte) ([]Site, error) {
	var resp SiteListResp
	err := json.Unmarshal(data, &resp)
	if err != nil {
		return []Site{}, fmt.Errorf("invalid API response %q: %w", data, err)
	}
	if len(resp.Locations.Location) < 1 {
		return []Site{}, fmt.Errorf("invalid API response %q: want at least one Weather element", data)
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

func GetSiteList(key string) ([]Site, error) {
	c := NewClient(key)

	sites, err := c.SiteReq()
	if err != nil {
		return []Site{}, err
	}
	return sites, err
}

func (c *Client) SiteReq() ([]Site, error) {
	URL := c.FormatURL("sitelist", 0)
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
