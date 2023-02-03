package weather

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Region struct {
	Id   string
	Name string
}

type RegionListResp struct {
	Locations struct {
		Location []struct {
			Id   string `json:"@id"`
			Name string `json:"@name"`
		} `json:"Location"`
	} `json:"Locations"`
}

func ParseRegionResponse(data []byte) ([]Region, error) {
	var resp RegionListResp
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

func GetRegionList(key string) ([]Region, error) {
	c := NewClient(key)

	regions, err := c.RegionReq()
	if err != nil {
		return []Region{}, err
	}
	return regions, err
}

func (c *Client) RegionReq() ([]Region, error) {
	URL := c.FormatURL(RegionList, "")
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
	regions, err := ParseRegionResponse(data)
	if err != nil {
		return []Region{}, err
	}
	return regions, nil
}
