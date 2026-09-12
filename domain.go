package govpsie

import (
	"context"
	"fmt"
	"net/http"
)

var domainsPath = "/apps/v2/domains"
var domainPath = "/apps/v2/domain"

type DomainService interface {
	ListDomainByProject(ctx context.Context, options *ListOptions, projectIdentifier string) ([]Domain, error)
	DnsRecord(ctx context.Context, domainIdentifier string, dnsRecord *DnsRecord) error
	ListDomains(ctx context.Context, options *ListOptions) ([]Domain, error)
	ListAllDomains(ctx context.Context) ([]Domain, error)
	ListDomainVpsies(ctx context.Context, options *ListOptions) ([]DomainVpsie, error)
	CreateDomain(ctx context.Context, createReq *CreateDomainRequest) error
	GetDomainByVpsie(ctx context.Context, domainIdentifier string) ([]Domain, error)
	UpdateReverse(ctx context.Context, reverseReq *ReverseRequest) error
	AddReverse(ctx context.Context, reverseReq *ReverseRequest) error
	UpdateDomain(ctx context.Context, dnsRecord *DnsRecord, domainIdentifier, vmIdentifier string) error
	DeleteReverse(ctx context.Context, ip, vmIdentifier string) error
	CreateDnsRecord(ctx context.Context, createReq CreateDnsRecordReq) error
	UpdateDnsRecord(ctx context.Context, updateReq *UpdateDnsRecordReq) error
	DeleteDomain(ctx context.Context, domainIdentifier, reason, note string) error
	DeleteDnsRecord(ctx context.Context, domainIdentifier string, record *Record) error
	ListReversePTRRecords(ctx context.Context) ([]ReversePTR, error)
}

type domainsServiceHandler struct {
	client *Client
}

var _ DomainService = &domainsServiceHandler{}

type ListDomainRoot struct {
	Error bool     `json:"error"`
	Data  []Domain `json:"data"`
	Total int      `json:"total "`
}

type Domain struct {
	DomainName  string `json:"domain_name"`
	Identifier  string `json:"identifier"`
	NsValidated int    `json:"ns_validated"`
	CreatedOn   string `json:"created_on"`
	LastCheck   string `json:"last_check"`
}

type DnsRecord struct {
	Name     string `json:"name"`
	Service  string `json:"service"`
	Protocol string `json:"protocol"`
	Content  string `json:"content"`
	Priority string `json:"priority"`
	Weight   string `json:"weight"`
	Port     string `json:"port"`
	Type     string `json:"type"`
	Ttl      int    `json:"ttl"`
}

type ReverseRequest struct {
	VmIdentifier     string `json:"vmIdentifier"`
	Ip               string `json:"ip"`
	DomainIdentifier string `json:"domainIdentifier"`
	HostName         string `json:"hostName"`
}

type DomainVpsie struct {
	HostName     string `json:"hostname"`
	IP           string `json:"ip"`
	Identifier   string `json:"identifier"`
	IPVersion    string `json:"ip_version"`
	MaskCIDR     int    `json:"mask_cidr"`
	DCIdentifier string `json:"dc_identifier"`
	Category     string `json:"category"`
	FullName     string `json:"fullname"`
	PCS          int    `json:"pcs"`
}

type ListDomainVpsieRoot struct {
	Error bool          `json:"error"`
	Vpsie []DomainVpsie `json:"data"`
	Total int           `json:"total"`
}

type CreateDomainRequest struct {
	ProjectIdentifier string   `json:"projectIdentifier"`
	Tags              []string `json:"tags,omitempty"`
	Domain            string   `json:"domain"`
}

type Record struct {
	Name    string `json:"name"`
	Content string `json:"content"`
	Type    string `json:"type"`
	TTL     int    `json:"ttl"`
}
type CreateDnsRecordReq struct {
	DomainIdentifier string `json:"domainIdentifier"`
	Record           Record `json:"record"`
}

type UpdateDnsRecordReq struct {
	DomainIdentifier string `json:"domainIdentifier"`
	Current          Record `json:"current"`
	New              Record `json:"new"`
}

