package weather_test

import (
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/rwx-yxu/weather"
)

func TestParseResponse(t *testing.T) {
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
	got, err := weather.ParseResponse(data)
	if err != nil {
		t.Fatal(err)
	}
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}
