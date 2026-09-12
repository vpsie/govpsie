package govpsie

import (
	"context"
	"net/http"
	"testing"
)

func TestTags_AddToEntity(t *testing.T) {
	var got map[string]interface{}
	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/apps/v2/tags/add" {
			t.Errorf("got %s %s", r.Method, r.URL.Path)
		}
		got = decodeBody(t, r)
		_, _ = w.Write([]byte(`{"error":false,"data":true}`))
	})

	if err := c.Tags.AddToEntity(context.Background(), "ssh_keys", "id-1", []string{"a", "b"}); err != nil {
		t.Fatalf("AddToEntity: %v", err)
	}
	if got["entity"] != "ssh_keys" || got["identifier"] != "id-1" {
		t.Errorf("unexpected body: %+v", got)
	}
	tags, ok := got["tags"].([]interface{})
	if !ok || len(tags) != 2 || tags[0] != "a" {
		t.Errorf("unexpected tags: %+v", got["tags"])
	}
}

func TestTags_ListForEntity(t *testing.T) {
	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("entity") != "ssh_keys" || r.URL.Query().Get("identifier") != "id-1" {
			t.Errorf("unexpected query: %s", r.URL.RawQuery)
		}
		// data is an object grouped by entity type; only the target group is populated.
		_, _ = w.Write([]byte(`{"error":false,"data":{"boxes":[],"ssh_keys":[{"name":"a","entity_id":1,"entity_type":"ssh_keys"},{"name":"b","entity_id":1,"entity_type":"ssh_keys"}]}}`))
	})

	names, err := c.Tags.ListForEntity(context.Background(), "ssh_keys", "id-1")
	if err != nil {
		t.Fatalf("ListForEntity: %v", err)
	}
	if len(names) != 2 || names[0] != "a" || names[1] != "b" {
		t.Fatalf("unexpected names: %+v", names)
	}
}

func TestTags_Delete(t *testing.T) {
	var got map[string]interface{}
	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete || r.URL.Path != "/apps/v2/tags/delete" {
			t.Errorf("got %s %s", r.Method, r.URL.Path)
		}
		got = decodeBody(t, r)
		_, _ = w.Write([]byte(`{"error":false,"data":true}`))
	})

	if err := c.Tags.Delete(context.Background(), "ssh_keys", "id-1"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if got["entity"] != "ssh_keys" || got["identifier"] != "id-1" {
		t.Errorf("unexpected body: %+v", got)
	}
}

func TestCertificate_AddAndList(t *testing.T) {
	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/apps/v2/certificates/add":
			got := decodeBody(t, r)
			if got["certName"] != "c1" || got["domainId"] != "d1" {
				t.Errorf("unexpected add body: %+v", got)
			}
			_, _ = w.Write([]byte(`{"error":false,"data":{}}`))
		case "/apps/v2/certificates/all":
			_, _ = w.Write([]byte(`{"error":false,"data":[{"identifier":"cert-1","certificateName":"c1","issuer":"LE"}],"total":1}`))
		default:
			t.Errorf("unexpected path %s", r.URL.Path)
		}
	})

	if err := c.Certificate.Add(context.Background(), "c1", "d1"); err != nil {
		t.Fatalf("Add: %v", err)
	}
	certs, err := c.Certificate.List(context.Background())
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(certs) != 1 || certs[0].Identifier != "cert-1" || certs[0].CertificateName != "c1" {
		t.Fatalf("unexpected certs: %+v", certs)
	}
}

