package govpsie

import (
	"context"
	"fmt"
	"net/http"
)

var dataCenterBasePath = "/apps/v2/datacenter"

type DataCenterService interface {
	List(ctx context.Context, options *ListOptions) ([]DataCenter, error)
}

type dataCenterServiceHandler struct {
	client *Client
}

var _ DataCenterService = &dataCenterServiceHandler{}

type DataCenter struct {
	DcName            string `json:"dc_name"`
	DcImage           string `json:"dc_image"`
	State             string `json:"state"`
	Country           string `json:"country"`
	IsActive          int    `json:"is_active"`
	Identifier        string `json:"identifier"`
	DefaultSelected   int    `json:"default_selected"`
	IsDeleted         int    `json:"is_deleted"`
	IsFipAvailable    int    `json:"is_fip_available"`
	IsBucketAvailable int    `json:"is_bucket_available"`
	IsPrivate         int    `json:"is_private"`
}

type DataCenterListRoot struct {
	Error bool         `json:"error"`
	Data  []DataCenter `json:"data"`
	Total int          `json:"total"`
}

func (d *dataCenterServiceHandler) List(ctx context.Context, options *ListOptions) ([]DataCenter, error) {
	path := fmt.Sprintf("%s?offset=%d&limit=%d", dataCenterBasePath, options.Page, options.PerPage)

	req, err := d.client.NewRequest(ctx, http.MethodGet, path, nil)

	if err != nil {
		return nil, err
	}

	root := new(DataCenterListRoot)
	if err = d.client.Do(ctx, req, root); err != nil {
		return nil, err
	}

	return root.Data, nil
}
