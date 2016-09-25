package promdata

import (
	"flag"
	"io/ioutil"
	"os"
	"testing"
)

var (
	dataFile string
)

func init() {
	flag.StringVar(&dataFile, "data", "test_data.txt", "test data file")
}

func TestParseJsonData(t *testing.T) {
	flag.Parse()
	f, err := os.Open(dataFile)
	if err != nil {
		t.Fatalf("Failed to open test data file: %s", err)
	}
	data, err := ioutil.ReadAll(f)
	if err != nil {
		t.Fatalf("Failed to read test data file: %s", err)
	}
	sds, err := ParseJsonData(data)
	if err != nil {
		t.Errorf("Failed to parse data: %s", err)
	}
	t.Logf("%#v", sds)

	sv := ExtractSeriesValues(sds)
	for _, s := range sv {
		t.Logf("%#v\n", s)
	}

}
