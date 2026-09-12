package govpsie

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

var gatewayPath = "/apps/v2/gateways"

type GatewayService interface {
	List(ctx context.Context, options *ListOptions) ([]Gateway, error)
	Delete(ctx context.Context, ipId int) error
	Create(ctx context.Context, createReq *CreateGatewayReq) error
	Get(ctx context.Context, id int64) (*Gateway, error)
	AttachVM(ctx context.Context, id int64, vms []string, ignoreLegacyVms int64) error
	DetachVM(ctx context.Context, id int64, mapping_id []int64) error
}

type gatewayServiceHandler struct {
	client *Client
}

var _ GatewayService = &gatewayServiceHandler{}

type ListGatewayRoot struct {
	Error bool `json:"error"`
	Data  struct {
		Rows                    []Gateway `json:"rows"`
		IsThereUserGatewayLevel bool      `json:"isThereUserGatewayLevel"`
		Count                   int64     `json:"count"`
	} `json:"data"`
}

type GetGatewayRoot struct {
	Error bool `json:"error"`
	Data  struct {
		Rows                    []Gateway `json:"rows"`
		IsThereUserGatewayLevel bool      `json:"isThereUserGatewayLevel"`
		Count                   int64     `json:"count"`
	} `json:"data"`
}
type Gateway struct {
	ID                   int64        `json:"id"`
	DatacenterID         int64        `json:"datacenter_id"`
	IPPropertiesID       int64        `json:"ip_properties_id"`
	IP                   string       `json:"ip"`
	IsReserved           int64        `json:"is_reserved"`
	IPVersion            string       `json:"ip_version"`
	BoxID                *int64       `json:"box_id,omitempty"`
	IsPrimary            int64        `json:"is_primary"`
	Notes                *string      `json:"notes,omitempty"`
	UserID               int64        `json:"user_id"`
	UpdatedAt            time.Time    `json:"updated_at"`
	IsGatewayReserved    int64        `json:"is_gateway_reserved"`
	IsUserAccountGateway int64        `json:"is_user_account_gateway"`
	DatacenterName       string       `json:"datacenterName"`
	State                string       `json:"state"`
	DcIdentifier         string       `json:"dcIdentifier"`
	CreatedBy            string       `json:"created_by"`
	AttachedVms          []AttachedVM `json:"attachedVms"`
}

type AttachedVM struct {
	Identifier       string `json:"identifier"`
	GatewayMappingID int64  `json:"gateway_mapping_id"`
}

type CreateGatewayReq struct {
	IPType       string   `json:"ipType"`
	Tags         []string `json:"tags,omitempty"`
	DcIdentifier string   `json:"dcIdentifier"`
}

func (s *gatewayServiceHandler) List(ctx context.Context, options *ListOptions) ([]Gateway, error) {
	path := fmt.Sprintf("%s/ips", gatewayPath)
	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var root ListGatewayRoot
	if err := s.client.Do(ctx, req, &root); err != nil {
		return nil, err
	}

	return root.Data.Rows, nil
}

func (s *gatewayServiceHandler) Create(ctx context.Context, createReq *CreateGatewayReq) error {
	path := fmt.Sprintf("%s/add/ip", gatewayPath)
	req, err := s.client.NewRequest(ctx, http.MethodPost, path, createReq)
	if err != nil {
		return err
	}

	return s.client.Do(ctx, req, nil)
}

func (s *gatewayServiceHandler) Delete(ctx context.Context, ipId int) error {
	path := fmt.Sprintf("%s/delete/ip", gatewayPath)

	delReq := struct {
		IPID int `json:"ipId"`
	}{
		IPID: ipId,
	}

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, delReq)
	if err != nil {
		return err
	}

	return s.client.Do(ctx, req, nil)
}

func (g *gatewayServiceHandler) Get(ctx context.Context, id int64) (*Gateway, error) {
	path := fmt.Sprintf("%s/ips?ipId=%d", gatewayPath, id)
	req, err := g.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var root GetGatewayRoot

	if err := g.client.Do(ctx, req, &root); err != nil {
		return nil, err
	}

	if len(root.Data.Rows) == 0 {
		return nil, fmt.Errorf("gateway with id %d not found", id)
	}

	return &root.Data.Rows[0], nil
}

func (g *gatewayServiceHandler) AttachVM(ctx context.Context, id int64, vms []string, ignoreLegacyVms int64) error {
	path := fmt.Sprintf("%s/attach/vms", gatewayPath)

	attReq := struct {
		Vms             []string `json:"vms"`
		IpId            int64    `json:"ipId"`
		IgnoreLegacyVms int64    `json:"ignoreLegacyVms"`
	}{
		Vms:             vms,
		IpId:            id,
		IgnoreLegacyVms: ignoreLegacyVms,
	}

	req, err := g.client.NewRequest(ctx, http.MethodPost, path, attReq)
	if err != nil {
		return err
	}

	return g.client.Do(ctx, req, nil)
}

func (g *gatewayServiceHandler) DetachVM(ctx context.Context, id int64, mapping_id []int64) error {
	path := fmt.Sprintf("%s/detach/vms", gatewayPath)

	detReq := struct {
		MapIds []int64 `json:"mapIds"`
		IpId   int64   `json:"ipId"`
	}{
		MapIds: mapping_id,
		IpId:   id,
	}

	req, err := g.client.NewRequest(ctx, http.MethodPost, path, detReq)
	if err != nil {
		return err
	}

	return g.client.Do(ctx, req, nil)
}
