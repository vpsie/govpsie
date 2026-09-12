package govpsie

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

var lbPath = "/api/v1/lb"

type LBsService interface {
	ListLBs(ctx context.Context, options *ListOptions) ([]LB, error)
	ListLBDataCenters(ctx context.Context, options *ListOptions) ([]LBDataCenter, error)
	ListOffers(ctx context.Context, dcIdentifier string) ([]LBOffers, error)
	GetLB(ctx context.Context, lbID string) (*LBDetails, error)
	CreateLB(ctx context.Context, createLBReq *CreateLBReq) error
	DeleteLB(ctx context.Context, lbID, reason, note string) error
	AddLBRule(ctx context.Context, addRuleReq *AddRuleReq) error
	DeleteLBRule(ctx context.Context, ruleID string) error
	AddLBDomain(ctx context.Context, domainAddReq *DomainAddReq) error
	ReplaceDomain(ctx context.Context, domainId, newDomainId string) error
	UpdateDomainBackend(ctx context.Context, domainId string, backends []Backend) error
	UpdateLBDomain(ctx context.Context, domainUpdateReq *DomainUpdateReq) error
	UpdateLBRules(ctx context.Context, ruleUpdateReq *RuleUpdateReq) error
	DeleteLBDomain(ctx context.Context, domainID string) error
	DeleteLBBackend(ctx context.Context, lbBackendID string) error
	ListPendingLBs(ctx context.Context) ([]PendingLB, error)
}

type lbsServiceHandler struct {
	client *Client
}

var _ LBsService = &lbsServiceHandler{}

type ListLBsRoot struct {
	Error bool `json:"error"`
	Data  []LB `json:"data"`
	Total int  `json:"total"`
}

type GetLBRoot struct {
	Error bool      `json:"error"`
	Data  LBDetails `json:"data"`
}

type ListLBDataCentersRoot struct {
	Error bool           `json:"error"`
	Data  []LBDataCenter `json:"data"`
	Total int            `json:"total"`
}

type LBDetails struct {
	LBName     string         `json:"lbName"`
	Identifier string         `json:"identifier"`
	Traffic    int            `json:"traffic"`
	BoxsizeID  int            `json:"boxsize_id"`
	DefaultIP  string         `json:"default_ip"`
	DcName     string         `json:"dc_name"`
	DcID       string         `json:"dcId"`
	CreatedBy  string         `json:"created_by"`
	UserID     int            `json:"user_id"`
	Rules      []LBRuleDetail `json:"rules"`
}

type LBRuleDetail struct {
	Scheme    string             `json:"scheme"`
	FrontPort int                `json:"frontPort"`
	BackPort  int                `json:"backPort"`
	CreatedOn time.Time          `json:"created_on"`
	RuleID    string             `json:"ruleId"`
	Domains   []LBDomainsDetail  `json:"domains,omitempty"`
	Backends  []LBBackendsDetail `json:"backends,omitempty"`
}

type LBBackendsDetail struct {
	IP           string    `json:"ip"`
	Identifier   string    `json:"identifier"`
	VMIdentifier string    `json:"vmIdentifier,omitempty"`
	CreatedOn    time.Time `json:"created_on"`
}

type LBDomainsDetail struct {
	DomainName      string             `json:"domainName"`
	BackendScheme   string             `json:"backendScheme"`
	Subdomain       *string            `json:"subdomain,omitempty"`
	Algorithm       string             `json:"algorithm"`
	RedirectHTTP    int                `json:"redirectHTTP"`
	HealthCheckPath string             `json:"healthCheckPath"`
	CookieCheck     int                `json:"cookieCheck"`
	CookieName      string             `json:"cookieName"`
	CreatedOn       time.Time          `json:"created_on"`
	BackPort        int                `json:"backPort"`
	DomainID        string             `json:"domainId"`
	CheckInterval   int                `json:"checkInterval"`
	FastInterval    int                `json:"fastInterval"`
	Rise            int                `json:"rise"`
	Fall            int                `json:"fall"`
	Backends        []LBBackendsDetail `json:"backends"`
}

