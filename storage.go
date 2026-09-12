package govpsie

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

var storageBasePath = "/apps/v2"

type StorageService interface {
	List(ctx context.Context, options *ListOptions) ([]Storage, error)
	Delete(ctx context.Context, storageIdentifier string) error
	AttachToServer(ctx context.Context, storageIdentifier, vmIdentifier string, vmType string) error
	DetachToServer(ctx context.Context, storageIdentifier, vmIdentifier string, vmType string) error
	CreateContainer(ctx context.Context, dcIdentifier string) error
	ListAll(ctx context.Context, options *ListOptions) ([]Storage, error)
	Create(ctx context.Context, createReq *StorageCreateRequest, vmIdentifier string, vmType string) error
	ListVmsToAttach(ctx context.Context) ([]VmToAttach, error)
	CreateVolume(ctx context.Context, creatReq *StorageCreateRequest) error
	CreateStorage(ctx context.Context, createReq *StorageCreateRequest) error
	DetachAllFromServer(ctx context.Context, vmIdentifier string, vmType string) error
	UpdateSize(ctx context.Context, storageIdentifier, size string) error
	UpdateName(ctx context.Context, storageIdentifier, name string) error
	CreateSnapshot(ctx context.Context, storageIdentifier, name, storageType string) error
	ListSnapshots(ctx context.Context, options *ListOptions) ([]StorageSnapShot, error)
	UpdateSnapshotName(ctx context.Context, snapshotIdentifier, name string) error
	RollbackSnapshot(ctx context.Context, snapshotIdentifier, snapType string) error
	CloneSnapshot(ctx context.Context, snapshotIdentifier, snapType string) error
	DeleteSnapshot(ctx context.Context, snapshotIdentifier string) error
	DeleteAllSnapshots(ctx context.Context, storageIdentifier string) error
	Get(ctx context.Context, identifier string) (*StorageDetail, error)
	ListStorageDataCenter(ctx context.Context) ([]DataCenter, error)
}

type storageServiceHandler struct {
	client *Client
}

var _ StorageService = &storageServiceHandler{}

type StorageDetail struct {
	ID                int       `json:"id"`
	Name              string    `json:"name"`
	Description       string    `json:"description"`
	UserID            int       `json:"user_id"`
	BoxID             int       `json:"box_id"`
	Identifier        string    `json:"identifier"`
	UserTemplateID    int       `json:"user_template_id"`
	StorageType       string    `json:"storage_type"`
	DiskFormat        string    `json:"disk_format"`
	BusDevice         string    `json:"bus_device"`
	IsAutomatic       int       `json:"is_automatic"`
	Size              int       `json:"size"`
	IsDeleted         int       `json:"is_deleted"`
	BusNumber         int       `json:"bus_number"`
	DatacenterID      int       `json:"datacenter_id"`
	UpdatedAt         time.Time `json:"updated_at"`
	StorageID         int       `json:"storage_id"`
	DiskKey           string    `json:"disk_key"`
	CreatedOn         time.Time `json:"created_on"`
	IndexTemplateProp string    `json:"index_template_prop"`
	EntityType        string    `json:"entity_type"`
	Hostname          string    `json:"hostname"`
	IsDeletedVM       int       `json:"is_deleted_vm"`
	State             string    `json:"state"`
	DcName            string    `json:"dcName"`
	DcIdentifier      string    `json:"dcIdentifier"`
	VMIdentifier      string    `json:"vmIdentifier"`
	OsIdentifier      string    `json:"osIdentifier"`
	OsFullName        string    `json:"osFullName"`
	VMCategory        string    `json:"vmCategory"`
	VMSSD             int       `json:"vmSSD"`
}

type GetStorageRoot struct {
	Error bool          `json:"error"`
	Data  StorageDetail `json:"data"`
}
type Storage struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	UserID         int    `json:"user_id"`
	BoxID          int    `json:"box_id"`
	Identifier     string `json:"identifier"`
	UserTemplateID int    `json:"user_template_id"`
	StorageType    string `json:"storage_type"`
	DiskFormat     string `json:"disk_format"`
	IsAutomatic    int    `json:"is_automatic"`
	Size           int    `json:"size"`
	StorageID      int    `json:"storage_id"`
	DiskKey        string `json:"disk_key"`
	CreatedOn      string `json:"created_on"`
	VmIdentifier   string `json:"vmIdentifier"`
	Hostname       string `json:"hostname"`
	OsIdentifier   string `json:"osIdentifier"`
	State          string `json:"state"`
	DcIdentifier   string `json:"dcIdentifier"`
	BusDevice      string `json:"bus_device"`
	BusNumber      int    `json:"bus_number"`
}

