package models

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

type AuthKind string

const (
	defaultAuthKind   AuthKind = `ServiceAccountKey`
	defaultDBEndpoint          = "grpc://localhost:2136"
	defaultDBLocation          = "/local"
	defaultTimeout             = "10"
)

// Settings - data loaded from grafana settings database
type Settings struct {
	AuthKind           AuthKind              `json:"authKind"`
	DBEndpoint         string                `json:"endpoint,omitempty"`
	DBLocation         string                `json:"dbLocation,omitempty"`
	User               string                `json:"user,omitempty"`
	Secrets            *SecretPluginSettings `json:"-"`
	Dsn                string
	IsSecureConnection bool
	Timeout            string
	TimeoutDuration    time.Duration
}

type SecretPluginSettings struct {
	ServiceAccAuthAccessKey string
	AccessToken             string
	Password                string
	Certificate             string
}

type SettingsOptionFunc func(settings *Settings)

// LoadSettings will read and validate Settings from the DataSourceConfig
func LoadSettings(source backend.DataSourceInstanceSettings) (*Settings, error) {
	if source.JSONData == nil || len(source.JSONData) < 1 {
		// If no settings have been saved return default values
		settings := Settings{
			AuthKind:   defaultAuthKind,
			DBEndpoint: defaultDBEndpoint,
			DBLocation: defaultDBLocation,
			Secrets:    loadSecretPluginSettings(source.DecryptedSecureJSONData),
			Timeout:    defaultTimeout,
		}

		return validateSettings(settings)
	}

	settings := Settings{
		AuthKind: defaultAuthKind,
		Timeout:  defaultTimeout,
	}

	err := json.Unmarshal(source.JSONData, &settings)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", err.Error(), ErrInvalidJSON)
	}

	settings.Secrets = loadSecretPluginSettings(source.DecryptedSecureJSONData)

	return validateSettings(settings)
}

func loadSecretPluginSettings(source map[string]string) *SecretPluginSettings {
	return &SecretPluginSettings{
		ServiceAccAuthAccessKey: source["serviceAccAuthAccessKey"],
		AccessToken:             source["accessToken"],
		Password:                source["password"],
		Certificate:             source["certificate"],
	}
}

func validateSettings(settings Settings) (*Settings, error) {
	if settings.DBEndpoint == "" {
		return nil, backend.DownstreamError(fmt.Errorf("%w", ErrEndpointEmpty))
	}
	if settings.DBLocation == "" {
		return nil, backend.DownstreamError(fmt.Errorf("%w", ErrDBLocationEmpty))
	}

	switch settings.AuthKind {
	case "ServiceAccountKey":
		if settings.Secrets.ServiceAccAuthAccessKey == "" {
			return nil, backend.DownstreamError(fmt.Errorf("%w", ErrServiceAccAuthAccessKeyEmpty))
		}
	case "AccessToken":
		if settings.Secrets.AccessToken == "" {
			return nil, backend.DownstreamError(fmt.Errorf("%w", ErrAccessTokenEmpty))
		}
	case "UserPassword":
		if settings.Secrets.Password == "" || settings.User == "" {
			return nil, backend.DownstreamError(fmt.Errorf("%w", ErrUserOrPasswordEmpty))
		}
	}

	settings.Dsn = fmt.Sprintf("%s%s", settings.DBEndpoint, settings.DBLocation)
	settings.IsSecureConnection = strings.HasPrefix(settings.DBEndpoint, "grpcs://")

	t, err := strconv.Atoi(settings.Timeout)
	if err != nil {
		return nil, backend.DownstreamError(fmt.Errorf("timeout %s invalid: %w", settings.Timeout, err))
	}

	settings.TimeoutDuration = time.Duration(t) * time.Second

	return &settings, nil
}