type LB struct {
	Cpu        int    `json:"cpu"`
	Ssd        int    `json:"ssd"`
	Ram        int    `json:"ram"`
	LBName     string `json:"lbName"`
	Traffic    int    `json:"traffic"`
	BoxsizeID  int    `json:"boxsize_id"`
	DefaultIP  string `json:"default_ip"`
	DCName     string `json:"dc_name"`
	Identifier string `json:"identifier"`
	CreatedOn  string `json:"created_on"`
	UpdatedAt  string `json:"updated_at"`
	Package    string `json:"package"`
	CreatedBy  string `json:"created_by"`
	UserID     int    `json:"user_id"`
}

type CreateLBReq struct {
	Algorithm          string `json:"algorithm"`
	CookieName         string `json:"cookieName"`
	HealthCheckPath    string `json:"healthCheckPath,omitempty"`
	CookieCheck        bool   `json:"cookieCheck"`
	RedirectHTTP       int    `json:"redirectHTTP"`
	LBName             string `json:"lbName"`
	ResourceIdentifier string `json:"resourceIdentifier"`
	DcIdentifier       string `json:"dcIdentifier"`
	Rule               []Rule `json:"rules"`
	CheckInterval      int    `json:"checkInterval,omitempty"`
	FastInterval       int    `json:"fastInterval,omitempty"`
	Rise               int    `json:"rise,omitempty"`
	Fall               int    `json:"fall,omitempty"`
	VpcID              int    `json:"vpcId,omitempty"`
}

type AddRuleReq struct {
	Scheme    string     `json:"scheme"`
	FrontPort string     `json:"frontPort"`
	BackPort  string     `json:"backPort"`
	LbId      string     `json:"lbId"`
	Domains   []LBDomain `json:"domains"`
}
type Rule struct {
	Scheme     string     `json:"scheme"`
	FrontPort  string     `json:"frontPort"`
	Domains    []LBDomain `json:"domains"`
	Backends   []Backend  `json:"backends"`
	BackPort   string     `json:"backPort"`
	DomainName string     `json:"domainName"`
}
type LBDomain struct {
	DomainID      string    `json:"domainId"`
	Backends      []Backend `json:"backends"`
	BackPort      string    `json:"backPort"`
	BackendScheme string    `json:"backendScheme"`
	DomainName    string    `json:"domainName"`
}
type Backend struct {
	Ip           string `json:"ip"`
	VmIdentifier string `json:"vmIdentifier"`
	Type         string `json:"type,omitempty"`
}

type LBDataCenter struct {
	DcName     string `json:"dc_name"`
	DcImage    string `json:"dc_image"`
	State      string `json:"state"`
	Country    string `json:"country"`
	Identifier string `json:"identifier"`
	IsActive   int    `json:"is_active"`
	IsDeleted  int    `json:"is_deleted"`
}

type RuleUpdateReq struct {
	RuleID    string    `json:"ruleId"`
	Backends  []Backend `json:"backends"`
	BackPort  int       `json:"backPort"`
	Scheme    string    `json:"scheme"`
	FrontPort int       `json:"frontPort"`
}

type LBOffers struct {
	Cpu         int    `json:"cpu"`
	Ram         int    `json:"ram"`
	Ssd         int    `json:"ssd"`
	Traffic     int    `json:"traffic"`
	Price       string `json:"price"`
	NickName    string `json:"nickname"`
	Identifier  string `json:"identifier"`
	Color       string `json:"color"`
	NetSpeed    int    `json:"net_speed"`
	Category    string `json:"category"`
	Description string `json:"description"`
}

type DomainAddReq struct {
	RuleID       string    `json:"ruleId"`
	DomainName   string    `json:"domainName"`
	DomainID     string    `json:"domainId"`
	Algorithm    string    `json:"algorithm"`
	RedirectHTTP int       `json:"redirectHTTP"`
	CookieCheck  bool      `json:"cookieCheck"`
	CookieName   string    `json:"cookieName"`
	BackPort     int       `json:"backPort"`
	Backends     []Backend `json:"backends"`
}

type DomainUpdateReq struct {
	DomainID      string `json:"domainId"`
	Subdomain     string `json:"subdomain"`
	Algorithm     string `json:"algorithm"`
	RedirectHTTP  int    `json:"redirectHTTP"`
	CookieCheck   bool   `json:"cookieCheck"`
	CookieName    string `json:"cookieName"`
	BackPort      int    `json:"backPort"`
	CheckInterval int    `json:"checkInterval"`
	FastInterval  int    `json:"fastInterval"`
	Rise          int    `json:"rise"`
	Fall          int    `json:"fall"`
}

