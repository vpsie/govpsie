package govpsie

import (
	"context"
	"fmt"
	"net/http"
)

var certificateBasePath = "/apps/v2/certificates"

// CertificateService is the interface for managing TLS certificates.
type CertificateService interface {
	List(ctx context.Context) ([]Certificate, error)
	Add(ctx context.Context, certName, domainID string) error
	Delete(ctx context.Context, certID string) error
	Validate(ctx context.Context, certName, cert, privKey string) error
}

type certificateServiceHandler struct {
	client *Client
}

var _ CertificateService = &certificateServiceHandler{}

// Certificate represents a managed TLS certificate.
type Certificate struct {
	Identifier         string   `json:"identifier"`
	CertificateName    string   `json:"certificateName"`
	DomainName         string   `json:"domainName"`
	CommonNames        []string `json:"commonNames"`
	Issuer             string   `json:"issuer"`
	Serial             string   `json:"serial"`
	PublicKeySize      string   `json:"publicKeySize"`
	PublicKeyAlgorithm string   `json:"publicKeyAlgorithm"`
	ValidFrom          string   `json:"validFrom"`
	ValidTo            string   `json:"validTo"`
	CertType           string   `json:"certType"`
	UserID             int64    `json:"user_id"`
	FullName           string   `json:"fullname"`
	CreatedOn          string   `json:"created_on"`
	UpdatedOn          string   `json:"updated_on"`
}

type certificatesListRoot struct {
	Error bool          `json:"error"`
	Data  []Certificate `json:"data"`
	Total int64         `json:"total"`
}

func (s *certificateServiceHandler) List(ctx context.Context) ([]Certificate, error) {
	path := fmt.Sprintf("%s/all", certificateBasePath)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	root := new(certificatesListRoot)
	if err = s.client.Do(ctx, req, root); err != nil {
		return nil, err
	}

	return root.Data, nil
}

// Add requests issuance of a certificate for the given domain.
func (s *certificateServiceHandler) Add(ctx context.Context, certName, domainID string) error {
	path := fmt.Sprintf("%s/add", certificateBasePath)

	addReq := struct {
		CertName string `json:"certName"`
		DomainID string `json:"domainId"`
	}{
		CertName: certName,
		DomainID: domainID,
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, &addReq)
	if err != nil {
		return err
	}

	return s.client.Do(ctx, req, nil)
}

func (s *certificateServiceHandler) Delete(ctx context.Context, certID string) error {
	path := fmt.Sprintf("%s/delete", certificateBasePath)

	delReq := struct {
		CertID string `json:"certId"`
	}{
		CertID: certID,
	}

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, &delReq)
	if err != nil {
		return err
	}

	return s.client.Do(ctx, req, nil)
}

// Validate uploads a custom certificate/key pair for validation.
func (s *certificateServiceHandler) Validate(ctx context.Context, certName, cert, privKey string) error {
	path := fmt.Sprintf("%s/validate", certificateBasePath)

	validateReq := struct {
		CertName string `json:"certName"`
		Cert     string `json:"cert"`
		PrivKey  string `json:"privKey"`
	}{
		CertName: certName,
		Cert:     cert,
		PrivKey:  privKey,
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, &validateReq)
	if err != nil {
		return err
	}

	return s.client.Do(ctx, req, nil)
}
