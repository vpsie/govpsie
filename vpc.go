package govpsie

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

var vpcPath = "/apps/v2"

type VPCService interface {
	List(ctx context.Context, options *ListOptions) ([]VPC, error)
	Get(ctx context.Context, id string) (*VPC, error)
	AssignServer(ctx context.Context, assignReq *AssignServerReq) error
	MoveServer(ctx context.Context, assignReq *AssignServerReq) error
	CreateVpc(ctx context.Context, createReq *CreateVpcReq) error
	ReleasePrivateIP(ctx context.Context, vmIdentifer string, privateIpId int) error
	ListServers(ctx context.Context, vpcID int) ([]VpcServer, error)
	DeleteVpc(ctx context.Context, vpcId, reason, note string) error
}

type vpcServiceHandler struct {
	client *Client
}

var _ VPCService = &vpcServiceHandler{}

type ListVPCRoot struct {
	Error bool  `json:"error"`
	Data  []VPC `json:"data"`
	Total int   `json:"total "`
}

type GetVPCRoot struct {
	Error      bool `json:"error"`
	Data       VPC  `json:"data"`
	VpcDetails struct {
		Rows  []interface{} `json:"rows"`
		Count int           `json:"count"`
	} `json:"vpcDetails"`
}

// VpcServer is a server attached to a VPC, as reported by
// GET /vpc/vms/{vpcId}. This is the authoritative source for a server's
// in-VPC address: VPC addresses are not returned by the private-IP listing.
type VpcServer struct {
	VmID           int    `json:"vmId"`
	Identifier     string `json:"identifier"`
	VpcIP          string `json:"vpc_ip"`
	Hostname       string `json:"hostname"`
	FullName       string `json:"fullname"`
	Category       string `json:"category"`
	EntityType     string `json:"entity_type"`
	IsDisconnected int    `json:"is_disconnected"`
	PrivateIPID    int    `json:"private_ip_id"`
}

type ListVpcServersRoot struct {
	Error bool        `json:"error"`
	Data  []VpcServer `json:"data"`
}

type AssignServerReq struct {
	VmIdentifier string `json:"vmIdentifier"`
	VpcID        int    `json:"vpcId"`
	DcIdentifier string `json:"dcIdentifier"`
}

type CreateVpcReq struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	DcIdentifier string `json:"dcIdentifier"`
	NetworkRange string `json:"networkRange"`
	NetworkSize  string `json:"networkSize"`
	AutoGenerate int    `json:"autoGenerate"`
}

type VPC struct {
	ID               int       `json:"id"`
	UserID           int       `json:"user_id"`
	OwnerID          int       `json:"owner_id"`
	DatacenterID     int       `json:"datacenter_id"`
	Name             string    `json:"name"`
	Description      string    `json:"description"`
	InterfaceNumber  int       `json:"interface_number"`
	NetworkTagNumber int       `json:"network_tag_number"`
	NetworkRange     string    `json:"network_range"`
	NetworkSize      int       `json:"network_size"`
	IsDefault        int       `json:"is_default"`
	CreatedBy        int       `json:"created_by"`
	UpdatedBy        int       `json:"updated_by"`
	CreatedOn        time.Time `json:"created_on"`
	LastUpdated      time.Time `json:"last_updated"`
	LowIPNum         int       `json:"low_ip_num"`
	HightIPNum       int       `json:"hight_ip_num"`
	IsUpcNetwork     int       `json:"is_upc_network"`
	Firstname        string    `json:"firstname"`
	Lastname         string    `json:"lastname"`
	Username         string    `json:"username"`
	State            string    `json:"state"`
	DcName           string    `json:"dc_name"`
	DcIdentifier     string    `json:"dc_identifier"`
}

func (s *vpcServiceHandler) List(ctx context.Context, options *ListOptions) ([]VPC, error) {
	path := fmt.Sprintf("%s/vpc", vpcPath)
	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	root := new(ListVPCRoot)
	if err := s.client.Do(ctx, req, &root); err != nil {
		return nil, err
	}

	return root.Data, nil
}

func (s *vpcServiceHandler) Get(ctx context.Context, id string) (*VPC, error) {
	path := fmt.Sprintf("%s/vpc/%s", vpcPath, id)
	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	root := new(GetVPCRoot)
	if err := s.client.Do(ctx, req, &root); err != nil {
		return nil, err
	}

	return &root.Data, nil
}

func (s *vpcServiceHandler) AssignServer(ctx context.Context, assignReq *AssignServerReq) error {
	path := fmt.Sprintf("%s/vm/add/vpc", vpcPath)

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, assignReq)
	if err != nil {
		return err
	}

	return s.client.Do(ctx, req, nil)
}

func (s *vpcServiceHandler) MoveServer(ctx context.Context, assignReq *AssignServerReq) error {
	path := fmt.Sprintf("%s/vm/vpc/move", vpcPath)

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, assignReq)
	if err != nil {
		return err
	}

	return s.client.Do(ctx, req, nil)
}

func (s *vpcServiceHandler) CreateVpc(ctx context.Context, createReq *CreateVpcReq) error {
	path := fmt.Sprintf("%s/vpc/add", vpcPath)

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, createReq)
	if err != nil {
		return err
	}

	return s.client.Do(ctx, req, nil)
}

func (s *vpcServiceHandler) ReleasePrivateIP(ctx context.Context, vmIdentifer string, privateIpId int) error {
	path := fmt.Sprintf("%s/vm/vpc", vpcPath)

	realseReq := struct {
		VmIdentifier string `json:"vmIdentifier"`
		PrivateIpId  int    `json:"privateIpId"`
	}{
		VmIdentifier: vmIdentifer,
		PrivateIpId:  privateIpId,
	}

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, realseReq)
	if err != nil {
		return err
	}

	return s.client.Do(ctx, req, nil)
}

func (s *vpcServiceHandler) DeleteVpc(ctx context.Context, vpcId, reason, note string) error {
	path := fmt.Sprintf("%s/vpc/%s", vpcPath, vpcId)

	deleteReq := struct {
		DeleteStatistic struct {
			Reason string `json:"reason"`
			Note   string `json:"note"`
		} `json:"deleteStatistic"`
	}{
		DeleteStatistic: struct {
			Reason string `json:"reason"`
			Note   string `json:"note"`
		}{
			Reason: reason,
			Note:   note,
		},
	}

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, &deleteReq)
	if err != nil {
		return err
	}

	return s.client.Do(ctx, req, nil)
}

// ListServers returns the servers currently attached to a VPC, including the
// private address and private-IP id each one was allocated inside it.
func (s *vpcServiceHandler) ListServers(ctx context.Context, vpcID int) ([]VpcServer, error) {
	path := fmt.Sprintf("%s/vpc/vms/%d", vpcPath, vpcID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	servers := new(ListVpcServersRoot)
	if err = s.client.Do(ctx, req, servers); err != nil {
		return nil, err
	}

	return servers.Data, nil
}
