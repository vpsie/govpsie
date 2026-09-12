package govpsie

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

var serverBasePath = "/apps/v2/vm"

type ServerService interface {
	ListServer(context.Context, *ListOptions, string) ([]VmData, error)
	List(context.Context, *ListOptions) ([]VmData, error)
	GetServerByIdentifier(context.Context, string) (*VmData, error)
	GetServerStatusByIdentifier(context.Context, string) (*Status, error)
	GetServerConsole(ctx context.Context, identifierId string) (*ServerConsole, error)
	CreateServer(context.Context, *CreateServerRequest) error
	DeleteServer(ctx context.Context, identifierId, password, reason, note string) error
	StartServer(ctx context.Context, identifierId string) error
	StopServer(ctx context.Context, identifierId string) error
	RestartServer(ctx context.Context, identifierId string) error
	ChangePassword(ctx context.Context, identifierId string, newPassword string) error
	ChangeHostName(ctx context.Context, identifierId string, newHostname string) error
	AddVPC(ctx context.Context, request *VpcRequest) error
	MoveVPC(ctx context.Context, request *VpcRequest) error
	AddTags(ctx context.Context, identifierId string, tags []string) error
	ResizeServer(ctx context.Context, identifierId, cpu, ram string) error
	AddSsh(ctx context.Context, identifierId, sshKeyIdentifier string) error
	AddScript(ctx context.Context, identifierId, scriptIdentifier string) error
	Lock(ctx context.Context, identifierId string) error
	UnLock(ctx context.Context, identifierId string) error
	DoMultiActions(ctx context.Context, vmsIdentifiers []string, actionType, sshKeyIdentifier string) error
	EnableIpv6(ctx context.Context, identifierId string) error
	EnableIpv4(ctx context.Context, identifierId string) error
	AddFip(ctx context.Context, identifierId, dcIdentifier string) error
	Resume(ctx context.Context, resumeReq *ResumeReq) error
	ResetNetwork(ctx context.Context, vmIdentifier string) error
	EditTag(ctx context.Context, tags []string, vmIdentifer string) error
	ResetAllFirewalls(ctx context.Context) error
	ListVirtualMachines(ctx context.Context) ([]VirtualMachine, error)
	ListAllNodesOfUser(ctx context.Context) ([]VmData, error)
	CheckAgentStatus(ctx context.Context, vmIdentifier string) (bool, error)
}

type serverServiceHandler struct {
	client *Client
}

var _ ServerService = &serverServiceHandler{}

type ListServerRoot struct {
	Error bool     `json:"error"`
	Data  []VmData `json:"data"`
	Total int64    `json:"total"`
}

type ListServerByIdentifierRoot struct {
	Error bool `json:"error"`
	Data  struct {
		VmData          VmData            `json:"vmData"`
		ImageCategories []ImageCategories `json:"imageCategories"`
		VmTags          []VmTags          `json:"vmTags"`
		PrivateIpData   []PrivateIpData   `json:"privateIpData"`
		FloatingIpData  []FloatingIpData  `json:"floatingIpData"`
	} `json:"data"`
	Total int64 `json:"total"`
}