type ListStorageRoot struct {
	Error bool      `json:"error"`
	Data  []Storage `json:"data"`
	Total int       `json:"total"`
}

type StorageUpdateRequest struct {
	StorageIdentifier string `json:"storageIdentifier"`
	Name              string `json:"name"`
	Description       string `json:"description"`
	Size              int    `json:"size"`
}

type StorageCreateRequest struct {
	Name         string `json:"name"`
	DcIdentifier string `json:"dcIdentifier"`
	Description  string `json:"description"`
	Size         int    `json:"size"`
	StorageType  string `json:"storageType"`
	DiskFormat   string `json:"diskFormat"`
	IsAutomatic  int    `json:"isAutomatic"`
	K8sTrials    int    `json:"k8sTrials"`
}

type VmToAttach struct {
	Hostname          string      `json:"hostname"`
	Identifier        string      `json:"identifier"`
	DatacenterID      int         `json:"datacenter_id"`
	DefaultIP         string      `json:"default_ip"`
	DefaultIpv6       string      `json:"default_ipv6"`
	PrivateIP         interface{} `json:"private_ip"`
	Ssd               int         `json:"ssd"`
	State             string      `json:"state"`
	IsFipAvailable    int         `json:"is_fip_available"`
	IsBucketAvailable int         `json:"is_bucket_available"`
	DcIdentifier      string      `json:"dcIdentifier"`
	Fullname          string      `json:"fullname"`
	Category          string      `json:"category"`
	Type              string      `json:"type"`
}

type ListVmToAttachRoot struct {
	Error bool         `json:"error"`
	Data  []VmToAttach `json:"data"`
	Total int          `json:"total"`
}

type StorageSnapShot struct {
	ID          int       `json:"id"`
	StorageID   int       `json:"storage_id"`
	Identifier  string    `json:"identifier"`
	Name        string    `json:"name"`
	Size        int       `json:"size"`
	CreatedOn   time.Time `json:"created_on"`
	UserID      int       `json:"user_id"`
	IsDeleted   int       `json:"is_deleted"`
	SnapshotKey string    `json:"snapshot_key"`
	StorageName string    `json:"storage_name"`
	StorageType string    `json:"storage_type"`
	DiskFormat  string    `json:"disk_format"`
	BoxID       int       `json:"box_id"`
	EntityType  string    `json:"entity_type"`
}

type ListStorageSnapShotRoot struct {
	Error bool              `json:"error"`
	Data  []StorageSnapShot `json:"data"`
	Total int               `json:"total"`
}

type ListStorageDataCenterRoot struct {
	Error bool         `json:"error"`
	Data  []DataCenter `json:"data"`
}

func (s *storageServiceHandler) List(ctx context.Context, options *ListOptions) ([]Storage, error) {
	path := fmt.Sprintf("%s/storages?offset=%d&limit%d", storageBasePath, options.Page, options.PerPage)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	storages := new(ListStorageRoot)
	if err = s.client.Do(ctx, req, &storages); err != nil {
		return nil, err
	}

	return storages.Data, nil
}

func (s *storageServiceHandler) Delete(ctx context.Context, storageIdentifier string) error {
	path := fmt.Sprintf("%s/storages", storageBasePath)

	deleteRequest := struct {
		StorageIdentifier string `json:"storageIdentifier"`
	}{
		StorageIdentifier: storageIdentifier,
	}

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, &deleteRequest)
	if err != nil {
		return err
	}

	return s.client.Do(ctx, req, nil)
}

func (s *storageServiceHandler) AttachToServer(ctx context.Context, storageIdentifier, vmIdentifier string, vmType string) error {

	path := fmt.Sprintf("%s/storages/vm/attach", storageBasePath)

	attachReq := struct {
		StorageIdentifier string `json:"storageIdentifier"`
		VmIdentifier      string `json:"vmIdentifier"`
		Type              string `json:"type"`
	}{
		StorageIdentifier: storageIdentifier,
		VmIdentifier:      vmIdentifier,
		Type:              vmType,
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, &attachReq)
	if err != nil {
		return err
	}

	return s.client.Do(ctx, req, nil)

}

