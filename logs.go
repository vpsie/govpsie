package govpsie

import (
	"context"
	"fmt"
	"net/http"
)

var logsPath = "/apps/v2/logs"

type LogsService interface {
	ListActivityLogs(ctx context.Context, options *ListOptions) ([]ActivityLog, error)
	ListBillingLogs(ctx context.Context, options *ListOptions) ([]BillingLog, error)
	ListAuditLogs(ctx context.Context, options *ListOptions) ([]AuditLog, error)
	ListVPSieLogs(ctx context.Context, options *ListOptions) ([]VmLog, error)
}

type logsServiceHandler struct {
	client *Client
}

type AuditLog struct {
	ID          int    `json:"id"`
	UserID      int    `json:"user_id"`
	UserAgent   string `json:"user_agent"`
	IPAddress   string `json:"ip_address"`
	Country     string `json:"country"`
	City        string `json:"city"`
	Action      string `json:"action"`
	Description string `json:"description"`
	CreatedOn   string `json:"created_on"`
	State       string `json:"state"`
}

type VmLog struct {
	ID          int    `json:"id"`
	UserID      int    `json:"user_id"`
	EntityID    int    `json:"entity_id"`
	EntityType  string `json:"entity_type"`
	Action      string `json:"action"`
	Description string `json:"description"`
	CreatedOn   string `json:"created_on"`
	UserAgent   string `json:"user_agent"`
	IPAddress   string `json:"ip_address"`
	Country     string `json:"country"`
	BoxID       int    `json:"box_id"`
}

type BillingLog struct {
	ID                    int    `json:"id"`
	UserID                int    `json:"user_id"`
	TransactionID         string `json:"transaction_id"`
	TransactionSerialized string `json:"transaction_serialized"`
	Message               string `json:"message"`
	Amount                int    `json:"amount"`
	CreatedOn             string `json:"created_on"`
	PpAccount             string `json:"pp_account"`
	PpFullname            string `json:"pp_full_name"`
	PpCountryCode         string `json:"pp_country_code"`
	TransactionOriginIP   string `json:"transaction_origin_ip"`
	RiskScore             string `json:"risk_score"`
	Purchaselogscol       string `json:"purchaselogscol"`
	MaxmindResponse       string `json:"maxmind_response"`
}

type ActivityLog struct {
	UserID      int    `json:"user_id"`
	CreatedBy   string `json:"created_by"`
	BoxID       int    `json:"box_id"`
	BackupID    int    `json:"backup_id"`
	Description string `json:"description"`
	CreatedOn   string `json:"created_on"`
}

type ListActivityLogsRoot struct {
	Error bool          `json:"error"`
	Data  []ActivityLog `json:"data"`
	Total int           `json:"total"`
}

type ListAuditLogsRoot struct {
	Error bool       `json:"error"`
	Data  []AuditLog `json:"data"`
	Total int        `json:"total"`
}

type ListBillingLogsRoot struct {
	Error bool         `json:"error"`
	Data  []BillingLog `json:"data"`
	Total int          `json:"total"`
}

type ListVmLogsRoot struct {
	Error bool    `json:"error"`
	Data  []VmLog `json:"data"`
	Total int     `json:"total"`
}

var _ LogsService = &logsServiceHandler{}

func (l *logsServiceHandler) ListActivityLogs(ctx context.Context, options *ListOptions) ([]ActivityLog, error) {
	path := fmt.Sprintf("%s/activity?offset=%d&limit=%d", logsPath, options.Page, options.PerPage)

	req, err := l.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	root := new(ListActivityLogsRoot)
	if err = l.client.Do(ctx, req, root); err != nil {
		return nil, err
	}

	return root.Data, nil
}

func (l *logsServiceHandler) ListBillingLogs(ctx context.Context, options *ListOptions) ([]BillingLog, error) {
	path := fmt.Sprintf("%s/billing?offset=%d&limit=%d", logsPath, options.Page, options.PerPage)

	req, err := l.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	root := new(ListBillingLogsRoot)
	if err = l.client.Do(ctx, req, root); err != nil {
		return nil, err
	}

	return root.Data, nil
}

func (l *logsServiceHandler) ListAuditLogs(ctx context.Context, options *ListOptions) ([]AuditLog, error) {
	path := fmt.Sprintf("%s/audit?offset=%d&limit=%d", logsPath, options.Page, options.PerPage)

	req, err := l.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	root := new(ListAuditLogsRoot)
	if err = l.client.Do(ctx, req, root); err != nil {
		return nil, err
	}

	return root.Data, nil
}
func (l *logsServiceHandler) ListVPSieLogs(ctx context.Context, options *ListOptions) ([]VmLog, error) {
	path := fmt.Sprintf("%s/vm?offset=%d&limit=%d", logsPath, options.Page, options.PerPage)

	req, err := l.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	root := new(ListVmLogsRoot)
	if err = l.client.Do(ctx, req, root); err != nil {
		return nil, err
	}

	return root.Data, nil
}
