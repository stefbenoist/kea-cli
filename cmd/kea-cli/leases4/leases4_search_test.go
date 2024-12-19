package leases4

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"gotest.tools/v3/assert"

	"kea-cli/api"
)

func TestSearch(t *testing.T) {
	respOK := []api.CommandResponse{
		{
			Result: 0,
			Text:   "IPv4 lease found.",
			Arguments: api.LeaseList{
				Count: 2,
				Leases: []api.Lease{
					{
						ClientID:  "42:42:42:42:42:42:42:42",
						Cltt:      12345678,
						FQDNFwd:   false,
						FQDNRev:   true,
						Hostname:  "myhost.example.com.",
						HwAddress: "08:08:08:08:08:08",
						IpAddress: "192.0.2.1",
						State:     0,
						SubnetID:  44,
						ValidLft:  3600,
					},
					{
						ClientID:  "01:00:0c:01:02:03:05",
						Cltt:      1600432234,
						FQDNFwd:   false,
						FQDNRev:   false,
						Hostname:  "",
						HwAddress: "00:0c:01:02:03:05",
						IpAddress: "192.168.1.151",
						State:     0,
						SubnetID:  1,
						ValidLft:  4000,
					},
				},
			},
		},
	}
	bRespOK, err := json.Marshal(respOK)
	assert.NilError(t, err, "expecting nil error")

	respKO := api.CommandResponse{
		Result: 1,
		Text:   "bad parameter",
	}
	bRespKO, err := json.Marshal(respKO)
	assert.NilError(t, err, "expecting nil error")

	// Start a local HTTP server simulating KEA REST API server
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		body, err := io.ReadAll(req.Body)
		assert.NilError(t, err, "expecting nil error")

		var cmd api.CommandRequest
		err = json.Unmarshal(body, &cmd)
		assert.NilError(t, err, "expecting nil error")

		if cmd.Arguments["hw-address"] == "bad" {
			rw.Write(bRespKO)
		} else {
			rw.Write(bRespOK)
		}
	}))
	defer server.Close()

	testClientConfig := api.Configuration{
		APIURL:     strings.TrimPrefix(server.URL, "http://"),
		HttpClient: server.Client(),
	}
	client, err := api.GetClient(testClientConfig)
	assert.NilError(t, err, "expecting nil error")

	type args struct {
		hwAddress string
		clientID  string
		hostname  string
	}

	tests := []struct {
		name        string
		args        args
		wantErr     bool
		errContains string
	}{
		{"SearchByHwAddress", args{hwAddress: "08:08:08:08:08:08"}, false, ""},
		{"SearchByHostname", args{hostname: "myhost.example.org"}, false, ""},
		{"SearchByClientID", args{clientID: "01:00:0c:01:02:03:04"}, false, ""},
		{"BadParameterTest", args{hwAddress: "bad"}, true, "bad parameter"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err = search(client, tt.args.hwAddress, tt.args.hostname, tt.args.clientID)
			if tt.wantErr != (err != nil) {
				t.Errorf("get() err= %v, wantError %v", err, tt.wantErr)
			}
			if tt.wantErr == true && tt.errContains != "" {
				assert.ErrorContains(t, err, tt.errContains)
			}
		})
	}
}
