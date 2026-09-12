package govpsie

import (
	"context"
	"crypto/rand"
	"fmt"
	"net/http"
)

// newProcessID returns a random RFC 4122 version 4 UUID string. The managed
// database endpoints require a client-supplied processId to track long-running
// operations; a fresh identifier is generated per call.
func newProcessID() (string, error) {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", err
	}
	b[6] = (b[6] & 0x0f) | 0x40 // version 4
	b[8] = (b[8] & 0x3f) | 0x80 // variant 10
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16]), nil
}

var managedDBBasePath = "/apps/v2/managed/db"

// ManagedDBService is the interface for managing managed database clusters.
type ManagedDBService interface {
	List(ctx context.Context) ([]ManagedDB, error)
	Get(ctx context.Context, clusterID string) (*ManagedDBDetails, error)
	Create(ctx context.Context, createReq *CreateManagedDBRequest) error
	Delete(ctx context.Context, identifier, reason string) error
	AddNode(ctx context.Context, identifier string) error
	ReduceNode(ctx context.Context, identifier string) error
	ListOffers(ctx context.Context, dcIdentifier string) ([]ManagedDBOffer, error)
}

type managedDBServiceHandler struct {
	client *Client
}

var _ ManagedDBService = &managedDBServiceHandler{}

// ManagedDB is a managed database cluster summary from the list endpoint.
type ManagedDB struct {
	Identifier  string `json:"identifier"`
	ClusterName string `json:"cluster_name"`
	Nickname    string `json:"nickname"`
	UserID      int64  `json:"user_id"`
	NodesCount  int64  `json:"nodes_count"`
	NodeCount   int64  `json:"nodeCount"`
	CPU         int64  `json:"cpu"`
	RAM         int64  `json:"ram"`
	Traffic     int64  `json:"traffic"`
	Color       string `json:"color"`
	Price       int64  `json:"price"`
	CreatedBy   string `json:"created_by"`
	CreatedOn   string `json:"created_on"`
	UpdatedOn   string `json:"updated_on"`
}

// ManagedDBNode is a single node within a managed database cluster.
type ManagedDBNode struct {
	ID           int64  `json:"id"`
	Identifier   string `json:"identifier"`
	Hostname     string `json:"hostname"`
	PrivateIP    string `json:"private_ip"`
	DefaultIP    string `json:"default_ip"`
	RAM          int64  `json:"ram"`
	CPU          int64  `json:"cpu"`
	SSD          int64  `json:"ssd"`
	Traffic      int64  `json:"traffic"`
	MainNode     int64  `json:"main_node"`
	Power        int64  `json:"power"`
	IsActive     int64  `json:"is_active"`
	IsDeleted    int64  `json:"is_deleted"`
	IsLocked     int64  `json:"is_locked"`
	IsSuspended  int64  `json:"is_suspended"`
	IsTerminated int64  `json:"is_terminated"`
	CreatedOn    string `json:"created_on"`
	LastUpdated  string `json:"last_updated"`
}

// ManagedDBDetails is the detailed view of a managed database cluster.
type ManagedDBDetails struct {
	ClusterID     int64           `json:"clusterId"`
	ClusterName   string          `json:"cluster_name"`
	Identifier    string          `json:"identifier"`
	Nickname      string          `json:"nickname"`
	NodesCount    int64           `json:"nodes_count"`
	AdminPassword string          `json:"admin_password"`
	CPU           int64           `json:"cpu"`
	RAM           int64           `json:"ram"`
	Traffic       int64           `json:"traffic"`
	Color         string          `json:"color"`
	Price         int64           `json:"price"`
	CreatedBy     string          `json:"created_by"`
	CreatedOn     string          `json:"created_on"`
	UpdatedOn     string          `json:"updated_on"`
	Nodes         []ManagedDBNode `json:"nodes"`
}

// ManagedDBOffer is a plan/offer available for a managed database in a datacenter.
type ManagedDBOffer struct {
	Identifier string `json:"identifier"`
	Name       string `json:"name"`
	PlanID     int64  `json:"planId"`
}