type ReversePTR struct {
	Ip           string `json:"ip"`
	HostName     string `json:"host_name"`
	VmIdentifier string `json:"vmIdentifier"`
	PtrRecord    string `json:"ptrRecord"`
}
type ListReversePTRRoot struct {
	Error bool         `json:"error"`
	Data  []ReversePTR `json:"data"`
}

func (d *domainsServiceHandler) ListDomainByProject(ctx context.Context, options *ListOptions, projectIdentifier string) ([]Domain, error) {
	path := fmt.Sprintf("%s/project/%s?offset=%d&limit=%d", domainsPath, projectIdentifier, options.Page, options.PerPage)

	req, err := d.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	domains := new(ListDomainRoot)
	if err = d.client.Do(ctx, req, domains); err != nil {
		return nil, err
	}

	return domains.Data, nil
}

func (d *domainsServiceHandler) CreateDnsRecord(ctx context.Context, createReq CreateDnsRecordReq) error {
	path := fmt.Sprintf("%s/dnsRecord", domainPath)

	req, err := d.client.NewRequest(ctx, http.MethodPost, path, createReq)
	if err != nil {
		return err
	}

	return d.client.Do(ctx, req, nil)
}

func (d *domainsServiceHandler) UpdateDnsRecord(ctx context.Context, updateReq *UpdateDnsRecordReq) error {
	path := fmt.Sprintf("%s/dnsRecord/update", domainPath)

	req, err := d.client.NewRequest(ctx, http.MethodPut, path, updateReq)
	if err != nil {
		return err
	}

	return d.client.Do(ctx, req, nil)
}

func (d *domainsServiceHandler) DeleteDnsRecord(ctx context.Context, domainIdentifier string, record *Record) error {
	path := fmt.Sprintf("%s/dnsRecord/delete", domainPath)

	updateReq := struct {
		DomainIdentifier string `json:"domainIdentifier"`
		Record           Record `json:"record"`
	}{
		DomainIdentifier: domainIdentifier,
		Record:           *record,
	}

	req, err := d.client.NewRequest(ctx, http.MethodDelete, path, updateReq)
	if err != nil {
		return nil
	}

	return d.client.Do(ctx, req, nil)
}

func (d *domainsServiceHandler) DnsRecord(ctx context.Context, domainIdentifier string, dnsRecord *DnsRecord) error {
	path := fmt.Sprintf("%s/dnsRecord", domainPath)

	dnsRecordReq := struct {
		DomainIdentifier string     `json:"domainIdentifier"`
		Record           *DnsRecord `json:"record"`
	}{
		DomainIdentifier: domainIdentifier,
		Record:           dnsRecord,
	}

	req, err := d.client.NewRequest(ctx, http.MethodPost, path, dnsRecordReq)
	if err != nil {
		return err
	}

	return d.client.Do(ctx, req, nil)
}

func (d *domainsServiceHandler) ListDomains(ctx context.Context, options *ListOptions) ([]Domain, error) {
	path := fmt.Sprintf("%s?offset=%d&limit=%d", domainsPath, options.Page, options.PerPage)

	req, err := d.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	domains := new(ListDomainRoot)
	if err = d.client.Do(ctx, req, &domains); err != nil {
		return nil, err
	}
	return domains.Data, nil
}

func (d *domainsServiceHandler) ListAllDomains(ctx context.Context) ([]Domain, error) {
	path := fmt.Sprintf("%s", domainsPath)

	req, err := d.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	domains := new(ListDomainRoot)
	if err = d.client.Do(ctx, req, &domains); err != nil {
		return nil, err
	}
	return domains.Data, nil
}

func (d *domainsServiceHandler) ListDomainVpsies(ctx context.Context, options *ListOptions) ([]DomainVpsie, error) {
	path := fmt.Sprintf("%s/vms?offset=%d&limit=%d", domainsPath, options.Page, options.PerPage)

	req, err := d.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	vpsies := new(ListDomainVpsieRoot)
	if err = d.client.Do(ctx, req, &vpsies); err != nil {
		return nil, err
	}

	return vpsies.Vpsie, nil
}

