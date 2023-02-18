package weather_test

import (
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/rwx-yxu/weather"
)

func TestParseSiteResponse(t *testing.T) {
	t.Parallel()
	data, err := os.ReadFile("testdata/locations.json")
	if err != nil {
		t.Fatal(err)
	}
	var want = []weather.Site{
		weather.Site{
			Id:     "14",
			Name:   "Carlisle Airport",
			Region: "nw",
			Area:   "Cumbria",
		},
		weather.Site{
			Id:     "26",
			Name:   "Liverpool John Lennon Airport",
			Region: "nw",
			Area:   "Merseyside",
		},
	}
	got, err := weather.ParseSiteResponse(data)
	if err != nil {
		t.Fatal(err)
	}
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestParseRegionResponse(t *testing.T) {
	t.Parallel()
	data, err := os.ReadFile("testdata/regions.json")
	if err != nil {
		t.Fatal(err)
	}
	var want = []weather.Region{
		weather.Region{
			Id:   "500",
			Name: "os",
		},
		weather.Region{
			Id:   "505",
			Name: "dg",
		},
	}
	got, err := weather.ParseRegionResponse(data)
	if err != nil {
		t.Fatal(err)
	}
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestParseSiteForecastResponse(t *testing.T) {
	t.Parallel()
	data, err := os.ReadFile("testdata/siteDailyForecast.json")
	if err != nil {
		t.Fatal(err)
	}
	var want = weather.SiteForecast{
		weather.Day{
			Temp:        "9",
			Description: "Overcast",
		},
		weather.Night{
			Temp:        "6",
			Description: "Overcast",
		},
	}
	got, err := weather.ParseForecastResponse(data)
	if err != nil {
		t.Fatal(err)
	}
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestParseRegionForecastResponse(t *testing.T) {
	t.Parallel()
	data, err := os.ReadFile("testdata/regionForecast.json")
	if err != nil {
		t.Fatal(err)
	}
	var want = []weather.RegionForecast{
		weather.RegionForecast{
			Title:   "Headline:",
			Content: "Patchy fog clearing, then dry and bright for most.",
		},
		weather.RegionForecast{
			Title:   "Today:",
			Content: "Mostly cloudy in Cumbria at first with the odd spot of light rain. Elsewhere, patchy mist and fog clearing to leave a dry day with bright or sunny spells. Turning breezier than of late. Maximum Temperature 9C.",
		},
		weather.RegionForecast{
			Title:   "Tonight:",
			Content: "Cloud thickening this evening with a band of patchy rain moving south across the region. Skies clearing overnight, allowing a patchy frost to develop. Winds becoming generally light. Minimum Temperature 1C.",
		},
		weather.RegionForecast{
			Title:   "Thursday:",
			Content: "After a frosty start in places, it will be a dry and fine day with sunny periods. Brisk westerly winds, especially in northern areas. Maximum Temperature 8C.",
		},
	}
	got, err := weather.ParseRegionForecastResponse(data)
	if err != nil {
		t.Fatal(err)
	}

	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}
