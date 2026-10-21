package govpsie

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

var bareMetalBasePath = "/apps/v2/bare-metal"

// BareMetalService manages dedicated (bare-metal) servers: the hardware plans
// and operating systems a datacenter offers, ordering a server, reading it,
// changing its hostname, project, password and private networks, power
// actions, and releasing it.
type BareMetalService interface {
	// ListPlans returns the active hardware plans of a datacenter. A
	// datacenter that sells no bare metal answers an empty list.
	ListPlans(ctx context.Context, dcIdentifier string, filter *BareMetalPlanFilter) ([]BareMetalPlan, error)
	// ListOS returns the operating systems a plan can be installed with,
	// grouped by distribution. planIdentifier may be empty.
	ListOS(ctx context.Context, dcIdentifier, planIdentifier string) ([]BareMetalOSGroup, error)
	// Deploy orders a server. A non-empty idempotencyKey (1-128 visible ASCII
	// characters) makes the order safe to retry: the API answers a repeat of
	// the same key with the first order's result for 24 hours.
	Deploy(ctx context.Context, req *BareMetalDeployRequest, idempotencyKey string) (*BareMetalDeployment, error)
	// List returns the account's servers.
	List(ctx context.Context) ([]BareMetalSummary, error)
	// Get returns one server. A server the account does not hold -- never
	// ordered, released, or a build-to-order server not delivered yet --
	// answers an *APIError with StatusCode 404.
	Get(ctx context.Context, identifier string) (*BareMetal, error)
	// UpdateHostname renames the server.
	UpdateHostname(ctx context.Context, identifier, hostname string) error
	// MoveToProject moves the server to another project of the account.
	MoveToProject(ctx context.Context, identifier, projectIdentifier string) error
	// ResetPassword sets a new administrator password. The server powers off
	// and boots a reset job; Get reports OSPhase "password_reset" until it
	// is back to "running" (or "failed").
	ResetPassword(ctx context.Context, identifier string, req *BareMetalPasswordResetRequest) (*BareMetalPasswordReset, error)
	// Power runs a power action: on, off, shutdown or reboot.
	Power(ctx context.Context, identifier string, action BareMetalPowerAction) error
	// PowerState reads the server's current power state from its BMC.
	PowerState(ctx context.Context, identifier string) (string, error)
	// ListPrivateNetworks returns the private networks attached to the server.
	ListPrivateNetworks(ctx context.Context, identifier string) (*BareMetalPrivateNetworks, error)
	// AttachPrivateNetwork attaches a VPC (by its numeric id). The attach
	// runs in the background: the returned item is "attaching" until the
	// switch port carries the network.
	AttachPrivateNetwork(ctx context.Context, identifier string, networkID int64) (*BareMetalPrivateNetwork, error)
	// DetachPrivateNetwork detaches a VPC. The detach runs in the background.
	DetachPrivateNetwork(ctx context.Context, identifier string, networkID int64) error
	// Release gives the server back. It is confirmed with the account
	// password and wipes the disks; Get answers 404 afterwards.
	Release(ctx context.Context, identifier string, req *BareMetalReleaseRequest) error
}

type bareMetalServiceHandler struct {
	client *Client
}

var _ BareMetalService = &bareMetalServiceHandler{}

// BareMetalPowerAction is a power action accepted by BareMetalService.Power.
type BareMetalPowerAction string

// Power actions.
const (
	BareMetalPowerOn       BareMetalPowerAction = "on"
	BareMetalPowerOff      BareMetalPowerAction = "off"
	BareMetalPowerShutdown BareMetalPowerAction = "shutdown"
	BareMetalPowerReboot   BareMetalPowerAction = "reboot"
)

// OS phases reported in BareMetal.OSPhase.
const (
	BareMetalOSPhaseInstalling    = "installing"
	BareMetalOSPhaseReinstalling  = "reinstalling"
	BareMetalOSPhasePasswordReset = "password_reset"
	BareMetalOSPhaseRunning       = "running"
	BareMetalOSPhaseFailed        = "failed"
)

