package govpsie

import (
	"context"
	"net/http"
	"testing"
)

func TestFirewallGroup_DeleteSendsDeleteStatistic(t *testing.T) {
	// Regression: the delete handler dereferences deleteStatistic.reason, so a
	// body of just {groupId} makes it fail with HTTP 500.
	var got map[string]interface{}
	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete || r.URL.Path != "/apps/v2/firewall/delete/group" {
			t.Errorf("got %s %s", r.Method, r.URL.Path)
		}
		got = decodeBody(t, r)
		_, _ = w.Write([]byte(`{"error":false,"data":true}`))
	})

	if err := c.FirewallGroup.Delete(context.Background(), "grp-1", "cleanup", "done"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if got["groupId"] != "grp-1" {
		t.Errorf("missing groupId: %+v", got)
	}
	stat, ok := got["deleteStatistic"].(map[string]interface{})
	if !ok || stat["reason"] != "cleanup" || stat["note"] != "done" {
		t.Errorf("unexpected deleteStatistic: %+v", got["deleteStatistic"])
	}
}

func TestFirewallGroup_CreateRuleBody(t *testing.T) {
	var got map[string]interface{}
	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/apps/v2/firewall/create/group" {
			t.Errorf("got %s %s", r.Method, r.URL.Path)
		}
		got = decodeBody(t, r)
		_, _ = w.Write([]byte(`{"error":false,"data":true}`))
	})

	rules := []FirewallUpdateReq{{
		Action: "ACCEPT", Type: "in", Proto: "tcp", Dport: "22",
		Source: []string{"0.0.0.0/0"}, Enable: 1, Comment: "ssh",
	}}
	if err := c.FirewallGroup.Create(context.Background(), "web", rules); err != nil {
		t.Fatalf("Create: %v", err)
	}
	if got["groupName"] != "web" {
		t.Errorf("missing groupName: %+v", got)
	}
	gotRules, ok := got["rules"].([]interface{})
	if !ok || len(gotRules) != 1 {
		t.Fatalf("unexpected rules: %+v", got["rules"])
	}
	rule := gotRules[0].(map[string]interface{})
	if rule["action"] != "ACCEPT" || rule["type"] != "in" || rule["dport"] != "22" {
		t.Errorf("unexpected rule: %+v", rule)
	}
}

func TestFirewallGroup_RenameGroup(t *testing.T) {
	var got map[string]interface{}
	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/apps/v2/firewall/group/rename" {
			t.Errorf("got %s %s", r.Method, r.URL.Path)
		}
		got = decodeBody(t, r)
		_, _ = w.Write([]byte(`{"error":false,"data":true}`))
	})

	if err := c.FirewallGroup.RenameGroup(context.Background(), "grp-1", "newname"); err != nil {
		t.Fatalf("RenameGroup: %v", err)
	}
	if got["groupId"] != "grp-1" || got["groupName"] != "newname" {
		t.Errorf("unexpected body: %+v", got)
	}
}

func TestFirewallGroup_UpdateAndDeleteRules(t *testing.T) {
	paths := map[string]map[string]interface{}{}
	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		paths[r.Method+" "+r.URL.Path] = decodeBody(t, r)
		_, _ = w.Write([]byte(`{"error":false,"data":true}`))
	})

	ctx := context.Background()
	if err := c.FirewallGroup.UpdateRules(ctx, "grp-1", []FirewallUpdateReq{{Action: "DROP", Type: "out", Enable: 1}}); err != nil {
		t.Fatalf("UpdateRules: %v", err)
	}
	if err := c.FirewallGroup.DeleteRules(ctx, "grp-1", []string{"rule-1", "rule-2"}); err != nil {
		t.Fatalf("DeleteRules: %v", err)
	}

	up, ok := paths["PUT /apps/v2/firewall/group/rules/update"]
	if !ok || up["groupId"] != "grp-1" {
		t.Errorf("update rules not sent correctly: %+v", up)
	}
	del, ok := paths["DELETE /apps/v2/firewall/group/rules/delete"]
	if !ok || del["groupId"] != "grp-1" {
		t.Errorf("delete rules not sent correctly: %+v", del)
	}
	ids, ok := del["rules"].([]interface{})
	if !ok || len(ids) != 2 || ids[0] != "rule-1" {
		t.Errorf("unexpected rule ids: %+v", del["rules"])
	}
}
