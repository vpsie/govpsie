package govpsie

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

const monitoringPath = "/apps/v2/monitoring"

type MonitoringService interface {
	ListMonitoringRule(ctx context.Context, options *ListOptions) ([]MonitoringRule, error)
	GetMonitoringRule(ctx context.Context, ruleIdentifier string) (*MonitoringRule, error)
	CreateRule(ctx context.Context, createReq *CreateMonitoringRuleReq) error
	ToggleMonitoringRuleStatus(ctx context.Context, status, ruleIdentifier string) error
	AttachVms(ctx context.Context, ruleIdentifier string, vms []string) error
	DetachVms(ctx context.Context, ruleIdentifier string, vms []string) error
	DeleteMonitoringRule(ctx context.Context, ruleIdentifier string) error
}

type monitoringServiceHandler struct {
	client *Client
}

var _ MonitoringService = &monitoringServiceHandler{}

// MonitoringAction is a single alert action attached to a metric rule.
type MonitoringAction struct {
	ID          int    `json:"id"`
	RuleID      int    `json:"rule_id"`
	ActionName  string `json:"action_name"`
	ActionKey   string `json:"action_key"`
	MonitorType string `json:"monitor_type"`
	Email       string `json:"email"`
	Value       string `json:"value"`
	Identifier  string `json:"identifier"`
}

// MonitoringMetric is a single metric/condition within a monitoring rule.
type MonitoringMetric struct {
	ID            int                `json:"id"`
	MetricType    string             `json:"metric_type"`
	Condition     string             `json:"condition"`
	Threshold     int                `json:"threshold"`
	ThresholdType string             `json:"threshold_type"`
	Period        int                `json:"period"`
	Status        int                `json:"status"`
	Actions       []MonitoringAction `json:"actions"`
}

// MonitoringVM is a VM attached to a monitoring rule.
type MonitoringVM struct {
	Identifier string `json:"identifier"`
	Hostname   string `json:"hostname"`
	Fullname   string `json:"fullname"`
}

// MonitoringRule is a monitoring rule as returned by list/get. The per-metric
// data lives in Metrics (the top-level metric_type/condition/threshold columns
// are legacy and null on current rules).
type MonitoringRule struct {
	ID            int                `json:"id"`
	UserId        int                `json:"user_id"`
	RuleName      string             `json:"rule_name"`
	Status        int                `json:"status"`
	CreatedOn     string             `json:"created_on"`
	Frequency     int                `json:"frequency"`
	LastAlertDate string             `json:"last_alert_date"`
	Identifier    string             `json:"identifier"`
	IsDeleted     int                `json:"is_deleted"`
	CreatedBy     string             `json:"created_by"`
	Metrics       []MonitoringMetric `json:"rules"`
	Vms           []MonitoringVM     `json:"vms"`
}

type ListMonitoringRuleRoot struct {
	Error bool             `json:"error"`
	Data  []MonitoringRule `json:"data"`
	// The monitoring endpoint returns total as an array (e.g. [{"count":0}]),
	// unlike other list endpoints that use a plain integer.
	Total []struct {
		Count int `json:"count"`
	} `json:"total"`
}

// CreateMonitoringActionReq is an action in a create/update request.
type CreateMonitoringActionReq struct {
	ActionName string `json:"actionName"`
	ActionKey  string `json:"actionKey"`
	Email      string `json:"email,omitempty"`
	Value      string `json:"value,omitempty"`
}

// CreateMonitoringMetricReq is a metric/condition in a create request.
type CreateMonitoringMetricReq struct {
	MetricType    string                      `json:"metricType"`
	Condition     string                      `json:"condition"`
	Threshold     string                      `json:"threshold"`
	ThresholdType string                      `json:"thresholdType"`
	Period        string                      `json:"period"`
	Status        string                      `json:"status"`
	Actions       []CreateMonitoringActionReq `json:"actions"`
}

