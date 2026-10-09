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
