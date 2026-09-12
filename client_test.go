package govpsie

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// setup returns a Client whose base URL points at a test server driven by
// handler. Service path constants are absolute, so they replace the base URL's
// path and arrive at the handler unchanged (e.g. /apps/v2/tags).
func setup(t *testing.T, handler http.HandlerFunc) *Client {
	t.Helper()
	ts := httptest.NewServer(handler)
	t.Cleanup(ts.Close)

	c := NewClient(nil)
	if err := c.SetBaseURL(ts.URL); err != nil {
		t.Fatalf("SetBaseURL: %v", err)
	}

	return c
}

func decodeBody(t *testing.T, r *http.Request) map[string]interface{} {
	t.Helper()
	var body map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		t.Fatalf("decode body: %v", err)
	}

	return body
}

func TestNewProcessID(t *testing.T) {
	seen := map[string]bool{}
	for i := 0; i < 100; i++ {
		id, err := newProcessID()
		if err != nil {
			t.Fatalf("newProcessID: %v", err)
		}
		if len(id) != 36 {
			t.Fatalf("expected 36-char UUID, got %q (%d)", id, len(id))
		}
		if id[14] != '4' {
			t.Errorf("expected version 4 nibble, got %q", id)
		}
		if seen[id] {
			t.Fatalf("duplicate process id %q", id)
		}
		seen[id] = true
	}
}

func TestDo_ErrorFallbackOnNonJSON(t *testing.T) {
	c := setup(t, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
		_, _ = w.Write([]byte("<html>502 Bad Gateway</html>"))
	})

	_, err := c.DataCenter.List(context.Background(), &ListOptions{})
	if err == nil {
		t.Fatal("expected an error for a non-JSON 502 response")
	}
	if want := "502"; !contains(err.Error(), want) {
		t.Errorf("expected error to mention status %s, got %q", want, err.Error())
	}
}

func TestDo_ErrorMessageFromEnvelope(t *testing.T) {
	c := setup(t, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error":true,"message":"Please authenticate"}`))
	})

	_, err := c.Tags.List(context.Background())
	if err == nil || err.Error() != "Please authenticate" {
		t.Fatalf("expected API message error, got %v", err)
	}
}

func TestDataCenterList_NilOptionsNoPanic(t *testing.T) {
	c := setup(t, func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"error":false,"data":[{"dc_name":"AMS","identifier":"ams1"}]}`))
	})

	// Passing nil options must not panic (regression: nil deref on options.Page).
	dcs, err := c.DataCenter.List(context.Background(), nil)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(dcs) != 1 || dcs[0].Identifier != "ams1" {
		t.Fatalf("unexpected datacenters: %+v", dcs)
	}
}

func contains(s, sub string) bool {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
