package plugin

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/stretchr/testify/assert"
)

const (
	endpoint   = "grpc://localhost:2136"
	dbLocation = "/local"
	authKind   = "anonymous"
)

func TestConnect(t *testing.T) {
	ydb := Ydb{}

	t.Run("should not error when valid settings passed", func(t *testing.T) {
		settings := backend.DataSourceInstanceSettings{JSONData: []byte(fmt.Sprintf(`{"authKind": "%s", "endpoint": "%s", "dbLocation": "%s"}`, authKind, endpoint, dbLocation))}
		_, err := ydb.Connect(settings, json.RawMessage{})
		assert.Equal(t, nil, err)
	})
}
