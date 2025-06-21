package macros_test

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/sqlds/v2"
	"github.com/stretchr/testify/assert"

	"github.com/ydb/grafana-ydb-datasource/pkg/macros"
	"github.com/ydb/grafana-ydb-datasource/pkg/plugin"
)

type YdbDriver struct {
	sqlds.Driver
}

type MockDB struct {
	YdbDriver
}

func (h *YdbDriver) Macros() sqlds.Macros {
	C := plugin.Ydb{}
	return C.Macros()
}

func TestMacroFromTimestampFilter(t *testing.T) {
	from, _ := time.Parse("2006-01-02T15:04:05.000Z", "2021-11-12T11:45:26.371Z")
	to, _ := time.Parse("2006-01-02T15:04:05.000Z", "2022-11-12T11:45:26.371Z")
	query := sqlds.Query{
		TimeRange: backend.TimeRange{
			From: from,
			To:   to,
		},
		RawSQL: "select foo from foo where bar > $__fromTimestamp",
	}
	got, err := macros.FromTimestampFilter(&query, []string{})
	assert.Nil(t, err)
	assert.Equal(t, "CAST(1636717526371000 AS TIMESTAMP)", got)
}

func TestMacroToTimestampFilter(t *testing.T) {
	from, _ := time.Parse("2006-01-02T15:04:05.000Z", "2021-11-12T11:45:26.371Z")
	to, _ := time.Parse("2006-01-02T15:04:05.000Z", "2022-11-12T11:45:26.371Z")
	query := sqlds.Query{
		TimeRange: backend.TimeRange{
			From: from,
			To:   to,
		},
		RawSQL: "select foo from foo where bar > $__toTimestamp",
	}
	got, err := macros.ToTimestampFilter(&query, []string{})
	assert.Nil(t, err)
	assert.Equal(t, "CAST(1668253526371000 AS TIMESTAMP)", got)
}

func TestMacroTimestampFilter(t *testing.T) {
	from, _ := time.Parse("2006-01-02T15:04:05.000Z", "2021-11-12T11:45:26.371Z")
	to, _ := time.Parse("2006-01-02T15:04:05.000Z", "2022-11-12T11:45:26.371Z")
	query := sqlds.Query{
		TimeRange: backend.TimeRange{
			From: from,
			To:   to,
		},
	}
	got, err := macros.TimestampFilter(&query, []string{"foo"})
	assert.Nil(t, err)
	assert.Equal(t, "foo >= CAST(1636717526371000 AS TIMESTAMP) AND foo <= CAST(1668253526371000 AS TIMESTAMP)", got)
}

func TestMacroDateFilter(t *testing.T) {
	from, _ := time.Parse("2006-01-02T15:04:05.000Z", "2021-11-12T11:45:26.371Z")
	to, _ := time.Parse("2006-01-02T15:04:05.000Z", "2022-11-12T11:45:26.371Z")
	query := sqlds.Query{
		TimeRange: backend.TimeRange{
			From: from,
			To:   to,
		},
	}
	got, err := macros.DateFilter(&query, []string{"foo"})
	assert.Nil(t, err)
	assert.Equal(t, "foo >= CAST(2021-11-12 AS DATE) AND foo <= CAST(2022-11-12 AS DATE)", got)
}

func TestMacroDateTimeFilter(t *testing.T) {
	from, _ := time.Parse("2006-01-02T15:04:05.000Z", "2021-11-12T11:45:26.371Z")
	to, _ := time.Parse("2006-01-02T15:04:05.000Z", "2022-11-12T11:45:26.371Z")
	query := sqlds.Query{
		TimeRange: backend.TimeRange{
			From: from,
			To:   to,
		},
	}
	got, err := macros.DateTimeFilter(&query, []string{"date_col", "time_col"})
	assert.Nil(t, err)
	assert.Equal(t, "date_col >= CAST(2021-11-12 AS DATE) AND date_col <= CAST(2022-11-12 AS DATE) "+
		"AND time_col >= CAST(1636717526 AS DATETIME) AND time_col <= CAST(1668253526 AS DATETIME)", got)
}

