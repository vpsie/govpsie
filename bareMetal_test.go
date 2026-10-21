package govpsie

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"
)

func TestBareMetal_ListPlans(t *testing.T) {
	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/apps/v2/bare-metal/dc/dc-1/hardware-plans" {
			t.Fatalf("unexpected request %s %s", r.Method, r.URL.Path)
		}
		q := r.URL.Query()
		if q.Get("minRamGb") != "64" || q.Get("minDiskGb") != "" || q.Get("gpu") != "false" {
			t.Fatalf("unexpected query %q", r.URL.RawQuery)
		}
		_, _ = w.Write([]byte(`{"error":false,"data":[{"identifier":"plan-1","name":"E-2388G","cpuCount":8,"ramGb":64,
			"price":"149.00","setupFee":0,"delivery":"instant","contractTerms":[{"term":"monthly","discountPct":10}],
			"linkSpeedOptions":[1000,10000],"availableServers":2,"isActive":1}]}`))
	})

	gpu := false
	plans, err := c.BareMetal.ListPlans(context.Background(), "dc-1", &BareMetalPlanFilter{MinRAMGB: 64, GPU: &gpu})
	if err != nil {
		t.Fatalf("ListPlans: %v", err)
	}
	if len(plans) != 1 {
		t.Fatalf("plans = %d, want 1", len(plans))
	}
	p := plans[0]
	if p.Identifier != "plan-1" || p.RAMGB != 64 || !p.IsActive || p.Price.String() != "149.00" || p.SetupFee.String() != "0" {
		t.Fatalf("unexpected plan %+v", p)
	}
	if len(p.ContractTerms) != 1 || p.ContractTerms[0].DiscountPct != 10 || len(p.LinkSpeedOptions) != 2 {
		t.Fatalf("unexpected plan terms %+v", p)
	}
}

func TestBareMetal_ListPlansNeedsDatacenter(t *testing.T) {
	c := NewClient(nil)
	if _, err := c.BareMetal.ListPlans(context.Background(), "", nil); err == nil {
		t.Fatal("expected an error without a datacenter")
	}
}

func TestBareMetal_ListOS(t *testing.T) {
	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/apps/v2/bare-metal/dc/dc-1/os" || r.URL.Query().Get("planId") != "plan-1" {
			t.Fatalf("unexpected request %s?%s", r.URL.Path, r.URL.RawQuery)
		}
		_, _ = w.Write([]byte(`{"error":false,"data":[{"name":"Ubuntu","is_default":1,"available_versions":[
			{"identifier":"bmos_abc","fullname":"Ubuntu 24.04","osFamily":"linux","remoteAccessConfigurable":true,"is_default":1},
			{"identifier":"bmos_def","fullname":"Ubuntu 22.04","osFamily":"linux","remoteAccessConfigurable":true,"is_default":0}]}]}`))
	})

	groups, err := c.BareMetal.ListOS(context.Background(), "dc-1", "plan-1")
	if err != nil {
		t.Fatalf("ListOS: %v", err)
	}
	if len(groups) != 1 || len(groups[0].Versions) != 2 || !groups[0].IsDefault {
		t.Fatalf("unexpected groups %+v", groups)
	}
	v := groups[0].Versions[0]
	if v.Identifier != "bmos_abc" || v.FullName != "Ubuntu 24.04" || v.OSFamily != "linux" || !v.IsDefault {
		t.Fatalf("unexpected version %+v", v)
	}
	if groups[0].Versions[1].IsDefault {
		t.Fatalf("is_default 0 decoded as true: %+v", groups[0].Versions[1])
	}
}

