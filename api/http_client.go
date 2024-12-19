package api

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/goware/urlx"
	"github.com/hashicorp/go-rootcerts"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// HTTPClient represents an HTTP client
type HTTPClient interface {
	NewRequest(ctx context.Context, commandRequest CommandRequest) (*http.Request, error)
	Do(req *http.Request) (*http.Response, error)

	Lease4() Lease4
}

// GetClient returns a HTTP Client
func GetClient(cc Configuration) (HTTPClient, error) {
	apiURL := strings.TrimRight(cc.APIURL, "/")
	caFile := cc.CAFile
	caPath := cc.CAPath
	certFile := cc.CertFile
	keyFile := cc.KeyFile

	httpClient := cc.HttpClient
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: cc.HTTPClientTimeout,
		}
	}
	client := &client{
		baseURL: "http://" + apiURL,
		Client:  httpClient,
	}
	if cc.SSLEnabled || cc.CAFile != "" || cc.CAPath != "" || (certFile != "" && keyFile != "") {
		url, err := urlx.Parse(apiURL)
		if err != nil {
			return nil, errors.Wrap(err, "Malformed API URL")
		}
		apiHost, _, err := urlx.SplitHostPort(url)
		if err != nil {
			return nil, errors.Wrap(err, "Malformed API URL")
		}

		tlsConfig := &tls.Config{ServerName: apiHost}
		if certFile != "" && keyFile != "" {
			cert, err := tls.LoadX509KeyPair(certFile, keyFile)
			if err != nil {
				return nil, errors.Wrap(err, "Failed to load TLS certificates")
			}
			tlsConfig.Certificates = []tls.Certificate{cert}
		}
		if caFile != "" || caPath != "" {
			cfg := &rootcerts.Config{
				CAFile: caFile,
				CAPath: caPath,
			}
			rootcerts.ConfigureTLS(tlsConfig, cfg)
			tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert
			tlsConfig.BuildNameToCertificate()
		}
		if cc.SkipTLSVerify {
			tlsConfig.InsecureSkipVerify = true
			fmt.Println("Warning : usage of skip_tls_verify is not recommended for production and may expose to MITM attack")
		}

		tr := &http.Transport{
			TLSClientConfig: tlsConfig,
		}
		client.baseURL = "https://" + apiURL
		httpClient.Transport = tr
		client.Client = httpClient
	}

	client.lease4 = &lease4{client: client}
	return client, nil
}

// client is the HTTP client structure
type client struct {
	*http.Client
	baseURL string

	lease4 *lease4
}

func (c *client) Lease4() Lease4 {
	return c.lease4
}

// NewRequest returns a new HTTP request
func (c *client) NewRequest(ctx context.Context, commandRequest CommandRequest) (*http.Request, error) {
	body, err := json.Marshal(commandRequest)
	if err != nil {
		return nil, errors.Wrapf(err, "Fail to marshall command request:%v", commandRequest)
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	request.Header.Add(contentTypeHeaderName, jsonAppHeader)
	request.Header.Add(acceptHeaderName, jsonAppHeader)

	log.Debugf("request:%+v\n", request)
	log.Debugf("request body:%s\n", string(body))
	return request, nil
}

// ReadResponse is a helper function that allow to fully read and close a KEA response body and
// unmarshal its json content into a provided data structure.
func ReadResponse(response *http.Response, data interface{}) error {
	defer response.Body.Close()
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return errors.Wrap(err, "Cannot read the response from KEA")
	}

	log.Debugf("response:%+v\n", response)
	log.Debugf("response body:%s\n", string(responseBody))

	// Looks like kea-control-agent returns an array instead of a single response as defined in doc:https://kea.readthedocs.io/en/latest/api.html#lease4-get-page
	var (
		commandResponseArray []CommandResponse
		commandResponse      CommandResponse
	)
	// Try firstly with array of commandResponse
	if err = json.Unmarshal(responseBody, &commandResponseArray); err != nil {
		// Retry with a single commandResponse
		if err = json.Unmarshal(responseBody, &commandResponse); err != nil {
			return errors.Wrap(err, "Unable to unmarshal content of the KEA response")
		}
	} else {
		commandResponse = commandResponseArray[0]
	}

	switch commandResponse.Result {
	case 0:
		if err = mapstructure.Decode(commandResponse.Arguments, &data); err != nil {
			return errors.Wrap(err, "Fail to decode specified data from KEA response")
		}
	case 1:
		return malformedError{mes: commandResponse.Text}
	case 2:
		return unsupportedError{mes: commandResponse.Text}
	case 3:
		return notFoundError{mes: commandResponse.Text}
	case 4:
		return conflictError{mes: commandResponse.Text}
	}
	return nil
}
