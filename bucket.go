package govpsie

import (
	"context"
	"fmt"
	"net/http"
)

var bucketsPath = "/apps/v2/buckets"
var bucketPath = "/apps/v2/bucket"

type BucketService interface {
	List(ctx context.Context, options *ListOptions) ([]Bucket, error)
	Get(ctx context.Context, id string) (*Bucket, error)
	Create(ctx context.Context, createReq *CreateBucketReq) error
	Delete(ctx context.Context, bucketId, reason, note string) error
	ToggleFileListing(ctx context.Context, bucketId string, fileListing bool) (bool, error)
	CheckFileListingStatus(ctx context.Context, bucketId string) (bool, error)
	GenerateKey(ctx context.Context, keyName string) error
	ListBucketKeys(ctx context.Context) ([]BucketKey, error)
}

type bucketServiceHandler struct {
	client *Client
}

var _ BucketService = &bucketServiceHandler{}

type ListBucketRoot struct {
	Error bool     `json:"error"`
	Data  []Bucket `json:"data"`
	Total int      `json:"total "`
}

type GetBucketRoot struct {
	Error bool   `json:"error"`
	Data  Bucket `json:"data"`
}

type CreateBucketReq struct {
	BucketName   string   `json:"bucketName"`
	FileListing  bool     `json:"fileListing"`
	ProjectId    string   `json:"projectId"`
	DataCenterId string   `json:"datacenterId"`
	Tags         []string `json:"tags"`
}

type FileListingStatusRoot struct {
	Error bool `json:"error"`
	Data  struct {
		FileListing bool `json:"fileListing"`
	} `json:"data"`
}

type Bucket struct {
	ID          int    `json:"id"`
	UserId      int    `json:"user_id"`
	AccessKey   string `json:"accessKey"`
	SecretKey   string `json:"secretKey"`
	BucketName  string `json:"bucketName"`
	ProjectName string `json:"projectName"`
	CreatedBy   string `json:"created_by"`
	EndPoint    string `json:"endPoint"`
	CreatedOn   string `json:"created_on"`
	Identifier  string `json:"identifier"`
	State       string `json:"state"`
	Country     string `json:"country"`
}

type BucketKey struct {
	AccessKey  string `json:"accessKey"`
	SecretKey  string `json:"secretKey"`
	KeyName    string `json:"key_name"`
	CreatedON  string `json:"created_on"`
	Identifier string `json:"identifier"`
}

type ListBucketKeysRoot struct {
	Error bool        `json:"error"`
	Data  []BucketKey `json:"data"`
}

func (s *bucketServiceHandler) List(ctx context.Context, options *ListOptions) ([]Bucket, error) {
	path := bucketsPath

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var root ListBucketRoot
	if err := s.client.Do(ctx, req, &root); err != nil {
		return nil, err
	}

	return root.Data, nil
}

func (s *bucketServiceHandler) Get(ctx context.Context, id string) (*Bucket, error) {
	path := bucketPath

	getReq := struct {
		BucketId string `json:"bucketId"`
	}{
		BucketId: id,
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, getReq)

	if err != nil {
		return nil, err
	}

	var root GetBucketRoot
	if err := s.client.Do(ctx, req, &root); err != nil {
		return nil, err
	}

	return &root.Data, nil
}

func (s *bucketServiceHandler) Create(ctx context.Context, createReq *CreateBucketReq) error {
	path := fmt.Sprintf("%s/create", bucketPath)

	// The API validates tags as an array; a nil slice marshals to null and is
	// rejected, so ensure it is sent as an empty array.
	if createReq.Tags == nil {
		createReq.Tags = []string{}
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, createReq)
	if err != nil {
		return err
	}

	return s.client.Do(ctx, req, nil)
}

func (s *bucketServiceHandler) Delete(ctx context.Context, buckId, reason, note string) error {
	path := fmt.Sprintf("%s/delete", bucketPath)

	deleteReq := struct {
		BucketId        string `json:"bucketId"`
		DeleteStatistic struct {
			Reason string `json:"reason"`
			Note   string `json:"note"`
		} `json:"deleteStatistic"`
	}{
		BucketId: buckId,
		DeleteStatistic: struct {
			Reason string `json:"reason"`
			Note   string `json:"note"`
		}{
			Reason: reason,
			Note:   note,
		},
	}

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, deleteReq)
	if err != nil {
		return err
	}

	return s.client.Do(ctx, req, nil)
}

func (s *bucketServiceHandler) GenerateKey(ctx context.Context, keyName string) error {
	path := fmt.Sprintf("%s/generate/keys", bucketPath)

	generateKeyReq := struct {
		KeyName string `json:"keyName"`
	}{
		KeyName: keyName,
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, generateKeyReq)
	if err != nil {
		return err
	}

	return s.client.Do(ctx, req, nil)
}

func (s *bucketServiceHandler) CheckFileListingStatus(ctx context.Context, bucketId string) (bool, error) {
	path := fmt.Sprintf("%s/listings/%s", bucketPath, bucketId)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return false, err
	}

	root := new(FileListingStatusRoot)
	if err := s.client.Do(ctx, req, root); err != nil {
		return false, err
	}

	return root.Data.FileListing, nil
}

func (s *bucketServiceHandler) ToggleFileListing(ctx context.Context, bucketId string, fileListing bool) (bool, error) {
	path := fmt.Sprintf("%s/listings/%s", bucketPath, bucketId)

	toogleReq := struct {
		FileListing bool `json:"fileListing"`
	}{
		FileListing: fileListing,
	}

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, toogleReq)
	if err != nil {
		return false, err
	}

	root := new(FileListingStatusRoot)
	if err := s.client.Do(ctx, req, root); err != nil {
		return false, err
	}

	return root.Data.FileListing, nil
}

func (s *bucketServiceHandler) ListBucketKeys(ctx context.Context) ([]BucketKey, error) {
	path := fmt.Sprintf("%s/keys", bucketPath)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	root := new(ListBucketKeysRoot)
	if err := s.client.Do(ctx, req, root); err != nil {
		return nil, err
	}

	return root.Data, nil
}
