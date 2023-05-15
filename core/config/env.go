package config

import (
	"fmt"
	"sync"

	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

const (
	// Dev env alias
	Dev = "dev"
	// Prod env alias
	Prod = "prod"
)

var (
	once sync.Once
	// Config is a global containing the application general configuration.
	// It can be considered safe for concurrent access,
	// but should not be modified anywhere other than in the sync.Once guard, in the package `init` function.
	Config *Manifest
)

func init() {
	once.Do(func() {
		var m Manifest
		if err := envconfig.Process("fizzbuzz", &m); err != nil {
			panic(errors.Wrapf(err, "ambiguous environment variables"))
		}
		Config = &m
	})
}

// Manifest contains global configuration.
type Manifest struct {
	// Mode defines the application environment, 'dev' or 'prod'.
	// This affects multiple internal subsystems such as logs, error handlers and other conveniences.
	Mode string `default:"dev"`
	// Port is the http port that should be used by the server, defaults to 8080
	Port uint16 `default:"8080"`
	// Addr if set will override Port config, expects a valid golang listener addr, such as ":8080", defaults to empty
	Addr string `splitwords:"true"`
	// CorsEnabled will enable CORS on all endpoints if set, defaults to true
	CorsEnabled bool `splitwords:"true" default:"true"`
	// AllowEmptyStr allows empty str to be used as words for the fizzbuzz endpoint, defaults to false
	AllowEmptyStr bool `splitwords:"true" default:"false"`
}

// IsProd check whether the application is configured for production.
func (m Manifest) IsProd() bool {
	return m.Mode == Prod
}

// ListenAddr returns a http.Listener compatible string address
func (m Manifest) ListenAddr() string {
	if m.Addr == "" {
		return fmt.Sprintf(":%d", m.Port)
	}
	return m.Addr
}