type ImageCategories struct {
}
type VmTags struct {
}
type PrivateIpData struct {
}
type FloatingIpData struct {
}
type VmData struct {
	ID                  int64   `json:"id"`
	UserID              int64   `json:"user_id"`
	BoxSizeID           int64   `json:"boxsize_id"`
	BoxImageID          int64   `json:"boximage_id"`
	DataCenterID        int64   `json:"datacenter_id"`
	NodeID              int64   `json:"node_id"`
	BoxdIsCountID       *int64  `json:"boxdiscount_id"`
	Hostname            string  `json:"hostname"`
	DefaultIP           string  `json:"default_ip"`
	DefaultIPv6         string  `json:"default_ipv6"`
	PrivateIP           string  `json:"private_ip"`
	IsAutoBackup        int64   `json:"is_autobackup"`
	BoxVirtualization   string  `json:"box_virtualization_id"`
	Ram                 int64   `json:"ram"`
	Cpu                 int64   `json:"cpu"`
	Ssd                 int64   `json:"ssd"`
	Traffic             int64   `json:"traffic"`
	AddedIpAddresses    *string `json:"added_ip_addresses"`
	InitialPassword     string  `json:"initial_password"`
	Notes               *string `json:"notes"`
	CreatedOn           string  `json:"created_on"`
	LastUpdated         string  `json:"last_updated"`
	DroppedOn           *string `json:"dropped_on"`
	IsActive            int64   `json:"is_active"`
	IsDeleted           int64   `json:"is_deleted"`
	Identifier          string  `json:"identifier"`
	Power               int64   `json:"power"`
	ProjectID           int64   `json:"project_id"`
	IsCustom            int64   `json:"is_custom"`
	NrAddedIps          int64   `json:"nr_added_ips"`
	InPcs               int64   `json:"in_pcs"`
	CustomPrice         *int64  `json:"custom_price"`
	PayableLicense      int64   `json:"payable_license"`
	LastLicensePay      *string `json:"last_license_pay"`
	ScriptID            *string `json:"script_id"`
	SshKeyID            *string `json:"sshkey_id"`
	IsLocked            int64   `json:"is_locked"`
	IsWorkWithNew       int64   `json:"is_work_with_new_version"`
	IsSuspended         int64   `json:"is_suspended"`
	IsTerminated        int64   `json:"is_terminated"`
	OldID               int64   `json:"old_id"`
	CustomIsoID         *int64  `json:"custom_iso_id"`
	IsIsoImageBootAble  int64   `json:"is_iso_image_bootable"`
	HasSsl              int64   `json:"has_ssl"`
	LastActionDate      *string `json:"last_action_date,omitempty"`
	IsCreatedFromLegacy int64   `json:"is_created_from_legacy"`
	IsSmtpAllowed       int64   `json:"is_smtp_allowed"`
	WeeklyBackup        int64   `json:"weekly_backup"`
	MonthlyBackup       int64   `json:"monthly_backup"`
	LibIsoID            *int64  `json:"lib_iso_id,omitempty"`
	DailySnapshot       int64   `json:"daily_snapshot"`
	WeeklySnapshot      int64   `json:"weekly_snapshot"`
	MonthlySnapshot     int64   `json:"monthly_snap"`
	LastActionInMin     int64   `json:"last_action_in_min"`
	FirstName           string  `json:"firstname"`
	LastName            string  `json:"lastname"`
	Username            string  `json:"username"`
	State               string  `json:"state"`
	IsFipAvailable      int64   `json:"is_fip_available"`
	IsBucketAvailable   int64   `json:"is_bucket_available"`
	DcIdentifier        string  `json:"dcIdentifier"`
	Category            string  `json:"category"`
	FullName            string  `json:"fullname"`
	VmDescription       string  `json:"vmDescription"`
	BoxesSuspended      int64   `json:"boxes_suspended"`
	IsSataAvailable     int64   `json:"is_sata_available"`
	IsSsdAvailable      int64   `json:"is_ssd_available"`
	PublicIp            *string `json:"publicIp,omitempty"`
	VMType              string  `json:"vmType,omitempty"`
}

type Status struct {
	Cpu            int64  `json:"cpu"`
	Ballon         int64  `json:"ballon"`
	Uptime         int64  `json:"uptime"`
	Pid            string `json:"pid"`
	Disk           int64  `json:"disk"`
	RunningMachine string `json:"running-machine"`
	RunningQemu    string `json:"running-qemu"`
	Status         string `json:"status"`
	DiskRead       string `json:"diskread"`
	DiskWrite      string `json:"diskwrite"`
	Fullname       string `json:"fullname"`
}

type GetStatusRoot struct {
	Error  bool   `json:"error"`
	Status Status `json:"data"`
}

