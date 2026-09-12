package govpsie

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

var backupsPath = "/apps/v2"

type BackupsService interface {
	List(ctx context.Context, options *ListOptions) ([]Backup, error)
	DeleteBackup(ctx context.Context, backupIdentifier, deleteReason, deleteNote string) error
	CreateBackups(ctx context.Context, vmIdentifier, name, notes string) error
	ListByServer(ctx context.Context, options *ListOptions, vmIdentifier string) ([]Backup, error)
	CreateServerByBackup(ctx context.Context, backupIdentifier string) error
	Get(ctx context.Context, identifer string) (*Backup, error)
	EnableAutoBackup(ctx context.Context, enableAutoReq *EnableAutoBackupReq) error
	Rename(ctx context.Context, backupIdentifier string, newName string) error
	GetBackupPolicy(ctx context.Context, identifier string) (*BackupPolicy, error)
	CreateBackupPolicy(ctx context.Context, createReq *CreateBackupPolicyReq) error
	DeleteBackupPolicy(ctx context.Context, policyId, identifier string) error
	ManageRetainBackupPolicy(ctx context.Context, policyId string, keep int) error
	AttachBackupPolicy(ctx context.Context, policyId string, vms []string) error
	DetachBackupPolicy(ctx context.Context, policyId string, vms []string) error
	ListBackupPolicies(ctx context.Context, options *ListOptions) ([]BackupPolicyListDetail, error)
}

type EnableAutoBackupReq struct {
	VmIdentifier  string   `json:"vmIdentifier"`
	VmID          int      `json:"vmId"`
	Period        string   `json:"period"`
	AutoBackup    int      `json:"autoBackup"`
	WeeklyBackup  int      `json:"weeklyBackup"`
	MonthlyBackup int      `json:"monthlyBackup"`
	Tags          []string `json:"tags"`
}

type backupsServiceHandler struct {
	client *Client
}

var _ BackupsService = &backupsServiceHandler{}

type ListBackupsRoot struct {
	Error bool     `json:"error"`
	Data  []Backup `json:"data"`
	Total int      `json:"total "`
}

type GetBackupsRoot struct {
	Error bool `json:"error"`
	Data  struct {
		Backup Backup `json:"backup"`
	} `json:"data"`
}

type Backup struct {
	HostName     string `json:"hostname"`
	Name         string `json:"name"`
	Identifier   string `json:"identifier"`
	Note         string `json:"note"`
	BackupKey    string `json:"backupKey"`
	State        string `json:"state"`
	DcIdentifier string `json:"dcIdentifier"`
	VMIdentifier string `json:"vmIdentifier"`
	BoxID        int    `json:"boxId"`
	BackupSHA1   string `json:"backupsha1"`
	OSFullName   string `json:"osFullName"`
	VMCategory   string `json:"vmCategory"`
	CreatedBy    string `json:"created_by"`
	CreatedOn    string `json:"created_on"`
}

func (b *backupsServiceHandler) List(ctx context.Context, options *ListOptions) ([]Backup, error) {
	path := fmt.Sprintf("%s/backups", backupsPath)

	req, err := b.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	backups := new(ListBackupsRoot)
	if err = b.client.Do(ctx, req, backups); err != nil {
		return nil, err
	}

	return backups.Data, nil
}

func (b *backupsServiceHandler) DeleteBackup(ctx context.Context, backupIdentifier, deleteReason, deleteNote string) error {
	path := fmt.Sprintf("%s/backup", backupsPath)

	deleteReq := struct {
		BackupIdentifier string `json:"backupIdentifier"`
		DeleteStatistic  struct {
			Reason string `json:"reason"`
			Note   string `json:"note"`
		} `json:"deleteStatistic"`
	}{
		BackupIdentifier: backupIdentifier,
		DeleteStatistic: struct {
			Reason string `json:"reason"`
			Note   string `json:"note"`
		}{
			Reason: deleteReason,
			Note:   deleteNote,
		},
	}

	req, err := b.client.NewRequest(ctx, http.MethodDelete, path, &deleteReq)
	if err != nil {
		return err
	}

	return b.client.Do(ctx, req, nil)
}