func (s *storageServiceHandler) DetachToServer(ctx context.Context, storageIdentifier, vmIdentifier string, vmType string) error {

	path := fmt.Sprintf("%s/storages/vm/detach", storageBasePath)

	detachReq := struct {
		StorageIdentifier string `json:"storageIdentifier"`
		VmIdentifier      string `json:"vmIdentifier"`
		Type              string `json:"type"`
	}{
		StorageIdentifier: storageIdentifier,
		VmIdentifier:      vmIdentifier,
		Type:              vmType,
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, &detachReq)
	if err != nil {
		return err
	}

	return s.client.Do(ctx, req, nil)

}
func (s *storageServiceHandler) CreateContainer(ctx context.Context, dcIdentifier string) error {

	path := fmt.Sprintf("%s/storages/create/container", storageBasePath)

	createContainerReq := struct {
		DcIdentifier string `json:"dcIdentifier"`
	}{
		DcIdentifier: dcIdentifier,
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, createContainerReq)
	if err != nil {
		return err
	}

	return s.client.Do(ctx, req, nil)
}

func (s *storageServiceHandler) UpdateSize(ctx context.Context, storageIdentifier, size string) error {
	path := fmt.Sprintf("%s/storages/edit", storageBasePath)

	updateReq := struct {
		StorageIdentifier string `json:"storageIdentifier"`
		Size              string `json:"size"`
	}{
		StorageIdentifier: storageIdentifier,
		Size:              size,
	}

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, &updateReq)
	if err != nil {
		return err
	}

	return s.client.Do(ctx, req, nil)
}

func (s *storageServiceHandler) UpdateName(ctx context.Context, storageIdentifier, name string) error {
	path := fmt.Sprintf("%s/storages/rename", storageBasePath)

	renameReq := struct {
		StorageIdentifier string `json:"storageIdentifier"`
		Name              string `json:"name"`
	}{
		StorageIdentifier: storageIdentifier,
		Name:              name,
	}
	req, err := s.client.NewRequest(ctx, http.MethodPut, path, &renameReq)
	if err != nil {
		return err
	}

	return s.client.Do(ctx, req, nil)
}

func (s *storageServiceHandler) ListAll(ctx context.Context, options *ListOptions) ([]Storage, error) {
	path := fmt.Sprintf("%s/storages", storageBasePath)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	storages := new(ListStorageRoot)
	if err = s.client.Do(ctx, req, &storages); err != nil {
		return nil, err
	}
	return storages.Data, nil
}

func (s *storageServiceHandler) Create(ctx context.Context, createReq *StorageCreateRequest, vmIdentifier string, vmType string) error {
	path := fmt.Sprintf("%s/storages/vm/attach/all", storageBasePath)

	createPayload := struct {
		Storages     []StorageCreateRequest `json:"storages"`
		VmIdentifier string                 `json:"vmIdentifier"`
		Type         string                 `json:"type"`
	}{
		Storages: []StorageCreateRequest{
			*createReq,
		},
		VmIdentifier: vmIdentifier,
		Type:         vmType,
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, createPayload)
	if err != nil {
		return err
	}

	return s.client.Do(ctx, req, nil)
}

func (s *storageServiceHandler) Get(ctx context.Context, identifier string) (*StorageDetail, error) {
	path := fmt.Sprintf("%s/storages/%s", storageBasePath, identifier)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	storage := new(GetStorageRoot)
	if err = s.client.Do(ctx, req, &storage); err != nil {
		return nil, err
	}

	return &storage.Data, nil
}

func (s *storageServiceHandler) ListVmsToAttach(ctx context.Context) ([]VmToAttach, error) {
	path := fmt.Sprintf("%s/storages/vms", storageBasePath)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		fmt.Printf("error in creating request: %v\n", err)
		return nil, err
	}

	vms := new(ListVmToAttachRoot)
	if err = s.client.Do(ctx, req, &vms); err != nil {
		fmt.Printf("error in making the request: %v\n", err)
		return nil, err
	}

	return vms.Data, nil
}

func (s *storageServiceHandler) CreateStorage(ctx context.Context, createReq *StorageCreateRequest) error {
	path := fmt.Sprintf("%s/storages/create/multiple", storageBasePath)
	fullReq := struct {
		Storages []StorageCreateRequest `json:"storages"`
	}{
		Storages: []StorageCreateRequest{
			*createReq,
		},
	}
	req, err := s.client.NewRequest(ctx, http.MethodPost, path, fullReq)
	if err != nil {
		return err
	}
	return s.client.Do(ctx, req, nil)
}

func (s *storageServiceHandler) CreateVolume(ctx context.Context, creatReq *StorageCreateRequest) error {
	path := fmt.Sprintf("%s/storage/create", storageBasePath)
	fullReq := struct {
		Storages []StorageCreateRequest `json:"storages"`
	}{
		Storages: []StorageCreateRequest{
			*creatReq,
		},
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, fullReq)
	if err != nil {
		return err
	}

	return s.client.Do(ctx, req, nil)
}