type ListOffersRoot struct {
	Error bool       `json:"error"`
	Data  []LBOffers `json:"Data"`
}

type PendingLB struct {
	ID   string `json:"id"`
	User struct {
	} `json:"user"`
	UserID int `json:"user_id"`
	Data   struct {
		Algorithm          string        `json:"algorithm"`
		LbName             string        `json:"lbName"`
		Rules              []interface{} `json:"rules"`
		DcIdentifier       string        `json:"dcIdentifier"`
		ResourceIdentifier string        `json:"resourceIdentifier"`
		CookieName         string        `json:"cookieName"`
		RedirectHTTP       int           `json:"redirectHTTP"`
		CookieCheck        bool          `json:"cookieCheck"`
		RequestIP          string        `json:"requestIp"`
		PrivateIps         []interface{} `json:"privateIps"`
	} `json:"data"`
	ResourceData struct {
	} `json:"resourceData"`
	Datacenter []interface{} `json:"datacenter"`
	OsData     struct {
	} `json:"osData"`
	Running int    `json:"running"`
	Type    string `json:"type"`
}

type PendingLBRoot struct {
	Error bool          `json:"error"`
	Data  [][]PendingLB `json:"data"`
}

func (l *lbsServiceHandler) ListLBs(ctx context.Context, options *ListOptions) ([]LB, error) {
	path := fmt.Sprintf("%s/all?sortField=created_on&sortDirection=DESC", lbPath)

	req, err := l.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	listLbsRoot := new(ListLBsRoot)
	if err := l.client.Do(ctx, req, listLbsRoot); err != nil {
		return nil, err
	}

	return listLbsRoot.Data, nil
}

func (l *lbsServiceHandler) GetLB(ctx context.Context, lbID string) (*LBDetails, error) {
	path := fmt.Sprintf("%s/%s", lbPath, lbID)

	req, err := l.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	lb := new(GetLBRoot)
	if err := l.client.Do(ctx, req, &lb); err != nil {
		return nil, err
	}

	return &lb.Data, nil
}

func (l *lbsServiceHandler) ListLBDataCenters(ctx context.Context, options *ListOptions) ([]LBDataCenter, error) {
	path := fmt.Sprintf("%s/datacenter", lbPath)

	req, err := l.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	lbDataCenters := new(ListLBDataCentersRoot)
	if err := l.client.Do(ctx, req, &lbDataCenters); err != nil {
		return nil, err
	}

	return lbDataCenters.Data, nil
}

func (l *lbsServiceHandler) CreateLB(ctx context.Context, createLBReq *CreateLBReq) error {
	path := fmt.Sprintf("%s/create", lbPath)

	log.Println("createLBReq", createLBReq)
	req, err := l.client.NewRequest(ctx, http.MethodPost, path, createLBReq)
	if err != nil {
		return err
	}

	return l.client.Do(ctx, req, nil)
}

func (l *lbsServiceHandler) DeleteLB(ctx context.Context, lbID, reason, note string) error {
	path := fmt.Sprintf("%s/%s", lbPath, lbID)

	deleteReq := struct {
		DeleteStatistic struct {
			Reason string `json:"reason"`
			Note   string `json:"note"`
		} `json:"deleteStatistic"`
	}{
		DeleteStatistic: struct {
			Reason string `json:"reason"`
			Note   string `json:"note"`
		}{
			Reason: reason,
			Note:   note,
		},
	}

	req, err := l.client.NewRequest(ctx, http.MethodDelete, path, &deleteReq)
	if err != nil {
		return err
	}

	return l.client.Do(ctx, req, nil)
}

func (l *lbsServiceHandler) AddLBRule(ctx context.Context, addRuleReq *AddRuleReq) error {
	path := fmt.Sprintf("%s/rule/add", lbPath)

	req, err := l.client.NewRequest(ctx, http.MethodPost, path, addRuleReq)
	if err != nil {
		return err
	}

	return l.client.Do(ctx, req, nil)
}

func (l *lbsServiceHandler) DeleteLBRule(ctx context.Context, ruleID string) error {
	path := fmt.Sprintf("%s/delete/rule", lbPath)
	delReq := struct {
		RuleID string `json:"ruleId"`
	}{
		RuleID: ruleID,
	}

	req, err := l.client.NewRequest(ctx, http.MethodDelete, path, delReq)
	if err != nil {
		return err
	}

	return l.client.Do(ctx, req, nil)
}

