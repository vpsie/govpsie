package govpsie

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

var snapshotBasePath = "/apps/v2/snapshot"
var snapshotsBasePath = "/apps/v2/snapshots"
var backupBasePath = "/apps/v2/backup"

type SnapshotService interface {
	List(ctx context.Context, options *ListOptions) ([]Snapshot, error)
	Create(ctx context.Context, name, vmIdentifier, note string) error
	ListByVm(ctx context.Context, options *ListOptions, vmIdentifier string) ([]Snapshot, error)
	Rollback(ctx context.Context, snapshotIdentifier string) error
	EnableAuto(ctx context.Context, enableReq *EnableAutoSnapshotReq) error
	Delete(ctx context.Context, snapshotIdentifier, reason, note string) error
	Update(ctx context.Context, snapshotIdentifier, newNote string) error
	Get(ctx context.Context, buckupIdentifier string) (*Snapshot, error)
	GetSnapShotPolicy(ctx context.Context, identifier string) (*SnapShotPolicy, error)
	CreateSnapShotPolicy(ctx context.Context, createReq *CreateSnapShotPolicyReq) error
	DeleteSnapShotPolicy(ctx context.Context, policyId, identifier string) error
	ManageRetainSnapShotPolicy(ctx context.Context, policyId string, keep int64) error
	AttachSnapShotPolicy(ctx context.Context, policyId string, vms []string) error
	DetachSnapShotPolicy(ctx context.Context, policyId string, vms []string) error
	ListSnapShotPolicies(ctx context.Context, options *ListOptions) ([]SnapShotPolicyListDetail, error)
}

type snapshotServiceHandler struct {
	client *Client
}

var _ SnapshotService = &snapshotServiceHandler{}

type Snapshot struct {
	Hostname     string    `json:"hostname"`
	Name         string    `json:"name"`
	Identifier   string    `json:"identifier"`
	BackupKey    string    `json:"backupKey"`
	State        string    `json:"state"`
	DcIdentifier string    `json:"dcIdentifier"`
	Daily        int64     `json:"daily"`
	IsSnapshot   int64     `json:"is_snapshot"`
	VmIdentifier string    `json:"vmIdentifier"`
	BackupSHA1   string    `json:"backupsha1"`
	IsDeletedVM  int64     `json:"is_deleted_vm"`
	CreatedOn    time.Time `json:"created_on"`
	Note         string    `json:"note"`
	BackupSize   int64     `json:"backup_size"`
	DcName       string    `json:"dcName"`
	Weekly       int64     `json:"weekly"`
	Monthly      int64     `json:"monthly"`
	BoxID        int64     `json:"box_id"`
	GlobalBackup int64     `json:"global_backup"`
	OsIdentifier string    `json:"osIdentifier"`
	OsFullName   string    `json:"osFullName"`
	VMCategory   string    `json:"vmCategory"`
	VMSSD        int64     `json:"vmSSD"`
}

type GetSnapshotRoot struct {
	Error bool `json:"error"`
	Data  struct {
		Backup Snapshot `json:"backup"`
	} `json:"data"`
}

type ListSnapshotsRoot struct {
	Error bool       `json:"error"`
	Data  []Snapshot `json:"data"`
	Total int64      `json:"total"`
}

type EnableAutoSnapshotReq struct {
	VMIdentifier    string   `json:"vmIdentifier"`
	VmId            int64    `json:"vmId"`
	Period          string   `json:"period"`
	DailySnapshot   int64    `json:"dailySnapshot"`
	WeeklySnapshot  int64    `json:"weeklySnapshot"`
	MonthlySnapshot int64    `json:"monthlySnapshot"`
	Tags            []string `json:"tags"`
}

func (s *snapshotServiceHandler) List(ctx context.Context, options *ListOptions) ([]Snapshot, error) {
	if options == nil {
		options = &ListOptions{}
	}
	path := fmt.Sprintf("%s?offset=%d&limit=%d", snapshotBasePath, options.Page, options.PerPage)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	snapshots := new(ListSnapshotsRoot)

	if err = s.client.Do(ctx, req, &snapshots); err != nil {
		return nil, err
	}

	return snapshots.Data, nil

}

func (s *snapshotServiceHandler) Create(ctx context.Context, name, vmIdentifier, note string) error {
	path := fmt.Sprintf("%s/add", snapshotBasePath)
	createSnapshotReq := struct {
		Name         string `json:"name"`
		VMIdentifier string `json:"vmIdentifier"`
		Note         string `json:"note"`
	}{
		Name:         name,
		VMIdentifier: vmIdentifier,
		Note:         note,
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, &createSnapshotReq)
	if err != nil {
		return err
	}

	return s.client.Do(ctx, req, nil)

}