func (s *storageServiceHandler) DetachAllFromServer(ctx context.Context, vmIdentifier string, vmType string) error {
	path := fmt.Sprintf("%s/storages/vm/detach/all", storageBasePath)

	detachReq := struct {
		VmIdentifier string `json:"vmIdentifier"`
		Type         string `json:"type"`
	}{
		VmIdentifier: vmIdentifier,
		Type:         vmType,
	}

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, &detachReq)
	if err != nil {
		return err
	}

	return s.client.Do(ctx, req, nil)

}

func (s *storageServiceHandler) CreateSnapshot(ctx context.Context, storageIdentifier, name, storageType string) error {
	path := fmt.Sprintf("%s/storages/snapshot", storageBasePath)

	snapshotReq := struct {
		StorageIdentifier string `json:"storageIdentifier"`
		Name              string `json:"name"`
		Type              string `json:"type"`
	}{
		StorageIdentifier: storageIdentifier,
		Name:              name,
		Type:              storageType,
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, &snapshotReq)
	if err != nil {
		return err
	}

	return s.client.Do(ctx, req, nil)
}

func (s *storageServiceHandler) ListSnapshots(ctx context.Context, options *ListOptions) ([]StorageSnapShot, error) {
	path := fmt.Sprintf("%s/storage/snapshots?offset=%d&limit=%d", storageBasePath, options.Page, options.PerPage)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	snapshots := new(ListStorageSnapShotRoot)
	if err = s.client.Do(ctx, req, &snapshots); err != nil {
		return nil, err
	}

	return snapshots.Data, nil
}

func (s *storageServiceHandler) UpdateSnapshotName(ctx context.Context, snapshotIdentifier, name string) error {
	path := fmt.Sprintf("%s/storages/snapshot/rename", storageBasePath)

	renameReq := struct {
		SnapshotIdentifier string `json:"snapshotIdentifier"`
		Name               string `json:"name"`
	}{
		SnapshotIdentifier: snapshotIdentifier,
		Name:               name,
	}

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, &renameReq)
	if err != nil {
		return err
	}

	return s.client.Do(ctx, req, nil)
}

func (s *storageServiceHandler) RollbackSnapshot(ctx context.Context, snapshotIdentifier, snapType string) error {
	path := fmt.Sprintf("%s/storages/snapshot/rollback", storageBasePath)

	rollbackReq := struct {
		SnapshotIdentifier string `json:"snapshotIdentifier"`
		Type               string `json:"type"`
	}{
		SnapshotIdentifier: snapshotIdentifier,
		Type:               snapType,
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, &rollbackReq)
	if err != nil {
		return err
	}

	return s.client.Do(ctx, req, nil)
}

func (s *storageServiceHandler) CloneSnapshot(ctx context.Context, snapshotIdentifier, snapType string) error {
	path := fmt.Sprintf("%s/storages/snapshot/clone", storageBasePath)

	cloneReq := struct {
		SnapshotIdentifier string `json:"snapshotIdentifier"`
		Type               string `json:"type"`
	}{
		SnapshotIdentifier: snapshotIdentifier,
		Type:               snapType,
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, &cloneReq)
	if err != nil {
		return err
	}

	return s.client.Do(ctx, req, nil)
}

func (s *storageServiceHandler) DeleteSnapshot(ctx context.Context, snapshotIdentifier string) error {
	path := fmt.Sprintf("%s/storages/snapshot/delete", storageBasePath)

	deleteReq := struct {
		SnapshotIdentifier string `json:"snapshotIdentifier"`
	}{
		SnapshotIdentifier: snapshotIdentifier,
	}

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, &deleteReq)
	if err != nil {
		return err
	}

	return s.client.Do(ctx, req, nil)
}

func (s *storageServiceHandler) DeleteAllSnapshots(ctx context.Context, storageIdentifier string) error {
	path := fmt.Sprintf("%s/storages/snapshots/delete/all", storageBasePath)

	deleteReq := struct {
		StorageIdentifier string `json:"storageIdentifier"`
	}{
		StorageIdentifier: storageIdentifier,
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, &deleteReq)
	if err != nil {
		return err
	}

	return s.client.Do(ctx, req, nil)
}

func (s *storageServiceHandler) ListStorageDataCenter(ctx context.Context) ([]DataCenter, error) {
	path := fmt.Sprintf("%s/storage/datacenter", storageBasePath)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	dataCenters := new(ListStorageDataCenterRoot)

	if err = s.client.Do(ctx, req, dataCenters); err != nil {
		return nil, err
	}

	return dataCenters.Data, nil
}