func (l *lbsServiceHandler) AddLBDomain(ctx context.Context, domainAddReq *DomainAddReq) error {
	path := fmt.Sprintf("%s/domain/add", lbPath)

	req, err := l.client.NewRequest(ctx, http.MethodPost, path, domainAddReq)
	if err != nil {
		return err
	}

	return l.client.Do(ctx, req, nil)

}

func (l *lbsServiceHandler) ReplaceDomain(ctx context.Context, domainId, newDomainId string) error {
	path := fmt.Sprintf("%s/domain/replace", lbPath)

	domainReplaceReq := struct {
		DomainID    string `json:"domainId"`
		NewDomainID string `json:"newDomainId"`
	}{
		DomainID:    domainId,
		NewDomainID: newDomainId,
	}

	req, err := l.client.NewRequest(ctx, http.MethodPost, path, domainReplaceReq)
	if err != nil {
		return err
	}

	return l.client.Do(ctx, req, nil)
}

func (l *lbsServiceHandler) UpdateDomainBackend(ctx context.Context, domainId string, backends []Backend) error {
	path := fmt.Sprintf("%s/backend/update", lbPath)

	updateDomainBackendReq := struct {
		DomainID string    `json:"domainId"`
		Backends []Backend `json:"backends"`
	}{
		DomainID: domainId,
		Backends: backends,
	}

	req, err := l.client.NewRequest(ctx, http.MethodPost, path, updateDomainBackendReq)
	if err != nil {
		return err
	}

	return l.client.Do(ctx, req, nil)
}

func (l *lbsServiceHandler) UpdateLBRules(ctx context.Context, ruleUpdateReq *RuleUpdateReq) error {
	path := fmt.Sprintf("%s/rule/update", lbPath)

	req, err := l.client.NewRequest(ctx, http.MethodPost, path, ruleUpdateReq)
	if err != nil {
		return err
	}

	return l.client.Do(ctx, req, nil)
}

func (l *lbsServiceHandler) UpdateLBDomain(ctx context.Context, domainUpdateReq *DomainUpdateReq) error {
	path := fmt.Sprintf("%s/domain/update", lbPath)

	req, err := l.client.NewRequest(ctx, http.MethodPost, path, domainUpdateReq)
	if err != nil {
		return err
	}

	return l.client.Do(ctx, req, nil)
}

func (l *lbsServiceHandler) DeleteLBDomain(ctx context.Context, domainID string) error {
	path := fmt.Sprintf("%s/delete/domain", lbPath)

	delReq := struct {
		DomainID string `json:"domainId"`
	}{
		DomainID: domainID,
	}

	req, err := l.client.NewRequest(ctx, http.MethodDelete, path, delReq)
	if err != nil {
		return err
	}

	return l.client.Do(ctx, req, nil)
}

func (l *lbsServiceHandler) DeleteLBBackend(ctx context.Context, lbBackendID string) error {
	path := fmt.Sprintf("%s/delete/backend", lbPath)

	delReq := struct {
		BackendID string `json:"backendId"`
	}{
		BackendID: lbBackendID,
	}

	req, err := l.client.NewRequest(ctx, http.MethodDelete, path, delReq)
	if err != nil {
		return err
	}

	return l.client.Do(ctx, req, nil)
}

func (l *lbsServiceHandler) ListOffers(ctx context.Context, dcIdentifier string) ([]LBOffers, error) {
	path := fmt.Sprintf("%s/offers", lbPath)

	offerReq := struct {
		DcIdentifier string `json:"dcIdentifier"`
	}{
		DcIdentifier: dcIdentifier,
	}

	req, err := l.client.NewRequest(ctx, http.MethodPost, path, &offerReq)
	if err != nil {
		return nil, err
	}

	offersRoot := new(ListOffersRoot)
	if err := l.client.Do(ctx, req, &offersRoot); err != nil {
		return nil, err
	}

	return offersRoot.Data, nil
}

func (l *lbsServiceHandler) ListPendingLBs(ctx context.Context) ([]PendingLB, error) {
	path := "/api/v2/lbs/pending"

	req, err := l.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	pendingLbs := new(PendingLBRoot)
	if err := l.client.Do(ctx, req, pendingLbs); err != nil {
		return nil, err
	}

	if len(pendingLbs.Data) == 0 {
		return []PendingLB{}, nil
	}

	return pendingLbs.Data[0], nil
}
