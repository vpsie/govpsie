package govpsie

import (
	"context"
	"fmt"
	"net/http"
)

var serverGroupBasePath = "/apps/v2/vm/group"

// ServerGroupService is the interface for managing server (VM) groups.
type ServerGroupService interface {
	List(ctx context.Context) ([]ServerGroup, error)
	Get(ctx context.Context, identifier string) (*ServerGroupDetails, error)
	Create(ctx context.Context, name, description string, isDistributed bool) error
	Delete(ctx context.Context, identifier string) error
	Join(ctx context.Context, groupIdentifier, vmIdentifier string) error
	Detach(ctx context.Context, groupIdentifier, vmIdentifier string) error
}

type serverGroupServiceHandler struct {
	client *Client
}

var _ ServerGroupService = &serverGroupServiceHandler{}

// ServerGroup represents a server group summary as returned by the list endpoint.
type ServerGroup struct {
	ID               int64  `json:"id"`
	Identifier       string `json:"identifier"`
	GroupName        string `json:"group_name"`
	GroupDescription string `json:"group_description"`
	IsDistributed    bool   `json:"is_distributed"`
	IsDeleted        bool   `json:"is_deleted"`
	UserID           int64  `json:"user_id"`
	VMCounts         int64  `json:"vmCounts"`
	CreatedOn        string `json:"created_on"`
	UpdatedAt        string `json:"updatedAt"`
}

// ServerGroupDetails is the detailed view of a single group with its attached VMs.
type ServerGroupDetails struct {
	GroupDetails struct {
		ID               int64  `json:"id"`
		Identifier       string `json:"identifier"`
		GroupName        string `json:"group_name"`
		GroupDescription string `json:"group_description"`
		IsDistributed    int64  `json:"is_distributed"`
		IsDeleted        int64  `json:"is_deleted"`
		UserID           int64  `json:"user_id"`
		CreatedOn        string `json:"created_on"`
		UpdatedOn        string `json:"updated_on"`
	} `json:"group_details"`
	AttachedVMs []ServerGroupVM `json:"attached_vms"`
}

// ServerGroupVM is a VM attached to a server group. Fields are mapped
// leniently; only the identifier is relied upon by callers.
type ServerGroupVM struct {
	Identifier string `json:"identifier"`
	Hostname   string `json:"hostname"`
}

type serverGroupListRoot struct {
	Error bool `json:"error"`
	Data  struct {
		Rows  []ServerGroup `json:"rows"`
		Total int64         `json:"total"`
	} `json:"data"`
}

type serverGroupGetRoot struct {
	Error bool                `json:"error"`
	Data  *ServerGroupDetails `json:"data"`
}

func (s *serverGroupServiceHandler) List(ctx context.Context) ([]ServerGroup, error) {
	req, err := s.client.NewRequest(ctx, http.MethodGet, serverGroupBasePath, nil)
	if err != nil {
		return nil, err
	}

	root := new(serverGroupListRoot)
	if err = s.client.Do(ctx, req, root); err != nil {
		return nil, err
	}

	return root.Data.Rows, nil
}

func (s *serverGroupServiceHandler) Get(ctx context.Context, identifier string) (*ServerGroupDetails, error) {
	path := fmt.Sprintf("%s/%s", serverGroupBasePath, identifier)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	root := new(serverGroupGetRoot)
	if err = s.client.Do(ctx, req, root); err != nil {
		return nil, err
	}

	return root.Data, nil
}

func (s *serverGroupServiceHandler) Create(ctx context.Context, name, description string, isDistributed bool) error {
	path := fmt.Sprintf("%s/new", serverGroupBasePath)

	createReq := struct {
		GroupName     string `json:"groupName"`
		GroupDesc     string `json:"groupDesc"`
		IsDistributed bool   `json:"isDistributed"`
	}{
		GroupName:     name,
		GroupDesc:     description,
		IsDistributed: isDistributed,
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, &createReq)
	if err != nil {
		return err
	}

	return s.client.Do(ctx, req, nil)
}

func (s *serverGroupServiceHandler) Delete(ctx context.Context, identifier string) error {
	path := fmt.Sprintf("%s/%s", serverGroupBasePath, identifier)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return err
	}

	return s.client.Do(ctx, req, nil)
}

func (s *serverGroupServiceHandler) Join(ctx context.Context, groupIdentifier, vmIdentifier string) error {
	path := fmt.Sprintf("%s/join", serverGroupBasePath)

	joinReq := struct {
		GroupIdentifier string `json:"groupIdentifier"`
		VMIdentifier    string `json:"vmIdentifier"`
	}{
		GroupIdentifier: groupIdentifier,
		VMIdentifier:    vmIdentifier,
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, &joinReq)
	if err != nil {
		return err
	}

	return s.client.Do(ctx, req, nil)
}

func (s *serverGroupServiceHandler) Detach(ctx context.Context, groupIdentifier, vmIdentifier string) error {
	path := fmt.Sprintf("%s/detach", serverGroupBasePath)

	detachReq := struct {
		GroupIdentifier string `json:"groupIdentifier"`
		VMIdentifier    string `json:"vmIdentifier"`
	}{
		GroupIdentifier: groupIdentifier,
		VMIdentifier:    vmIdentifier,
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, &detachReq)
	if err != nil {
		return err
	}

	return s.client.Do(ctx, req, nil)
}