func (s *snapshotServiceHandler) Get(ctx context.Context, buckupIdentifier string) (*Snapshot, error) {
	path := fmt.Sprintf("%s/%s", backupBasePath, buckupIdentifier)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	backup := new(GetSnapshotRoot)
	if err = s.client.Do(ctx, req, backup); err != nil {
		return nil, err
	}

	return &backup.Data.Backup, nil
}

func (s *snapshotServiceHandler) Update(ctx context.Context, snapshotIdentifier, newNote string) error {
	path := "/apps/v2/backups/update"

	updateReq := struct {
		Identifier string `json:"identifier"`
		Note       string `json:"note"`
	}{
		Identifier: snapshotIdentifier,
		Note:       newNote,
	}

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, &updateReq)
	if err != nil {
		return err
	}

	return s.client.Do(ctx, req, nil)
}

func (s *snapshotServiceHandler) ListByVm(ctx context.Context, options *ListOptions, vmIdentifier string) ([]Snapshot, error) {
	if options == nil {
		options = &ListOptions{}
	}
	path := fmt.Sprintf("/apps/v2/vm/snapshot/%s?offset=%d&limit=%d", vmIdentifier, options.Page, options.PerPage)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	snapshots := new(ListSnapshotsRoot)

	if err = s.client.Do(ctx, req, &snapshots); err != nil {
		return nil, err
	}

	return snapshots.Data, nil

}

func (s *snapshotServiceHandler) Delete(ctx context.Context, snapshotIdentifier, reason, note string) error {
	deleteReq := struct {
		SnapshotIdentifier string `json:"snapshotIdentifier"`
		DeleteStatistic    struct {
			Reason string `json:"reason"`
			Note   string `json:"note"`
		} `json:"deleteStatistic"`
	}{
		SnapshotIdentifier: snapshotIdentifier,
		DeleteStatistic: struct {
			Reason string `json:"reason"`
			Note   string `json:"note"`
		}{
			Reason: reason,
			Note:   note,
		},
	}

	req, err := s.client.NewRequest(ctx, http.MethodDelete, snapshotBasePath, &deleteReq)
	if err != nil {
		return err
	}

	return s.client.Do(ctx, req, nil)
}

func (s *snapshotServiceHandler) Rollback(ctx context.Context, snapshotIdentifier string) error {
	path := fmt.Sprintf("%s/rollback", snapshotBasePath)

	rollbackReq := struct {
		SnapshotIdentifier string `json:"snapshotIdentifier"`
	}{
		SnapshotIdentifier: snapshotIdentifier,
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, &rollbackReq)
	if err != nil {
		return err
	}

	return s.client.Do(ctx, req, nil)
}

func (s *snapshotServiceHandler) EnableAuto(ctx context.Context, enableReq *EnableAutoSnapshotReq) error {
	path := fmt.Sprintf("%s/enable/auto", snapshotBasePath)

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, enableReq)
	if err != nil {
		return err
	}

	return s.client.Do(ctx, req, nil)
}

// snap shot policy
type SnapShotVms struct {
	Name       string `json:"name"`
	Identifier string `json:"identifier"`
	Category   string `json:"category"`
	Fullname   string `json:"fullname"`
	Type       string `json:"type"`
}
type SnapShotPolicy struct {
	Name       string        `json:"name"`
	Identifier string        `json:"identifier"`
	CreatedOn  string        `json:"created_on"`
	CreatedBy  string        `json:"created_by"`
	BackupPlan string        `json:"backupPlan"`
	PlanEvery  int64         `json:"planEvery"`
	Keep       int64         `json:"keep"`
	Disabled   int64         `json:"disabled"`
	UserId     int64         `json:"userId"`
	Vms        []SnapShotVms `json:"vms"`
}

type SnapShotPolicyListDetail struct {
	Name       string `json:"name"`
	Identifier string `json:"identifier"`
	CreatedOn  string `json:"created_on"`
	CreatedBy  string `json:"created_by"`
	BackupPlan string `json:"backupPlan"`
	PlanEvery  int64  `json:"planEvery"`
	Keep       int64  `json:"keep"`
	Disabled   int64  `json:"disabled"`
	VmsCount   int64  `json:"vmsCount"`
	UserId     int64  `json:"userId"`
}
type ListSnapShotPoliciesRoot struct {
	Error bool `json:"error"`
	Data  struct {
		Rows []SnapShotPolicyListDetail `json:"rows"`
	} `json:"data"`
}

