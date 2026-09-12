package govpsie

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

var firewallGroupBasePath = "/apps/v2/firewall"

type FirewallGroupService interface {
	Create(ctx context.Context, groupName string, firewallUpdateReq []FirewallUpdateReq) error
	List(ctx context.Context, options *ListOptions) ([]FirewallGroupListData, error)
	Get(ctx context.Context, fwGroupId string) (*FirewallGroupDetailData, error)
	Delete(ctx context.Context, fwGroupId, reason, note string) error
	RenameGroup(ctx context.Context, fwGroupId, groupName string) error
	UpdateRules(ctx context.Context, fwGroupId string, rules []FirewallUpdateReq) error
	DeleteRules(ctx context.Context, fwGroupId string, ruleIdentifiers []string) error
	AssignToVpsie(ctx context.Context, groupId, vmId string) error
	DetachFromVpsie(ctx context.Context, groupId, vmId string) error
	AttachToVpsie(ctx context.Context, groupId, vmId string) error
	DeleteFirewallGroupOfServer(ctx context.Context, groupId, vmId string) error
	GetFirewallGroup(ctx context.Context, fwGroupId string) (*FirewallGroupDetailData, error)
}

type firewallGroupServiceHandler struct {
	client *Client
}

var _ FirewallGroupService = &firewallGroupServiceHandler{}

type ListFirewallGroupsRoot struct {
	Error bool                    `json:"error"`
	Data  []FirewallGroupListData `json:"data"`
	Total int64                   `json:"total"`
}

type GetFirewallGroupRoot struct {
	Error bool                    `json:"error"`
	Data  FirewallGroupDetailData `json:"data"`
}

type FirewallGroupListData struct {
	UserName      string `json:"user_name"`
	ID            int64  `json:"id"`
	GroupName     string `json:"group_name"`
	Identifier    string `json:"identifier"`
	CreatedOn     string `json:"created_on"`
	UpdatedOn     string `json:"updated_on"`
	InboundCount  int64  `json:"inbound_count"`
	OutboundCount int64  `json:"outbound_count"`
	Vms           int64  `json:"vms"`
	CreatedBy     int64  `json:"created_by"`

	// Rules is a single object grouped into inBound/outBound arrays, not a slice.
	Rules   FirewallRules `json:"rules"`
	VmsData []VmsData     `json:"vmsData"`
}
type FirewallGroupDetailData struct {
	Group FirewallGroup `json:"group"`
	Rules FirewallRules `json:"rules"`
	Vms   []VmsData     `json:"vms"`
	Count int64         `json:"count"`
}

type VmsData struct {
	Hostname   string `json:"hostname"`
	Identifier string `json:"identifier"`
	Fullname   string `json:"fullname"`
	Category   string `json:"category"`
}

type FirewallRules struct {
	InBound  []InBoundFirewallRules  `json:"inBound"`
	OutBound []OutBoundFirewallRules `json:"outBound"`
}

type OutBoundFirewallRules struct {
	ID         int64     `json:"id"`
	GroupID    int64     `json:"group_id"`
	UserID     int64     `json:"user_id"`
	Action     string    `json:"action"`
	Type       string    `json:"type"`
	Comment    string    `json:"comment"`
	Dest       []string  `json:"dest,omitempty"`
	Dport      string    `json:"dport"`
	Proto      string    `json:"proto"`
	Source     []string  `json:"source"`
	Sport      string    `json:"sport"`
	Enable     int64     `json:"enable"`
	Iface      string    `json:"iface,omitempty"`
	Log        string    `json:"log,omitempty"`
	Macro      string    `json:"macro,omitempty"`
	Identifier string    `json:"identifier"`
	CreatedOn  time.Time `json:"created_on"`
	UpdatedOn  time.Time `json:"updated_on"`
}

type InBoundFirewallRules struct {
	ID         int64     `json:"id"`
	GroupID    int64     `json:"group_id"`
	UserID     int64     `json:"user_id"`
	Action     string    `json:"action"`
	Type       string    `json:"type"`
	Comment    string    `json:"comment"`
	Dest       []string  `json:"dest"`
	Dport      string    `json:"dport"`
	Proto      string    `json:"proto"`
	Source     []string  `json:"source,omitempty"`
	Sport      string    `json:"sport"`
	Enable     int64     `json:"enable"`
	Iface      string    `json:"iface,omitempty"`
	Log        string    `json:"log,omitempty"`
	Macro      string    `json:"macro,omitempty"`
	Identifier string    `json:"identifier"`
	CreatedOn  time.Time `json:"created_on"`
	UpdatedOn  time.Time `json:"updated_on"`
}