func (b *backupsServiceHandler) CreateBackups(ctx context.Context, vmIdentifier, name, notes string) error {
	path := fmt.Sprintf("%s/backup/add", backupsPath)

	createBackupReq := struct {
		VmIdentifier string `json:"vmIdentifier"`
		Name         string `json:"name"`
	}{
		VmIdentifier: vmIdentifier,
		Name:         name,
	}

	req, err := b.client.NewRequest(ctx, http.MethodPost, path, &createBackupReq)

	if err != nil {
		return err
	}

	return b.client.Do(ctx, req, nil)

}

func (b *backupsServiceHandler) ListByServer(ctx context.Context, options *ListOptions, serverId string) ([]Backup, error) {
	path := fmt.Sprintf("%s/vm/backups/%s", backupsPath, serverId)

	req, err := b.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	backups := new(ListBackupsRoot)
	if err = b.client.Do(ctx, req, backups); err != nil {
		return nil, err
	}

	return backups.Data, nil
}

func (b *backupsServiceHandler) CreateServerByBackup(ctx context.Context, backupIdentifier string) error {
	path := fmt.Sprintf("%s/backups/create", backupsPath)

	createServerReq := struct {
		BackupIdentifier string `json:"backupIdentifier"`
	}{
		BackupIdentifier: backupIdentifier,
	}

	req, err := b.client.NewRequest(ctx, http.MethodPost, path, &createServerReq)
	if err != nil {
		return err
	}

	return b.client.Do(ctx, req, nil)
}

func (b *backupsServiceHandler) Get(ctx context.Context, identifer string) (*Backup, error) {
	path := fmt.Sprintf("%s/backup/%s", backupsPath, identifer)

	req, err := b.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	backup := new(GetBackupsRoot)
	if err = b.client.Do(ctx, req, backup); err != nil {
		return nil, err
	}

	return &backup.Data.Backup, nil
}

func (b *backupsServiceHandler) EnableAutoBackup(ctx context.Context, enableAutoReq *EnableAutoBackupReq) error {
	path := fmt.Sprintf("%s/backups/enable/auto", backupsPath)

	req, err := b.client.NewRequest(ctx, http.MethodPost, path, enableAutoReq)
	if err != nil {
		return err
	}

	return b.client.Do(ctx, req, nil)

}

func (b *backupsServiceHandler) Rename(ctx context.Context, backupIdentifier string, newName string) error {
	path := fmt.Sprintf("%s/backups/update/name", backupsPath)

	renameReq := struct {
		Identifier string `json:"identifier"`
		NewName    string `json:"new_name"`
	}{
		Identifier: backupIdentifier,
		NewName:    newName,
	}

	req, err := b.client.NewRequest(ctx, http.MethodPut, path, &renameReq)
	if err != nil {
		return err
	}

	return b.client.Do(ctx, req, nil)
}

// backup policies apis

type BackupPolicyVm struct {
	Name          string `json:"name"`
	Identifier    string `json:"identifier"`
	Category      string `json:"category"`
	Fullname      string `json:"fullname"`
	ImgIdentifier string `json:"imgIdentifier"`
	Type          string `json:"type"`
}

type BackupPolicy struct {
	Name       string           `json:"name"`
	Identifier string           `json:"identifier"`
	CreatedOn  string           `json:"created_on"`
	CreatedBy  string           `json:"created_by"`
	BackupPlan string           `json:"backupPlan"`
	PlanEvery  int              `json:"planEvery"`
	Keep       int              `json:"keep"`
	Disabled   int              `json:"disabled"`
	UserId     int              `json:"userId"`
	Vms        []BackupPolicyVm `json:"vms"`
}

type BackupPolicyListDetail struct {
	Name       string `json:"name"`
	Identifier string `json:"identifier"`
	CreatedOn  string `json:"created_on"`
	CreatedBy  string `json:"created_by"`
	BackupPlan string `json:"backupPlan"`
	PlanEvery  int    `json:"planEvery"`
	Keep       int    `json:"keep"`
	Disabled   int    `json:"disabled"`
	VmsCount   int    `json:"vmsCount"`
	UserId     int    `json:"userId"`
}
type ListBackupPoliciesRoot struct {
	Error bool `json:"error"`
	Data  struct {
		Rows []BackupPolicyListDetail `json:"rows"`
	} `json:"data"`
}

type GetBackupPolicyRoot struct {
	Error bool            `json:"error"`
	Data  json.RawMessage `json:"data"`
}

