package govpsie

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

var projectsBasePath = "/apps/v2/projects"

type ProjectsService interface {
	List(context.Context, *ListOptions) ([]Project, error)
	SetDefault(context.Context, string) error
	Get(ctx context.Context, identifer string) (*Project, error)
	Create(context.Context, *CreateProjectRequest) error
	ListAnotherVms(context.Context, string) ([]VmData, error)
	MoveVms(context.Context, string, string) error
	AssignToVms(ctx context.Context, projectIdentifier, projectId string) error
	ListDomains(ctx context.Context, projectIdentifier string) ([]Domain, error)
	Delete(ctx context.Context, id string) error
	ListUserLimits(ctx context.Context) (*UserLimit, error)
}

type projectsServiceHandler struct {
	client *Client
}

var _ ProjectsService = &projectsServiceHandler{}

type Project struct {
	ID          uint64 `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatedOn   string `json:"created_on"`
	UpdatedAt   string `json:"updated_at"`
	Identifier  string `json:"identifier"`
	CreatedBy   uint64 `json:"created_by"`
	IsDefault   int    `json:"is_default"`
}

// CreateProjectRequest represents the request to create a new project.
type CreateProjectRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ProjectsRoot struct {
	Error bool `json:"error"`
	Data  Data `json:"data"`
}

type Data struct {
	Rows  []Project `json:"rows"`
	Count int       `json:"count"`
}

type ProjectRoot struct {
	Error   bool     `json:"error"`
	Project *Project `json:"data"`
}

func (r *ProjectRoot) UnmarshalJSON(b []byte) error {
	var raw struct {
		Error bool            `json:"error"`
		Data  json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	r.Error = raw.Error
	if len(raw.Data) == 0 || string(raw.Data) == "false" || string(raw.Data) == "null" {
		return nil
	}
	var p Project
	if err := json.Unmarshal(raw.Data, &p); err != nil {
		return err
	}
	r.Project = &p
	return nil
}

type ListUserLimitRoot struct {
	Error bool      `json:"error"`
	Data  UserLimit `json:"data"`
}

type UserLimit struct {
	BackupsLimit      int `json:"backups_Limit"`
	SnapshotLimit     int `json:"snapshot_Limit"`
	FirewallLimit     int `json:"firewall_Limit"`
	BucketsLimit      int `json:"buckets_Limit"`
	CertificatesLimit int `json:"certificates_Limit"`
	DNSDomainsLimit   int `json:"dns_domains_Limit"`
}

func (p *projectsServiceHandler) List(ctx context.Context, options *ListOptions) ([]Project, error) {
	req, err := p.client.NewRequest(ctx, http.MethodGet, projectsBasePath, nil)
	if err != nil {
		return nil, err
	}

	projects := new(ProjectsRoot)
	if err = p.client.Do(ctx, req, projects); err != nil {
		return nil, err
	}

	return projects.Data.Rows, nil
}

func (p *projectsServiceHandler) SetDefault(ctx context.Context, projectIdentifier string) error {
	path := fmt.Sprintf("%s/set/default", projectsBasePath)

	projectReq := struct {
		ProjectIdentifier string `json:"projectIdentifier"`
	}{
		ProjectIdentifier: projectIdentifier,
	}
	req, err := p.client.NewRequest(ctx, http.MethodPost, path, &projectReq)
	if err != nil {
		return err
	}

	return p.client.Do(ctx, req, nil)
}

func (p *projectsServiceHandler) Get(ctx context.Context, identifer string) (*Project, error) {
	path := fmt.Sprintf("%s/%s", projectsBasePath, identifer)
	req, err := p.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	project := new(ProjectRoot)
	if err := p.client.Do(ctx, req, project); err != nil {
		return nil, err
	}

	if project.Project == nil {
		return nil, fmt.Errorf("project not found")
	}

	return project.Project, nil
}

func (p *projectsServiceHandler) Create(ctx context.Context, projectReq *CreateProjectRequest) error {
	path := fmt.Sprintf("%s/add", projectsBasePath)

	req, err := p.client.NewRequest(ctx, http.MethodPost, path, projectReq)
	if err != nil {
		return err
	}
	return p.client.Do(ctx, req, nil)

}

func (p *projectsServiceHandler) ListAnotherVms(ctx context.Context, projectId string) ([]VmData, error) {
	path := fmt.Sprintf("%s/another/vms?projectId=%s", projectsBasePath, projectId)

	vms := new(ListServerRoot)
	req, err := p.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	if err = p.client.Do(ctx, req, vms); err != nil {
		return nil, err
	}

	return vms.Data, nil
}

func (p *projectsServiceHandler) MoveVms(ctx context.Context, projectIdentifier, projectId string) error {
	path := fmt.Sprintf("%s/move/vms", projectsBasePath)

	moveReq := struct {
		VmsIdentifiers []string `json:"vmsIdentifiers"`
		ProjectId      string   `json:"projectId"`
	}{
		VmsIdentifiers: []string{projectIdentifier},
		ProjectId:      projectId,
	}
	req, err := p.client.NewRequest(ctx, http.MethodPost, path, &moveReq)
	if err != nil {
		return err
	}

	return p.client.Do(ctx, req, nil)
}

func (p *projectsServiceHandler) AssignToVms(ctx context.Context, projectIdentifier, projectId string) error {
	path := fmt.Sprintf("%s/vm", projectsBasePath)

	assignReq := struct {
		VmIdentifier string `json:"vmIdentifier"`
		ProjectId    string `json:"projectId"`
	}{
		VmIdentifier: projectIdentifier,
		ProjectId:    projectId,
	}

	req, err := p.client.NewRequest(ctx, http.MethodPost, path, &assignReq)
	if err != nil {
		return err
	}

	return p.client.Do(ctx, req, nil)
}

func (p *projectsServiceHandler) ListDomains(ctx context.Context, projectIdentifier string) ([]Domain, error) {
	path := fmt.Sprintf("/apps/v2/domains/project/%s", projectIdentifier)

	domains := new(ListDomainRoot)
	req, err := p.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	if err = p.client.Do(ctx, req, domains); err != nil {
		return nil, err
	}

	return domains.Data, nil
}

func (p *projectsServiceHandler) Delete(ctx context.Context, id string) error {
	path := fmt.Sprintf("%s/%s", projectsBasePath, id)

	req, err := p.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return err
	}

	return p.client.Do(ctx, req, nil)
}

func (p *projectsServiceHandler) ListUserLimits(ctx context.Context) (*UserLimit, error) {
	path := "/apps/v2/profile/product/limits"

	req, err := p.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	limits := new(ListUserLimitRoot)

	if err = p.client.Do(ctx, req, limits); err != nil {
		return nil, err
	}

	return &limits.Data, nil
}
