package api

const (
	jsonAppHeader         = "application/json"
	contentTypeHeaderName = "Content-Type"
	acceptHeaderName      = "Accept"

	httpRequestFailureErrMsg = "failed to create http request"
)

type Lease struct {
	ClientID  string `json:"client-id" mapstructure:"client-id"`
	Cltt      int64  `json:"cltt" mapstructure:"cltt"`
	FQDNFwd   bool   `json:"fqdn-fwd" mapstructure:"fqdn-fwd"`
	FQDNRev   bool   `json:"fqdn-rev" mapstructure:"fqdn-rev"`
	Hostname  string `json:"hostname" mapstructure:"hostname"`
	HwAddress string `json:"hw-address" mapstructure:"hw-address"`
	IpAddress string `json:"ip-address" mapstructure:"ip-address"`
	State     int    `json:"state" mapstructure:"state"`
	SubnetID  int    `json:"subnet-id" mapstructure:"subnet-id"`
	ValidLft  int    `json:"valid-lft"  mapstructure:"valid-lft"`
}

type LeaseList struct {
	Count  int     `json:"count" mapstructure:"count"`
	Leases []Lease `json:"leases" mapstructure:"leases"`
}

type LeaseIdentifier struct {
	IdentifierType  string `json:"identifier-type,omitempty" mapstructure:"identifier-type"`
	IdentifierValue string `json:"identifier,omitempty" mapstructure:"identifier"`
	IpAddress       string `json:"ip-address,omitempty" mapstructure:"ip-address"`
	SubnetID        int    `json:"subnet-id,omitempty" mapstructure:"subnet-id"`
}

type CommandRequest struct {
	Command   string                 `json:"command"`
	Services  []string               `json:"service"`
	Arguments map[string]interface{} `json:"arguments,omitempty"`
}

type CommandResponse struct {
	Arguments interface{} `json:"arguments,omitempty"`
	Result    int         `json:"result"`
	Text      string      `json:"text"`
}