// CreateMonitoringRuleReq is the body for creating a monitoring rule. The
// metric/condition data lives inside Rules; the legacy top-level flat shape is
// no longer accepted by the API.
type CreateMonitoringRuleReq struct {
	RuleName  string                      `json:"ruleName"`
	Status    string                      `json:"status"`
	Frequency string                      `json:"frequency"`
	Rules     []CreateMonitoringMetricReq `json:"rules"`
	Vms       []string                    `json:"vms"`
	Tags      []string                    `json:"tags,omitempty"`
}

func (s *monitoringServiceHandler) ListMonitoringRule(ctx context.Context, options *ListOptions) ([]MonitoringRule, error) {
	path := fmt.Sprintf("%s/rules", monitoringPath)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var root ListMonitoringRuleRoot
	if err := s.client.Do(ctx, req, &root); err != nil {
		return nil, err
	}

	return root.Data, nil
}

// GetMonitoringRule fetches a single rule (with its metrics, actions and VMs) by
// its UUID identifier. There is no dedicated single-GET route, so the list
// endpoint is filtered by ruleIdentifier. Returns (nil, nil) when not found.
func (s *monitoringServiceHandler) GetMonitoringRule(ctx context.Context, ruleIdentifier string) (*MonitoringRule, error) {
	path := fmt.Sprintf("%s/rules?ruleIdentifier=%s", monitoringPath, url.QueryEscape(ruleIdentifier))

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var root ListMonitoringRuleRoot
	if err := s.client.Do(ctx, req, &root); err != nil {
		return nil, err
	}

	for i := range root.Data {
		if root.Data[i].Identifier == ruleIdentifier {
			return &root.Data[i], nil
		}
	}
	return nil, nil
}

func (s *monitoringServiceHandler) CreateRule(ctx context.Context, createReq *CreateMonitoringRuleReq) error {
	path := fmt.Sprintf("%s/rules/add", monitoringPath)

	if createReq.Vms == nil {
		createReq.Vms = []string{}
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, createReq)
	if err != nil {
		return err
	}

	return s.client.Do(ctx, req, nil)
}

func (s *monitoringServiceHandler) ToggleMonitoringRuleStatus(ctx context.Context, status, ruleIdentifier string) error {
	path := fmt.Sprintf("%s/rules/edit", monitoringPath)

	toggleReq := struct {
		Status         string `json:"status"`
		RuleIdentifier string `json:"ruleIdentifier"`
	}{
		Status:         status,
		RuleIdentifier: ruleIdentifier,
	}

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, &toggleReq)
	if err != nil {
		return err
	}

	return s.client.Do(ctx, req, nil)
}

func (s *monitoringServiceHandler) AttachVms(ctx context.Context, ruleIdentifier string, vms []string) error {
	return s.vmsAction(ctx, "attach", ruleIdentifier, vms)
}

func (s *monitoringServiceHandler) DetachVms(ctx context.Context, ruleIdentifier string, vms []string) error {
	return s.vmsAction(ctx, "detach", ruleIdentifier, vms)
}

func (s *monitoringServiceHandler) vmsAction(ctx context.Context, action, ruleIdentifier string, vms []string) error {
	path := fmt.Sprintf("%s/rules/%s", monitoringPath, action)

	if vms == nil {
		vms = []string{}
	}
	body := struct {
		Vms            []string `json:"vms"`
		RuleIdentifier string   `json:"ruleIdentifier"`
	}{Vms: vms, RuleIdentifier: ruleIdentifier}

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, &body)
	if err != nil {
		return err
	}

	return s.client.Do(ctx, req, nil)
}

func (s *monitoringServiceHandler) DeleteMonitoringRule(ctx context.Context, ruleIdentifier string) error {
	path := fmt.Sprintf("%s/rules/%s", monitoringPath, ruleIdentifier)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return err
	}

	return s.client.Do(ctx, req, nil)
}
