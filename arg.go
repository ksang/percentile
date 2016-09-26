/*
percentile is a command line tool to fetch x% percentile
sample value from prometheus time series query.
Please read README.md for instructions
*/
package percentile

type Arg struct {
	PromURL    string
	FilePath   string
	Percent    int
	PrintTable bool
}
