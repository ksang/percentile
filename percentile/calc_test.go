package main

import (
	"math"
	"testing"

	"github.com/ksang/percentile/promdata"
)

func TestCalc(t *testing.T) {
	var tests = []struct {
		svs []promdata.SeriesValue
		p   int
		// expected[avg, percentile]
		expected []float64
	}{
		{
			svs: []promdata.SeriesValue{
				promdata.SeriesValue{
					Value: "0",
				},
				promdata.SeriesValue{
					Value: "50",
				},
				promdata.SeriesValue{
					Value: "100",
				},
			},
			p:        95,
			expected: []float64{50, 95},
		},
		{
			svs:      []promdata.SeriesValue{},
			p:        95,
			expected: []float64{0, 0},
		},
		{
			svs: []promdata.SeriesValue{
				promdata.SeriesValue{
					Value: "abc",
				},
			},
			p:        95,
			expected: []float64{0, 0},
		},
	}

	for caseid, test := range tests {
		avg := calcAverage(test.svs)
		percentile := calcPercentile(test.svs, test.p)
		t.Logf("avg: %f percentile: %f", avg, percentile)
		if math.Abs(avg-test.expected[0]) > 0.0001 || math.Abs(percentile-test.expected[1]) > 0.0001 {
			t.Errorf("Case #%d, Actual: %f/%f, Expected: %f/%f",
				caseid+1, avg, percentile,
				test.expected[0], test.expected[1])
		}
	}
}