type FirewallGroup struct {
	UserName      string `json:"user_name"`
	ID            int64  `json:"id"`
	GroupName     string `json:"group_name"`
	Identifier    string `json:"identifier"`
	CreatedOn     string `json:"created_on"`
	UpdatedOn     string `json:"updated_on"`
	InboundCount  int64  `json:"inbound_count"`
	OutboundCount int64  `json:"outbound_count"`
	Vms           int64  `json:"vms"`
	CreatedBy     int64  `json:"created_by"`
}

type FirewallUpdateReq struct {
	Action  string   `json:"action"`
	Type    string   `json:"type"`
	Dport   string   `json:"dport"`
	Proto   string   `json:"proto"`
	Source  []string `json:"source,omitempty"`
	Sport   string   `json:"sport"`
	Enable  int64    `json:"enable"`
	Macro   string   `json:"macro"`
	Comment string   `json:"comment"`
	Dest    []string `json:"dest,omitempty"`
	Log     string   `json:"log,omitempty"`
	// Identifier targets an existing rule on an update; leave empty to add a new rule.
	Identifier string `json:"identifier,omitempty"`
}

func (f *firewallGroupServiceHandler) Create(ctx context.Context, groupName string, firewallUpdateReq []FirewallUpdateReq) error {
	fwGroupReq := struct {
		GroupName string              `json:"groupName"`
		Rules     []FirewallUpdateReq `json:"rules,omitempty"`
	}{
		GroupName: groupName,
		Rules:     firewallUpdateReq,
	}

	path := fmt.Sprintf("%s/create/group", firewallGroupBasePath)

	req, err := f.client.NewRequest(ctx, http.MethodPost, path, &fwGroupReq)
	if err != nil {
		return err
	}

	return f.client.Do(ctx, req, nil)
}

func (f *firewallGroupServiceHandler) List(ctx context.Context, options *ListOptions) ([]FirewallGroupListData, error) {
	path := fmt.Sprintf("%s/groups", firewallGroupBasePath)

	req, err := f.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	fwGroups := new(ListFirewallGroupsRoot)

	if err = f.client.Do(ctx, req, &fwGroups); err != nil {
		return nil, err
	}

	return fwGroups.Data, nil

}

func (f *firewallGroupServiceHandler) Get(ctx context.Context, fwGroupId string) (*FirewallGroupDetailData, error) {
	path := fmt.Sprintf("%s/group/%s", firewallGroupBasePath, fwGroupId)

	req, err := f.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	fwGroup := new(GetFirewallGroupRoot)

	if err = f.client.Do(ctx, req, fwGroup); err != nil {
		return nil, err
	}

	return &fwGroup.Data, nil
}

// Delete removes a firewall group by its UUID identifier. The API requires a
// deleteStatistic payload; omitting it makes the handler dereference an
// undefined field and fail with HTTP 500.
func (f *firewallGroupServiceHandler) Delete(ctx context.Context, fwGroupId, reason, note string) error {
	path := fmt.Sprintf("%s/delete/group", firewallGroupBasePath)

	delReq := struct {
		GroupId         string `json:"groupId"`
		DeleteStatistic struct {
			Reason string `json:"reason"`
			Note   string `json:"note"`
		} `json:"deleteStatistic"`
	}{GroupId: fwGroupId}
	delReq.DeleteStatistic.Reason = reason
	delReq.DeleteStatistic.Note = note

	req, err := f.client.NewRequest(ctx, http.MethodDelete, path, &delReq)
	if err != nil {
		return err
	}

	return f.client.Do(ctx, req, nil)
}

// RenameGroup changes a firewall group's name (keyed on the group UUID).
func (f *firewallGroupServiceHandler) RenameGroup(ctx context.Context, fwGroupId, groupName string) error {
	path := fmt.Sprintf("%s/group/rename", firewallGroupBasePath)

	renameReq := struct {
		GroupId   string `json:"groupId"`
		GroupName string `json:"groupName"`
	}{GroupId: fwGroupId, GroupName: groupName}

	req, err := f.client.NewRequest(ctx, http.MethodPost, path, &renameReq)
	if err != nil {
		return err
	}

	return f.client.Do(ctx, req, nil)
}

