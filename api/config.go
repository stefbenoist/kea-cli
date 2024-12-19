package api

import (
	"net/http"
	"time"
)

const (
	DefaultKeaAPIURL         = "127.0.0.1:8000"
	DefaultHTTPClientTimeout = 30 * time.Second
	DefaultPageLimit         = 50
)

// Configuration holds the api client configuration
type Configuration struct {
	APIURL            string        `mapstructure:"api_url"`
	HTTPClientTimeout time.Duration `mapstructure:"http_client_timeout"`
	SSLEnabled        bool          `mapstructure:"ssl_enabled"`
	SkipTLSVerify     bool          `mapstructure:"skip_tls_verify"`
	KeyFile           string        `mapstructure:"key_file"`
	CertFile          string        `mapstructure:"cert_file"`
	CAFile            string        `mapstructure:"ca_file"`
	CAPath            string        `mapstructure:"ca_path"`
	LogLevel          string        `mapstructure:"log_level"`

	HttpClient *http.Client
}
