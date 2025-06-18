package macros

import (
	"errors"
	"fmt"
	"math"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data/sqlutil"

	"github.com/grafana/sqlds/v2"
)

var ErrInvalidVariableFallback = errors.New("fallback should contain at least one character")

type timeQueryType string

const (
	timeQueryTypeFrom timeQueryType = "from"
	timeQueryTypeTo   timeQueryType = "to"
)

func newTimeFilter(queryType timeQueryType, query *sqlds.Query) (string, error) {
	date := query.TimeRange.From
	if queryType == timeQueryTypeTo {
		date = query.TimeRange.To
	}
	micros := date.UTC().UnixMicro()
	return fmt.Sprintf("CAST(%d AS TIMESTAMP)", micros), nil
}

// FromTimestampFilter return time filter query based on grafana's timepicker's from time in microseconds
func FromTimestampFilter(query *sqlds.Query, args []string) (string, error) {
	return newTimeFilter(timeQueryTypeFrom, query)
}

// ToTimestampFilter return time filter query based on grafana's timepicker's to time in microseconds
func ToTimestampFilter(query *sqlds.Query, args []string) (string, error) {
	return newTimeFilter(timeQueryTypeTo, query)
}

func TimestampFilter(query *sqlds.Query, args []string) (string, error) {
	if len(args) != 1 {
		return "", backend.DownstreamError(fmt.Errorf("%w: expected 1 argument, received %d", sqlutil.ErrorBadArgumentCount, len(args)))
	}
	var (
		column = args[0]
		from   = query.TimeRange.From.UTC().UnixMicro()
		to     = query.TimeRange.To.UTC().UnixMicro()
	)
	return fmt.Sprintf("%s >= CAST(%d AS TIMESTAMP) AND %s <=  CAST(%d AS TIMESTAMP)", column, from, column, to), nil
}

func DateFilter(query *sqlds.Query, args []string) (string, error) {
	if len(args) != 1 {
		return "", backend.DownstreamError(fmt.Errorf("%w: expected 1 argument, received %d", sqlutil.ErrorBadArgumentCount, len(args)))
	}
	var (
		column = args[0]
		from   = query.TimeRange.From.UTC().UnixMicro()
		to     = query.TimeRange.To.UTC().UnixMicro()
	)
	return fmt.Sprintf("%s >= CAST(%d AS DATE) AND %s <=  CAST(%d AS DATE)", column, from, column, to), nil
}

func DateTimeFilter(query *sqlds.Query, args []string) (string, error) {
	if len(args) != 2 {
		return "", backend.DownstreamError(fmt.Errorf("%w: expected 2 arguments, received %d", sqlutil.ErrorBadArgumentCount, len(args)))
	}
	var (
		dateColumn = args[0]
		timeColumn = args[1]
		from       = query.TimeRange.From.UTC().Unix()
		to         = query.TimeRange.To.UTC().Unix()
	)

	dateFilter := fmt.Sprintf("%s >= CAST(%d AS DATE) AND %s <=  CAST(%d AS DATE)", dateColumn, from, dateColumn, to)
	timeFilter := fmt.Sprintf("%s >= CAST(%d AS DATETIME) AND %s <=  CAST(%d AS DATETIME)", timeColumn, from, timeColumn, to)
	return fmt.Sprintf("%s AND %s", dateFilter, timeFilter), nil
}

func IntervalSeconds(query *sqlds.Query, args []string) (string, error) {
	seconds := math.Max(query.Interval.Seconds(), 1)
	return fmt.Sprintf("%d", int(seconds)), nil
}

func VariableFallback(query *sqlds.Query, args []string) (string, error) {
	if len(args) != 2 {
		return "", backend.DownstreamError(fmt.Errorf("%w: expected 2 arguments, received %d", sqlutil.ErrorBadArgumentCount, len(args)))
	}
	value := args[1]
	if len(value) == 0 {
		fallback := args[0]
		if len(fallback) == 0 {
			return "", fmt.Errorf("%w", ErrInvalidVariableFallback)
		}
		value = fallback
	}
	return value, nil
}

// Macros is a map of all macro functions
var Macros = map[string]sqlds.MacroFunc{
	"fromTime":       FromTimestampFilter,
	"toTime":         ToTimestampFilter,
	"timeFilter":     TimestampFilter,
	"dateFilter":     DateFilter,
	"dateTimeFilter": DateTimeFilter,
	"dt":             DateTimeFilter,
	"interval_s":     IntervalSeconds,
	"varFallback":    VariableFallback,
}