type ServerConsoleRoot struct {
	Error         bool          `json:"error"`
	ServerConsole ServerConsole `json:"data"`
}

type ServerConsole struct {
	Upid     string `json:"upid"`
	Ticket   string `json:"ticket"`
	Cert     string `json:"cert"`
	User     string `json:"user"`
	Port     string `json:"port"`
	Token    string `json:"token"`
	OS       string `json:"os"`
	Fullname string `json:"fullname"`
}

type CreateServerRequest struct {
	ResourceIdentifier  string    `json:"resourceIdentifier"`
	OsIdentifier        string    `json:"osIdentifier"`
	DcIdentifier        string    `json:"dcIdentifier"`
	Hostname            string    `json:"hostname"`
	Notes               *string   `json:"notes,omitempty"`
	BackupEnabled       *int64    `json:"backupEnabled,omitempty"`
	AddPublicIpV4       *int64    `json:"addPublicIpV4,omitempty"`
	AddPublicIpV6       *int64    `json:"addPublicIpV6,omitempty"`
	AddPrivateIp        *int64    `json:"addPrivateIp,omitempty"`
	SshKeyIdentifier    *string   `json:"sshKeyIdentifier,omitempty"`
	ProjectIdentifier   string    `json:"projectId"`
	VmPassword          string    `json:"vmPassword,omitempty"`
	IsGeneratedPassword bool      `json:"isGeneratedPassword"`
	Tags                []*string `json:"tags,omitempty"`
	ScriptIdentifier    *string   `json:"scriptIdentifier,omitempty"`
}

type ActionRequest struct {
	VmIdentifier string `json:"vmIdentifier"`
}

type VpcRequest struct {
	VmIdentifier string `json:"vmIdentifier"`
	VpcId        string `json:"vpcId"`
	DcIdentifier string `json:"dcIdentifier"`
}

type ResumeReq struct {
	VmIdentifier     string `json:"vmIdentifier"`
	OsIdentifier     string `json:"osIdentifier"`
	ScriptIdentifier string `json:"scriptIdentifier"`
	SshKeyIdentifier string `json:"sshKeyIdentifier"`
	HostName         string `json:"hostname"`
	Password         string `json:"password"`
	IsOnPremise      bool   `json:"isOnPremise"`
}

type ListAllNodesOfUserRoot struct {
	Error bool     `json:"error"`
	Data  []VmData `json:"data"`
}

type VirtualMachine struct {
	Hostname          string      `json:"hostname"`
	ID                int         `json:"id"`
	Identifier        string      `json:"identifier"`
	DefaultIP         string      `json:"default_ip"`
	DefaultIpv6       string      `json:"default_ipv6"`
	RAM               int         `json:"ram"`
	CPU               int         `json:"cpu"`
	Ssd               int         `json:"ssd"`
	IsSuspended       int         `json:"is_suspended"`
	IsLocked          int         `json:"is_locked"`
	IsActive          int         `json:"is_active"`
	CreatedOn         time.Time   `json:"created_on"`
	UserID            int         `json:"user_id"`
	PrivateIP         interface{} `json:"private_ip"`
	Power             int         `json:"power"`
	Traffic           int         `json:"traffic"`
	IsAgentActive     int         `json:"is_agent_active"`
	Firstname         string      `json:"firstname"`
	Lastname          string      `json:"lastname"`
	Username          string      `json:"username"`
	Category          string      `json:"category"`
	Fullname          string      `json:"fullname"`
	APIConfID         string      `json:"api_conf_id"`
	VMDescription     string      `json:"vmDescription"`
	State             string      `json:"state"`
	IsFipAvailable    int         `json:"is_fip_available"`
	IsBucketAvailable int         `json:"is_bucket_available"`
	DcIdentifier      string      `json:"dcIdentifier"`
}

type ListVirtualMachineRoot struct {
	Error bool             `json:"error"`
	Data  []VirtualMachine `json:"data"`
}

