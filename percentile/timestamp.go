package main

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

func ParseStartEnd(start string, duration string) (string, error) {
	s, err := time.Parse(time.RFC3339, start)
	if err != nil {
		var (
			sec int64
		)
		parts := strings.Split(start, ".")
		sec, err = strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			return "", errors.New("unsupported time format")
		}
		// nanosec is not necessary
		s = time.Unix(sec, 0)
	}
	d, err := time.ParseDuration(duration)
	if err != nil {
		return "", err
	}
	e := s.Add(d)
	res := "&start=" + strconv.FormatInt(s.Unix(), 10) + "&end=" + strconv.FormatInt(e.Unix(), 10)
	return res, nil
}
