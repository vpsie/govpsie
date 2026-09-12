package govpsie

import (
	"context"
	"fmt"
	"net/http"
)

const monitoringPath = "/apps/v2/monitoring"

type MonitoringService interface {
	ListMonitoringRule(ctx context.Context, options *ListOptions) ([]MonitoringRule, error)
	CreateRule(ctx context.Context, createReq *CreateMonitoringRuleReq) error
	ToggleMonitoringRuleStatus(ctx context.Context, status, ruleIdentifier string) error
	DeleteMonitoringRule(ctx context.Context, ruleIdentifier string) error
}

type monitoringServiceHandler struct {
	client *Client
}

var _ MonitoringService = &monitoringServiceHandler{}

type MonitoringRule struct {
	ID            int    `json:"id"`
	UserId        int    `json:"user_id"`
	MetricType    string `json:"metric_type"`
	RuleName      string `json:"rule_name"`
	Condition     string `json:"condition"`
	Email         string `json:"email"`
	Threshold     int    `json:"threshold"`
	ThresholdType string `json:"threshold_type"`
	Period        int    `json:"period"`
	Status        int    `json:"status"`
	CreatedOn     string `json:"created_on"`
	Frequency     int    `json:"frequency"`
	LastAlertDate string `json:"last_alert_date"`
	Identifier    string `json:"identifier"`
	IsDeleted     int    `json:"is_deleted"`
	CreatedBY     string `json:"created_by"`
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

type CreateMonitoringRuleReq struct {
	MetricType    string `json:"metricType"`
	RuleName      string `json:"ruleName"`
	Condition     string `json:"condition"`
	ThresholdType string `json:"thresholdType"`
	ThresholdId   string `json:"thresholdId"`
	Period        string `json:"period"`
	Frequency     string `json:"frequency"`
	Status        string `json:"status"`
	Threshold     string `json:"threshold"`
	Actions       struct {
		Email      string `json:"email"`
		ActionKey  string `json:"actionKey"`
		ActionName string `json:"actionName"`
	} `json:"actions"`
	Vms  []string `json:"vms"`
	Tags []string `json:"tags"`
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

func (s *monitoringServiceHandler) CreateRule(ctx context.Context, createReq *CreateMonitoringRuleReq) error {
	path := fmt.Sprintf("%s/rules/add", monitoringPath)

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, &createReq)
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

func (s *monitoringServiceHandler) DeleteMonitoringRule(ctx context.Context, ruleIdentifier string) error {
	path := fmt.Sprintf("%s/rules/%s", monitoringPath, ruleIdentifier)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return err
	}

	return s.client.Do(ctx, req, nil)
}
