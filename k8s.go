package govpsie

import (
	"context"
	"fmt"
	"net/http"
)

var k8sPath = "/apps/v2/k8s"

type K8sService interface {
	List(ctx context.Context, options *ListOptions) ([]ListK8s, error)
	Delete(ctx context.Context, identifier, reason, note string) error
	Create(ctx context.Context, createReq *CreateK8sReq) error
	Get(ctx context.Context, identifier string) (*K8s, error)
	AddSlave(ctx context.Context, identifier string) error
	RemoveSlave(ctx context.Context, identifier string) error
}

type k8sServiceHandler struct {
	client *Client
}

var _ K8sService = &k8sServiceHandler{}

type ListK8sRoot struct {
	Error bool `json:"error"`
	Data  []ListK8s `json:"data"`
}

type GetK8sRoot struct {
	Error bool `json:"error"`
	Data K8s `json:"data"`
}

type K8s struct {
	ClusterName string `json:"cluster_name"`
	Identifier string `json:"identifier"`
	Count int `json:"count"`
	Nodes []Node `json:"nodes"`
	CreatedOn string `json:"created_on"`
	UpdatedOn string `json:"updated_on"`
	CreatedBy string `json:"created_by"`
	NickName string `json:"nickname"`
	Cpu int `json:"cpu"`
	Ram int `json:"ram"`
	Traffic int `json:"traffic"`
	Color string `json:"color"`
	Price float64 `json:"price"`
}

type ListK8s struct {
	ClusterName string `json:"cluster_name"`
	Identifier string `json:"identifier"`
	Count int `json:"count"`
	CreatedOn string `json:"created_on"`
	UpdatedOn string `json:"updated_on"`
	CreatedBy string `json:"created_by"`
	NickName string `json:"nickname"`
	Cpu int `json:"cpu"`
	Ram int `json:"ram"`
	Traffic int `json:"traffic"`
	Color string `json:"color"`
	Price float64 `json:"price"`
	ManagerCount int `json:"managerCount"`
	SlaveCount int `json:"slaveCount"`
}

type Node struct {
	Id int `json:"id"`
	UserId int `json:"user_id"`
	HostName string `json:"hostname"`
	DefaultIP string `json:"default_ip"`
	PrivateIP string `json:"private_ip"`
	NodeType int `json:"node_type"`
	NodeId int `json:"node_id"`
	DatacenterId int `json:"datacenter_id"`
	CreatedOn string `json:"created_on"`
}

type CreateK8sReq struct {
	ClusterName string `json:"clusterName"`
	DcIdentifier string `json:"dcIdentifier"`
	NodesCountMaster int `json:"nodesCountMaster"`
	NodesCountSlave int `json:"nodesCountSlave"`
	VpcId int `json:"vpcId"`
	KuberVer int `json:"kuberVer"`
	ResourceIdentifier string `json:"resourceIdentifier"`
	ProjectIdentifier string `json:"projectIdentifier"`
}

func (s *k8sServiceHandler) List(ctx context.Context, options *ListOptions) ([]ListK8s, error) {
	path := fmt.Sprintf("%s/cluster/all", k8sPath)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var root ListK8sRoot
	if err := s.client.Do(ctx, req, &root); err != nil {
		return nil, err
	}

	return root.Data, nil
}

func (s *k8sServiceHandler) Delete(ctx context.Context, identifier, reason, note string) error {
	path := fmt.Sprintf("%s/cluster/byId/%s", k8sPath, identifier)


	deleteStat := struct {
		DeleteStatistic struct {
			Reason string `json:"reason"`
			Note string `json:"note"`
		} `json:"deleteStatistic"`
	}{
		DeleteStatistic: struct {
			Reason string `json:"reason"`
			Note string `json:"note"`
		}{
			Reason: reason,
			Note: note,
		},
	}

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, &deleteStat)
	if err != nil {
		return err
	}

	return s.client.Do(ctx, req, nil)
}

func (s *k8sServiceHandler) Create(ctx context.Context, createReq *CreateK8sReq) error {
	path := fmt.Sprintf("%s/create/cluster", k8sPath)

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, createReq)
	if err != nil {
		return err
	}

	return s.client.Do(ctx, req, nil)
}

func (s *k8sServiceHandler) Get(ctx context.Context, identifier string) (*K8s, error) {
	path := fmt.Sprintf("%s/cluster/byId/%s", k8sPath, identifier)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var root GetK8sRoot
	if err := s.client.Do(ctx, req, &root); err != nil {
		return nil, err
	}

	return &root.Data, nil
}

func (s *k8sServiceHandler) AddSlave(ctx context.Context, identifier string) error {
	path := fmt.Sprintf("%s/cluster/byId/%s/add/slave", k8sPath, identifier)

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, nil)
	if err != nil {
		return err
	}

	return s.client.Do(ctx, req, nil)
}

func (s *k8sServiceHandler) RemoveSlave(ctx context.Context, identifier string) error {
	path := fmt.Sprintf("%s/cluster/byId/%s/reduce", k8sPath, identifier)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return err
	}

	return s.client.Do(ctx, req, nil)
}

