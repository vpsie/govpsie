package govpsie

import (
	"context"
	"fmt"
	"net/http"
)

var snapshotBasePath = "/apps/v2/snapshot"

type SnapshotService interface {
	List(ctx context.Context, options *ListOptions) ([]Snapshot, error)
	Create(ctx context.Context, name, vmIdentifier string) error
	ListByVm(ctx context.Context, options *ListOptions, vmIdentifier string) ([]Snapshot, error)
	Rollback(ctx context.Context, snapshotIdentifier string) error
	EnableAuto(ctx context.Context, enableReq *EnableAutoSnapshotReq) error
	Delete(ctx context.Context, snapshotIdentifier, reason, note string) error
}

type snapshotServiceHandler struct {
	client *Client
}

var _ SnapshotService = &snapshotServiceHandler{}

type Snapshot struct {
	Hostname     string `json:"hostname"`
	Name         string `json:"name"`
	Identifier   string `json:"identifier"`
	BackupKey    string `json:"backupKey"`
	State        string `json:"state"`
	DcIdentifier string `json:"dcIdentifier"`
	Daily        int    `json:"daily"`
	IsSnapshot   int    `json:"is_snapshot"`
	VmIdentifier string `json:"vmIdentifier"`
	BackupSHA1   string `json:"backupsha1"`
	OSIdentifier string `json:"os_identifier"`
	UserID       int    `json:"user_id"`
}

type ListSnapshotsRoot struct {
	Error bool       `json:"error"`
	Data  []Snapshot `json:"data"`
	Total int        `json:"total"`
}

type EnableAutoSnapshotReq struct {
	VMIdentifier    string   `json:"vmIdentifier"`
	VmId            int      `json:"vmId"`
	Period          string   `json:"period"`
	DailySnapshot   int      `json:"dailySnapshot"`
	WeeklySnapshot  int      `json:"weeklySnapshot"`
	MonthlySnapshot int      `json:"monthlySnapshot"`
	Tags            []string `json:"tags"`
}

func (s *snapshotServiceHandler) List(ctx context.Context, options *ListOptions) ([]Snapshot, error) {
	path := fmt.Sprintf("%s?offset=%d&limit%d", snapshotBasePath, options.Page, options.PerPage)

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

func (s *snapshotServiceHandler) Create(ctx context.Context, name, vmIdentifier string) error {
	path := fmt.Sprintf("%s/add", snapshotBasePath)
	createSnapshotReq := struct {
		Name         string `json:"name"`
		VMIdentifier string `json:"vmIdentifier"`
	}{
		Name:         name,
		VMIdentifier: vmIdentifier,
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, &createSnapshotReq)
	if err != nil {
		return err
	}

	return s.client.Do(ctx, req, nil)

}

func (s *snapshotServiceHandler) ListByVm(ctx context.Context, options *ListOptions, vmIdentifier string) ([]Snapshot, error) {
	path := fmt.Sprintf("/apps/v2/vm/snapshot/%s?offset=%d&limit%d", vmIdentifier, options.Page, options.PerPage)

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
		DeleteStatistic struct {
			Reason string `json:"reason"`
			Note  string `json:"note"`
		} `json:"deleteStatistic"`
	}{
		SnapshotIdentifier: snapshotIdentifier,
		DeleteStatistic: struct {
			Reason string `json:"reason"`
			Note  string `json:"note"`
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