type CreateBackupPolicyReq struct {
	Name       string   `json:"name"`
	BackupPlan string   `json:"backupPlan"`
	PlanEvery  string   `json:"planEvery"`
	Keep       string   `json:"keep"`
	Vms        []string `json:"vms"`
	Tags       []string `json:"tags"`
}

func (b *backupsServiceHandler) ListBackupPolicies(ctx context.Context, options *ListOptions) ([]BackupPolicyListDetail, error) {
	path := fmt.Sprintf("%s/backups/policy/all", backupsPath)

	req, err := b.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	policies := new(ListBackupPoliciesRoot)
	if err = b.client.Do(ctx, req, policies); err != nil {
		return nil, err
	}

	return policies.Data.Rows, nil
}

// GetBackupPolicy fetches a single backup policy by its UUID identifier. The
// route is singular (/backup/policy/:id) even though the create and list routes
// are plural. When the policy does not exist the API returns HTTP 200 with a
// literal `"data": false`, so (nil, nil) is returned to signal not-found.
func (b *backupsServiceHandler) GetBackupPolicy(ctx context.Context, identifier string) (*BackupPolicy, error) {
	path := fmt.Sprintf("%s/backup/policy/%s", backupsPath, identifier)

	req, err := b.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	root := new(GetBackupPolicyRoot)
	if err = b.client.Do(ctx, req, root); err != nil {
		return nil, err
	}

	if isNullData(root.Data) {
		return nil, nil
	}

	policy := new(BackupPolicy)
	if err = json.Unmarshal(root.Data, policy); err != nil {
		return nil, err
	}

	return policy, nil
}

func (b *backupsServiceHandler) CreateBackupPolicy(ctx context.Context, createReq *CreateBackupPolicyReq) error {
	path := fmt.Sprintf("%s/backups/policy/create", backupsPath)

	// The API validates tags and vms as arrays; nil slices marshal to null and
	// are rejected, so ensure they are sent as empty arrays.
	if createReq.Tags == nil {
		createReq.Tags = []string{}
	}
	if createReq.Vms == nil {
		createReq.Vms = []string{}
	}

	req, err := b.client.NewRequest(ctx, http.MethodPost, path, createReq)
	if err != nil {
		return err
	}

	return b.client.Do(ctx, req, nil)
}

func (b *backupsServiceHandler) DeleteBackupPolicy(ctx context.Context, policyId, identifier string) error {
	path := fmt.Sprintf("%s/backup/policy/%s", backupsPath, identifier)

	deleteBackup := struct {
		PolicyId string `json:"policyId"`
	}{
		PolicyId: policyId,
	}

	req, err := b.client.NewRequest(ctx, http.MethodDelete, path, deleteBackup)
	if err != nil {
		return err
	}

	return b.client.Do(ctx, req, nil)
}

func (b *backupsServiceHandler) ManageRetainBackupPolicy(ctx context.Context, policyId string, keep int) error {
	path := fmt.Sprintf("%s/backups/policy/keep", backupsPath)

	manageBackup := struct {
		PolicyId string `json:"policyId"`
		Keep     int    `json:"keep"`
	}{
		PolicyId: policyId,
		Keep:     keep,
	}

	req, err := b.client.NewRequest(ctx, http.MethodPost, path, manageBackup)
	if err != nil {
		return err
	}

	return b.client.Do(ctx, req, nil)
}

func (b *backupsServiceHandler) AttachBackupPolicy(ctx context.Context, policyId string, vms []string) error {
	path := fmt.Sprintf("%s/backups/policy/attach", backupsPath)

	attachBackup := struct {
		PolicyId string   `json:"policyId"`
		Vms      []string `json:"vms"`
	}{
		PolicyId: policyId,
		Vms:      vms,
	}

	req, err := b.client.NewRequest(ctx, http.MethodPost, path, attachBackup)
	if err != nil {
		return err
	}

	return b.client.Do(ctx, req, nil)
}

func (b *backupsServiceHandler) DetachBackupPolicy(ctx context.Context, policyId string, vms []string) error {
	path := fmt.Sprintf("%s/backups/policy/detach", backupsPath)

	detachBackup := struct {
		PolicyId string   `json:"policyId"`
		Vms      []string `json:"vms"`
	}{
		PolicyId: policyId,
		Vms:      vms,
	}

	req, err := b.client.NewRequest(ctx, http.MethodPost, path, detachBackup)
	if err != nil {
		return err
	}

	return b.client.Do(ctx, req, nil)
}