// Order statuses reported in BareMetalDeployment.Status.
const (
	// BareMetalOrderSold means a stocked server was assigned and is
	// installing.
	BareMetalOrderSold = "sold"
	// BareMetalOrderPendingBuild means the plan is built to order: the
	// server is assembled first (see DeliveryHours) and only then appears in
	// Get and List.
	BareMetalOrderPendingBuild = "pending_build"
)

// Private network attachment statuses.
const (
	BareMetalNetworkAttaching = "attaching"
	BareMetalNetworkAttached  = "attached"
	BareMetalNetworkDetaching = "detaching"
	BareMetalNetworkFailed    = "failed"
)

// BareMetalPlanFilter narrows ListPlans. Zero values do not filter.
type BareMetalPlanFilter struct {
	MinRAMGB  int
	MinDiskGB int
	// GPU, when set, keeps only plans with (true) or without (false) GPUs.
	GPU *bool
}

// BareMetalContractTerm is a billing term a plan can be ordered on.
type BareMetalContractTerm struct {
	Term        string  `json:"term"`
	DiscountPct float64 `json:"discountPct"`
}

// BareMetalPlan is a hardware plan.
type BareMetalPlan struct {
	Identifier       string                  `json:"identifier"`
	Name             string                  `json:"name"`
	ProcessorModel   string                  `json:"processorModel"`
	CPUCount         int                     `json:"cpuCount"`
	RAMGB            int                     `json:"ramGb"`
	DiskGB           int                     `json:"diskGb"`
	DiskTypeSummary  string                  `json:"diskTypeSummary"`
	GPUModel         string                  `json:"gpuModel"`
	GPUCount         int                     `json:"gpuCount"`
	GPUMemoryGB      int                     `json:"gpuMemoryGb"`
	Price            FlexString              `json:"price"`
	Currency         string                  `json:"currency"`
	SetupFee         FlexString              `json:"setupFee"`
	Delivery         string                  `json:"delivery"`
	DeliveryHours    int                     `json:"deliveryHours"`
	ContractTerms    []BareMetalContractTerm `json:"contractTerms"`
	LinkSpeedOptions []int                   `json:"linkSpeedOptions"`
	AvailableServers int                     `json:"availableServers"`
	IsActive         bool                    `json:"isActive"`
}

// BareMetalOSVersion is one installable operating system.
type BareMetalOSVersion struct {
	Identifier               string `json:"identifier"`
	FullName                 string `json:"fullname"`
	Category                 string `json:"category"`
	Description              string `json:"description"`
	GPU                      bool   `json:"gpu"`
	OSFamily                 string `json:"osFamily"`
	RemoteAccessConfigurable bool   `json:"remoteAccessConfigurable"`
	IsDefault                bool   `json:"is_default"`
}

// BareMetalOSGroup is a distribution and its installable versions.
type BareMetalOSGroup struct {
	Name      string               `json:"name"`
	IsDefault bool                 `json:"is_default"`
	Versions  []BareMetalOSVersion `json:"available_versions"`
}

// BareMetalRemoteAccess configures the remote-access firewall rules written
// at install time: SSH (or RDP) open to the internet, WinRM on Windows, or a
// firewall group of the account.
type BareMetalRemoteAccess struct {
	Enabled            bool    `json:"enabled"`
	WinRM              bool    `json:"winrm"`
	FirewallIdentifier *string `json:"firewallIdentifier"`
}