func TestServerGroup_CreateJoinDetach(t *testing.T) {
	paths := map[string]bool{}
	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		paths[r.Method+" "+r.URL.Path] = true
		switch r.URL.Path {
		case "/apps/v2/vm/group/new", "/apps/v2/vm/group/join", "/apps/v2/vm/group/detach":
			_, _ = w.Write([]byte(`{"error":false,"data":true}`))
		case "/apps/v2/vm/group":
			_, _ = w.Write([]byte(`{"error":false,"data":{"rows":[{"identifier":"g1","group_name":"web","is_distributed":true}],"total":1}}`))
		default:
			t.Errorf("unexpected path %s", r.URL.Path)
		}
	})

	ctx := context.Background()
	if err := c.ServerGroup.Create(ctx, "web", "desc", true); err != nil {
		t.Fatalf("Create: %v", err)
	}
	groups, err := c.ServerGroup.List(ctx)
	if err != nil || len(groups) != 1 || groups[0].Identifier != "g1" {
		t.Fatalf("List: %v %+v", err, groups)
	}
	if err := c.ServerGroup.Join(ctx, "g1", "vm1"); err != nil {
		t.Fatalf("Join: %v", err)
	}
	if err := c.ServerGroup.Detach(ctx, "g1", "vm1"); err != nil {
		t.Fatalf("Detach: %v", err)
	}
	for _, want := range []string{"POST /apps/v2/vm/group/new", "GET /apps/v2/vm/group", "POST /apps/v2/vm/group/join", "POST /apps/v2/vm/group/detach"} {
		if !paths[want] {
			t.Errorf("expected request %q", want)
		}
	}
}

func TestRegistry_CreateAndCheckName(t *testing.T) {
	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/apps/v2/registries":
			if r.Method != http.MethodPost {
				t.Errorf("expected POST, got %s", r.Method)
			}
			got := decodeBody(t, r)
			if got["name"] != "reg1" || got["dcIdentifier"] != "ams1" {
				t.Errorf("unexpected body: %+v", got)
			}
			_, _ = w.Write([]byte(`{"error":false,"data":{}}`))
		case "/apps/v2/registries/check":
			_, _ = w.Write([]byte(`{"error":false,"data":true}`))
		default:
			t.Errorf("unexpected path %s", r.URL.Path)
		}
	})

	ctx := context.Background()
	if err := c.Registry.Create(ctx, "reg1", "ams1", "plan1"); err != nil {
		t.Fatalf("Create: %v", err)
	}
	taken, err := c.Registry.CheckName(ctx, "reg1")
	if err != nil || !taken {
		t.Fatalf("CheckName: %v taken=%v", err, taken)
	}
}

func TestManagedDB_DeleteSendsProcessID(t *testing.T) {
	var got map[string]interface{}
	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete || r.URL.Path != "/apps/v2/managed/db/cluster/byId/mdb-1" {
			t.Errorf("got %s %s", r.Method, r.URL.Path)
		}
		got = decodeBody(t, r)
		_, _ = w.Write([]byte(`{"error":false,"data":true}`))
	})

	if err := c.ManagedDB.Delete(context.Background(), "mdb-1", "no longer needed"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if got["mdbIdentifier"] != "mdb-1" {
		t.Errorf("missing mdbIdentifier: %+v", got)
	}
	if pid, _ := got["processId"].(string); pid == "" {
		t.Errorf("expected a generated processId, got %+v", got["processId"])
	}
	stat, ok := got["deleteStatistic"].(map[string]interface{})
	if !ok || stat["reason"] != "no longer needed" {
		t.Errorf("unexpected deleteStatistic: %+v", got["deleteStatistic"])
	}
}

func TestDomain_DeleteDnsRecordPath(t *testing.T) {
	// Regression: the delete must hit DELETE /domain/dnsRecord, not the
	// non-existent /domain/dnsRecord/delete (which returned "Not found").
	var got map[string]interface{}
	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete || r.URL.Path != "/apps/v2/domain/dnsRecord" {
			t.Errorf("got %s %s, want DELETE /apps/v2/domain/dnsRecord", r.Method, r.URL.Path)
		}
		got = decodeBody(t, r)
		_, _ = w.Write([]byte(`{"error":false,"data":true}`))
	})

	rec := &Record{Name: "www", Content: "1.2.3.4", Type: "A", TTL: 3600}
	if err := c.Domain.DeleteDnsRecord(context.Background(), "dom-1", rec); err != nil {
		t.Fatalf("DeleteDnsRecord: %v", err)
	}
	if got["domainIdentifier"] != "dom-1" {
		t.Errorf("missing domainIdentifier: %+v", got)
	}
	record, ok := got["record"].(map[string]interface{})
	if !ok || record["name"] != "www" || record["content"] != "1.2.3.4" {
		t.Errorf("unexpected record body: %+v", got["record"])
	}
}