type AgentStatus struct {
	Error bool `json:"error"`
	Data  bool `json:"data"`
}

func (v *serverServiceHandler) ListServer(ctx context.Context, options *ListOptions, projectId string) ([]VmData, error) {
	path := fmt.Sprintf("%s?projectId=%s", serverBasePath, projectId)
	req, err := v.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	servers := new(ListServerRoot)
	if err = v.client.Do(ctx, req, servers); err != nil {
		return nil, err
	}

	return servers.Data, nil
}

func (v *serverServiceHandler) List(ctx context.Context, options *ListOptions) ([]VmData, error) {
	req, err := v.client.NewRequest(ctx, http.MethodGet, serverBasePath, nil)
	if err != nil {
		return nil, err
	}

	Servers := new(ListServerRoot)
	if err = v.client.Do(ctx, req, Servers); err != nil {
		return nil, err
	}

	return Servers.Data, nil
}

func (v *serverServiceHandler) GetServerByIdentifier(ctx context.Context, identifierId string) (*VmData, error) {
	path := fmt.Sprintf("%s/%s", serverBasePath, identifierId)
	req, err := v.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	Servers := new(ListServerByIdentifierRoot)
	if err = v.client.Do(ctx, req, Servers); err != nil {
		return nil, err
	}

	return &Servers.Data.VmData, nil
}

func (v *serverServiceHandler) GetServerStatusByIdentifier(ctx context.Context, identifierId string) (*Status, error) {
	path := fmt.Sprintf("%s/status/%s", serverBasePath, identifierId)
	req, err := v.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	status := new(GetStatusRoot)
	if err = v.client.Do(ctx, req, status); err != nil {
		return nil, err
	}

	return &status.Status, nil
}

func (v *serverServiceHandler) GetServerConsole(ctx context.Context, identifierId string) (*ServerConsole, error) {
	path := fmt.Sprintf("%s/console/%s", serverBasePath, identifierId)
	req, err := v.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	console := new(ServerConsoleRoot)
	if err = v.client.Do(ctx, req, console); err != nil {
		return nil, err
	}

	return &console.ServerConsole, nil
}

func (v *serverServiceHandler) CreateServer(ctx context.Context, server *CreateServerRequest) error {
	req, err := v.client.NewRequest(ctx, http.MethodPost, serverBasePath, server)
	if err != nil {
		return err
	}

	if err = v.client.Do(ctx, req, nil); err != nil {
		return err
	}

	return nil
}

func (v *serverServiceHandler) DeleteServer(ctx context.Context, identifierId, password, reason, note string) error {

	deleteReq := struct {
		VMIdentifier    string `json:"vmIdentifier"`
		Password        string `json:"password"`
		DeleteStatistic struct {
			Reason string `json:"reason"`
			Note   string `json:"note"`
		} `json:"deleteStatistic"`
	}{
		VMIdentifier: identifierId,
		Password:     password,
		DeleteStatistic: struct {
			Reason string `json:"reason"`
			Note   string `json:"note"`
		}{
			Reason: reason,
			Note:   note,
		},
	}

	req, err := v.client.NewRequest(ctx, http.MethodDelete, serverBasePath, deleteReq)
	if err != nil {
		return err
	}

	if err = v.client.Do(ctx, req, nil); err != nil {
		return err
	}

	return nil
}

func (v *serverServiceHandler) StartServer(ctx context.Context, identifierId string) error {
	vmIdentifier := &ActionRequest{
		VmIdentifier: identifierId,
	}

	path := fmt.Sprintf("%s/start", serverBasePath)
	req, err := v.client.NewRequest(ctx, http.MethodPost, path, vmIdentifier)
	if err != nil {
		return err
	}

	if err = v.client.Do(ctx, req, nil); err != nil {
		return err
	}

	return nil
}