// BareMetalDeployRequest orders a server. Exactly one of OSIdentifier and
// TemplateName is required.
type BareMetalDeployRequest struct {
	DCIdentifier      string                 `json:"dcId"`
	PlanIdentifier    string                 `json:"planId"`
	Hostname          string                 `json:"hostname"`
	Password          string                 `json:"password"`
	OSIdentifier      string                 `json:"osIdentifier,omitempty"`
	TemplateName      string                 `json:"templateName,omitempty"`
	SSHKeyIdentifiers []string               `json:"sshKeyIdentifiers,omitempty"`
	RemoteAccess      *BareMetalRemoteAccess `json:"remoteAccess,omitempty"`
	DNSZoneIdentifier string                 `json:"dnsZoneIdentifier,omitempty"`
	ProjectIdentifier string                 `json:"projectId,omitempty"`
	// ContractTerm is "hourly" or "monthly"; empty takes the plan default.
	ContractTerm string `json:"contractTerm,omitempty"`
	// LinkSpeedMbps is one of the plan's LinkSpeedOptions; 0 takes the
	// plan default.
	LinkSpeedMbps int `json:"linkSpeedMbps,omitempty"`
	// UserData is a cloud-config document (starting with "#cloud-config")
	// or a script (starting with "#!") run on first boot, at most 64 KB.
	UserData string `json:"userData,omitempty"`
	// ScriptIdentifiers are saved scripts of the account (at most 5) run on
	// first boot.
	ScriptIdentifiers []string `json:"scriptIds,omitempty"`
}

// BareMetalDeployment is the answer to an order.
type BareMetalDeployment struct {
	Success             bool       `json:"success"`
	Identifier          string     `json:"identifier"`
	OSPhase             string     `json:"osPhase"`
	OperationID         FlexString `json:"operationId"`
	IPAddress           string     `json:"ipAddress"`
	Status              string     `json:"status"`
	Delivery            string     `json:"delivery"`
	DeliveryHours       int        `json:"deliveryHours"`
	ContractTerm        string     `json:"contractTerm"`
	ContractDiscountPct float64    `json:"contractDiscountPct"`
	LinkSpeedMbps       int        `json:"linkSpeedMbps"`
	SetupFee            FlexString `json:"setupFee"`
}

// BareMetalSummary is a server as List returns it.
type BareMetalSummary struct {
	Identifier           string `json:"identifier"`
	MachineName          string `json:"machine_name"`
	Status               string `json:"status"`
	CurrentOSName        string `json:"current_os_name"`
	CurrentOSTitle       string `json:"current_os_title"`
	Location             string `json:"location"`
	CreatedOn            string `json:"created_on"`
	DatacenterIdentifier string `json:"datacenter_identifier"`
	PlanName             string `json:"plan_name"`
	IPAddress            string `json:"ip_address"`
	OSPhase              string `json:"osPhase"`
	DisplayState         string `json:"displayState"`
	ProjectIdentifier    string `json:"projectIdentifier"`
	ProjectName          string `json:"projectName"`
}

// BareMetalHardware describes the server's hardware.
type BareMetalHardware struct {
	ProcessorModel string `json:"processor_model"`
	Cores          int    `json:"cores"`
	RAMGB          int    `json:"ram_gb"`
	StorageGB      int    `json:"storage_gb"`
	GPUModel       string `json:"gpu_model"`
	GPUCount       int    `json:"gpu_count"`
	GPUMemoryGB    int    `json:"gpu_memory_gb"`
}

// BareMetalContract is the billing term the server runs on.
type BareMetalContract struct {
	Term        string  `json:"term"`
	DiscountPct float64 `json:"discountPct"`
}

// BareMetalAdminHold reports holds staff placed on the server.
type BareMetalAdminHold struct {
	Suspended    bool `json:"suspended"`
	Disconnected bool `json:"disconnected"`
	Locked       bool `json:"locked"`
}

// BareMetalPrivateNetwork is a VPC attached to a server.
type BareMetalPrivateNetwork struct {
	NetworkID    int64  `json:"networkId"`
	Name         string `json:"name"`
	NetworkRange string `json:"networkRange"`
	VLANID       int    `json:"vlanId"`
	Status       string `json:"status"`
	ErrorCode    string `json:"errorCode"`
	AttachedAt   string `json:"attachedAt"`
}

// BareMetalPrivateNetworks lists a server's private networks. Available is
// false (with Reason) where the server's site cannot attach them.
type BareMetalPrivateNetworks struct {
	Available   bool                      `json:"available"`
	Reason      string                    `json:"reason"`
	MaxNetworks int                       `json:"maxNetworks"`
	Items       []BareMetalPrivateNetwork `json:"items"`
}

