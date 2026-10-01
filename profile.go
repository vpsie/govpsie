package govpsie

import (
	"context"
	"fmt"
	"net/http"
)

var profilePath = "/apps/v2/profile"

type ProfilesService interface {
	ListQuickActionOfUser(ctx context.Context, options *ListOptions) ([]QuickActions, error)
	ListQuickActionOfAccount(ctx context.Context, options *ListOptions) ([]QuickActions, error)
	SaveQuickActions(ctx context.Context, actions []int) error
	GetProfile(ctx context.Context) (*Profile, error)
	UpdateProfile(context.Context, UpdateProfileRequest) error
	GetPermissionGroups(ctx context.Context) ([]PermissionGroup, error)
	DeletePermissionGroup(ctx context.Context, groupId string) error
	CreatePermissionGroup(ctx context.Context, groupName string) error
	ChangePassword(ctx context.Context, oldPassword string, newPassword string) error
	UpdateBilling(ctx context.Context, billing BillingAddress) error
	ValidatePhone(ctx context.Context, phone string) error
	VerifyPhone(ctx context.Context, code string) error
	EnableTwofa(ctx context.Context) error
	DisableTwofa(ctx context.Context) error
}

type profilesServiceHandler struct {
	client *Client
}

var _ ProfilesService = &profilesServiceHandler{}

type ListActionOfUserRoot struct {
	Error bool           `json:"error"`
	Data  []QuickActions `json:"data"`
	Total int            `json:"total "`
}

type QuickActions struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Url         string `json:"url"`
	MappingID   int    `json:"mapping_id"`
	AutoEnabled int    `json:"auto_enabled"`
}

type GetPermissionGroupRoot struct {
	Error bool `json:"error"`
	Data  struct {
		OwnGroups    []PermissionGroup `json:"ownGroups"`
		InvitedGroup []PermissionGroup `json:"InvitedGroups"`
	} `json:"data"`
}

type PermissionGroup struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	CreatedOn string `json:"create_on"`
	UserId    int    `json:"user_id"`
	Perms     string `json:"perms"`
}

type BillingAddress struct {
	Address string `json:"address"`
	City    string `json:"city"`
	State   string `json:"state"`
	Zip     string `json:"zip"`
	Country string `json:"country"`
}

type GetProfileRoot struct {
	Error bool    `json:"error"`
	Data  Profile `json:"data"`
}

type Profile struct {
	Id             int        `json:"id"`
	Name           string     `json:"name"`
	FirstName      string     `json:"firstname"`
	LastName       string     `json:"lastname"`
	Email          string     `json:"email"`
	CurrentBalance float64    `json:"current_balance"`
	Identifier     string     `json:"identifier"`
	Status         string     `json:"status"`
	MonthlyCharge  FlexString `json:"monthly_charge"`
	Invitations    int        `json:"invitations"`
}

type UpdateProfileRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Title     string `json:"title"`
	Company   string `json:"company"`
	VatNumber string `json:"vatNumber"`
	Industry  string `json:"industry"`
	TimeZone  string `json:"timeZone"`
}

func (p *profilesServiceHandler) ListQuickActionOfUser(ctx context.Context, options *ListOptions) ([]QuickActions, error) {
	path := fmt.Sprintf("%s/user/quick/actions", profilePath)

	req, err := p.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	actions := new(ListActionOfUserRoot)
	if err = p.client.Do(ctx, req, actions); err != nil {
		return nil, err
	}

	return actions.Data, nil
}

func (p *profilesServiceHandler) ListQuickActionOfAccount(ctx context.Context, options *ListOptions) ([]QuickActions, error) {
	path := fmt.Sprintf("%s/quick/actions", profilePath)

	req, err := p.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	actions := new(ListActionOfUserRoot)
	if err = p.client.Do(ctx, req, actions); err != nil {
		return nil, err
	}

	return actions.Data, nil
}

