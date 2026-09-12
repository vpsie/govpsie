package govpsie

import (
	"context"
	"fmt"
	"net/http"
)

var tagBasePath = "/apps/v2/tags"

// TagsService is the interface for managing resource tags.
type TagsService interface {
	List(ctx context.Context) ([]Tag, error)
	Get(ctx context.Context, identifier string) (*Tag, error)
	Create(ctx context.Context, name, color string) (string, error)
	Delete(ctx context.Context, entity, identifier string) error
	Edit(ctx context.Context, entity, identifier string, tags []string) error
}

type tagsServiceHandler struct {
	client *Client
}

var _ TagsService = &tagsServiceHandler{}

// Tag represents a resource tag.
type Tag struct {
	Identifier string `json:"identifier"`
	Name       string `json:"name"`
	Color      string `json:"color"`
}

type tagsListRoot struct {
	Error bool  `json:"error"`
	Data  []Tag `json:"data"`
}

type tagGetRoot struct {
	Error bool `json:"error"`
	Data  *Tag `json:"data"`
}

type tagCreateRoot struct {
	Error bool `json:"error"`
	Data  struct {
		Identifier string `json:"identifier"`
	} `json:"data"`
}

func (s *tagsServiceHandler) List(ctx context.Context) ([]Tag, error) {
	req, err := s.client.NewRequest(ctx, http.MethodGet, tagBasePath, nil)
	if err != nil {
		return nil, err
	}

	root := new(tagsListRoot)
	if err = s.client.Do(ctx, req, root); err != nil {
		return nil, err
	}

	return root.Data, nil
}

func (s *tagsServiceHandler) Get(ctx context.Context, identifier string) (*Tag, error) {
	path := fmt.Sprintf("%s/%s", tagBasePath, identifier)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	root := new(tagGetRoot)
	if err = s.client.Do(ctx, req, root); err != nil {
		return nil, err
	}

	return root.Data, nil
}

// Create adds a new resource tag and returns its identifier.
func (s *tagsServiceHandler) Create(ctx context.Context, name, color string) (string, error) {
	path := fmt.Sprintf("%s/add", tagBasePath)

	createReq := struct {
		Name  string `json:"name"`
		Color string `json:"color"`
	}{
		Name:  name,
		Color: color,
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, &createReq)
	if err != nil {
		return "", err
	}

	root := new(tagCreateRoot)
	if err = s.client.Do(ctx, req, root); err != nil {
		return "", err
	}

	return root.Data.Identifier, nil
}

// Delete detaches/removes a tag from a resource entity.
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

// Edit replaces the set of tags attached to a resource entity.
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
