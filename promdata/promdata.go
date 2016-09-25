package promdata

import (
	"encoding/json"
	"reflect"
	"sort"
	"strconv"
)

// SeriesDataSet is the full json format prometheus series returns
type SeriesDataSet struct {
	Status string `json:"status"`
	Data   struct {
		ResultType string `json:"resultType"`
		Result     []struct {
			Metric struct{}        `json:"metric"`
			Values [][]interface{} `json:"values"`
		} `json:"result"`
	} `json:"data"`
}

// SeriesValue are the actual data returned by prometheus series,
// one record represents a sample.
// To get percentile of all the samples, we need to sort them.
type SeriesValue struct {
	TimeStamp float64
	Value     interface{}
}

type ByValue []SeriesValue

func (s ByValue) Len() int      { return len(s) }
func (s ByValue) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s ByValue) Less(i, j int) bool {
	var (
		v1, v2 float64
		err    error
	)
	t1 := reflect.TypeOf(s[i].Value)
	switch t1.Kind() {
	case reflect.String:
		v1, err = strconv.ParseFloat(s[i].Value.(string), 64)
		if err != nil {
			panic("failed to parse data")
		}
	}

	t2 := reflect.TypeOf(s[j].Value)
	switch t2.Kind() {
	case reflect.String:
		v2, err = strconv.ParseFloat(s[j].Value.(string), 64)
		if err != nil {
			panic("failed to parse data")
		}
	}
	return v1 < v2
}

// Parse raw series data result to SeriesDataSet
func ParseJsonData(data []byte) (*SeriesDataSet, error) {
	var sds SeriesDataSet
	if err := json.Unmarshal(data, &sds); err != nil {
		return nil, err
	}
	return &sds, nil
}

// extract data SeriesDataSet and generate sorted slice for value pairs
func ExtractSeriesValues(sds *SeriesDataSet) []SeriesValue {
	var res []SeriesValue
	for _, result := range sds.Data.Result {
		for _, record := range result.Values {
			// valid record should be [timestamp, value]
			if len(record) == 2 {
				sv := SeriesValue{
					TimeStamp: record[0].(float64),
					Value:     record[1],
				}
				res = append(res, sv)
			}
		}
	}
	sort.Sort(ByValue(res))
	return res
}
