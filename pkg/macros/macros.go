package macros

import (
	"errors"
	"fmt"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data/sqlutil"
)

var ErrInvalidVariableFallback = errors.New("fallback should contain at least one character")

type timeQueryType string

const (
	timeQueryTypeFrom timeQueryType = "from"
	timeQueryTypeTo   timeQueryType = "to"
)

// Converts a time.Time to a Date
func timeToDate(t time.Time) string {
	return fmt.Sprintf("%s", t.Format("2006-01-02"))
}

func newTimeFilter(queryType timeQueryType, query *sqlutil.Query) (string, error) {
	date := query.TimeRange.From
	if queryType == timeQueryTypeTo {
		date = query.TimeRange.To
	}
	micros := date.UTC().UnixMicro()
	return fmt.Sprintf("CAST(%d AS TIMESTAMP)", micros), nil
}

// FromTimestampFilter return time filter query based on grafana's timepicker's from time in microseconds
func FromTimestampFilter(query *sqlutil.Query, args []string) (string, error) {
	return newTimeFilter(timeQueryTypeFrom, query)
}

// ToTimestampFilter return time filter query based on grafana's timepicker's to time in microseconds
func ToTimestampFilter(query *sqlutil.Query, args []string) (string, error) {
	return newTimeFilter(timeQueryTypeTo, query)
}

func TimestampFilter(query *sqlutil.Query, args []string) (string, error) {
	if len(args) != 1 {
		return "", backend.DownstreamError(fmt.Errorf("%w: expected 1 argument, received %d", sqlutil.ErrorBadArgumentCount, len(args)))
	}
	var (
		column = args[0]
		from   = query.TimeRange.From.UTC().UnixMicro()
		to     = query.TimeRange.To.UTC().UnixMicro()
	)
	return fmt.Sprintf("%s >= CAST(%d AS TIMESTAMP) AND %s <= CAST(%d AS TIMESTAMP)", column, from, column, to), nil
}

func DateFilter(query *sqlutil.Query, args []string) (string, error) {
	if len(args) != 1 {
		return "", backend.DownstreamError(fmt.Errorf("%w: expected 1 argument, received %d", sqlutil.ErrorBadArgumentCount, len(args)))
	}
	var (
		column = args[0]
		from   = query.TimeRange.From.UTC()
		to     = query.TimeRange.To.UTC()
	)
	return fmt.Sprintf("%s >= CAST(%s AS DATE) AND %s <= CAST(%s AS DATE)", column, timeToDate(from), column, timeToDate(to)), nil
}

func DateTimeFilter(query *sqlutil.Query, args []string) (string, error) {
	if len(args) != 2 {
		return "", backend.DownstreamError(fmt.Errorf("%w: expected 2 arguments, received %d", sqlutil.ErrorBadArgumentCount, len(args)))
	}
	var (
		dateColumn = args[0]
		timeColumn = args[1]
		from       = query.TimeRange.From.UTC()
		to         = query.TimeRange.To.UTC()
	)

	dateFilter := fmt.Sprintf("%s >= CAST(%s AS DATE) AND %s <= CAST(%s AS DATE)", dateColumn, timeToDate(from), dateColumn, timeToDate(to))
	timeFilter := fmt.Sprintf("%s >= CAST(%d AS DATETIME) AND %s <= CAST(%d AS DATETIME)", timeColumn, from.Unix(), timeColumn, to.Unix())
	return fmt.Sprintf("%s AND %s", dateFilter, timeFilter), nil
}

func VariableFallback(query *sqlutil.Query, args []string) (string, error) {
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
var Macros = map[string]sqlutil.MacroFunc{
	"fromTimestamp":  FromTimestampFilter,
	"toTimestamp":    ToTimestampFilter,
	"timeFilter":     TimestampFilter,
	"dateFilter":     DateFilter,
	"dateTimeFilter": DateTimeFilter,
	"varFallback":    VariableFallback,
}
