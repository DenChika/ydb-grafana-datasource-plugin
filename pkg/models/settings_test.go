package models

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/stretchr/testify/assert"
)

func TestLoadSettings(t *testing.T) {
	t.Run("should parse settings correctly", func(t *testing.T) {
		type args struct {
			config backend.DataSourceInstanceSettings
		}
		tests := []struct {
			name         string
			args         args
			wantSettings Settings
		}{
			{
				name: "should set all json fields correctly by default",
				args: args{
					config: backend.DataSourceInstanceSettings{
						UID:      "ds-uid",
						JSONData: []byte(""),
						DecryptedSecureJSONData: map[string]string{
							"serviceAccAuthAccessKey": "test-key",
							"accessToken":             "test-token",
							"password":                "test-pswd",
							"certificate":             "test-cert",
						},
					},
				},
				wantSettings: Settings{
					AuthKind:   "ServiceAccountKey",
					DBEndpoint: "grpc://localhost:2136",
					DBLocation: "/local",
					Secrets: &SecretPluginSettings{
						ServiceAccAuthAccessKey: "test-key",
						AccessToken:             "test-token",
						Password:                "test-pswd",
						Certificate:             "test-cert",
					},
					Dsn:                "grpc://localhost:2136/local",
					IsSecureConnection: false,
					Timeout:            "10",
					TimeoutDuration:    time.Duration(10) * time.Second,
				},
			},
			{
				name: "anonymous auth: should parse and set all json fields correctly",
				args: args{
					config: backend.DataSourceInstanceSettings{
						UID: "ds-uid",
						JSONData: []byte(`{
							"authKind": "anonymous",
							"endpoint": "grpc://foo:2136",
							"dbLocation": "/local"
						}`),
					},
				},
				wantSettings: Settings{
					AuthKind:   "anonymous",
					DBEndpoint: "grpc://foo:2136",
					DBLocation: "/local",
					Secrets: &SecretPluginSettings{
						ServiceAccAuthAccessKey: "",
						AccessToken:             "",
						Password:                "",
						Certificate:             "",
					},
					Dsn:                "grpc://foo:2136/local",
					IsSecureConnection: false,
					Timeout:            "10",
					TimeoutDuration:    time.Duration(10) * time.Second,
				},
			},
			{
				name: "ServiceAccountKey auth: should parse and set all json fields correctly",
				args: args{
					config: backend.DataSourceInstanceSettings{
						UID: "ds-uid",
						JSONData: []byte(`{
							"authKind": "ServiceAccountKey",
							"endpoint": "grpc://foo:2136",
							"dbLocation": "/local"
						}`),
						DecryptedSecureJSONData: map[string]string{
							"serviceAccAuthAccessKey": "test-key",
						},
					},
				},
				wantSettings: Settings{
					AuthKind:   "ServiceAccountKey",
					DBEndpoint: "grpc://foo:2136",
					DBLocation: "/local",
					Secrets: &SecretPluginSettings{
						ServiceAccAuthAccessKey: "test-key",
						AccessToken:             "",
						Password:                "",
						Certificate:             "",
					},
					Dsn:                "grpc://foo:2136/local",
					IsSecureConnection: false,
					Timeout:            "10",
					TimeoutDuration:    time.Duration(10) * time.Second,
				},
			},
			{
				name: "AccessToken auth: should parse and set all json fields correctly",
				args: args{
					config: backend.DataSourceInstanceSettings{
						UID: "ds-uid",
						JSONData: []byte(`{
							"authKind": "AccessToken",
							"endpoint": "grpc://foo:2136",
							"dbLocation": "/local"
						}`),
						DecryptedSecureJSONData: map[string]string{
							"accessToken": "test-token",
						},
					},
				},
				wantSettings: Settings{
					AuthKind:   "AccessToken",
					DBEndpoint: "grpc://foo:2136",
					DBLocation: "/local",
					Secrets: &SecretPluginSettings{
						ServiceAccAuthAccessKey: "",
						AccessToken:             "test-token",
						Password:                "",
						Certificate:             "",
					},
					Dsn:                "grpc://foo:2136/local",
					IsSecureConnection: false,
					Timeout:            "10",
					TimeoutDuration:    time.Duration(10) * time.Second,
				},
			},
			{
				name: "UserPassword auth: should parse and set all json fields correctly",
				args: args{
					config: backend.DataSourceInstanceSettings{
						UID: "ds-uid",
						JSONData: []byte(`{
							"authKind": "UserPassword",
							"endpoint": "grpc://foo:2136",
							"dbLocation": "/local",
							"user": "test-user"
						}`),
						DecryptedSecureJSONData: map[string]string{
							"password": "test-pswd",
						},
					},
				},
				wantSettings: Settings{
					AuthKind:   "UserPassword",
					DBEndpoint: "grpc://foo:2136",
					DBLocation: "/local",
					User:       "test-user",
					Secrets: &SecretPluginSettings{
						ServiceAccAuthAccessKey: "",
						AccessToken:             "",
						Password:                "test-pswd",
						Certificate:             "",
					},
					Dsn:                "grpc://foo:2136/local",
					IsSecureConnection: false,
					Timeout:            "10",
					TimeoutDuration:    time.Duration(10) * time.Second,
				},
			},
			{
				name: "should parse and set secureConnection correctly",
				args: args{
					config: backend.DataSourceInstanceSettings{
						UID: "ds-uid",
						JSONData: []byte(`{
							"authKind": "anonymous",
							"endpoint": "grpcs://foo:2136",
							"dbLocation": "/local"
						}`),
					},
				},
				wantSettings: Settings{
					AuthKind:   "anonymous",
					DBEndpoint: "grpcs://foo:2136",
					DBLocation: "/local",
					Secrets: &SecretPluginSettings{
						ServiceAccAuthAccessKey: "",
						AccessToken:             "",
						Password:                "",
						Certificate:             "",
					},
					Dsn:                "grpcs://foo:2136/local",
					IsSecureConnection: true,
					Timeout:            "10",
					TimeoutDuration:    time.Duration(10) * time.Second,
				},
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				gotSettings, err := LoadSettings(tt.args.config)
				assert.Equal(t, nil, err)
				if !reflect.DeepEqual(*gotSettings, tt.wantSettings) {
					t.Errorf("LoadSettings() = %v, want %v", *gotSettings, tt.wantSettings)
				}
			})
		}
	})
	t.Run("should capture invalid settings", func(t *testing.T) {
		type args struct {
			config backend.DataSourceInstanceSettings
		}
		tests := []struct {
			name    string
			args    args
			wantErr error
		}{
			{
				name: "should capture empty endpoint",
				args: args{
					config: backend.DataSourceInstanceSettings{
						UID: "ds-uid",
						JSONData: []byte(`{
							"authKind": "anonymous",
							"dbLocation": "/local"
						}`),
					},
				},
				wantErr: backend.DownstreamError(fmt.Errorf("%w", ErrEndpointEmpty)),
			},
			{
				name: "should capture empty db location",
				args: args{
					config: backend.DataSourceInstanceSettings{
						UID: "ds-uid",
						JSONData: []byte(`{
							"authKind": "anonymous",
							"endpoint": "grpc://foo:2136"
						}`),
					},
				},
				wantErr: backend.DownstreamError(fmt.Errorf("%w", ErrDBLocationEmpty)),
			},
			{
				name: "ServiceAccountKey auth: should capture empty key",
				args: args{
					config: backend.DataSourceInstanceSettings{
						UID: "ds-uid",
						JSONData: []byte(`{
							"authKind": "ServiceAccountKey",
							"endpoint": "grpc://foo:2136",
							"dbLocation": "/local"
						}`),
					},
				},
				wantErr: backend.DownstreamError(fmt.Errorf("%w", ErrServiceAccAuthAccessKeyEmpty)),
			},
			{
				name: "AccessTokenAuth auth: should capture empty token",
				args: args{
					config: backend.DataSourceInstanceSettings{
						UID: "ds-uid",
						JSONData: []byte(`{
							"authKind": "AccessToken",
							"endpoint": "grpc://foo:2136",
							"dbLocation": "/local"
						}`),
					},
				},
				wantErr: backend.DownstreamError(fmt.Errorf("%w", ErrAccessTokenEmpty)),
			},
			{
				name: "UserPassword auth: should capture empty user",
				args: args{
					config: backend.DataSourceInstanceSettings{
						UID: "ds-uid",
						JSONData: []byte(`{
							"authKind": "UserPassword",
							"endpoint": "grpc://foo:2136",
							"dbLocation": "/local"
						}`),
						DecryptedSecureJSONData: map[string]string{
							"password": "test-pswd",
						},
					},
				},
				wantErr: backend.DownstreamError(fmt.Errorf("%w", ErrUserOrPasswordEmpty)),
			},
			{
				name: "UserPassword auth: should capture empty password",
				args: args{
					config: backend.DataSourceInstanceSettings{
						UID: "ds-uid",
						JSONData: []byte(`{
							"authKind": "UserPassword",
							"endpoint": "grpc://foo:2136",
							"dbLocation": "/local",
							"user": "test-user"
						}`),
					},
				},
				wantErr: backend.DownstreamError(fmt.Errorf("%w", ErrUserOrPasswordEmpty)),
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				_, err := LoadSettings(tt.args.config)
				assert.Equal(t, tt.wantErr, err)
			})
		}
	})
}
