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

	svm := ExtractSeriesValues(sds)
	for k, s := range svm {
		t.Logf("Metric: %s\n", k)
		for _, it := range s {
			t.Logf("%#v\n", it)
		}
	}
}

func TestParseMetricToString(t *testing.T) {
	var tests = []struct {
		input    map[string]interface{}
		expected string
	}{
		{
			input: map[string]interface{}{
				"__name__": "node_cpu",
				"node":     "node1",
				"cpu":      "cpu1",
			},
			expected: `node_cpu{node1,cpu1}`,
		},
		{
			input: map[string]interface{}{
				"node": "node1",
				"cpu":  "cpu1",
			},
			expected: `{node1,cpu1}`,
		},
		{
			input: map[string]interface{}{
				"__name__": "node_cpu",
			},
			expected: `node_cpu{}`,
		},
		{
			input:    map[string]interface{}{},
			expected: `{}`,
		},
	}

	for caseid, test := range tests {
		res := ParseMetricToString(test.input)
		t.Logf("ParseMetricToString(): %s", res)
		if len(res) != len(test.expected) {
			t.Errorf("Case #%d, Actual: %s, Expected: %s", caseid+1, res, test.expected)
		}
	}
}
