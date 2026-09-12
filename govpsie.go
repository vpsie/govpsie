package govpsie

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	libraryVersion = "2.0"
	defaultBaseURL = "https://api.vpsie.com/apps/v2"
	userAgent      = "vpsiecli/" + libraryVersion
	mediaType      = "application/json"
)

type Client struct {
	// HTTP client used to communicate with the VPSIE API.
	client *http.Client

	// Base URL for API requests.
	BaseURL *url.URL

	// User agent for client
	UserAgent string
	headers   map[string]string

	// AccountPassword is the password of the account the access token belongs
	// to. A few destructive endpoints -- deleting a server, notably -- confirm
	// the operation against the account password rather than the resource's own
	// credentials, and reject the call without it.
	AccountPassword string

	// services
	Account       AccountService
	Project       ProjectsService
	Server        ServerService
	Image         ImagesService
	SShKey        SshkeysService
	Profile       ProfilesService
	Backup        BackupsService
	IP            IPsService
	Domain        DomainService
	Fip           FipService
	FirewallGroup FirewallGroupService
	Firewall      FirewallService
	Storage       StorageService
	Snapshot      SnapshotService
	Logs          LogsService
	DataCenter    DataCenterService
	LB            LBsService
	Scripts       ScriptsService
	Pending       PendingService
	Gateway       GatewayService
	VPC           VPCService
	Bucket        BucketService
	K8s           K8sService
	AccessToken   AccessTokenService
	Billing       BillingService
	Monitoring    MonitoringService
	Tags          TagsService
	Certificate   CertificateService
	ServerGroup   ServerGroupService
	Registry      RegistryService
	ManagedDB     ManagedDBService
}

type ErrorRsp struct {
	Error   bool   `json:"error"`
	Code    int    `json:"code"`
	Message string `json:"message"`
	Stack   string `json:"stack"`
}

type GeneralRspRoot struct {
	Error bool        `json:"error"`
	Data  interface{} `json:"data"`
}

// ListOptions specifies the optional parameters to various List methods that support pagination.
type ListOptions struct {
	// For paginated result sets, page of results to retrieve.
	Page int `url:"page,omitempty"`

	// For paginated result sets, the number of results to include per page.
	PerPage int `url:"per_page,omitempty"`
}

func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	baseURL, _ := url.Parse(defaultBaseURL)

	c := &Client{
		client:    httpClient,
		BaseURL:   baseURL,
		UserAgent: userAgent,
	}

	c.Account = &accountServiceHandler{client: c}
	c.Project = &projectsServiceHandler{client: c}
	c.Server = &serverServiceHandler{client: c}
	c.Image = &imagesServiceHandler{client: c}
	c.SShKey = &sshkeysServiceHandler{client: c}
	c.Profile = &profilesServiceHandler{client: c}
	c.Backup = &backupsServiceHandler{client: c}
	c.IP = &iPsServiceHandler{client: c}
	c.Domain = &domainsServiceHandler{client: c}
	c.Fip = &fipServiceHandler{client: c}
	c.FirewallGroup = &firewallGroupServiceHandler{client: c}
	c.Firewall = &firewallServiceHandler{client: c}
	c.Storage = &storageServiceHandler{client: c}
	c.Snapshot = &snapshotServiceHandler{client: c}
	c.Logs = &logsServiceHandler{client: c}
	c.DataCenter = &dataCenterServiceHandler{client: c}
	c.LB = &lbsServiceHandler{client: c}
	c.Pending = &pendingServiceHandler{client: c}
	c.Scripts = &scriptsServiceHandler{client: c}
	c.Gateway = &gatewayServiceHandler{client: c}
	c.VPC = &vpcServiceHandler{client: c}
	c.Bucket = &bucketServiceHandler{client: c}
	c.K8s = &k8sServiceHandler{client: c}
	c.AccessToken = &accessTokenServiceHandler{client: c}
	c.Billing = &billingServiceHandler{client: c}
	c.Monitoring = &monitoringServiceHandler{client: c}
	c.Tags = &tagsServiceHandler{client: c}
	c.Certificate = &certificateServiceHandler{client: c}
	c.ServerGroup = &serverGroupServiceHandler{client: c}
	c.Registry = &registryServiceHandler{client: c}
	c.ManagedDB = &managedDBServiceHandler{client: c}

	c.headers = make(map[string]string)
	return c
}

// SetAccountPassword records the account password used to confirm destructive
// operations. See Client.AccountPassword.
func (c *Client) SetAccountPassword(password string) {
	c.AccountPassword = password
}

func (c *Client) SetRequestHeaders(headers map[string]string) {

	for k, v := range headers {
		c.headers[k] = v
	}
}

// SetUserAgent Overrides the default UserAgent
func (c *Client) SetUserAgent(ua string) {
	c.UserAgent = ua
}