func (p *profilesServiceHandler) SaveQuickActions(ctx context.Context, actions []int) error {
	path := fmt.Sprintf("%s/quick/actions/save", profilePath)

	req, err := p.client.NewRequest(ctx, http.MethodPut, path, actions)
	if err != nil {
		return err
	}

	return p.client.Do(ctx, req, nil)
}
func (p *profilesServiceHandler) GetProfile(ctx context.Context) (*Profile, error) {
	path := fmt.Sprintf("%s/user", profilePath)

	req, err := p.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	profile := new(GetProfileRoot)
	if err = p.client.Do(ctx, req, profile); err != nil {
		return nil, err
	}
	return &profile.Data, nil
}
func (p *profilesServiceHandler) UpdateProfile(ctx context.Context, updateReq UpdateProfileRequest) error {
	path := fmt.Sprintf("%s/user", profilePath)

	req, err := p.client.NewRequest(context.TODO(), http.MethodPut, path, updateReq)
	if err != nil {
		return err
	}

	return p.client.Do(ctx, req, nil)
}
func (p *profilesServiceHandler) GetPermissionGroups(ctx context.Context) ([]PermissionGroup, error) {
	path := fmt.Sprintf("%s/permission/group", profilePath)

	req, err := p.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	groups := new(GetPermissionGroupRoot)
	if err = p.client.Do(ctx, req, groups); err != nil {
		return nil, err
	}

	return groups.Data.OwnGroups, err
}

func (p *profilesServiceHandler) DeletePermissionGroup(ctx context.Context, groupId string) error {
	path := fmt.Sprintf("%s/permission/group/%s", profilePath, groupId)

	req, err := p.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return err
	}

	return p.client.Do(ctx, req, nil)
}

func (p *profilesServiceHandler) CreatePermissionGroup(ctx context.Context, groupName string) error {
	path := fmt.Sprintf("%s/permission/group", profilePath)
	createReq := struct {
		GroupName string `json:"groupName"`
	}{
		GroupName: groupName,
	}

	req, err := p.client.NewRequest(ctx, http.MethodPost, path, groupName)
	if err != nil {
		return err
	}

	return p.client.Do(ctx, req, createReq)
}

func (p *profilesServiceHandler) ChangePassword(ctx context.Context, oldPassword string, newPassword string) error {
	changePassReq := struct {
		OldPassword string `json:"oldPassword"`
		NewPassword string `json:"newPassword"`
	}{
		OldPassword: oldPassword,
		NewPassword: newPassword,
	}

	path := fmt.Sprintf("%s/changepass", profilePath)
	req, err := p.client.NewRequest(ctx, http.MethodPost, path, changePassReq)
	if err != nil {
		return err
	}

	return p.client.Do(ctx, req, nil)

}

func (p *profilesServiceHandler) UpdateBilling(ctx context.Context, billing BillingAddress) error {
	path := fmt.Sprintf("%s/billing", profilePath)

	req, err := p.client.NewRequest(ctx, http.MethodPut, path, billing)
	if err != nil {
		return err
	}

	return p.client.Do(ctx, req, nil)

}

func (p *profilesServiceHandler) ValidatePhone(ctx context.Context, phone string) error {
	path := fmt.Sprintf("%s/phone/validate", profilePath)
	validateReq := struct {
		PhoneNumber string `json:"phoneNumber"`
	}{
		PhoneNumber: phone,
	}

	req, err := p.client.NewRequest(ctx, http.MethodPost, path, validateReq)
	if err != nil {
		return err
	}

	return p.client.Do(ctx, req, nil)
}

func (p *profilesServiceHandler) VerifyPhone(ctx context.Context, code string) error {
	path := fmt.Sprintf("%s/phone/verify", profilePath)
	verifyReq := struct {
		Code string `json:"code"`
	}{
		Code: code,
	}

	req, err := p.client.NewRequest(ctx, http.MethodPost, path, verifyReq)
	if err != nil {
		return err
	}

	return p.client.Do(ctx, req, nil)
}

func (p *profilesServiceHandler) EnableTwofa(ctx context.Context) error {
	path := fmt.Sprintf("%s/twoFa/enable", profilePath)

	req, err := p.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return err
	}

	return p.client.Do(ctx, req, nil)
}

func (p *profilesServiceHandler) DisableTwofa(ctx context.Context) error {
	path := fmt.Sprintf("%s/twoFa/disable", profilePath)

	req, err := p.client.NewRequest(ctx, http.MethodPost, path, nil)
	if err != nil {
		return err
	}

	return p.client.Do(ctx, req, nil)
}
