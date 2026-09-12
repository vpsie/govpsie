package govpsie

import (
	"context"
	"fmt"
	"net/http"
)

var ipsPath = "/apps/v2/ips"

type IPsService interface {
	ListPrivateIPs(ctx context.Context, options *ListOptions) ([]IP, error)
	ListPublicIPs(ctx context.Context, options *ListOptions) ([]IP, error)
	ListAllIPs(ctx context.Context, options *ListOptions) ([]IP, error)
	DeleteIP(ctx context.Context, ip, vmIdentifier string) error
	CreateIps(ctx context.Context, ipType, vmIdentifier string) error
}

type iPsServiceHandler struct {
	client *Client
}

var _ IPsService = &iPsServiceHandler{}

type ListIPsRoot struct {
	Error bool `json:"error"`
	Data  []IP `json:"data"`
	Total int  `json:"total"`
}

type IP struct {
	ID            int    `json:"id"`
	DcName        string `json:"dcName"`
	DcIdentifier  string `json:"dcIdentifier"`
	RangeID       int    `json:"range_id"`
	IP            string `json:"ip"`
	IPVersion     string `json:"ip_version"`
	IsPrimary     int    `json:"is_primary"`
	Hostname      string `json:"hostname"`
	BoxID         int    `json:"box_id"`
	BoxIdentifier string `json:"box_identifier"`
	FullName      string `json:"fullName"`
	Category      string `json:"category"`
	UserID        int    `json:"user_id"`
	OwnerID       int    `json:"owner_id"`
	CreatedBy     string `json:"created_by"`
	Type          string `json:"type"`
	UpdatedAt     string `json:"updated_at"`
}

func (i *iPsServiceHandler) ListPrivateIPs(ctx context.Context, options *ListOptions) ([]IP, error) {
	path := fmt.Sprintf("%s/private", ipsPath)

	req, err := i.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	ips := new(ListIPsRoot)
	if err = i.client.Do(ctx, req, ips); err != nil {
		return nil, err
	}

	return ips.Data, nil
}

func (i *iPsServiceHandler) ListPublicIPs(ctx context.Context, options *ListOptions) ([]IP, error) {
	if options == nil {
		options = &ListOptions{}
	}
	path := fmt.Sprintf("%s/public?offset=%d&limit=%d", ipsPath, options.Page, options.PerPage)

	req, err := i.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	ips := new(ListIPsRoot)
	if err = i.client.Do(ctx, req, ips); err != nil {
		return nil, err
	}
	return ips.Data, nil
}

func (i *iPsServiceHandler) ListAllIPs(ctx context.Context, options *ListOptions) ([]IP, error) {
	if options == nil {
		options = &ListOptions{}
	}
	path := fmt.Sprintf("%s?offset=%d&limit=%d", ipsPath, options.Page, options.PerPage)

	req, err := i.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	ips := new(ListIPsRoot)
	if err = i.client.Do(ctx, req, ips); err != nil {
		return nil, err

	}
	return ips.Data, nil
}

func (i *iPsServiceHandler) DeleteIP(ctx context.Context, ip, vmIdentifier string) error {
	path := fmt.Sprintf("%s/delete", ipsPath)

	deleteRequest := struct {
		IP           string `json:"ip"`
		VmIdentifier string `json:"vmIdentifier"`
	}{
		IP:           ip,
		VmIdentifier: vmIdentifier,
	}

	req, err := i.client.NewRequest(ctx, http.MethodPost, path, &deleteRequest)
	if err != nil {
		return err
	}
	return i.client.Do(ctx, req, nil)
}

func (i *iPsServiceHandler) CreateIps(ctx context.Context, ipType, vmIdentifier string) error {
	path := fmt.Sprintf("%s/add", ipsPath)

	createRequest := struct {
		IPType       string `json:"ipType"`
		VMIdentifier string `json:"vmIdentifier"`
	}{
		IPType:       ipType,
		VMIdentifier: vmIdentifier,
	}

	req, err := i.client.NewRequest(ctx, http.MethodPost, path, &createRequest)
	if err != nil {
		return err
	}
	return i.client.Do(ctx, req, nil)
}
