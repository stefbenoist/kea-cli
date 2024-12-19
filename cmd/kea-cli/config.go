package main

import (
	"fmt"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"kea-cli/api"
)

// initFlags allows to init client configuration

func initFlags() {
	rootCmd.Version = fmt.Sprintf("%s (Git commit: %s)", version, commit)
	rootCmd.SetHelpCommand(&cobra.Command{Hidden: true})
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "Config file (default is /etc/kea-cli/kea-cli.[json|yaml])")
	rootCmd.PersistentFlags().StringP("api_url", "", api.DefaultKeaAPIURL, "Specify the host and port used to join the Kea REST API")
	rootCmd.PersistentFlags().Duration("http_client_timeout", api.DefaultHTTPClientTimeout, "Timeout for HTTP Client used to join the Kea REST API")
	rootCmd.PersistentFlags().BoolP("ssl_enabled", "", false, "Use HTTPS to connect to the Kea REST API")
	rootCmd.PersistentFlags().BoolP("skip_tls_verify", "", false, "Controls whether a client verifies the server's certificate chain and host name. If set to true, TLS accepts any certificate presented by the server and any host name in that certificate. In this mode, TLS is susceptible to man-in-the-middle attacks. This should be used only for testing. This implies the use of HTTPS to connect to the Kea REST API.")
	rootCmd.PersistentFlags().StringP("cert_file", "", "", "File path to a PEM-encoded client certificate used to authenticate to the Kea REST API. This must be provided along with key-file. If one of key-file or cert-file is not provided then SSL authentication is disabled. If both cert-file and key-file are provided this implies the use of HTTPS to connect to the Kea REST API.")
	rootCmd.PersistentFlags().StringP("key_file", "", "", "File path to a PEM-encoded client private key used to authenticate to the Kea REST API. This must be provided along with cert-file. If one of key-file or cert-file is not provided then SSL authentication is disabled. If both cert-file and key-file are provided this implies the use of HTTPS to connect to the Kea REST API.")
	rootCmd.PersistentFlags().StringP("ca_file", "", "", "This provides a file path to a PEM-encoded certificate authority. This implies the use of HTTPS to connect to the Kea REST API.")
	rootCmd.PersistentFlags().StringP("ca_path", "", "", "Path to a directory of PEM-encoded certificates authorities. This implies the use of HTTPS to connect to the Kea REST API.")
	rootCmd.PersistentFlags().StringP("log_level", "", log.InfoLevel.String(), "Log level: can be 'trace', 'debug', 'info', 'warning', 'error', 'fatal', 'panic'.")

	viper.BindPFlag("api_url", rootCmd.PersistentFlags().Lookup("api_url"))
	viper.BindPFlag("http_client_timeout", rootCmd.PersistentFlags().Lookup("http_client_timeout"))
	viper.BindPFlag("ssl_enabled", rootCmd.PersistentFlags().Lookup("ssl_enabled"))
	viper.BindPFlag("ca_file", rootCmd.PersistentFlags().Lookup("ca_file"))
	viper.BindPFlag("ca_path", rootCmd.PersistentFlags().Lookup("ca_path"))
	viper.BindPFlag("key_file", rootCmd.PersistentFlags().Lookup("key_file"))
	viper.BindPFlag("cert_file", rootCmd.PersistentFlags().Lookup("cert_file"))
	viper.BindPFlag("skip_tls_verify", rootCmd.PersistentFlags().Lookup("skip_tls_verify"))
	viper.BindPFlag("log_level", rootCmd.PersistentFlags().Lookup("log_level"))

	viper.SetEnvPrefix("kea_cli")
	viper.AutomaticEnv()
	viper.BindEnv("api_url")
	viper.BindEnv("http_client_timeout")
	viper.BindEnv("ssl_enabled")
	viper.BindEnv("ca_file")
	viper.BindEnv("ca_path")
	viper.BindEnv("key_file")
	viper.BindEnv("cert_file")
	viper.BindEnv("skip_tls_verify")
	viper.BindEnv("log_level")

	viper.SetDefault("api_url", api.DefaultKeaAPIURL)
	viper.SetDefault("http_client_timeout", api.DefaultHTTPClientTimeout)
	viper.SetDefault("ssl_enabled", false)
	viper.SetDefault("log_level", log.InfoLevel.String())

	//Configuration file directories
	viper.SetConfigName("kea-cli")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.kea-cli")
	viper.AddConfigPath("/etc/kea-cli/")
}

// getClientConfig retrieves client configuration
func getClientConfig(clientConfig *api.Configuration) error {
	if cfgFile != "" {
		// enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		ok := errors.As(err, &configFileNotFoundError)
		if !ok {
			fmt.Println("Can't use config file:", err)
		}
	}
	err := viper.Unmarshal(&clientConfig)
	if err != nil {
		return errors.Wrapf(err, "bad configuration")
	}
	setupLogs(clientConfig)
	return nil
}

func setupLogs(c *api.Configuration) {
	level := log.InfoLevel.String()
	if c.LogLevel != "" {
		level = c.LogLevel
	}
	l, err := log.ParseLevel(level)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Fatal("failed to parse log level")
	}
	log.SetLevel(l)
	log.SetFormatter(&log.TextFormatter{})
}