// BareMetal is a server as Get returns it.
type BareMetal struct {
	Identifier           string                    `json:"identifier"`
	Name                 string                    `json:"name"`
	PowerState           string                    `json:"powerState"`
	OSFamily             string                    `json:"osFamily"`
	OSName               string                    `json:"osName"`
	OSPhase              string                    `json:"osPhase"`
	DisplayState         string                    `json:"displayState"`
	DisplayStateReason   string                    `json:"displayStateReason"`
	DeployedAt           string                    `json:"deployedAt"`
	RemoteAccess         *BareMetalRemoteAccess    `json:"remoteAccess"`
	Hardware             *BareMetalHardware        `json:"hardwareSpecs"`
	DatacenterIdentifier string                    `json:"datacenterIdentifier"`
	IPAddress            string                    `json:"ip_address"`
	Netmask              string                    `json:"netmask"`
	Gateway              string                    `json:"gateway"`
	Contract             *BareMetalContract        `json:"contract"`
	LinkSpeedMbps        int                       `json:"linkSpeedMbps"`
	ProjectIdentifier    string                    `json:"projectIdentifier"`
	ProjectName          string                    `json:"projectName"`
	AdminHold            BareMetalAdminHold        `json:"adminHold"`
	PrivateNetworks      *BareMetalPrivateNetworks `json:"privateNetworks"`
}

// BareMetalPasswordResetRequest sets a new administrator password. OSType
// ("linux" or "windows") and Username are detected from the installed OS
// when empty.
type BareMetalPasswordResetRequest struct {
	Password string `json:"password"`
	OSType   string `json:"osType,omitempty"`
	Username string `json:"username,omitempty"`
}

// BareMetalPasswordReset is the answer to a password reset.
type BareMetalPasswordReset struct {
	Success  bool   `json:"success"`
	Started  bool   `json:"started"`
	OSPhase  string `json:"osPhase"`
	OSType   string `json:"osType"`
	Username string `json:"username"`
}

// BareMetalReleaseRequest releases a server. Password is the ACCOUNT
// password (see Client.AccountPassword), not the server's.
type BareMetalReleaseRequest struct {
	Password string `json:"password"`
	Reason   string `json:"reason"`
	Note     string `json:"note,omitempty"`
}

type bareMetalPlansRoot struct {
	Error bool            `json:"error"`
	Data  []BareMetalPlan `json:"data"`
}

type bareMetalOSRoot struct {
	Error bool               `json:"error"`
	Data  []BareMetalOSGroup `json:"data"`
}

type bareMetalDeployRoot struct {
	Error bool                `json:"error"`
	Data  BareMetalDeployment `json:"data"`
}

type bareMetalListRoot struct {
	Success bool               `json:"success"`
	Data    []BareMetalSummary `json:"data"`
}

type bareMetalRoot struct {
	Error bool            `json:"error"`
	Data  json.RawMessage `json:"data"`
}

type bareMetalPasswordResetRoot struct {
	Error bool                   `json:"error"`
	Data  BareMetalPasswordReset `json:"data"`
}

type bareMetalPowerStateRoot struct {
	Error bool `json:"error"`
	Data  struct {
		PowerState string `json:"powerState"`
	} `json:"data"`
}

type bareMetalNetworksRoot struct {
	Error bool                     `json:"error"`
	Data  BareMetalPrivateNetworks `json:"data"`
}

type bareMetalNetworkRoot struct {
	Error bool                    `json:"error"`
	Data  BareMetalPrivateNetwork `json:"data"`
}

// errBareMetalIdentifier is returned when a method is called without the
// server identifier, which would otherwise address the collection route.
var errBareMetalIdentifier = errors.New("govpsie: bare-metal identifier is required")

func bareMetalPath(identifier string, parts ...string) string {
	p := bareMetalBasePath + "/" + url.PathEscape(identifier)
	for _, part := range parts {
		p += "/" + part
	}

	return p
}