func (v *serverServiceHandler) StopServer(ctx context.Context, identifierId string) error {
	vmIdentifier := &ActionRequest{
		VmIdentifier: identifierId,
	}

	path := fmt.Sprintf("%s/stop", serverBasePath)
	req, err := v.client.NewRequest(ctx, http.MethodPost, path, vmIdentifier)
	if err != nil {
		return err
	}

	if err = v.client.Do(ctx, req, nil); err != nil {
		return err
	}

	return nil
}

func (v *serverServiceHandler) RestartServer(ctx context.Context, identifierId string) error {
	vmIdentifier := &ActionRequest{
		VmIdentifier: identifierId,
	}
	path := fmt.Sprintf("%s/restart", serverBasePath)
	req, err := v.client.NewRequest(ctx, http.MethodPost, path, vmIdentifier)
	if err != nil {
		return err
	}

	if err = v.client.Do(ctx, req, nil); err != nil {
		return err
	}

	return nil
}

func (v *serverServiceHandler) ChangePassword(ctx context.Context, identifierId string, newPassword string) error {
	changePassReq := struct {
		VmIdentifier string `json:"vmIdentifier"`
		NewPassword  string `json:"newpassword"`
	}{
		VmIdentifier: identifierId,
		NewPassword:  newPassword,
	}
	path := fmt.Sprintf("%s/changepass", serverBasePath)
	req, err := v.client.NewRequest(ctx, http.MethodPost, path, changePassReq)
	if err != nil {
		return err
	}

	if err = v.client.Do(ctx, req, nil); err != nil {
		return err
	}

	return nil
}

func (v *serverServiceHandler) ChangeHostName(ctx context.Context, identifierId string, newHostname string) error {
	changeHostNameReq := struct {
		VmIdentifier string `json:"vmIdentifier"`
		Hostname     string `json:"hostname"`
	}{
		VmIdentifier: identifierId,
		Hostname:     newHostname,
	}
	path := fmt.Sprintf("%s/changehostname", serverBasePath)
	req, err := v.client.NewRequest(ctx, http.MethodPost, path, changeHostNameReq)
	if err != nil {
		return err
	}

	if err = v.client.Do(ctx, req, nil); err != nil {
		return err
	}

	return nil
}

func (v *serverServiceHandler) AddVPC(ctx context.Context, request *VpcRequest) error {
	path := fmt.Sprintf("%s/add/vpc", serverBasePath)

	req, err := v.client.NewRequest(ctx, http.MethodPost, path, request)
	if err != nil {
		return err
	}

	if err = v.client.Do(ctx, req, nil); err != nil {
		return err
	}

	return nil
}

func (v *serverServiceHandler) MoveVPC(ctx context.Context, request *VpcRequest) error {
	path := fmt.Sprintf("%s/vpc/move", serverBasePath)

	req, err := v.client.NewRequest(ctx, http.MethodPost, path, request)
	if err != nil {
		return err
	}

	if err = v.client.Do(ctx, req, nil); err != nil {
		return err
	}

	return nil
}

func (v *serverServiceHandler) ResizeServer(ctx context.Context, identifierId, cpu, ram string) error {
	path := fmt.Sprintf("%s/resize", serverBasePath)

	resizeServer := struct {
		VmIdentifier string `json:"vmIdentifier"`
		Ram          string `json:"ram"`
		Cpu          string `json:"cup"`
	}{
		VmIdentifier: identifierId,
		Ram:          ram,
		Cpu:          cpu,
	}

	req, err := v.client.NewRequest(ctx, http.MethodPost, path, resizeServer)
	if err != nil {
		return err
	}

	if err = v.client.Do(ctx, req, nil); err != nil {
		return err
	}

	return nil
}

func (v *serverServiceHandler) AddTags(ctx context.Context, identifierId string, tags []string) error {
	path := fmt.Sprintf("%s/addtags", serverBasePath)

	addTagsRequest := struct {
		VmIdentifier string   `json:"vmIdentifier"`
		Tags         []string `json:"tags"`
	}{
		VmIdentifier: identifierId,
		Tags:         tags,
	}

	req, err := v.client.NewRequest(ctx, http.MethodPost, path, addTagsRequest)
	if err != nil {
		return err
	}

	if err = v.client.Do(ctx, req, nil); err != nil {
		return err
	}

	return nil
}