// TestBareMetal_FlagsAcceptNumbersAndBooleans pins the flag encodings: the
// API sends isActive and is_default as TINYINT 1 / 0, and a boolean must
// still decode should an endpoint switch to it.
func TestBareMetal_FlagsAcceptNumbersAndBooleans(t *testing.T) {
	for _, tt := range []struct {
		in   string
		want bool
	}{
		{in: `1`, want: true},
		{in: `0`, want: false},
		{in: `true`, want: true},
		{in: `false`, want: false},
		{in: `"1"`, want: true},
	} {
		var plan BareMetalPlan
		if err := json.Unmarshal([]byte(`{"isActive":`+tt.in+`}`), &plan); err != nil {
			t.Fatalf("plan isActive %s: %v", tt.in, err)
		}
		var group BareMetalOSGroup
		if err := json.Unmarshal([]byte(`{"is_default":`+tt.in+`,"available_versions":[{"is_default":`+tt.in+`}]}`), &group); err != nil {
			t.Fatalf("os is_default %s: %v", tt.in, err)
		}
		if bool(plan.IsActive) != tt.want || bool(group.IsDefault) != tt.want || bool(group.Versions[0].IsDefault) != tt.want {
			t.Fatalf("%s decoded as plan=%v group=%v version=%v, want %v", tt.in, plan.IsActive, group.IsDefault, group.Versions[0].IsDefault, tt.want)
		}
	}

	var plan BareMetalPlan
	if err := json.Unmarshal([]byte(`{"isActive":2}`), &plan); err == nil {
		t.Fatal("isActive 2 was accepted")
	}
}

func TestBareMetal_DeploySendsBodyAndIdempotencyKey(t *testing.T) {
	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/apps/v2/bare-metal/deploy" {
			t.Fatalf("unexpected request %s %s", r.Method, r.URL.Path)
		}
		if got := r.Header.Get("Idempotency-Key"); got != "key-1" {
			t.Fatalf("Idempotency-Key = %q", got)
		}
		body := decodeBody(t, r)
		want := map[string]interface{}{
			"dcId": "dc-1", "planId": "plan-1", "hostname": "web-1", "password": "Secret123!",
			"osIdentifier": "bmos_abc", "contractTerm": "monthly", "linkSpeedMbps": float64(10000),
			"userData": "#cloud-config\n", "projectId": "proj-1",
		}
		for k, v := range want {
			if body[k] != v {
				t.Fatalf("body[%s] = %v, want %v (body %v)", k, body[k], v, body)
			}
		}
		for _, absent := range []string{"templateName", "dnsZoneIdentifier"} {
			if _, ok := body[absent]; ok {
				t.Fatalf("body carries empty %s: %v", absent, body)
			}
		}
		keys, ok := body["sshKeyIdentifiers"].([]interface{})
		if !ok || len(keys) != 1 || keys[0] != "key-a" {
			t.Fatalf("sshKeyIdentifiers = %v", body["sshKeyIdentifiers"])
		}
		ra, ok := body["remoteAccess"].(map[string]interface{})
		if !ok || ra["enabled"] != true || ra["winrm"] != false || ra["firewallIdentifier"] != nil {
			t.Fatalf("remoteAccess = %v", body["remoteAccess"])
		}
		scripts, ok := body["scriptIds"].([]interface{})
		if !ok || len(scripts) != 1 || scripts[0] != "script-1" {
			t.Fatalf("scriptIds = %v", body["scriptIds"])
		}
		_, _ = w.Write([]byte(`{"error":false,"data":{"success":true,"identifier":"bm-1","osPhase":"installing",
			"operationId":42,"ipAddress":"203.0.113.10","status":"sold","contractTerm":"monthly","linkSpeedMbps":10000}}`))
	})

	got, err := c.BareMetal.Deploy(context.Background(), &BareMetalDeployRequest{
		DCIdentifier: "dc-1", PlanIdentifier: "plan-1", Hostname: "web-1", Password: "Secret123!",
		OSIdentifier: "bmos_abc", SSHKeyIdentifiers: []string{"key-a"},
		RemoteAccess: &BareMetalRemoteAccess{Enabled: true}, ProjectIdentifier: "proj-1",
		ContractTerm: "monthly", LinkSpeedMbps: 10000, UserData: "#cloud-config\n",
		ScriptIdentifiers: []string{"script-1"},
	}, "key-1")
	if err != nil {
		t.Fatalf("Deploy: %v", err)
	}
	if got.Identifier != "bm-1" || got.Status != BareMetalOrderSold || got.OperationID.String() != "42" || got.IPAddress != "203.0.113.10" {
		t.Fatalf("unexpected deployment %+v", got)
	}
}

func TestBareMetal_DeployIsNotReplayedOnGatewayError(t *testing.T) {
	calls := 0
	c := setup(t, func(w http.ResponseWriter, _ *http.Request) {
		calls++
		w.WriteHeader(http.StatusBadGateway)
	})

	_, err := c.BareMetal.Deploy(context.Background(), &BareMetalDeployRequest{DCIdentifier: "dc-1"}, "key-1")
	if err == nil {
		t.Fatal("expected an error")
	}
	if calls != 1 {
		t.Fatalf("deploy sent %d times, want 1 (POST must not be replayed by the client)", calls)
	}
}