func (s *bareMetalServiceHandler) ListPlans(ctx context.Context, dcIdentifier string, filter *BareMetalPlanFilter) ([]BareMetalPlan, error) {
	if dcIdentifier == "" {
		return nil, errors.New("govpsie: datacenter identifier is required")
	}

	q := url.Values{}
	if filter != nil {
		if filter.MinRAMGB > 0 {
			q.Set("minRamGb", strconv.Itoa(filter.MinRAMGB))
		}
		if filter.MinDiskGB > 0 {
			q.Set("minDiskGb", strconv.Itoa(filter.MinDiskGB))
		}
		if filter.GPU != nil {
			q.Set("gpu", strconv.FormatBool(*filter.GPU))
		}
	}

	path := fmt.Sprintf("%s/dc/%s/hardware-plans", bareMetalBasePath, url.PathEscape(dcIdentifier))
	if len(q) > 0 {
		path += "?" + q.Encode()
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	root := new(bareMetalPlansRoot)
	if err := s.client.Do(ctx, req, root); err != nil {
		return nil, err
	}

	return root.Data, nil
}

func (s *bareMetalServiceHandler) ListOS(ctx context.Context, dcIdentifier, planIdentifier string) ([]BareMetalOSGroup, error) {
	if dcIdentifier == "" {
		return nil, errors.New("govpsie: datacenter identifier is required")
	}

	path := fmt.Sprintf("%s/dc/%s/os", bareMetalBasePath, url.PathEscape(dcIdentifier))
	if planIdentifier != "" {
		path += "?" + url.Values{"planId": {planIdentifier}}.Encode()
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	root := new(bareMetalOSRoot)
	if err := s.client.Do(ctx, req, root); err != nil {
		return nil, err
	}

	return root.Data, nil
}

func (s *bareMetalServiceHandler) Deploy(ctx context.Context, deploy *BareMetalDeployRequest, idempotencyKey string) (*BareMetalDeployment, error) {
	if deploy == nil {
		return nil, errors.New("govpsie: deploy request is required")
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, bareMetalBasePath+"/deploy", deploy)
	if err != nil {
		return nil, err
	}
	if idempotencyKey != "" {
		req.Header.Set("Idempotency-Key", idempotencyKey)
	}

	root := new(bareMetalDeployRoot)
	if err := s.client.Do(ctx, req, root); err != nil {
		return nil, err
	}

	return &root.Data, nil
}

func (s *bareMetalServiceHandler) List(ctx context.Context) ([]BareMetalSummary, error) {
	req, err := s.client.NewRequest(ctx, http.MethodGet, bareMetalBasePath, nil)
	if err != nil {
		return nil, err
	}

	root := new(bareMetalListRoot)
	if err := s.client.Do(ctx, req, root); err != nil {
		return nil, err
	}

	return root.Data, nil
}

func (s *bareMetalServiceHandler) Get(ctx context.Context, identifier string) (*BareMetal, error) {
	if identifier == "" {
		return nil, errBareMetalIdentifier
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, bareMetalPath(identifier), nil)
	if err != nil {
		return nil, err
	}

	root := new(bareMetalRoot)
	if err := s.client.Do(ctx, req, root); err != nil {
		return nil, err
	}
	if isNullData(root.Data) {
		return nil, &APIError{StatusCode: http.StatusNotFound, Message: fmt.Sprintf("bare-metal server %s not found", identifier)}
	}

	server := new(BareMetal)
	if err := json.Unmarshal(root.Data, server); err != nil {
		return nil, err
	}

	return server, nil
}

func (s *bareMetalServiceHandler) put(ctx context.Context, path string, body interface{}) error {
	req, err := s.client.NewRequest(ctx, http.MethodPut, path, body)
	if err != nil {
		return err
	}

	return s.client.Do(ctx, req, nil)
}

func (s *bareMetalServiceHandler) UpdateHostname(ctx context.Context, identifier, hostname string) error {
	if identifier == "" {
		return errBareMetalIdentifier
	}

	return s.put(ctx, bareMetalPath(identifier, "hostname"), map[string]string{"hostname": hostname})
}

func (s *bareMetalServiceHandler) MoveToProject(ctx context.Context, identifier, projectIdentifier string) error {
	if identifier == "" {
		return errBareMetalIdentifier
	}

	return s.put(ctx, bareMetalPath(identifier, "project"), map[string]string{"projectId": projectIdentifier})
}

func (s *bareMetalServiceHandler) ResetPassword(ctx context.Context, identifier string, reset *BareMetalPasswordResetRequest) (*BareMetalPasswordReset, error) {
	if identifier == "" {
		return nil, errBareMetalIdentifier
	}
	if reset == nil {
		return nil, errors.New("govpsie: password reset request is required")
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, bareMetalPath(identifier, "reset-password"), reset)
	if err != nil {
		return nil, err
	}

	root := new(bareMetalPasswordResetRoot)
	if err := s.client.Do(ctx, req, root); err != nil {
		return nil, err
	}

	return &root.Data, nil
}

func (s *bareMetalServiceHandler) Power(ctx context.Context, identifier string, action BareMetalPowerAction) error {
	if identifier == "" {
		return errBareMetalIdentifier
	}

	switch action {
	case BareMetalPowerOn, BareMetalPowerOff, BareMetalPowerShutdown, BareMetalPowerReboot:
	default:
		return fmt.Errorf("govpsie: unknown bare-metal power action %q", action)
	}

	return s.put(ctx, bareMetalPath(identifier, string(action)), struct{}{})
}

func (s *bareMetalServiceHandler) PowerState(ctx context.Context, identifier string) (string, error) {
	if identifier == "" {
		return "", errBareMetalIdentifier
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, bareMetalPath(identifier, "power-state"), nil)
	if err != nil {
		return "", err
	}

	root := new(bareMetalPowerStateRoot)
	if err := s.client.Do(ctx, req, root); err != nil {
		return "", err
	}

	return root.Data.PowerState, nil
}

func (s *bareMetalServiceHandler) ListPrivateNetworks(ctx context.Context, identifier string) (*BareMetalPrivateNetworks, error) {
	if identifier == "" {
		return nil, errBareMetalIdentifier
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, bareMetalPath(identifier, "private-networks"), nil)
	if err != nil {
		return nil, err
	}

	root := new(bareMetalNetworksRoot)
	if err := s.client.Do(ctx, req, root); err != nil {
		return nil, err
	}

	return &root.Data, nil
}

func (s *bareMetalServiceHandler) AttachPrivateNetwork(ctx context.Context, identifier string, networkID int64) (*BareMetalPrivateNetwork, error) {
	if identifier == "" {
		return nil, errBareMetalIdentifier
	}
	if networkID <= 0 {
		return nil, errors.New("govpsie: private network id must be positive")
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, bareMetalPath(identifier, "private-networks"), map[string]int64{"networkId": networkID})
	if err != nil {
		return nil, err
	}

	root := new(bareMetalNetworkRoot)
	if err := s.client.Do(ctx, req, root); err != nil {
		return nil, err
	}

	return &root.Data, nil
}

func (s *bareMetalServiceHandler) DetachPrivateNetwork(ctx context.Context, identifier string, networkID int64) error {
	if identifier == "" {
		return errBareMetalIdentifier
	}
	if networkID <= 0 {
		return errors.New("govpsie: private network id must be positive")
	}

	path := bareMetalPath(identifier, "private-networks", strconv.FormatInt(networkID, 10))
	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return err
	}

	return s.client.Do(ctx, req, nil)
}

func (s *bareMetalServiceHandler) Release(ctx context.Context, identifier string, release *BareMetalReleaseRequest) error {
	if identifier == "" {
		return errBareMetalIdentifier
	}
	if release == nil || release.Password == "" {
		return errors.New("govpsie: releasing a bare-metal server needs the account password")
	}
	if release.Reason == "" {
		return errors.New("govpsie: releasing a bare-metal server needs a reason")
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, bareMetalPath(identifier, "release"), release)
	if err != nil {
		return err
	}

	return s.client.Do(ctx, req, nil)
}
