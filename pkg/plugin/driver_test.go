package plugin

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/stretchr/testify/assert"
)

const (
	endpoint   = "grpc://localhost:2136"
	dbLocation = "/local"
)

func TestConnect(t *testing.T) {
	ydb := Ydb{}

	ctx := context.Background()

	t.Run("anonymous: should not error when valid settings passed", func(t *testing.T) {
		settings := backend.DataSourceInstanceSettings{
			JSONData: []byte(
				fmt.Sprintf(`{
					"authKind": "anonymous", 
					"endpoint": "%s", 
					"dbLocation": "%s"
				}`, endpoint, dbLocation)),
		}
		_, err := ydb.Connect(ctx, settings, json.RawMessage{})
		assert.Equal(t, nil, err)
	})
}
