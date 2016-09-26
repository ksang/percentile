package main

import (
	"testing"
)

func TestParseTS(t *testing.T) {
	var tests = []struct {
		start    string
		duration string
		expected string
		err      bool
	}{
		{
			start:    "1435781430.781",
			duration: "2s",
			expected: "&start=1435781430&end=1435781432",
			err:      false,
		},
		{
			start:    "1435781430.781",
			duration: "2",
			expected: "&start=1435781430&end=1435781432",
			err:      true,
		},
		{
			start:    "abcd",
			duration: "2s",
			expected: "&start=1435781430.781&end=1435781432.781",
			err:      true,
		},
		{
			start:    "2016-01-02T15:04:05Z",
			duration: "60m",
			expected: "&start=1451747045&end=1451750645",
			err:      false,
		},
		{
			start:    "1996-12-19T16:39:57-08:00",
			duration: "60m",
			expected: "&start=851042397&end=851045997",
			err:      false,
		},
	}

	for caseid, test := range tests {
		res, err := ParseStartEnd(test.start, test.duration)
		t.Logf("Res: %s Error: %s", res, err)
		if err != nil {
			if !test.err {
				t.Errorf("Case #%d, error:%s", caseid+1, err)
			}
			continue
		}
		if res != test.expected {
			t.Errorf("Case #%d, Actual: %s, Expected: %s",
				caseid+1, res, test.expected)
		}
	}
}
