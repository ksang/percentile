/*
Package Promdata parsing and extracing metadata/data from prometheus HTTP API.
It will sort the data value in groups for easy fetching percentile/max/min/avg
sample value.
*/
package promdata

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

// SeriesDataSet is the full json format prometheus series returns
type SeriesDataSet struct {
	Status string `json:"status"`
	Data   struct {
		ResultType string `json:"resultType"`
		Result     []struct {
			Metric interface{}     `json:"metric"`
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

// ByValue defines sort methods of SeriesValue
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

// ParseJsonData parse raw series data result to SeriesDataSet
func ParseJsonData(data []byte) (*SeriesDataSet, error) {
	var sds SeriesDataSet
	if err := json.Unmarshal(data, &sds); err != nil {
		return nil, err
	}
	return &sds, nil
}

// ExtractSeriesValues extracts data SeriesDataSet and generate sorted slice for value pairs
func ExtractSeriesValues(sds *SeriesDataSet) map[string][]SeriesValue {
	res := make(map[string][]SeriesValue)
	for _, result := range sds.Data.Result {
		var svs []SeriesValue
		metricName := ParseMetricToString(result.Metric)
		for _, record := range result.Values {
			// valid record should be [timestamp, value]
			if len(record) == 2 {
				sv := SeriesValue{
					TimeStamp: record[0].(float64),
					Value:     record[1],
				}
				svs = append(svs, sv)
			}
		}
		sort.Sort(ByValue(svs))
		res[metricName] = svs
	}

	return res
}

// ExtractSVM is a helper function to parse json then extract SeriesValue map
func ExtractSVM(data []byte) (map[string][]SeriesValue, error) {
	sds, err := ParseJsonData(data)
	if err != nil {
		return nil, err
	}
	return ExtractSeriesValues(sds), nil
}

// parse a uknown metric struct to string
func ParseMetricToString(metric interface{}) string {
	var res string
	if name, ok := metric.(map[string]interface{})["__name__"]; ok {
		res += fmt.Sprintf("%s", name)
	}
	res += "{"
	for k, v := range metric.(map[string]interface{}) {
		if k == "__name__" {
			continue
		}
		//res = res + k + "="
		res = res + fmt.Sprintf("%s", v) + ","
	}
	res = strings.TrimSuffix(res, ",")
	res += "}"
	return res
}
