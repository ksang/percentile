package main

import (
	"fmt"
	"math"
	"strconv"

	"github.com/ksang/percentile/promdata"
)

func calcAverage(svs []promdata.SeriesValue) float64 {
	var sum float64
	for _, v := range svs {
		u, err := strconv.ParseFloat(v.Value.(string), 64)
		if err != nil {
			// ignore the value
			continue
		}
		sum += u
	}
	return sum / float64(len(svs))
}

// assume svs is already sorted
func calcPercentile(svs []promdata.SeriesValue, p int) float64 {
	if len(svs) == 0 {
		return math.NaN()
	}
	n := float64(len(svs))
	// When the percentile lies between two samples,
	// we use a weighted average of the two samples.
	rank := float64(p) / float64(100) * (n - 1)

	lowerIndex := math.Max(0, math.Floor(rank))
	upperIndex := math.Min(n-1, lowerIndex+1)
	lowerValue, err := strconv.ParseFloat(svs[int(lowerIndex)].Value.(string), 64)
	if err != nil {
		fmt.Printf("ERROR: Failed to get lower for percentile value: %s", err)
		return math.NaN()
	}
	upperValue, err := strconv.ParseFloat(svs[int(upperIndex)].Value.(string), 64)
	if err != nil {
		fmt.Printf("ERROR: Failed to get upper for percentile value: %s", err)
		return math.NaN()
	}
	weight := rank - math.Floor(rank)
	return float64(lowerValue)*(1-weight) + float64(upperValue)*weight
}