func TestBareMetal_DeployErrorCarriesStatusAndCode(t *testing.T) {
	c := setup(t, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusConflict)
		_, _ = w.Write([]byte(`{"error":true,"code":"bareMetalDeployInProgress","message":"An order with this key is still running"}`))
	})

	_, err := c.BareMetal.Deploy(context.Background(), &BareMetalDeployRequest{DCIdentifier: "dc-1"}, "")
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *APIError, got %T %v", err, err)
	}
	if apiErr.StatusCode != http.StatusConflict || apiErr.Code != "bareMetalDeployInProgress" {
		t.Fatalf("unexpected APIError %+v", apiErr)
	}
	if err.Error() != "An order with this key is still running" || ErrorCode(err) != "bareMetalDeployInProgress" {
		t.Fatalf("unexpected error text/code %q %q", err.Error(), ErrorCode(err))
	}
}

func TestBareMetal_Get(t *testing.T) {
	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/apps/v2/bare-metal/bm-1" {
			t.Fatalf("unexpected path %s", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"error":false,"data":{"identifier":"bm-1","name":"web-1","powerState":"on","osFamily":"linux",
			"osName":"Ubuntu 24.04","osPhase":"running","displayState":"ready","ip_address":"203.0.113.10",
			"netmask":"255.255.255.0","gateway":"203.0.113.1","contract":{"term":"monthly","discountPct":10},
			"linkSpeedMbps":10000,"projectIdentifier":"proj-1","datacenterIdentifier":"dc-1",
			"remoteAccess":{"enabled":true,"winrm":false,"firewallIdentifier":null},
			"hardwareSpecs":{"processor_model":"E-2388G","cores":8,"ram_gb":64,"storage_gb":960},
			"adminHold":{"suspended":false},
			"privateNetworks":{"available":true,"items":[{"networkId":7,"status":"attached","vlanId":120}]}}}`))
	})

	bm, err := c.BareMetal.Get(context.Background(), "bm-1")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if bm.Name != "web-1" || bm.OSPhase != BareMetalOSPhaseRunning || bm.Contract == nil || bm.Contract.Term != "monthly" {
		t.Fatalf("unexpected server %+v", bm)
	}
	if bm.RemoteAccess == nil || !bm.RemoteAccess.Enabled || bm.Hardware == nil || bm.Hardware.RAMGB != 64 {
		t.Fatalf("unexpected nested fields %+v", bm)
	}
	if bm.PrivateNetworks == nil || len(bm.PrivateNetworks.Items) != 1 || bm.PrivateNetworks.Items[0].NetworkID != 7 {
		t.Fatalf("unexpected private networks %+v", bm.PrivateNetworks)
	}
}

func TestBareMetal_GetNotFound(t *testing.T) {
	c := setup(t, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":true,"message":"Machine not found"}`))
	})

	_, err := c.BareMetal.Get(context.Background(), "bm-gone")
	if !IsNotFound(err) {
		t.Fatalf("IsNotFound(%v) = false", err)
	}
}

func TestBareMetal_GetNullDataIsNotFound(t *testing.T) {
	c := setup(t, func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"error":false,"data":null}`))
	})

	_, err := c.BareMetal.Get(context.Background(), "bm-1")
	if !IsNotFound(err) {
		t.Fatalf("IsNotFound(%v) = false", err)
	}
}

func TestBareMetal_GetNeedsIdentifier(t *testing.T) {
	c := NewClient(nil)
	if _, err := c.BareMetal.Get(context.Background(), ""); !errors.Is(err, errBareMetalIdentifier) {
		t.Fatalf("expected errBareMetalIdentifier, got %v", err)
	}
}

func TestBareMetal_List(t *testing.T) {
	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/apps/v2/bare-metal" {
			t.Fatalf("unexpected path %s", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"success":true,"data":[{"identifier":"bm-1","machine_name":"web-1","status":"deployed",
			"datacenter_identifier":"dc-1","ip_address":"203.0.113.10","displayState":"ready"}],"timestamp":"x"}`))
	})

	list, err := c.BareMetal.List(context.Background())
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(list) != 1 || list[0].MachineName != "web-1" || list[0].DatacenterIdentifier != "dc-1" {
		t.Fatalf("unexpected list %+v", list)
	}
}

