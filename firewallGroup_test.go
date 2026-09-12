package govpsie

import (
	"context"
	"os"
	"testing"

	"golang.org/x/oauth2"
)

// TestFirewallGroupServiceHandlerDelete exercises a live Delete call. It talks
// to the real VPSie API, so it only runs when VPSIE_ACCESS_TOKEN is set;
// otherwise it is skipped to keep `go test ./...` green without credentials.
func TestFirewallGroupServiceHandlerDelete(t *testing.T) {
	token := os.Getenv("VPSIE_ACCESS_TOKEN")
	if token == "" {
		t.Skip("VPSIE_ACCESS_TOKEN not set; skipping live firewall group delete test")
	}

	identifier := os.Getenv("VPSIE_FIREWALL_GROUP_IDENTIFIER")
	if identifier == "" {
		t.Skip("VPSIE_FIREWALL_GROUP_IDENTIFIER not set; skipping live firewall group delete test")
	}

	client := NewClient(oauth2.NewClient(context.Background(), nil))

	client.SetUserAgent(userAgent)
	client.SetRequestHeaders(map[string]string{
		"Vpsie-Auth": token,
	})

	if err := client.FirewallGroup.Delete(context.Background(), identifier); err != nil {
		t.Error(err)
	}
}
