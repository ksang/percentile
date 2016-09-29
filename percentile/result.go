package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/ksang/percentile"
	"github.com/ksang/percentile/gcd"
	"github.com/ksang/percentile/promdata"

	"github.com/olekukonko/tablewriter"
)

// Result represents a percentile result for a label set
type Result struct {
	Labels     string
	Enough     bool
	Percentile float64
	Max        float64
	Min        float64
	Average    float64
}

// GenResult is the portal of running command arguments
func GenResult(arg *percentile.Arg) {
	data, err := GetSeriesOutput(arg)
	if err != nil {
		fmt.Printf("Failed to get data, error: %s\n", err)
		return
	}
	svm, err := promdata.ExtractSVM(data)
	if err != nil {
		fmt.Printf("Failed to extract data, error: %s\n", err)
		return
	}
	var rs []*Result
	for labels, svs := range svm {
		if verbose {
			fmt.Printf("Labels: %s got %d value", labels, len(svs))
		}
		if len(svs) < 100/int(gcd.GCD(uint64(arg.Percent), 100)) {
			fmt.Printf("WARN: %s only has %d sample(s), may not enough to get good percentile value.\n",
				labels, len(svs))
		}
		r := genResult(svs, labels, arg.PrintTable, arg.Percent)
		rs = append(rs, r)
	}
	printResult(rs, arg)
}

func printResult(rs []*Result, arg *percentile.Arg) {
	if arg.PrintTable {
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{fmt.Sprintf("%d%%", arg.Percent),
			"Max", "Min", "Avg", "Labels"})
		for _, r := range rs {
			table.Append([]string{
				strconv.FormatFloat(r.Percentile, 'f', 2, 64),
				strconv.FormatFloat(r.Max, 'f', 2, 64),
				strconv.FormatFloat(r.Min, 'f', 2, 64),
				strconv.FormatFloat(r.Average, 'f', 2, 64),
				r.Labels})
		}
		table.Render()
	} else {
		for _, r := range rs {
			fmt.Printf("%f\t%s\n", r.Percentile, r.Labels)
		}
	}
}

func genResult(svs []promdata.SeriesValue, labels string, table bool, p int) *Result {
	res := &Result{
		Labels: labels,
	}
	res.Enough = checkEnough(len(svs), p)
	if len(svs) == 0 {
		return res
	}
	res.Percentile = calcPercentile(svs, p)
	if table {
		var err error
		res.Max, err = strconv.ParseFloat(svs[len(svs)-1].Value.(string), 64)
		if err != nil {
			fmt.Printf("ERROR: Failed to get Max: %s", err)
		}
		res.Min, err = strconv.ParseFloat(svs[0].Value.(string), 64)
		if err != nil {
			fmt.Printf("ERROR: Failed to get Min: %s", err)
		}
		res.Average = calcAverage(svs)
	}
	return res
}

func checkEnough(count int, p int) bool {
	return true
}

// GetSeriesOutput get rawdata from prometheus http api or file
func GetSeriesOutput(arg *percentile.Arg) ([]byte, error) {
	if len(arg.PromURL) == 0 {
		f, err := os.Open(arg.FilePath)
		if err != nil {
			return nil, err
		}
		data, err := ioutil.ReadAll(f)
		if err != nil {
			return nil, err
		}
		return data, nil
	}

	client := &http.Client{}
	if verbose {
		fmt.Printf("Requesting: %s\n", arg.PromURL)
	}
	resp, err := client.Get(arg.PromURL)
	if err != nil || resp == nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}
