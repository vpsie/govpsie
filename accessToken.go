package govpsie

import (
	"context"
	"fmt"
	"net/http"
)

var accessTokenBasePath = "/apps/v2/profile/security"

type AccessTokenService interface {
	List(ctx context.Context, options *ListOptions) ([]AccessToken, error)
	Create(ctx context.Context, name, accessToken, expirationDate, status string) error
	Delete(ctx context.Context, accessTokenIdentifier string) error
	Update(ctx context.Context, accessTokenIdentifier, name, expirationDate, status string) error
}

type accessTokenServiceHandler struct {
	client *Client
}

var _ AccessTokenService = &accessTokenServiceHandler{}

type AccessToken struct {
	AccessTokenIdentifier string `json:"identifier"`
	Name                  string `json:"name"`
	CreatedOn             string `json:"created_on"`
	ExpirationDate        string `json:"expiration_date"`
}

type ListAccessTokensRoot struct {
	Error bool          `json:"error"`
	Data  []AccessToken `json:"data"`
}

func (s *accessTokenServiceHandler) List(ctx context.Context, options *ListOptions) ([]AccessToken, error) {
	path := fmt.Sprintf("%s/access/token", accessTokenBasePath)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	accessTokens := new(ListAccessTokensRoot)
	if err := s.client.Do(ctx, req, &accessTokens); err != nil {
		return nil, err
	}

	return accessTokens.Data, nil
}

func (s *accessTokenServiceHandler) Create(ctx context.Context, name, accessToken, expirationDate, status string) error {
	path := fmt.Sprintf("%s/access/token", accessTokenBasePath)

	createAccessTokenReq := struct {
		AccessTokenName string `json:"accessTokenName"`
		AccessToken     string `json:"accessToken"`
		ExpirationDate  string `json:"expirationDate"`
		Status          string `json:"status"`
	}{
		AccessTokenName: name,
		AccessToken:     accessToken,
		ExpirationDate:  expirationDate,
		Status:          status,
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, &createAccessTokenReq)
	if err != nil {
		return err
	}

	return s.client.Do(ctx, req, nil)
}

func (s *accessTokenServiceHandler) Delete(ctx context.Context, accessTokenIdentifier string) error {
	path := fmt.Sprintf("%s/access/token/%s", accessTokenBasePath, accessTokenIdentifier)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return err
	}

	return s.client.Do(ctx, req, nil)
}

func (s *accessTokenServiceHandler) Update(ctx context.Context, accessTokenIdentifier, name, expirationDate, status string) error {
	path := fmt.Sprintf("%s/access/token/%s", accessTokenBasePath, accessTokenIdentifier)

	updateAccessTokenReq := struct {
		AccessTokenName string `json:"accessTokenName"`
		ExpirationDate  string `json:"expirationDate"`
		Status          string `json:"status"`
	}{
		AccessTokenName: name,
		ExpirationDate:  expirationDate,
		Status:          status,
	}

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, &updateAccessTokenReq)
	if err != nil {
		return err
	}

	return s.client.Do(ctx, req, nil)
}