func TestBareMetal_UpdateHostnameAndProject(t *testing.T) {
	seen := map[string]map[string]interface{}{}
	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Fatalf("unexpected method %s", r.Method)
		}
		seen[r.URL.Path] = decodeBody(t, r)
		_, _ = w.Write([]byte(`{"error":false,"data":{}}`))
	})

	ctx := context.Background()
	if err := c.BareMetal.UpdateHostname(ctx, "bm-1", "web-2"); err != nil {
		t.Fatalf("UpdateHostname: %v", err)
	}
	if err := c.BareMetal.MoveToProject(ctx, "bm-1", "proj-2"); err != nil {
		t.Fatalf("MoveToProject: %v", err)
	}
	if seen["/apps/v2/bare-metal/bm-1/hostname"]["hostname"] != "web-2" {
		t.Fatalf("hostname body %v", seen)
	}
	if seen["/apps/v2/bare-metal/bm-1/project"]["projectId"] != "proj-2" {
		t.Fatalf("project body %v", seen)
	}
}

func TestBareMetal_ResetPassword(t *testing.T) {
	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/apps/v2/bare-metal/bm-1/reset-password" {
			t.Fatalf("unexpected request %s %s", r.Method, r.URL.Path)
		}
		body := decodeBody(t, r)
		if body["password"] != "NewSecret1!" {
			t.Fatalf("body %v", body)
		}
		if _, ok := body["osType"]; ok {
			t.Fatalf("empty osType must be omitted: %v", body)
		}
		_, _ = w.Write([]byte(`{"error":false,"data":{"success":true,"started":true,"osPhase":"password_reset"}}`))
	})

	got, err := c.BareMetal.ResetPassword(context.Background(), "bm-1", &BareMetalPasswordResetRequest{Password: "NewSecret1!"})
	if err != nil {
		t.Fatalf("ResetPassword: %v", err)
	}
	if !got.Started || got.OSPhase != BareMetalOSPhasePasswordReset {
		t.Fatalf("unexpected answer %+v", got)
	}
}

func TestBareMetal_Power(t *testing.T) {
	var paths []string
	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Fatalf("unexpected method %s", r.Method)
		}
		paths = append(paths, r.URL.Path)
		_, _ = w.Write([]byte(`{"error":false,"data":{}}`))
	})

	ctx := context.Background()
	for _, a := range []BareMetalPowerAction{BareMetalPowerOn, BareMetalPowerOff, BareMetalPowerShutdown, BareMetalPowerReboot} {
		if err := c.BareMetal.Power(ctx, "bm-1", a); err != nil {
			t.Fatalf("Power(%s): %v", a, err)
		}
	}
	want := []string{"/apps/v2/bare-metal/bm-1/on", "/apps/v2/bare-metal/bm-1/off", "/apps/v2/bare-metal/bm-1/shutdown", "/apps/v2/bare-metal/bm-1/reboot"}
	for i := range want {
		if paths[i] != want[i] {
			t.Fatalf("paths = %v, want %v", paths, want)
		}
	}

	if err := c.BareMetal.Power(ctx, "bm-1", "release"); err == nil {
		t.Fatal("expected an unknown action to be refused before any request")
	}
	if len(paths) != len(want) {
		t.Fatalf("an unknown action reached the API: %v", paths)
	}
}

