package main

import (
	"errors"
	"flag"
	"log/slog"
	"thermpro_exporter"
	"time"
)

var (
	start    = flag.String("start", "", "specifies the start datetime in RFC3339 format")
	end      = flag.String("end", "", "specifies the end datetime in RFC3339 format")
	interval = flag.String("interval", "1m", "specifies the interval in Go duration format")
)

func main() {
	flag.Parse()
	var startDate time.Time
	var endDate = time.Now()
	toSkip := 1 * time.Minute

	var allErr []error
	var err error
	if *start != "" {
		startDate, err = time.Parse(time.RFC3339, *start)
		if err != nil {
			allErr = append(allErr, err)
		}
	}
	if *end != "" {
		endDate, err = time.Parse(time.RFC3339, *end)
		if err != nil {
			allErr = append(allErr, err)
		}
	}

	if *interval != "" {
		toSkip, err = time.ParseDuration(*interval)
		if err != nil {
			allErr = append(allErr, err)
		}
	}

	if len(allErr) > 0 {
		slog.Error("bad input", slog.Any("error", errors.Join(allErr...)))
	}
	err = thermpro_exporter.GenerateCSV(startDate, endDate, toSkip)
	if err != nil {
		slog.Error("failed due to error", slog.Any("error", err))
	}
}
