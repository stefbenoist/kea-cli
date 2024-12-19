package api

import (
	"context"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

var DHCPv4Services = []string{"dhcp4"}

type Lease4 interface {
	List(ctx context.Context, from string, limit int) ([]Lease, error)
	Search(ctx context.Context, criteria, value string) ([]Lease, error)
	Get(ctx context.Context, identifier *LeaseIdentifier) (*Lease, error)
	Del(ctx context.Context, identifier *LeaseIdentifier) error
	Add(ctx context.Context, ipAddress, hwAddress string, subnetID int) error
}

type lease4 struct {
	client *client
}

func (l *lease4) List(ctx context.Context, from string, limit int) ([]Lease, error) {
	command := CommandRequest{
		Command:  "lease4-get-page",
		Services: DHCPv4Services,
		Arguments: map[string]interface{}{
			"from":  from,
			"limit": limit,
		},
	}
	request, err := l.client.NewRequest(ctx, command)
	if err != nil {
		return nil, errors.Wrap(err, httpRequestFailureErrMsg)
	}
	response, err := l.client.Do(request)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to send http request to get all leases from:%s and limited to:%d", from, limit)
	}

	var leaseList LeaseList
	if err = ReadResponse(response, &leaseList); err != nil {
		return nil, err
	}
	return leaseList.Leases, nil
}

func (l *lease4) Search(ctx context.Context, criteria, value string) ([]Lease, error) {
	command := CommandRequest{
		Command:  fmt.Sprintf("lease4-get-by-%s", criteria),
		Services: DHCPv4Services,
		Arguments: map[string]interface{}{
			criteria: value,
		},
	}
	request, err := l.client.NewRequest(ctx, command)
	if err != nil {
		return nil, errors.Wrap(err, httpRequestFailureErrMsg)
	}
	response, err := l.client.Do(request)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to send http request to search leases with %s=%s", criteria, value)
	}

	var leaseList LeaseList
	if err = ReadResponse(response, &leaseList); err != nil {
		return nil, err
	}
	return leaseList.Leases, nil
}

func (l *lease4) Get(ctx context.Context, identifier *LeaseIdentifier) (*Lease, error) {
	arguments, err := fillArgsFromIdentifier(identifier)
	if err != nil {
		return nil, errors.Wrap(err, "failed to fill command arguments")
	}
	command := CommandRequest{
		Command:   "lease4-get",
		Services:  DHCPv4Services,
		Arguments: arguments,
	}
	request, err := l.client.NewRequest(ctx, command)
	if err != nil {
		return nil, errors.Wrap(err, httpRequestFailureErrMsg)
	}
	response, err := l.client.Do(request)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to send http request to get lease for identifier:%+v", *identifier)
	}

	var lease Lease
	if err = ReadResponse(response, &lease); err != nil {
		return nil, err
	}

	return &lease, nil
}

func (l *lease4) Add(ctx context.Context, ipAddress, hwAddress string, subnetID int) error {
	command := CommandRequest{
		Command:  "lease4-add",
		Services: DHCPv4Services,
		Arguments: map[string]interface{}{
			"ip-address": ipAddress,
			"hw-address": hwAddress,
			"subnet-id":  subnetID,
		},
	}
	request, err := l.client.NewRequest(ctx, command)
	if err != nil {
		return errors.Wrap(err, httpRequestFailureErrMsg)
	}
	response, err := l.client.Do(request)
	if err != nil {
		return errors.Wrapf(err, "failed to send http request to add lease for ip:%s and MAC:%s", ipAddress, hwAddress)
	}
	return ReadResponse(response, nil)
}

func (l *lease4) Del(ctx context.Context, identifier *LeaseIdentifier) error {
	arguments, err := fillArgsFromIdentifier(identifier)
	if err != nil {
		return errors.Wrap(err, "failed to fill command arguments")
	}
	command := CommandRequest{
		Command:   "lease4-del",
		Services:  DHCPv4Services,
		Arguments: arguments,
	}
	request, err := l.client.NewRequest(ctx, command)
	if err != nil {
		return errors.Wrap(err, httpRequestFailureErrMsg)
	}
	response, err := l.client.Do(request)
	if err != nil {
		return errors.Wrapf(err, "failed to send http request to delete lease for indentifier:%+v", *identifier)
	}
	return ReadResponse(response, nil)
}

func fillArgsFromIdentifier(identifier *LeaseIdentifier) (map[string]interface{}, error) {
	if identifier == nil {
		return nil, errors.New("missing identifier for requesting leases")
	}
	var m map[string]interface{}
	if err := mapstructure.Decode(identifier, &m); err != nil {
		return nil, err
	}
	return m, nil
}
