package main

import (
	"errors"
	"flag"
	"fmt"
	"net/url"
	"strings"

	"github.com/ksang/percentile"
)

var (
	full      string
	table     bool
	percent   int
	starttime string
	duration  string
	step      string
	query     string
	promhost  string
	verbose   bool
)

func init() {
	flag.StringVar(&full, "f", "", "full URI to prometheus series api or result file output")
	flag.BoolVar(&table, "t", false, "print table result, including max/avg/min")
	flag.BoolVar(&verbose, "v", false, "verbose output")
	flag.IntVar(&percent, "p", -1, "percentile to get, a int number in range 0-100")
	flag.StringVar(&starttime, "start", "", "start timestamp, int RFC3339 or unix timestamp")
	flag.StringVar(&duration, "duration", "", "duration, e.g 2s, 3m, 4h.")
	flag.StringVar(&step, "step", "30s", "step, e.g 10s, 15s, 1m")
	flag.StringVar(&query, "query", "", "query used for query_range, without time series")
	flag.StringVar(&promhost, "prom", "", "promethues host url, e.g http://1.1.1.1:9090")
}

func main() {
	flag.Parse()
	arg := &percentile.Arg{}
	if err := getArg(arg); err != nil {
		fmt.Println("ERROR: failed to parse arguments:", err)
		return
	}
	GenResult(arg)
}

// GetArg parse and check command line arguments
func getArg(arg *percentile.Arg) error {
	if percent < 0 || percent > 100 {
		return errors.New("percent value is not valid")
	}
	if len(full) == 0 && len(promhost) == 0 {
		return errors.New("no data source provided, need prometheus info or data file")
	}
	if len(full) > 0 {
		_, err := url.ParseRequestURI(full)
		if err != nil {
			arg.FilePath = full
		} else {
			arg.PromURL = full
		}
	}
	arg.PrintTable = table
	arg.Percent = percent

	if len(full) == 0 {
		if err := genFullURI(arg); err != nil {
			return err
		}
	}
	return nil
}

func genFullURI(arg *percentile.Arg) error {
	res := strings.TrimSuffix(promhost, "/") + "/api/v1/query_range?query="
	res += url.QueryEscape(query)
	se, err := ParseStartEnd(starttime, duration)
	if err != nil {
		return err
	}
	arg.PromURL = res + se + "&step=" + step
	return nil
}
