package govpsie

import (
	"encoding/json"
	"strings"
	"testing"
)

// The API answers passThrough=true with 400 (TLS passthrough is not
// supported). A domain built without the field must still send false
// explicitly, so a request never depends on the API's default.
func TestLBDomainSendsPassThroughFalse(t *testing.T) {
	body, err := json.Marshal(LBDomain{DomainName: "example.com"})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(body), `"passThrough":false`) {
		t.Errorf("expected passThrough:false in %s", body)
	}
}

// corsHeaders is sent only when the caller sets it: omitted, the API keeps the
// permissive CORS headers on (its contract for API clients).
func TestLBDomainSendsCorsHeadersOnlyWhenSet(t *testing.T) {
	body, err := json.Marshal(LBDomain{DomainName: "example.com"})
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(body), "corsHeaders") {
		t.Errorf("corsHeaders must be omitted when unset, got %s", body)
	}
	off := false
	body, err = json.Marshal(LBDomain{DomainName: "example.com", CorsHeaders: &off})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(body), `"corsHeaders":false`) {
		t.Errorf("expected corsHeaders:false in %s", body)
	}
}

// GET /lb/{id} returns corsHeaders as 0 / 1; an older API does not return it.
func TestLBDomainsDetailReadsCorsHeaders(t *testing.T) {
	cases := []struct {
		json string
		want *bool
	}{
		{`{"corsHeaders":1}`, boolPtr(true)},
		{`{"corsHeaders":0}`, boolPtr(false)},
		{`{"corsHeaders":true}`, boolPtr(true)},
		{`{"corsHeaders":"0"}`, boolPtr(false)},
		{`{}`, nil},
		{`{"corsHeaders":null}`, nil},
	}
	for _, tc := range cases {
		var d LBDomainsDetail
		if err := json.Unmarshal([]byte(tc.json), &d); err != nil {
			t.Fatalf("%s: %v", tc.json, err)
		}
		got := d.CorsHeaders.Bool()
		if (got == nil) != (tc.want == nil) || (got != nil && *got != *tc.want) {
			t.Errorf("%s: got %v, want %v", tc.json, got, tc.want)
		}
	}
	var bad LBDomainsDetail
	if err := json.Unmarshal([]byte(`{"corsHeaders":"maybe"}`), &bad); err == nil {
		t.Errorf("an unknown value must be an error")
	}
}

func boolPtr(v bool) *bool { return &v }
