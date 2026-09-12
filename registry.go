package govpsie

import (
	"context"
	"fmt"
	"net/http"
)

var registryBasePath = "/apps/v2/registries"

// RegistryService is the interface for managing container registries.
type RegistryService interface {
	List(ctx context.Context) ([]Registry, error)
	Create(ctx context.Context, name, datacenterID, resourceIdentifier, projectID string) error
	Delete(ctx context.Context, registryID string) error
	CheckName(ctx context.Context, name string) (bool, error)
	AssignProject(ctx context.Context, registryIdentifier, projectIdentifier string) error
	ListPlans(ctx context.Context) ([]RegistryPlan, error)
	ListDatacenters(ctx context.Context) ([]RegistryDatacenter, error)
}

type registryServiceHandler struct {
	client *Client
}

var _ RegistryService = &registryServiceHandler{}

// Registry represents a container registry. The list endpoint's response is
// only loosely specified upstream, so fields are mapped best-effort; unknown
// API fields are ignored rather than causing an error.
type Registry struct {
	Identifier     string `json:"registry_id"`
	Name           string `json:"registry_name"`
	DcIdentifier   string `json:"dcIdentifier"`
	DatacenterName string `json:"dc_name"`
	RegistryPlanID int64  `json:"registry_plan_id"`
	DatacenterID   int64  `json:"datacenter_id"`
	ProjectID      int64  `json:"project_id"`
	PlanIdentifier string `json:"planIdentifier"`
	Status         string `json:"status"`
	UserID         int64  `json:"user_id"`
	CreatedOn      string `json:"created_on"`
	UpdatedOn      string `json:"updated_at"`
	CreatedBy      string `json:"created_by"`
}

// RegistryPlan is a resource plan available for registries.
type RegistryPlan struct {
	Identifier string `json:"identifier"`
	Name       string `json:"name"`
}

// RegistryDatacenter is a datacenter available for registries.
type RegistryDatacenter struct {
	Identifier string `json:"identifier"`
	Name       string `json:"name"`
}

type registryListRoot struct {
	Error bool       `json:"error"`
	Data  []Registry `json:"data"`
}

type registryCheckRoot struct {
	Error bool `json:"error"`
	Data  bool `json:"data"`
}

type registryPlansRoot struct {
	Error bool           `json:"error"`
	Data  []RegistryPlan `json:"data"`
}

type registryDatacentersRoot struct {
	Error bool                 `json:"error"`
	Data  []RegistryDatacenter `json:"data"`
}

func (s *registryServiceHandler) List(ctx context.Context) ([]Registry, error) {
	path := fmt.Sprintf("%s/all", registryBasePath)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	root := new(registryListRoot)
	if err = s.client.Do(ctx, req, root); err != nil {
		return nil, err
	}

	return root.Data, nil
}

func (s *registryServiceHandler) Create(ctx context.Context, name, datacenterID, resourceIdentifier, projectID string) error {
	createReq := struct {
		Name               string `json:"name"`
		DatacenterID       string `json:"datacenterId"`
		ResourceIdentifier string `json:"resourceIdentifier"`
		ProjectID          string `json:"projectId"`
	}{
		Name:               name,
		DatacenterID:       datacenterID,
		ResourceIdentifier: resourceIdentifier,
		ProjectID:          projectID,
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, registryBasePath, &createReq)
	if err != nil {
		return err
	}

	return s.client.Do(ctx, req, nil)
}

func (s *registryServiceHandler) Delete(ctx context.Context, registryID string) error {
	path := fmt.Sprintf("%s/%s", registryBasePath, registryID)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return err
	}

	return s.client.Do(ctx, req, nil)
}

// CheckName reports whether a registry name is already taken. It returns true
// when the name already exists (is not available).
func (s *registryServiceHandler) CheckName(ctx context.Context, name string) (bool, error) {
	path := fmt.Sprintf("%s/check", registryBasePath)

	checkReq := struct {
		Name string `json:"name"`
	}{
		Name: name,
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, &checkReq)
	if err != nil {
		return false, err
	}

	root := new(registryCheckRoot)
	if err = s.client.Do(ctx, req, root); err != nil {
		return false, err
	}

	return root.Data, nil
}

// AssignProject assigns a registry to a project.
func (s *registryServiceHandler) AssignProject(ctx context.Context, registryIdentifier, projectIdentifier string) error {
	path := fmt.Sprintf("%s/project", registryBasePath)

	assignReq := struct {
		RegistryIdentifier string `json:"registryIdentifier"`
		ProjectIdentifier  string `json:"projectIdentifier"`
	}{
		RegistryIdentifier: registryIdentifier,
		ProjectIdentifier:  projectIdentifier,
	}

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, &assignReq)
	if err != nil {
		return err
	}

	return s.client.Do(ctx, req, nil)
}

func (s *registryServiceHandler) ListPlans(ctx context.Context) ([]RegistryPlan, error) {
	path := fmt.Sprintf("%s/resources/plan", registryBasePath)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	root := new(registryPlansRoot)
	if err = s.client.Do(ctx, req, root); err != nil {
		return nil, err
	}

	return root.Data, nil
}

func (s *registryServiceHandler) ListDatacenters(ctx context.Context) ([]RegistryDatacenter, error) {
	path := fmt.Sprintf("%s/datacenter", registryBasePath)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	root := new(registryDatacentersRoot)
	if err = s.client.Do(ctx, req, root); err != nil {
		return nil, err
	}

	return root.Data, nil
}