type GetSnapShotPolicyRoot struct {
	Error bool            `json:"error"`
	Data  json.RawMessage `json:"data"`
}

type CreateSnapShotPolicyReq struct {
	Name       string   `json:"name"`
	BackupPlan string   `json:"backupPlan"`
	PlanEvery  string   `json:"planEvery"`
	Keep       string   `json:"keep"`
	Vms        []string `json:"vms"`
	Tags       []string `json:"tags"`
}

func (s *snapshotServiceHandler) ListSnapShotPolicies(ctx context.Context, options *ListOptions) ([]SnapShotPolicyListDetail, error) {
	path := fmt.Sprintf("%s/policy/all", snapshotBasePath)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	policies := new(ListSnapShotPoliciesRoot)
	if err = s.client.Do(ctx, req, policies); err != nil {
		return nil, err
	}

	return policies.Data.Rows, nil
}

// GetSnapShotPolicy fetches a single snapshot policy by its UUID identifier.
// When the policy does not exist the API returns HTTP 200 with a literal
// `"data": false`, so (nil, nil) is returned to signal not-found.
func (s *snapshotServiceHandler) GetSnapShotPolicy(ctx context.Context, identifier string) (*SnapShotPolicy, error) {
	path := fmt.Sprintf("%s/policy/%s", snapshotBasePath, identifier)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	root := new(GetSnapShotPolicyRoot)
	if err = s.client.Do(ctx, req, root); err != nil {
		return nil, err
	}

	if isNullData(root.Data) {
		return nil, nil
	}

	policy := new(SnapShotPolicy)
	if err = json.Unmarshal(root.Data, policy); err != nil {
		return nil, err
	}

	return policy, nil
}

func (s *snapshotServiceHandler) CreateSnapShotPolicy(ctx context.Context, createReq *CreateSnapShotPolicyReq) error {
	path := fmt.Sprintf("%s/policy/create", snapshotBasePath)

	// The API validates tags and vms as arrays; nil slices marshal to null and
	// are rejected, so ensure they are sent as empty arrays.
	if createReq.Tags == nil {
		createReq.Tags = []string{}
	}
	if createReq.Vms == nil {
		createReq.Vms = []string{}
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, createReq)
	if err != nil {
		return err
	}

	return s.client.Do(ctx, req, nil)
}

func (s *snapshotServiceHandler) DeleteSnapShotPolicy(ctx context.Context, policyId, identifier string) error {
	path := fmt.Sprintf("%s/policy/%s", snapshotBasePath, identifier)

	deleteSnapShot := struct {
		PolicyId string `json:"policyId"`
	}{
		PolicyId: policyId,
	}

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, deleteSnapShot)
	if err != nil {
		return err
	}

	return s.client.Do(ctx, req, nil)
}

func (s *snapshotServiceHandler) ManageRetainSnapShotPolicy(ctx context.Context, policyId string, keep int64) error {
	path := fmt.Sprintf("%s/policy/keep", snapshotBasePath)

	manageSnapShot := struct {
		PolicyId string `json:"policyId"`
		Keep     int64  `json:"keep"`
	}{
		PolicyId: policyId,
		Keep:     keep,
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, manageSnapShot)
	if err != nil {
		return err
	}

	return s.client.Do(ctx, req, nil)
}

func (s *snapshotServiceHandler) AttachSnapShotPolicy(ctx context.Context, policyId string, vms []string) error {
	path := fmt.Sprintf("%s/policy/attach", snapshotsBasePath)

	attachSnapShot := struct {
		PolicyId string   `json:"policyId"`
		Vms      []string `json:"vms"`
	}{
		PolicyId: policyId,
		Vms:      vms,
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, attachSnapShot)
	if err != nil {
		return err
	}

	return s.client.Do(ctx, req, nil)
}

func (s *snapshotServiceHandler) DetachSnapShotPolicy(ctx context.Context, policyId string, vms []string) error {
	path := fmt.Sprintf("%s/policy/detach", snapshotsBasePath)

	detachSnapShot := struct {
		PolicyId string   `json:"policyId"`
		Vms      []string `json:"vms"`
	}{
		PolicyId: policyId,
		Vms:      vms,
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, detachSnapShot)
	if err != nil {
		return err
	}

	return s.client.Do(ctx, req, nil)
}