func TestBareMetal_PowerState(t *testing.T) {
	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/apps/v2/bare-metal/bm-1/power-state" {
			t.Fatalf("unexpected path %s", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"error":false,"data":{"powerState":"off"}}`))
	})

	state, err := c.BareMetal.PowerState(context.Background(), "bm-1")
	if err != nil || state != "off" {
		t.Fatalf("PowerState = %q, %v", state, err)
	}
}

func TestBareMetal_PrivateNetworks(t *testing.T) {
	var calls []string
	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		calls = append(calls, r.Method+" "+r.URL.Path)
		switch r.Method {
		case http.MethodGet:
			_, _ = w.Write([]byte(`{"error":false,"data":{"available":true,"maxNetworks":8,"items":[{"networkId":7,"status":"attached"}]}}`))
		case http.MethodPost:
			if body := decodeBody(t, r); body["networkId"] != float64(9) {
				t.Fatalf("attach body %v", body)
			}
			w.WriteHeader(http.StatusAccepted)
			_, _ = w.Write([]byte(`{"error":false,"data":{"networkId":9,"status":"attaching"}}`))
		case http.MethodDelete:
			w.WriteHeader(http.StatusAccepted)
			_, _ = w.Write([]byte(`{"error":false,"data":{"networkId":7,"status":"detaching"}}`))
		}
	})

	ctx := context.Background()
	list, err := c.BareMetal.ListPrivateNetworks(ctx, "bm-1")
	if err != nil || !list.Available || list.MaxNetworks != 8 || len(list.Items) != 1 {
		t.Fatalf("ListPrivateNetworks = %+v, %v", list, err)
	}
	item, err := c.BareMetal.AttachPrivateNetwork(ctx, "bm-1", 9)
	if err != nil || item.NetworkID != 9 || item.Status != BareMetalNetworkAttaching {
		t.Fatalf("AttachPrivateNetwork = %+v, %v", item, err)
	}
	if err := c.BareMetal.DetachPrivateNetwork(ctx, "bm-1", 7); err != nil {
		t.Fatalf("DetachPrivateNetwork: %v", err)
	}
	want := []string{
		"GET /apps/v2/bare-metal/bm-1/private-networks",
		"POST /apps/v2/bare-metal/bm-1/private-networks",
		"DELETE /apps/v2/bare-metal/bm-1/private-networks/7",
	}
	for i := range want {
		if calls[i] != want[i] {
			t.Fatalf("calls = %v, want %v", calls, want)
		}
	}
	if _, err := c.BareMetal.AttachPrivateNetwork(ctx, "bm-1", 0); err == nil {
		t.Fatal("expected a non-positive network id to be refused")
	}
}

func TestBareMetal_PrivateNetworkBusyCode(t *testing.T) {
	c := setup(t, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusConflict)
		_, _ = w.Write([]byte(`{"error":true,"success":false,"code":"private_network_busy","message":"Another change is running"}`))
	})

	_, err := c.BareMetal.AttachPrivateNetwork(context.Background(), "bm-1", 9)
	if ErrorCode(err) != "private_network_busy" {
		t.Fatalf("ErrorCode(%v) = %q", err, ErrorCode(err))
	}
}

func TestBareMetal_Release(t *testing.T) {
	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/apps/v2/bare-metal/bm-1/release" {
			t.Fatalf("unexpected request %s %s", r.Method, r.URL.Path)
		}
		body := decodeBody(t, r)
		if body["password"] != "account-pw" || body["reason"] != "Destroyed by Terraform" {
			t.Fatalf("body %v", body)
		}
		if _, ok := body["note"]; ok {
			t.Fatalf("empty note must be omitted: %v", body)
		}
		_, _ = w.Write([]byte(`{"error":false,"data":{"success":true}}`))
	})

	err := c.BareMetal.Release(context.Background(), "bm-1", &BareMetalReleaseRequest{Password: "account-pw", Reason: "Destroyed by Terraform"})
	if err != nil {
		t.Fatalf("Release: %v", err)
	}
}

func TestBareMetal_ReleaseNeedsPasswordAndReason(t *testing.T) {
	c := setup(t, func(_ http.ResponseWriter, r *http.Request) {
		t.Fatalf("no request expected, got %s %s", r.Method, r.URL.Path)
	})

	ctx := context.Background()
	if err := c.BareMetal.Release(ctx, "bm-1", &BareMetalReleaseRequest{Reason: "x"}); err == nil {
		t.Fatal("expected an error without the account password")
	}
	if err := c.BareMetal.Release(ctx, "bm-1", &BareMetalReleaseRequest{Password: "pw"}); err == nil {
		t.Fatal("expected an error without a reason")
	}
}

func TestAPIError_NumericCodeAndNonJSONBody(t *testing.T) {
	numeric := newAPIError(http.StatusBadRequest, []byte(`{"error":true,"code":400,"message":"Bad input"}`))
	if numeric.Code != "400" || numeric.Error() != "Bad input" {
		t.Fatalf("numeric code: %+v", numeric)
	}

	html := newAPIError(http.StatusBadGateway, []byte(`<html>bad gateway</html>`))
	if html.Error() != "vpsie: unexpected response: 502 Bad Gateway" || html.Code != "" {
		t.Fatalf("non-JSON body: %+v", html)
	}
}
