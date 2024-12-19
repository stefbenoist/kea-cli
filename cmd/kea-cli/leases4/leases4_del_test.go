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

func TestDel(t *testing.T) {
	respOK := api.CommandResponse{
		Result: 0,
		Text:   "Lease added.",
	}
	bRespOK, err := json.Marshal(respOK)
	assert.NilError(t, err, "expecting nil error")

	respKO := api.CommandResponse{
		Result: 1,
		Text:   "bad parameter 'ip-address'",
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

		if cmd.Arguments["ip-address"] == "bad" {
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
		ipAddress string
		hwAddress string
		clientID  string
		subnetID  int
	}

	tests := []struct {
		name        string
		args        args
		wantErr     bool
		errContains string
	}{
		{"SimpleTestWithIpAddress", args{ipAddress: "192.0.2.202"}, false, ""},
		{"SimpleTestWithMACAddress", args{hwAddress: "1a:1b:1c:1d:1e:1f", clientID: "42:42:42:42:42:42:42:42"}, false, ""},
		{"SimpleTestWithClientID", args{hwAddress: "08:08:08:08:08:08:42", clientID: "42:42:42:42:42:42:42:42"}, false, ""},
		{"BadIpAddressTest", args{ipAddress: "bad"}, true, "bad parameter 'ip-address'"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err = del(client, tt.args.ipAddress, tt.args.hwAddress, tt.args.clientID, tt.args.subnetID)
			if tt.wantErr != (err != nil) {
				t.Errorf("del() err= %v, wantError %v", err, tt.wantErr)
			}
			if tt.wantErr == true && tt.errContains != "" {
				assert.ErrorContains(t, err, tt.errContains)
			}
		})
	}
}
