package govpsie

import (
	"context"
	"net/http"
	"testing"
)

func TestAccessTokenCreateReturnsServerMintedToken(t *testing.T) {
	var body map[string]interface{}
	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/apps/v2/profile/security/access/token" {
			t.Errorf("unexpected request %s %s", r.Method, r.URL.Path)
		}
		body = decodeBody(t, r)
		_, _ = w.Write([]byte(`{"error":false,"data":{"accessToken":"minted-by-server"}}`))
	})

	token, err := c.AccessToken.Create(context.Background(), "ci-token", "2030-01-01T00:00:00Z", "allowed")
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if token != "minted-by-server" {
		t.Errorf("token = %q", token)
	}
	if _, sent := body["accessToken"]; sent {
		t.Error("the token value must not be sent; the API generates it")
	}
	if body["accessTokenName"] != "ci-token" || body["expirationDate"] != "2030-01-01T00:00:00Z" || body["status"] != "allowed" {
		t.Errorf("body = %v", body)
	}
}

func TestAccessTokenCreateWithoutValueInResponse(t *testing.T) {
	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"error":false,"data":{}}`))
	})

	if _, err := c.AccessToken.Create(context.Background(), "ci-token", "2030-01-01T00:00:00Z", "allowed"); err == nil {
		t.Fatal("expected an error when the response carries no token")
	}
}

func TestCreateScriptSendsTagsArray(t *testing.T) {
	var body map[string]interface{}
	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		body = decodeBody(t, r)
		_, _ = w.Write([]byte(`{"error":false,"data":true}`))
	})

	req := &CreateScriptRequest{Name: "boot", ScriptContent: "echo hi", ScriptType: "bash"}
	if err := c.Scripts.CreateScript(context.Background(), req); err != nil {
		t.Fatalf("CreateScript: %v", err)
	}

	tags, ok := body["tags"].([]interface{})
	if !ok || len(tags) != 0 {
		t.Errorf("tags = %#v, want an empty array", body["tags"])
	}
	if req.Tags != nil {
		t.Error("CreateScript must not modify the caller's request")
	}
}
