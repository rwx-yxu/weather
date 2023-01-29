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