func TestMacroVariableFallback(t *testing.T) {
	query := sqlds.Query{
		RawSQL: "select $__varFallback(fallback, value)",
	}
	got, err := macros.VariableFallback(&query, []string{"fallback", "value"})
	assert.Nil(t, err)
	assert.Equal(t, "value", got)
}

func TestMacroVariableFallbackNoValue(t *testing.T) {
	query := sqlds.Query{
		RawSQL: "select $__varFallback(fallback, '')",
	}
	got, err := macros.VariableFallback(&query, []string{"fallback", ""})
	assert.Nil(t, err)
	assert.Equal(t, "fallback", got)
}

// test sqlds query interpolation with YDB filters used
func TestInterpolate(t *testing.T) {
	from, _ := time.Parse("2006-01-02T15:04:05.000Z", "2021-11-12T11:45:26.371Z")
	to, _ := time.Parse("2006-01-02T15:04:05.000Z", "2022-11-12T11:45:26.371Z")

	tableName := "my_table"
	tableColumn := "my_col"

	type test struct {
		name   string
		input  string
		output string
	}

	tests := []test{
		{
			input:  "SELECT * FROM foo WHERE (DATE >= $__fromTimestamp AND DATE <= $__toTimestamp)",
			output: "SELECT * FROM foo WHERE (DATE >= CAST(1636717526371000 AS TIMESTAMP) AND DATE <= CAST(1668253526371000 AS TIMESTAMP))",
			name:   "YDB fromTimestamp and toTimestamp",
		},
		{
			input:  "SELECT * FROM foo WHERE $__timeFilter(sth)",
			output: "SELECT * FROM foo WHERE sth >= CAST(1636717526371000 AS TIMESTAMP) AND sth <= CAST(1668253526371000 AS TIMESTAMP)",
			name:   "YDB timeFilter",
		},
		{
			input:  "SELECT * FROM foo WHERE $__timeFilter(sth )",
			output: "SELECT * FROM foo WHERE sth >= CAST(1636717526371000 AS TIMESTAMP) AND sth <= CAST(1668253526371000 AS TIMESTAMP)",
			name:   "YDB timeFilter with spaces",
		},
		{
			input:  "SELECT * FROM foo WHERE $__dateFilter(sth)",
			output: "SELECT * FROM foo WHERE sth >= CAST(2021-11-12 AS DATE) AND sth <= CAST(2022-11-12 AS DATE)",
			name:   "YDB dateFilter",
		},
		{
			input:  "SELECT * FROM foo WHERE $__dateFilter(sth )",
			output: "SELECT * FROM foo WHERE sth >= CAST(2021-11-12 AS DATE) AND sth <= CAST(2022-11-12 AS DATE)",
			name:   "YDB dateFilter with spaces",
		},
		{
			input:  "SELECT * FROM foo WHERE $__dateTimeFilter(d_col, t_col)",
			output: "SELECT * FROM foo WHERE d_col >= CAST(2021-11-12 AS DATE) AND d_col <= CAST(2022-11-12 AS DATE) AND t_col >= CAST(1636717526 AS DATETIME) AND t_col <= CAST(1668253526 AS DATETIME)",
			name:   "YDB dateTimeFilter",
		},
	}

	for i, tc := range tests {
		driver := MockDB{}
		t.Run(fmt.Sprintf("[%d/%d] %s", i+1, len(tests), tc.name), func(t *testing.T) {
			query := &sqlds.Query{
				RawSQL: tc.input,
				Table:  tableName,
				Column: tableColumn,
				TimeRange: backend.TimeRange{
					From: from,
					To:   to,
				},
			}
			interpolatedQuery, err := sqlds.Interpolate(&driver, query)
			require.Nil(t, err)
			assert.Equal(t, tc.output, interpolatedQuery)
		})
	}
}