// CreateManagedDBRequest is the payload for creating a managed database cluster.
type CreateManagedDBRequest struct {
	Name         string `json:"name"`
	DBType       string `json:"dbType"`
	DcIdentifier string `json:"dcIdentifier"`
	NodeCount    int64  `json:"nodeCount"`
	PlanID       int64  `json:"planId"`
	ProjectID    string `json:"projectId,omitempty"`
}

type managedDBListRoot struct {
	Error bool        `json:"error"`
	Data  []ManagedDB `json:"data"`
	Total int64       `json:"total"`
}

type managedDBGetRoot struct {
	Error bool              `json:"error"`
	Data  *ManagedDBDetails `json:"data"`
}

type managedDBOffersRoot struct {
	Error bool             `json:"error"`
	Data  []ManagedDBOffer `json:"data"`
}

func (s *managedDBServiceHandler) List(ctx context.Context) ([]ManagedDB, error) {
	path := fmt.Sprintf("%s/cluster/all", managedDBBasePath)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	root := new(managedDBListRoot)
	if err = s.client.Do(ctx, req, root); err != nil {
		return nil, err
	}

	return root.Data, nil
}

func (s *managedDBServiceHandler) Get(ctx context.Context, clusterID string) (*ManagedDBDetails, error) {
	path := fmt.Sprintf("%s/cluster/byId/%s", managedDBBasePath, clusterID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	root := new(managedDBGetRoot)
	if err = s.client.Do(ctx, req, root); err != nil {
		return nil, err
	}

	return root.Data, nil
}

func (s *managedDBServiceHandler) Create(ctx context.Context, createReq *CreateManagedDBRequest) error {
	path := fmt.Sprintf("%s/create/cluster", managedDBBasePath)

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, createReq)
	if err != nil {
		return err
	}

	return s.client.Do(ctx, req, nil)
}

func (s *managedDBServiceHandler) Delete(ctx context.Context, identifier, reason string) error {
	path := fmt.Sprintf("%s/cluster/byId/%s", managedDBBasePath, identifier)

	processID, err := newProcessID()
	if err != nil {
		return err
	}

	delReq := struct {
		MdbIdentifier   string `json:"mdbIdentifier"`
		ProcessID       string `json:"processId"`
		DeleteStatistic struct {
			Reason string `json:"reason"`
		} `json:"deleteStatistic"`
	}{
		MdbIdentifier: identifier,
		ProcessID:     processID,
	}
	delReq.DeleteStatistic.Reason = reason

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, &delReq)
	if err != nil {
		return err
	}

	return s.client.Do(ctx, req, nil)
}

// AddNode adds a single node to a managed database cluster.
func (s *managedDBServiceHandler) AddNode(ctx context.Context, identifier string) error {
	path := fmt.Sprintf("%s/cluster/byId/%s/add", managedDBBasePath, identifier)

	processID, err := newProcessID()
	if err != nil {
		return err
	}

	addReq := struct {
		CreateFromPool string `json:"createFromPool"`
		ProcessID      string `json:"processId"`
	}{
		CreateFromPool: "0",
		ProcessID:      processID,
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, &addReq)
	if err != nil {
		return err
	}

	return s.client.Do(ctx, req, nil)
}

// ReduceNode removes a single node from a managed database cluster.
func (s *managedDBServiceHandler) ReduceNode(ctx context.Context, identifier string) error {
	path := fmt.Sprintf("%s/cluster/byId/%s/reduce", managedDBBasePath, identifier)

	processID, err := newProcessID()
	if err != nil {
		return err
	}

	reduceReq := struct {
		ProcessID string `json:"processId"`
	}{
		ProcessID: processID,
	}

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, &reduceReq)
	if err != nil {
		return err
	}

	return s.client.Do(ctx, req, nil)
}

func (s *managedDBServiceHandler) ListOffers(ctx context.Context, dcIdentifier string) ([]ManagedDBOffer, error) {
	path := fmt.Sprintf("%s/offers", managedDBBasePath)

	offersReq := struct {
		DcIdentifier string `json:"dcIdentifier"`
	}{
		DcIdentifier: dcIdentifier,
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, &offersReq)
	if err != nil {
		return nil, err
	}

	root := new(managedDBOffersRoot)
	if err = s.client.Do(ctx, req, root); err != nil {
		return nil, err
	}

	return root.Data, nil
}