// SetBaseURL Overrides the default BaseUrl
func (c *Client) SetBaseURL(baseURL string) error {
	updatedURL, err := url.Parse(baseURL)

	if err != nil {
		return err
	}

	c.BaseURL = updatedURL
	return nil
}

// value pointed to by body is JSON encoded and included in as the request body.
func (c *Client) NewRequest(ctx context.Context, method, urlStr string, body interface{}) (*http.Request, error) {
	u, err := c.BaseURL.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	var req *http.Request
	switch method {
	case http.MethodGet, http.MethodHead, http.MethodOptions:
		req, err = http.NewRequestWithContext(ctx, method, u.String(), nil)
		if err != nil {
			return nil, err
		}

	default:
		buf := new(bytes.Buffer)
		if body != nil {
			err = json.NewEncoder(buf).Encode(body)
			if err != nil {
				return nil, err
			}
		}

		req, err = http.NewRequestWithContext(ctx, method, u.String(), buf)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", mediaType)
	}

	for k, v := range c.headers {
		req.Header.Add(k, v)
	}

	req.Header.Set("Accept", mediaType)
	req.Header.Set("User-Agent", c.UserAgent)

	return req, nil
}

// retryableStatus reports whether an HTTP status represents a transient
// gateway/upstream failure that is worth retrying.
func retryableStatus(code int) bool {
	switch code {
	case http.StatusBadGateway, http.StatusServiceUnavailable, http.StatusGatewayTimeout:
		return true
	default:
		return false
	}
}

// retryableMethod reports whether a request may be safely replayed.
//
// Only methods that are idempotent by HTTP semantics are retried. POST and
// PATCH are deliberately excluded: the VPSIE API provisions billable resources
// on POST, and a gateway error gives no way to tell whether the upstream
// already accepted the request, so replaying one risks creating duplicates.
func retryableMethod(method string) bool {
	switch method {
	case http.MethodGet, http.MethodHead, http.MethodOptions, http.MethodPut, http.MethodDelete:
		return true
	default:
		return false
	}
}

// Do sends an API request and decodes the response into v.
//
// Transient gateway failures (502/503/504) are retried with exponential
// backoff for idempotent methods, because the API sits behind a proxy that
// intermittently returns them while an upstream is rolling or unhealthy.
// Without this a momentary blip fails the caller mid-operation, which for a
// create can leave a billable resource behind that nothing is tracking.
func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) error {
	const maxAttempts = 4

	backoff := 500 * time.Millisecond
	var lastErr error

	for attempt := 1; ; attempt++ {
		if attempt > 1 {
			// Rewind the body so the request can be replayed.
			if req.GetBody != nil {
				body, err := req.GetBody()
				if err != nil {
					return lastErr
				}
				req.Body = body
			}

			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(backoff):
			}
			backoff *= 2
		}

		err, retryable := c.do(ctx, req, v)
		if err == nil {
			return nil
		}

		lastErr = err
		if !retryable || attempt == maxAttempts || !retryableMethod(req.Method) {
			return err
		}
	}
}

// do performs a single attempt and reports whether the failure is transient.
func (c *Client) do(ctx context.Context, req *http.Request, v interface{}) (error, bool) {
	res, err := c.client.Do(req.WithContext(ctx))
	if err != nil {
		// A transport-level failure means no response was produced.
		return err, true
	}

	defer func() { _ = res.Body.Close() }()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err, true
	}

	if res.StatusCode == http.StatusNoContent {
		return nil, false
	}

	if res.StatusCode < http.StatusOK || res.StatusCode >= 300 {
		retryable := retryableStatus(res.StatusCode)

		var errRsp ErrorRsp
		if jsonErr := json.Unmarshal(body, &errRsp); jsonErr != nil || errRsp.Message == "" {
			// The body was not the expected JSON envelope (e.g. an HTML error
			// page from a proxy) or carried no message; fall back to the status.
			return fmt.Errorf("vpsie: unexpected response: %d %s", res.StatusCode, http.StatusText(res.StatusCode)), retryable
		}

		return errors.New(errRsp.Message), retryable
	}

	if v != nil {
		if err := json.Unmarshal(body, v); err != nil {
			return err, false
		}
	}

	return nil, false
}

// StreamToString converts a reader to a string
func StreamToString(stream io.Reader) string {
	buf := new(bytes.Buffer)
	_, _ = buf.ReadFrom(stream)
	return buf.String()
}

// isNullData reports whether a JSON `data` payload is absent. Several VPSie
// endpoints return HTTP 200 with a literal `"data": false` (or `null`) to mean
// "not found" rather than a 404, so callers use this to distinguish a missing
// resource from a real one.
func isNullData(raw json.RawMessage) bool {
	trimmed := bytes.TrimSpace(raw)
	return len(trimmed) == 0 ||
		bytes.Equal(trimmed, []byte("null")) ||
		bytes.Equal(trimmed, []byte("false"))
}