func TestDomain_GetByIdentifierGroupsRecords(t *testing.T) {
	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/apps/v2/domain/dom-1" {
			t.Errorf("got %s %s", r.Method, r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"error":false,"data":{"domainIdentifier":"dom-1","domain":"ex.com","records":{"A":[{"name":"www.ex.com","content":"1.2.3.4","type":"A","ttl":3600}],"AAAA":[],"CNAME":[],"MX":[],"TXT":[],"CAA":[],"SRV":[],"NS":[]}}}`))
	})

	dom, err := c.Domain.GetDomainByIdentifier(context.Background(), "dom-1")
	if err != nil {
		t.Fatalf("GetDomainByIdentifier: %v", err)
	}
	if dom.Domain != "ex.com" {
		t.Errorf("unexpected domain: %q", dom.Domain)
	}
	aRecords := dom.RecordsByType("A")
	if len(aRecords) != 1 || aRecords[0].Name != "www.ex.com" || aRecords[0].TTL != 3600 {
		t.Errorf("unexpected A records: %+v", aRecords)
	}
	if len(dom.RecordsByType("MX")) != 0 {
		t.Errorf("expected no MX records")
	}
}

func TestBackupPolicy_GetSingularPathAndVms(t *testing.T) {
	// Regression: GET must use the singular /backup/policy/:id route (the plural
	// /backups/policy/:id returned "Not found"), and vms are objects.
	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/apps/v2/backup/policy/bp-1" {
			t.Errorf("got %s, want /apps/v2/backup/policy/bp-1", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"error":false,"data":{"name":"nightly","identifier":"bp-1","backupPlan":"day","planEvery":1,"keep":3,"vms":[{"name":"web","identifier":"vm-1","type":"vm"}]}}`))
	})

	policy, err := c.Backup.GetBackupPolicy(context.Background(), "bp-1")
	if err != nil {
		t.Fatalf("GetBackupPolicy: %v", err)
	}
	if policy == nil || policy.Name != "nightly" || len(policy.Vms) != 1 || policy.Vms[0].Identifier != "vm-1" {
		t.Fatalf("unexpected policy: %+v", policy)
	}
}

func TestBackupPolicy_GetNotFound(t *testing.T) {
	c := setup(t, func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"error":false,"data":false}`))
	})
	policy, err := c.Backup.GetBackupPolicy(context.Background(), "gone")
	if err != nil {
		t.Fatalf("GetBackupPolicy: %v", err)
	}
	if policy != nil {
		t.Fatalf("expected nil policy for data:false, got %+v", policy)
	}
}

func TestSnapshotPolicy_GetNotFound(t *testing.T) {
	c := setup(t, func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"error":false,"data":false}`))
	})
	policy, err := c.Snapshot.GetSnapShotPolicy(context.Background(), "gone")
	if err != nil {
		t.Fatalf("GetSnapShotPolicy: %v", err)
	}
	if policy != nil {
		t.Fatalf("expected nil policy for data:false, got %+v", policy)
	}
}

func TestManagedDB_AddNodeSendsRequiredFields(t *testing.T) {
	var got map[string]interface{}
	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/apps/v2/managed/db/cluster/byId/mdb-1/add" {
			t.Errorf("unexpected path %s", r.URL.Path)
		}
		got = decodeBody(t, r)
		_, _ = w.Write([]byte(`{"error":false,"data":true}`))
	})

	if err := c.ManagedDB.AddNode(context.Background(), "mdb-1"); err != nil {
		t.Fatalf("AddNode: %v", err)
	}
	if _, ok := got["createFromPool"]; !ok {
		t.Errorf("expected createFromPool in body: %+v", got)
	}
	if pid, _ := got["processId"].(string); pid == "" {
		t.Errorf("expected processId in body: %+v", got)
	}
}