// UpdateRules adds or updates rules on a group (keyed on the group UUID). Rules
// carrying an Identifier are updated in place; rules without one are added.
func (f *firewallGroupServiceHandler) UpdateRules(ctx context.Context, fwGroupId string, rules []FirewallUpdateReq) error {
	path := fmt.Sprintf("%s/group/rules/update", firewallGroupBasePath)

	updateReq := struct {
		GroupId string              `json:"groupId"`
		Rules   []FirewallUpdateReq `json:"rules"`
	}{GroupId: fwGroupId, Rules: rules}

	req, err := f.client.NewRequest(ctx, http.MethodPut, path, &updateReq)
	if err != nil {
		return err
	}

	return f.client.Do(ctx, req, nil)
}

// DeleteRules removes rules from a group by their UUID identifiers.
func (f *firewallGroupServiceHandler) DeleteRules(ctx context.Context, fwGroupId string, ruleIdentifiers []string) error {
	path := fmt.Sprintf("%s/group/rules/delete", firewallGroupBasePath)

	if ruleIdentifiers == nil {
		ruleIdentifiers = []string{}
	}
	delReq := struct {
		GroupId string   `json:"groupId"`
		Rules   []string `json:"rules"`
	}{GroupId: fwGroupId, Rules: ruleIdentifiers}

	req, err := f.client.NewRequest(ctx, http.MethodDelete, path, &delReq)
	if err != nil {
		return err
	}

	return f.client.Do(ctx, req, nil)
}

func (f *firewallGroupServiceHandler) AssignToVpsie(ctx context.Context, groupId string, vmId string) error {
	path := fmt.Sprintf("%s/setGroupVm", firewallGroupBasePath)

	assignReq := struct {
		VmID    string `json:"vmId"`
		GroupID string `json:"groupId"`
	}{
		VmID:    vmId,
		GroupID: groupId,
	}

	req, err := f.client.NewRequest(ctx, http.MethodPost, path, assignReq)
	if err != nil {
		return err
	}

	return f.client.Do(ctx, req, nil)
}

func (f *firewallGroupServiceHandler) AttachToVpsie(ctx context.Context, groupId, vmId string) error {
	path := fmt.Sprintf("%s/attach/group", firewallGroupBasePath)

	assignReq := struct {
		VmID    string `json:"vmId"`
		GroupID string `json:"groupId"`
	}{
		VmID:    vmId,
		GroupID: groupId,
	}

	req, err := f.client.NewRequest(ctx, http.MethodPost, path, assignReq)
	if err != nil {
		return err
	}

	return f.client.Do(ctx, req, nil)
}

func (f *firewallGroupServiceHandler) DetachFromVpsie(ctx context.Context, groupId, vmId string) error {
	path := fmt.Sprintf("%s/detach/group", firewallGroupBasePath)

	assignReq := struct {
		VmID    string `json:"vmId"`
		GroupID string `json:"groupId"`
	}{
		VmID:    vmId,
		GroupID: groupId,
	}

	req, err := f.client.NewRequest(ctx, http.MethodPost, path, assignReq)
	if err != nil {
		return err
	}

	return f.client.Do(ctx, req, nil)
}

func (f *firewallGroupServiceHandler) GetFirewallGroup(ctx context.Context, fwGroupId string) (*FirewallGroupDetailData, error) {
	path := fmt.Sprintf("%s/group/%s", firewallGroupBasePath, fwGroupId)

	req, err := f.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	fwGroup := new(GetFirewallGroupRoot)

	if err = f.client.Do(ctx, req, fwGroup); err != nil {
		return nil, err
	}

	return &fwGroup.Data, nil
}

func (f *firewallGroupServiceHandler) DeleteFirewallGroupOfServer(ctx context.Context, groupId, vmId string) error {
	path := fmt.Sprintf("%s/firewall/removeGroupVm", firewallGroupBasePath)

	delReq := struct {
		GroupId string `json:"groupId"`
		VmId    string `json:"vmId"`
	}{
		GroupId: groupId,
		VmId:    vmId,
	}

	req, err := f.client.NewRequest(ctx, http.MethodDelete, path, &delReq)
	if err != nil {
		return err
	}

	return f.client.Do(ctx, req, nil)
}