func (v *serverServiceHandler) AddSsh(ctx context.Context, identifierId, sshKeyIdentifier string) error {
	path := fmt.Sprintf("%s/sshkey", serverBasePath)
	addSshReq := struct {
		VmIdentifier     string `json:"vmIdentifier"`
		SshKeyIdentifier string `json:"sshKeyIdentifier"`
	}{
		VmIdentifier:     identifierId,
		SshKeyIdentifier: sshKeyIdentifier,
	}

	req, err := v.client.NewRequest(ctx, http.MethodPost, path, addSshReq)
	if err != nil {
		return err
	}

	if err = v.client.Do(ctx, req, nil); err != nil {
		return err
	}

	return nil
}

func (v *serverServiceHandler) AddScript(ctx context.Context, identifierId, scriptIdentifier string) error {
	path := fmt.Sprintf("%s/script", serverBasePath)
	addScriptReq := struct {
		VmIdentifier     string `json:"vmIdentifier"`
		ScriptIdentifier string `json:"scriptIdentifier"`
	}{
		VmIdentifier:     identifierId,
		ScriptIdentifier: scriptIdentifier,
	}

	req, err := v.client.NewRequest(ctx, http.MethodPost, path, addScriptReq)
	if err != nil {
		return err
	}

	if err = v.client.Do(ctx, req, nil); err != nil {
		return err
	}

	return nil
}

func (v *serverServiceHandler) Lock(ctx context.Context, identifierId string) error {
	path := fmt.Sprintf("%s/lock", serverBasePath)
	lockReq := struct {
		VmIdentifier string `json:"vmIdentifier"`
	}{
		VmIdentifier: identifierId,
	}

	req, err := v.client.NewRequest(ctx, http.MethodPost, path, lockReq)
	if err != nil {
		return err
	}

	if err = v.client.Do(ctx, req, nil); err != nil {
		return err
	}

	return nil
}

func (v *serverServiceHandler) UnLock(ctx context.Context, identifierId string) error {
	path := fmt.Sprintf("%s/unlock", serverBasePath)
	unLockReq := struct {
		VmIdentifier string `json:"vmIdentifier"`
	}{
		VmIdentifier: identifierId,
	}

	req, err := v.client.NewRequest(ctx, http.MethodPost, path, unLockReq)
	if err != nil {
		return err
	}

	if err = v.client.Do(ctx, req, nil); err != nil {
		return err
	}

	return nil
}

func (v *serverServiceHandler) DoMultiActions(ctx context.Context, vmsIdentifiers []string, actionType, sshKeyIdentifier string) error {
	path := fmt.Sprintf("%s/actions", serverBasePath)
	doMultiActionsReq := struct {
		VmsIdentifiers   []string `json:"vmsIdentifiers"`
		ActionType       string   `json:"actionType"`
		SshKeyIdentifier string   `json:"sshKeyIdentifier"`
	}{
		VmsIdentifiers:   vmsIdentifiers,
		ActionType:       actionType,
		SshKeyIdentifier: sshKeyIdentifier,
	}

	req, err := v.client.NewRequest(ctx, http.MethodPost, path, doMultiActionsReq)
	if err != nil {
		return err
	}

	if err = v.client.Do(ctx, req, nil); err != nil {
		return err
	}

	return nil
}

func (v *serverServiceHandler) AddFip(ctx context.Context, identifierId, dcIdentifier string) error {
	path := fmt.Sprintf("%s/fip/add", serverBasePath)
	addFipReq := struct {
		VmIdentifier string `json:"vmIdentifier"`
		DcIdentifier string `json:"dcIdentifier"`
	}{
		VmIdentifier: identifierId,
		DcIdentifier: dcIdentifier,
	}

	req, err := v.client.NewRequest(ctx, http.MethodPost, path, addFipReq)
	if err != nil {
		return err
	}

	if err = v.client.Do(ctx, req, nil); err != nil {
		return err
	}

	return nil
}

