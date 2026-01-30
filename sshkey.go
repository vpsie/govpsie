package govpsie

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

var sshkeysBasePath = "/apps/v2/sshkeys"
var sshkeyBasePath = "/apps/v2/sshkey"

type SshkeysService interface {
	List(context.Context) ([]SShKey, error)
	Delete(context.Context, string) error
	Get(context.Context, string) (*SShKey, error)
	Create(context.Context, string, string) error
}

type sshkeysServiceHandler struct {
	client *Client
}

type SshKeysGetRoot struct {
	Error bool    `json:"error"`
	Data  *SShKey `json:"data"`
}

func (r *SshKeysGetRoot) UnmarshalJSON(b []byte) error {
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
	var k SShKey
	if err := json.Unmarshal(raw.Data, &k); err != nil {
		return err
	}
	r.Data = &k
	return nil
}

type SshKeysListRoot struct {
	Error bool     `json:"error"`
	Data  []SShKey `json:"data"`
	Total int64    `json:"total"`
}

type SShKey struct {
	Id         int64  `json:"id"`
	UserId     int64  `json:"user_id"`
	Name       string `json:"name"`
	PrivateKey string `json:"private_key"`
	CreatedOn  string `json:"created_on"`
	Identifier string `json:"identifier"`
	CreatedBy  string `json:"created_by"`
}

var _ SshkeysService = &sshkeysServiceHandler{}

func (s *sshkeysServiceHandler) List(ctx context.Context) ([]SShKey, error) {
	req, err := s.client.NewRequest(ctx, http.MethodGet, sshkeysBasePath, nil)
	if err != nil {
		return nil, err
	}

	sshKeys := new(SshKeysListRoot)
	if err = s.client.Do(ctx, req, sshKeys); err != nil {
		return nil, err
	}

	return sshKeys.Data, nil
}
func (s *sshkeysServiceHandler) Delete(ctx context.Context, sshKeyIdentifier string) error {
	delReq := struct {
		Identifier string `json:"identifier"`
	}{
		Identifier: sshKeyIdentifier,
	}

	req, err := s.client.NewRequest(ctx, http.MethodDelete, sshkeyBasePath, delReq)
	if err != nil {
		return err
	}

	if err = s.client.Do(ctx, req, nil); err != nil {
		return err
	}

	return nil
}
func (s *sshkeysServiceHandler) Get(ctx context.Context, sshKeyIdentifier string) (*SShKey, error) {
	path := fmt.Sprintf("%s/%s", sshkeyBasePath, sshKeyIdentifier)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	sshKeys := new(SshKeysGetRoot)
	if err = s.client.Do(ctx, req, sshKeys); err != nil {
		return nil, err
	}

	if sshKeys.Data == nil {
		return nil, fmt.Errorf("ssh key not found")
	}

	return sshKeys.Data, nil
}

func (s *sshkeysServiceHandler) Create(ctx context.Context, privateKey, name string) error {
	path := fmt.Sprintf("%s/add", sshkeyBasePath)

	createReq := struct {
		PrivateKey string `json:"privateKey"`
		Name       string `json:"name"`
	}{
		PrivateKey: privateKey,
		Name:       name,
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, createReq)
	if err != nil {
		return err
	}

	if err = s.client.Do(ctx, req, nil); err != nil {
		return err
	}

	return nil
}
