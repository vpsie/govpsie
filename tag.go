package govpsie

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

var tagBasePath = "/apps/v2/tags"

// TagsService manages resource tags. In the VPSie API a tag is a label attached
// to an entity (a server, VPC, storage, ssh key, …); tags are not standalone
// objects, they are applied to a specific resource identified by its entity type
// and identifier.
type TagsService interface {
	// List returns all tag definitions visible to the user.
	List(ctx context.Context) ([]Tag, error)
	// ListForEntity returns the tag names applied to a specific entity.
	ListForEntity(ctx context.Context, entity, identifier string) ([]string, error)
	// AddToEntity applies the given tag names to an entity.
	AddToEntity(ctx context.Context, entity, identifier string, tags []string) error
	// Edit replaces the set of tag names applied to an entity.
	Edit(ctx context.Context, entity, identifier string, tags []string) error
	// Delete removes all tags applied to an entity.
	Delete(ctx context.Context, entity, identifier string) error
}

type tagsServiceHandler struct {
	client *Client
}

var _ TagsService = &tagsServiceHandler{}

// Tag represents a tag applied to a resource. The API returns applied tags
// grouped by entity type; each row carries the tag name and the entity it is
// attached to.
type Tag struct {
	Name       string `json:"name"`
	EntityType string `json:"entity_type"`
	EntityID   int64  `json:"entity_id"`
}

// tagsRoot matches the API response, whose `data` is an object keyed by entity
// group, each value a slice of applied tags.
type tagsRoot struct {
	Error bool             `json:"error"`
	Data  map[string][]Tag `json:"data"`
}

func flattenTags(groups map[string][]Tag) []Tag {
	var out []Tag
	for _, group := range groups {
		out = append(out, group...)
	}

	return out
}

func (s *tagsServiceHandler) List(ctx context.Context) ([]Tag, error) {
	req, err := s.client.NewRequest(ctx, http.MethodGet, tagBasePath, nil)
	if err != nil {
		return nil, err
	}

	root := new(tagsRoot)
	if err = s.client.Do(ctx, req, root); err != nil {
		return nil, err
	}

	return flattenTags(root.Data), nil
}

func (s *tagsServiceHandler) ListForEntity(ctx context.Context, entity, identifier string) ([]string, error) {
	query := url.Values{}
	query.Set("entity", entity)
	query.Set("identifier", identifier)
	path := fmt.Sprintf("%s?%s", tagBasePath, query.Encode())

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	root := new(tagsRoot)
	if err = s.client.Do(ctx, req, root); err != nil {
		return nil, err
	}

	// The query filters by entity+identifier, so only the matching group is
	// populated; flatten all groups to collect the applied tag names.
	var names []string
	for _, t := range flattenTags(root.Data) {
		names = append(names, t.Name)
	}

	return names, nil
}

func (s *tagsServiceHandler) AddToEntity(ctx context.Context, entity, identifier string, tags []string) error {
	path := fmt.Sprintf("%s/add", tagBasePath)

	addReq := struct {
		Entity     string   `json:"entity"`
		Identifier string   `json:"identifier"`
		Tags       []string `json:"tags"`
	}{
		Entity:     entity,
		Identifier: identifier,
		Tags:       tags,
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, &addReq)
	if err != nil {
		return err
	}

	return s.client.Do(ctx, req, nil)
}

func (s *tagsServiceHandler) Edit(ctx context.Context, entity, identifier string, tags []string) error {
	path := fmt.Sprintf("%s/edit", tagBasePath)

	editReq := struct {
		Entity     string   `json:"entity"`
		Identifier string   `json:"identifier"`
		Tags       []string `json:"tags"`
	}{
		Entity:     entity,
		Identifier: identifier,
		Tags:       tags,
	}

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, &editReq)
	if err != nil {
		return err
	}

	return s.client.Do(ctx, req, nil)
}

func (s *tagsServiceHandler) Delete(ctx context.Context, entity, identifier string) error {
	path := fmt.Sprintf("%s/delete", tagBasePath)

	delReq := struct {
		Entity     string `json:"entity"`
		Identifier string `json:"identifier"`
	}{
		Entity:     entity,
		Identifier: identifier,
	}

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, &delReq)
	if err != nil {
		return err
	}

	return s.client.Do(ctx, req, nil)
}