func (v *serverServiceHandler) EnableIpv6(ctx context.Context, identifierId string) error {
	path := fmt.Sprintf("%s/enable/ipv6", serverBasePath)
	enableIpv6Req := struct {
		VmIdentifier string `json:"vmIdentifier"`
	}{
		VmIdentifier: identifierId,
	}

	req, err := v.client.NewRequest(ctx, http.MethodPost, path, enableIpv6Req)
	if err != nil {
		return err
	}

	if err = v.client.Do(ctx, req, nil); err != nil {
		return err
	}

	return nil
}

func (v *serverServiceHandler) EnableIpv4(ctx context.Context, identifierId string) error {
	path := fmt.Sprintf("%s/enable/ipv4", serverBasePath)
	enableIpv4Req := struct {
		VmIdentifier string `json:"vmIdentifier"`
	}{
		VmIdentifier: identifierId,
	}

	req, err := v.client.NewRequest(ctx, http.MethodPost, path, enableIpv4Req)
	if err != nil {
		return err
	}

	if err = v.client.Do(ctx, req, nil); err != nil {
		return err
	}

	return nil
}

func (v *serverServiceHandler) Resume(ctx context.Context, resumeReq *ResumeReq) error {
	path := fmt.Sprintf("%s/resume", serverBasePath)

	req, err := v.client.NewRequest(ctx, http.MethodPost, path, resumeReq)
	if err != nil {
		return err
	}

	if err = v.client.Do(ctx, req, nil); err != nil {
		return err
	}

	return nil
}

func (v *serverServiceHandler) ResetNetwork(ctx context.Context, vmIdentifier string) error {
	path := fmt.Sprintf("/apps/v2/refresh/vm/ips/%s", vmIdentifier)

	req, err := v.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return err
	}

	if err = v.client.Do(ctx, req, nil); err != nil {
		return err
	}

	return nil
}

func (v *serverServiceHandler) EditTag(ctx context.Context, tags []string, vmIdentifer string) error {
	path := fmt.Sprintf("%s/tags/edit", serverBasePath)

	req, err := v.client.NewRequest(ctx, http.MethodPost, path, nil)
	if err != nil {
		return err
	}

	if err = v.client.Do(ctx, req, nil); err != nil {
		return err
	}

	return nil
}

func (v *serverServiceHandler) ListVirtualMachines(ctx context.Context) ([]VirtualMachine, error) {
	path := "/apps/v2/virtualMachines"
	req, err := v.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	vms := new(ListVirtualMachineRoot)

	if err = v.client.Do(ctx, req, &vms); err != nil {
		return nil, err
	}

	return vms.Data, nil
}

func (v *serverServiceHandler) ListAllNodesOfUser(ctx context.Context) ([]VmData, error) {
	path := "/apps/v2/vms/all/user"
	req, err := v.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	nodes := new(ListAllNodesOfUserRoot)

	if err = v.client.Do(ctx, req, &nodes); err != nil {
		return nil, err
	}

	return nodes.Data, nil
}

func (v *serverServiceHandler) CheckAgentStatus(ctx context.Context, vmIdentifier string) (bool, error) {
	path := fmt.Sprintf("%s/live/agent/status/%s", serverBasePath, vmIdentifier)
	req, err := v.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return false, err
	}

	status := new(AgentStatus)

	if err = v.client.Do(ctx, req, &status); err != nil {
		return false, err
	}

	return status.Data, nil
}

func (v *serverServiceHandler) ResetAllFirewalls(ctx context.Context) error {
	path := fmt.Sprintf("%s/reset/firewall/for/all", serverBasePath)
	req, err := v.client.NewRequest(ctx, http.MethodPost, path, nil)
	if err != nil {
		return err
	}

	return v.client.Do(ctx, req, nil)
}