func (d *domainsServiceHandler) CreateDomain(ctx context.Context, createReq *CreateDomainRequest) error {
	path := fmt.Sprintf("%s/add", domainPath)

	req, err := d.client.NewRequest(ctx, http.MethodPost, path, createReq)
	if err != nil {
		return err
	}

	return d.client.Do(ctx, req, nil)
}

func (d *domainsServiceHandler) GetDomainByVpsie(ctx context.Context, domainIdentifier string) ([]Domain, error) {
	path := fmt.Sprintf("%s/vm/%s", domainsPath, domainIdentifier)

	req, err := d.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	domains := new(ListDomainRoot)
	if err = d.client.Do(ctx, req, &domains); err != nil {
		return nil, err
	}

	return domains.Data, nil
}

func (d *domainsServiceHandler) DeleteDomain(ctx context.Context, domainIdentifier, reason, note string) error {
	path := fmt.Sprintf("%s/delete", domainPath)

	deleteReq := struct {
		DomainIdentifier string `json:"domainIdentifier"`
		DeleteStatistic  struct {
			Reason string `json:"reason"`
			Note   string `json:"note"`
		}
	}{
		DomainIdentifier: domainIdentifier,
		DeleteStatistic: struct {
			Reason string `json:"reason"`
			Note   string `json:"note"`
		}{
			Reason: reason,
			Note:   note,
		},
	}

	req, err := d.client.NewRequest(ctx, http.MethodDelete, path, deleteReq)
	if err != nil {
		return err
	}
	return d.client.Do(ctx, req, nil)
}

func (d *domainsServiceHandler) AddReverse(ctx context.Context, reverseReq *ReverseRequest) error {
	path := fmt.Sprintf("%s/addreverse", domainPath)

	req, err := d.client.NewRequest(ctx, http.MethodPost, path, reverseReq)
	if err != nil {
		return err
	}

	return d.client.Do(ctx, req, nil)
}

func (d *domainsServiceHandler) UpdateReverse(ctx context.Context, reverseReq *ReverseRequest) error {
	path := fmt.Sprintf("%s/reverse/update", domainPath)

	req, err := d.client.NewRequest(ctx, http.MethodPut, path, reverseReq)
	if err != nil {
		return err
	}

	return d.client.Do(ctx, req, nil)
}

func (d *domainsServiceHandler) UpdateDomain(ctx context.Context, dnsRecord *DnsRecord, domainIdentifier, vmIdentifier string) error {
	path := fmt.Sprintf("%s/update", domainPath)

	dnsRecordReq := struct {
		DomainIdentifier string     `json:"domainIdentifier"`
		VmIdentifier     string     `json:"vmIdentifier"`
		Record           *DnsRecord `json:"record"`
	}{
		DomainIdentifier: domainIdentifier,
		VmIdentifier:     vmIdentifier,
		Record:           dnsRecord,
	}

	req, err := d.client.NewRequest(ctx, http.MethodPut, path, dnsRecordReq)
	if err != nil {
		return err
	}

	return d.client.Do(ctx, req, nil)
}

func (d *domainsServiceHandler) DeleteReverse(ctx context.Context, ip, vmIdentifier string) error {
	path := fmt.Sprintf("%s/reverse/delete", domainPath)

	rvDeleteReq := struct {
		VmIdentifier string `json:"vmIdentifier"`
		Ip           string `json:"ip"`
	}{
		VmIdentifier: vmIdentifier,
		Ip:           ip,
	}

	req, err := d.client.NewRequest(ctx, http.MethodDelete, path, rvDeleteReq)
	if err != nil {
		return err
	}

	return d.client.Do(ctx, req, nil)
}

func (d *domainsServiceHandler) ListReversePTRRecords(ctx context.Context) ([]ReversePTR, error) {
	path := fmt.Sprintf("%s/reverse/all", domainPath)

	req, err := d.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	rv := new(ListReversePTRRoot)
	if err = d.client.Do(ctx, req, rv); err != nil {
		return nil, err
	}

	return rv.Data, nil
}
