/*
Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testPoolAttachmentResourceData(t *testing.T, pool, node string) *schema.ResourceData {
	t.Helper()
	d := schema.TestResourceDataRaw(t, resourceBigipLtmPoolAttachment().Schema, map[string]interface{}{
		"pool":                  pool,
		"node":                  node,
		"ratio":                 1,
		"priority_group":        0,
		"connection_limit":      0,
		"connection_rate_limit": "0",
		"monitor":               "default",
		"state":                 "enabled",
		"dynamic_ratio":         1,
	})
	return d
}

// ---- SplitNodePort ----

func TestSplitNodePort(t *testing.T) {
	cases := []struct {
		in       string
		expected []string
	}{
		{"10.10.10.10:80", []string{"10.10.10.10", "80"}},
		{"/Common/node1:80", []string{"/Common/node1", "80"}},
		{"web-server1", []string{"web-server1"}},
		{"node.example.com:443", []string{"node.example.com", "443"}},
	}
	for _, c := range cases {
		assert.Equal(t, c.expected, SplitNodePort(c.in), "input=%s", c.in)
	}
}

// ---- Create ----

func TestResourceBigipLtmPoolAttachmentCreate_matchNonFQDNSuccess(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	poolName := "/Common/test-pool"
	nodeName := "/Common/node1:80"

	mux.HandleFunc("/mgmt/tm/ltm/node/~Common~node1", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		_, _ = fmt.Fprintf(w, `{"name":"node1","fullPath":"/Common/node1"}`)
	})
	mux.HandleFunc("/mgmt/tm/ltm/pool/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			_, _ = fmt.Fprintf(w, `{}`)
		case "PATCH":
			_, _ = fmt.Fprintf(w, `{}`)
		case "GET":
			if strings.HasSuffix(r.URL.Path, "/members") {
				_, _ = fmt.Fprintf(w, `{"items":[{"name":"node1:80","fullPath":"/Common/node1:80"}]}`)
				return
			}
			_, _ = fmt.Fprintf(w, `{"name":"test-pool"}`)
		default:
			t.Fatalf("unexpected method %s on %s", r.Method, r.URL.Path)
		}
	})

	client := newDatasourceTestClient(server.URL)
	d := testPoolAttachmentResourceData(t, poolName, nodeName)

	diags := resourceBigipLtmPoolAttachmentCreate(context.Background(), d, client)
	require.False(t, diags.HasError(), "%v", diags)
	assert.Equal(t, fmt.Sprintf("%s-%s", poolName, nodeName), d.Id())
}

func TestResourceBigipLtmPoolAttachmentCreate_matchFQDNSuccess(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	poolName := "/Common/test-pool"
	nodeName := "/Common/node1:80"

	mux.HandleFunc("/mgmt/tm/ltm/node/~Common~node1", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"name":"node1","fullPath":"/Common/node1","fqdn":{"tmName":"node1.example.com","interval":"3600","addressFamily":"ipv4","autopopulate":"enabled"}}`)
	})
	mux.HandleFunc("/mgmt/tm/ltm/pool/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			_, _ = fmt.Fprintf(w, `{}`)
		case "PATCH":
			_, _ = fmt.Fprintf(w, `{}`)
		case "GET":
			if strings.HasSuffix(r.URL.Path, "/members") {
				_, _ = fmt.Fprintf(w, `{"items":[{"name":"node1:80","fullPath":"/Common/node1:80"}]}`)
				return
			}
			_, _ = fmt.Fprintf(w, `{"name":"test-pool"}`)
		default:
			t.Fatalf("unexpected method %s on %s", r.Method, r.URL.Path)
		}
	})

	client := newDatasourceTestClient(server.URL)
	d := testPoolAttachmentResourceData(t, poolName, nodeName)

	diags := resourceBigipLtmPoolAttachmentCreate(context.Background(), d, client)
	require.False(t, diags.HasError(), "%v", diags)
	assert.Equal(t, fmt.Sprintf("%s-%s", poolName, nodeName), d.Id())
}

func TestResourceBigipLtmPoolAttachmentCreate_getNodeError(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	poolName := "/Common/test-pool"
	nodeName := "/Common/node1:80"

	mux.HandleFunc("/mgmt/tm/ltm/node/~Common~node1", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(w, `{"code":500,"message":"boom"}`)
	})

	client := newDatasourceTestClient(server.URL)
	d := testPoolAttachmentResourceData(t, poolName, nodeName)

	diags := resourceBigipLtmPoolAttachmentCreate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

// A 404 from GetNode means the node no longer exists on the device; Create
// clears the resource ID rather than returning an error.
func TestResourceBigipLtmPoolAttachmentCreate_getNodeNotFound(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	poolName := "/Common/test-pool"
	nodeName := "/Common/node1:80"

	mux.HandleFunc("/mgmt/tm/ltm/node/~Common~node1", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = fmt.Fprintf(w, `{"code":404,"message":"not found"}`)
	})

	client := newDatasourceTestClient(server.URL)
	d := testPoolAttachmentResourceData(t, poolName, nodeName)
	d.SetId("preexisting")

	diags := resourceBigipLtmPoolAttachmentCreate(context.Background(), d, client)
	require.False(t, diags.HasError(), "%v", diags)
	require.Empty(t, d.Id())
}

func TestResourceBigipLtmPoolAttachmentCreate_addPoolMemberFQDNError(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	poolName := "/Common/test-pool"
	nodeName := "/Common/node1:80"

	mux.HandleFunc("/mgmt/tm/ltm/node/~Common~node1", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"name":"node1","fullPath":"/Common/node1","fqdn":{"tmName":"node1.example.com"}}`)
	})
	mux.HandleFunc("/mgmt/tm/ltm/pool/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(w, `{"code":500,"message":"add failed"}`)
	})

	client := newDatasourceTestClient(server.URL)
	d := testPoolAttachmentResourceData(t, poolName, nodeName)

	diags := resourceBigipLtmPoolAttachmentCreate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestResourceBigipLtmPoolAttachmentCreate_addPoolMemberNodeError(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	poolName := "/Common/test-pool"
	nodeName := "/Common/node1:80"

	mux.HandleFunc("/mgmt/tm/ltm/node/~Common~node1", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"name":"node1","fullPath":"/Common/node1"}`)
	})
	mux.HandleFunc("/mgmt/tm/ltm/pool/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(w, `{"code":500,"message":"add failed"}`)
	})

	client := newDatasourceTestClient(server.URL)
	d := testPoolAttachmentResourceData(t, poolName, nodeName)

	diags := resourceBigipLtmPoolAttachmentCreate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestResourceBigipLtmPoolAttachmentCreate_noMatchIPSuccess(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	poolName := "/Common/test-pool"
	nodeName := "10.10.10.10:80"

	mux.HandleFunc("/mgmt/tm/ltm/pool/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			_, _ = fmt.Fprintf(w, `{}`)
		case "PATCH":
			_, _ = fmt.Fprintf(w, `{}`)
		case "GET":
			if strings.HasSuffix(r.URL.Path, "/members") {
				_, _ = fmt.Fprintf(w, `{"items":[{"name":"10.10.10.10:80","fullPath":"/Common/10.10.10.10:80"}]}`)
				return
			}
			_, _ = fmt.Fprintf(w, `{"name":"test-pool"}`)
		default:
			t.Fatalf("unexpected method %s on %s", r.Method, r.URL.Path)
		}
	})

	client := newDatasourceTestClient(server.URL)
	d := testPoolAttachmentResourceData(t, poolName, nodeName)
	d.SetId(poolName)

	diags := resourceBigipLtmPoolAttachmentCreate(context.Background(), d, client)
	require.False(t, diags.HasError(), "%v", diags)
	assert.Equal(t, poolName, d.Id())
}

func TestResourceBigipLtmPoolAttachmentCreate_noMatchFQDNAutopopulate(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	poolName := "/Common/test-pool"
	nodeName := "myhost.example.com:80"

	mux.HandleFunc("/mgmt/tm/ltm/pool/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			_, _ = fmt.Fprintf(w, `{}`)
		case "PATCH":
			_, _ = fmt.Fprintf(w, `{}`)
		case "GET":
			_, _ = fmt.Fprintf(w, `{"name":"test-pool"}`)
		default:
			t.Fatalf("unexpected method %s on %s", r.Method, r.URL.Path)
		}
	})

	client := newDatasourceTestClient(server.URL)
	d := testPoolAttachmentResourceData(t, poolName, nodeName)
	d.SetId(poolName)
	_ = d.Set("fqdn_autopopulate", "disabled")

	diags := resourceBigipLtmPoolAttachmentCreate(context.Background(), d, client)
	require.False(t, diags.HasError(), "%v", diags)
}

func TestResourceBigipLtmPoolAttachmentCreate_noMatchAddPoolMemberError(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	poolName := "/Common/test-pool"
	nodeName := "10.10.10.10:80"

	mux.HandleFunc("/mgmt/tm/ltm/pool/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(w, `{"code":500,"message":"add failed"}`)
	})

	client := newDatasourceTestClient(server.URL)
	d := testPoolAttachmentResourceData(t, poolName, nodeName)

	diags := resourceBigipLtmPoolAttachmentCreate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

// ---- Update ----

func TestResourceBigipLtmPoolAttachmentUpdate_matchSuccess(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	poolName := "/Common/test-pool"
	nodeName := "/Common/node1:80"

	mux.HandleFunc("/mgmt/tm/ltm/node/~Common~node1", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"name":"node1","fullPath":"/Common/node1"}`)
	})
	mux.HandleFunc("/mgmt/tm/ltm/pool/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "PATCH":
			_, _ = fmt.Fprintf(w, `{}`)
		case "GET":
			_, _ = fmt.Fprintf(w, `{"name":"test-pool"}`)
		default:
			t.Fatalf("unexpected method %s on %s", r.Method, r.URL.Path)
		}
	})

	client := newDatasourceTestClient(server.URL)
	d := testPoolAttachmentResourceData(t, poolName, nodeName)
	d.SetId(fmt.Sprintf("%s-%s", poolName, nodeName))
	_ = d.Set("state", "disabled")

	diags := resourceBigipLtmPoolAttachmentUpdate(context.Background(), d, client)
	require.False(t, diags.HasError(), "%v", diags)
}

func TestResourceBigipLtmPoolAttachmentUpdate_matchFQDNForcedOffline(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	poolName := "/Common/test-pool"
	nodeName := "/Common/node1:80"

	mux.HandleFunc("/mgmt/tm/ltm/node/~Common~node1", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"name":"node1","fullPath":"/Common/node1","fqdn":{"tmName":"node1.example.com"}}`)
	})
	mux.HandleFunc("/mgmt/tm/ltm/pool/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "PATCH":
			_, _ = fmt.Fprintf(w, `{}`)
		case "GET":
			_, _ = fmt.Fprintf(w, `{"name":"test-pool"}`)
		default:
			t.Fatalf("unexpected method %s on %s", r.Method, r.URL.Path)
		}
	})

	client := newDatasourceTestClient(server.URL)
	d := testPoolAttachmentResourceData(t, poolName, nodeName)
	d.SetId(fmt.Sprintf("%s-%s", poolName, nodeName))
	_ = d.Set("state", "forced_offline")

	diags := resourceBigipLtmPoolAttachmentUpdate(context.Background(), d, client)
	require.False(t, diags.HasError(), "%v", diags)
}

func TestResourceBigipLtmPoolAttachmentUpdate_matchGetNodeError(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	poolName := "/Common/test-pool"
	nodeName := "/Common/node1:80"

	mux.HandleFunc("/mgmt/tm/ltm/node/~Common~node1", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(w, `{"code":500,"message":"boom"}`)
	})

	client := newDatasourceTestClient(server.URL)
	d := testPoolAttachmentResourceData(t, poolName, nodeName)
	d.SetId(fmt.Sprintf("%s-%s", poolName, nodeName))

	diags := resourceBigipLtmPoolAttachmentUpdate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

// A 404 from GetNode means the node no longer exists on the device; Update
// clears the resource ID rather than returning an error.
func TestResourceBigipLtmPoolAttachmentUpdate_matchGetNodeNotFound(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	poolName := "/Common/test-pool"
	nodeName := "/Common/node1:80"

	mux.HandleFunc("/mgmt/tm/ltm/node/~Common~node1", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = fmt.Fprintf(w, `{"code":404,"message":"not found"}`)
	})

	client := newDatasourceTestClient(server.URL)
	d := testPoolAttachmentResourceData(t, poolName, nodeName)
	d.SetId(fmt.Sprintf("%s-%s", poolName, nodeName))

	diags := resourceBigipLtmPoolAttachmentUpdate(context.Background(), d, client)
	require.False(t, diags.HasError(), "%v", diags)
	require.Empty(t, d.Id())
}

func TestResourceBigipLtmPoolAttachmentUpdate_matchModifyPoolMemberError(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	poolName := "/Common/test-pool"
	nodeName := "/Common/node1:80"

	mux.HandleFunc("/mgmt/tm/ltm/node/~Common~node1", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"name":"node1","fullPath":"/Common/node1"}`)
	})
	mux.HandleFunc("/mgmt/tm/ltm/pool/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(w, `{"code":500,"message":"modify failed"}`)
	})

	client := newDatasourceTestClient(server.URL)
	d := testPoolAttachmentResourceData(t, poolName, nodeName)
	d.SetId(fmt.Sprintf("%s-%s", poolName, nodeName))

	diags := resourceBigipLtmPoolAttachmentUpdate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestResourceBigipLtmPoolAttachmentUpdate_noMatchSuccess(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	poolName := "/Common/test-pool"
	nodeName := "10.10.10.10:80"

	mux.HandleFunc("/mgmt/tm/ltm/pool/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "PATCH":
			_, _ = fmt.Fprintf(w, `{}`)
		case "GET":
			_, _ = fmt.Fprintf(w, `{"name":"test-pool"}`)
		default:
			t.Fatalf("unexpected method %s on %s", r.Method, r.URL.Path)
		}
	})

	client := newDatasourceTestClient(server.URL)
	d := testPoolAttachmentResourceData(t, poolName, nodeName)
	d.SetId(poolName)
	_ = d.Set("state", "disabled")

	diags := resourceBigipLtmPoolAttachmentUpdate(context.Background(), d, client)
	require.False(t, diags.HasError(), "%v", diags)
}

func TestResourceBigipLtmPoolAttachmentUpdate_noMatchFQDNAutopopulate(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	poolName := "/Common/test-pool"
	nodeName := "myhost.example.com:80"

	mux.HandleFunc("/mgmt/tm/ltm/pool/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "PATCH":
			_, _ = fmt.Fprintf(w, `{}`)
		case "GET":
			_, _ = fmt.Fprintf(w, `{"name":"test-pool"}`)
		default:
			t.Fatalf("unexpected method %s on %s", r.Method, r.URL.Path)
		}
	})

	client := newDatasourceTestClient(server.URL)
	d := testPoolAttachmentResourceData(t, poolName, nodeName)
	d.SetId(poolName)
	_ = d.Set("state", "forced_offline")

	diags := resourceBigipLtmPoolAttachmentUpdate(context.Background(), d, client)
	require.False(t, diags.HasError(), "%v", diags)
}

func TestResourceBigipLtmPoolAttachmentUpdate_noMatchModifyError(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	poolName := "/Common/test-pool"
	nodeName := "10.10.10.10:80"

	mux.HandleFunc("/mgmt/tm/ltm/pool/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(w, `{"code":500,"message":"modify failed"}`)
	})

	client := newDatasourceTestClient(server.URL)
	d := testPoolAttachmentResourceData(t, poolName, nodeName)
	d.SetId(poolName)

	diags := resourceBigipLtmPoolAttachmentUpdate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

// ---- Read ----

func TestResourceBigipLtmPoolAttachmentRead_matchFound(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	poolName := "/Common/test-pool"
	nodeName := "/Common/node1:80"

	mux.HandleFunc("/mgmt/tm/ltm/pool/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/mgmt/tm/ltm/pool/~Common~test-pool/members" {
			_, _ = fmt.Fprintf(w, `{"items":[{"name":"node1:80","fullPath":"/Common/node1:80","priorityGroup":0,"ratio":1,"connectionLimit":0,"rateLimit":"disabled","dynamicRatio":1,"monitor":"default"}]}`)
			return
		}
		_, _ = fmt.Fprintf(w, `{"name":"test-pool"}`)
	})

	client := newDatasourceTestClient(server.URL)
	d := testPoolAttachmentResourceData(t, poolName, nodeName)
	d.SetId(fmt.Sprintf("%s-%s", poolName, nodeName))

	diags := resourceBigipLtmPoolAttachmentRead(context.Background(), d, client)
	require.False(t, diags.HasError(), "%v", diags)
	assert.NotEqual(t, "", d.Id())
}

func TestResourceBigipLtmPoolAttachmentRead_matchNotFound(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	poolName := "/Common/test-pool"
	nodeName := "/Common/node1:80"

	mux.HandleFunc("/mgmt/tm/ltm/pool/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/mgmt/tm/ltm/pool/~Common~test-pool/members" {
			_, _ = fmt.Fprintf(w, `{"items":[{"name":"node2:80","fullPath":"/Common/node2:80"}]}`)
			return
		}
		_, _ = fmt.Fprintf(w, `{"name":"test-pool"}`)
	})

	client := newDatasourceTestClient(server.URL)
	d := testPoolAttachmentResourceData(t, poolName, nodeName)
	d.SetId(fmt.Sprintf("%s-%s", poolName, nodeName))

	diags := resourceBigipLtmPoolAttachmentRead(context.Background(), d, client)
	require.False(t, diags.HasError(), "%v", diags)
	assert.Equal(t, "", d.Id())
}

func TestResourceBigipLtmPoolAttachmentRead_noMatchFound(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	poolName := "/Common/test-pool"
	nodeName := "10.10.10.10:80"

	mux.HandleFunc("/mgmt/tm/ltm/pool/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/mgmt/tm/ltm/pool/~Common~test-pool/members" {
			_, _ = fmt.Fprintf(w, `{"items":[{"name":"10.10.10.10:80","fullPath":"/Common/10.10.10.10:80"}]}`)
			return
		}
		_, _ = fmt.Fprintf(w, `{"name":"test-pool"}`)
	})

	client := newDatasourceTestClient(server.URL)
	d := testPoolAttachmentResourceData(t, poolName, nodeName)
	d.SetId(poolName)

	diags := resourceBigipLtmPoolAttachmentRead(context.Background(), d, client)
	require.False(t, diags.HasError(), "%v", diags)
}

func TestResourceBigipLtmPoolAttachmentRead_getPoolError(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	poolName := "/Common/test-pool"
	nodeName := "/Common/node1:80"

	mux.HandleFunc("/mgmt/tm/ltm/pool/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(w, `{"code":500,"message":"boom"}`)
	})

	client := newDatasourceTestClient(server.URL)
	d := testPoolAttachmentResourceData(t, poolName, nodeName)
	d.SetId(fmt.Sprintf("%s-%s", poolName, nodeName))

	diags := resourceBigipLtmPoolAttachmentRead(context.Background(), d, client)
	require.True(t, diags.HasError())
}

// A 404 from GetPool means the pool no longer exists on the device; Read
// clears the resource ID rather than returning an error.
func TestResourceBigipLtmPoolAttachmentRead_poolNotFound(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	poolName := "/Common/test-pool"
	nodeName := "/Common/node1:80"

	mux.HandleFunc("/mgmt/tm/ltm/pool/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = fmt.Fprintf(w, `{"code":404,"message":"not found"}`)
	})

	client := newDatasourceTestClient(server.URL)
	d := testPoolAttachmentResourceData(t, poolName, nodeName)
	d.SetId(fmt.Sprintf("%s-%s", poolName, nodeName))

	diags := resourceBigipLtmPoolAttachmentRead(context.Background(), d, client)
	require.False(t, diags.HasError(), "%v", diags)
	require.Empty(t, d.Id())
}

func TestResourceBigipLtmPoolAttachmentRead_poolMembersError(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	poolName := "/Common/test-pool"
	nodeName := "/Common/node1:80"

	mux.HandleFunc("/mgmt/tm/ltm/pool/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/mgmt/tm/ltm/pool/~Common~test-pool/members" {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprintf(w, `{"code":500,"message":"boom"}`)
			return
		}
		_, _ = fmt.Fprintf(w, `{"name":"test-pool"}`)
	})

	client := newDatasourceTestClient(server.URL)
	d := testPoolAttachmentResourceData(t, poolName, nodeName)
	d.SetId(fmt.Sprintf("%s-%s", poolName, nodeName))

	diags := resourceBigipLtmPoolAttachmentRead(context.Background(), d, client)
	require.True(t, diags.HasError())
}

// ---- Delete ----

func TestResourceBigipLtmPoolAttachmentDelete_success(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	poolName := "/Common/test-pool"
	nodeName := "/Common/node1:80"

	mux.HandleFunc("/mgmt/tm/ltm/pool/", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "DELETE", r.Method)
		w.WriteHeader(http.StatusOK)
	})

	client := newDatasourceTestClient(server.URL)
	d := testPoolAttachmentResourceData(t, poolName, nodeName)
	d.SetId(fmt.Sprintf("%s-%s", poolName, nodeName))

	diags := resourceBigipLtmPoolAttachmentDelete(context.Background(), d, client)
	require.False(t, diags.HasError(), "%v", diags)
	assert.Equal(t, "", d.Id())
}

func TestResourceBigipLtmPoolAttachmentDelete_error(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	poolName := "/Common/test-pool"
	nodeName := "/Common/node1:80"

	mux.HandleFunc("/mgmt/tm/ltm/pool/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(w, `{"code":500,"message":"delete failed"}`)
	})

	client := newDatasourceTestClient(server.URL)
	d := testPoolAttachmentResourceData(t, poolName, nodeName)
	d.SetId(fmt.Sprintf("%s-%s", poolName, nodeName))

	diags := resourceBigipLtmPoolAttachmentDelete(context.Background(), d, client)
	require.True(t, diags.HasError())
}

// ---- Import ----

func TestResourceBigipLtmPoolAttachmentImport_success(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	poolName := "/Common/test-pool"
	nodeName := "/Common/node1:80"

	mux.HandleFunc("/mgmt/tm/ltm/pool/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/mgmt/tm/ltm/pool/~Common~test-pool/members" {
			_, _ = fmt.Fprintf(w, `{"items":[{"name":"node1:80","fullPath":"/Common/node1:80","priorityGroup":0,"ratio":1,"connectionLimit":0,"rateLimit":"disabled","dynamicRatio":1,"monitor":"default"}]}`)
			return
		}
		_, _ = fmt.Fprintf(w, `{"name":"test-pool"}`)
	})

	client := newDatasourceTestClient(server.URL)
	d := testPoolAttachmentResourceData(t, "", "")
	d.SetId(fmt.Sprintf(`{"pool":"%s","node":"%s"}`, poolName, nodeName))

	results, err := resourceBigipLtmPoolAttachmentImport(context.Background(), d, client)
	require.NoError(t, err)
	require.Len(t, results, 1)
	assert.Equal(t, fmt.Sprintf("%s-%s", poolName, nodeName), results[0].Id())
}

func TestResourceBigipLtmPoolAttachmentImport_badJSON(t *testing.T) {
	client := newDatasourceTestClient("http://127.0.0.1:0")
	d := testPoolAttachmentResourceData(t, "", "")
	d.SetId("not-json")

	_, err := resourceBigipLtmPoolAttachmentImport(context.Background(), d, client)
	require.Error(t, err)
}

func TestResourceBigipLtmPoolAttachmentImport_missingPool(t *testing.T) {
	client := newDatasourceTestClient("http://127.0.0.1:0")
	d := testPoolAttachmentResourceData(t, "", "")
	d.SetId(`{"node":"/Common/node1:80"}`)

	_, err := resourceBigipLtmPoolAttachmentImport(context.Background(), d, client)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "missing pool name")
}

func TestResourceBigipLtmPoolAttachmentImport_missingNode(t *testing.T) {
	client := newDatasourceTestClient("http://127.0.0.1:0")
	d := testPoolAttachmentResourceData(t, "", "")
	d.SetId(`{"pool":"/Common/test-pool"}`)

	_, err := resourceBigipLtmPoolAttachmentImport(context.Background(), d, client)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "missing node name")
}

func TestResourceBigipLtmPoolAttachmentImport_getPoolError(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	mux.HandleFunc("/mgmt/tm/ltm/pool/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(w, `{"code":500,"message":"boom"}`)
	})

	client := newDatasourceTestClient(server.URL)
	d := testPoolAttachmentResourceData(t, "", "")
	d.SetId(`{"pool":"/Common/test-pool","node":"/Common/node1:80"}`)

	_, err := resourceBigipLtmPoolAttachmentImport(context.Background(), d, client)
	require.Error(t, err)
}

// A 404 from GetPool means the pool does not exist on the device; Import
// surfaces this as an "unable to find the pool" error.
func TestResourceBigipLtmPoolAttachmentImport_poolNotFound(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	mux.HandleFunc("/mgmt/tm/ltm/pool/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = fmt.Fprintf(w, `{"code":404,"message":"not found"}`)
	})

	client := newDatasourceTestClient(server.URL)
	d := testPoolAttachmentResourceData(t, "", "")
	d.SetId(`{"pool":"/Common/test-pool","node":"/Common/node1:80"}`)

	_, err := resourceBigipLtmPoolAttachmentImport(context.Background(), d, client)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unable to find the pool")
}

func TestResourceBigipLtmPoolAttachmentImport_poolMembersError(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	mux.HandleFunc("/mgmt/tm/ltm/pool/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/mgmt/tm/ltm/pool/~Common~test-pool/members" {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprintf(w, `{"code":500,"message":"boom"}`)
			return
		}
		_, _ = fmt.Fprintf(w, `{"name":"test-pool"}`)
	})

	client := newDatasourceTestClient(server.URL)
	d := testPoolAttachmentResourceData(t, "", "")
	d.SetId(`{"pool":"/Common/test-pool","node":"/Common/node1:80"}`)

	_, err := resourceBigipLtmPoolAttachmentImport(context.Background(), d, client)
	require.Error(t, err)
}

func TestResourceBigipLtmPoolAttachmentImport_notFoundInMembers(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	mux.HandleFunc("/mgmt/tm/ltm/pool/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/mgmt/tm/ltm/pool/~Common~test-pool/members" {
			_, _ = fmt.Fprintf(w, `{"items":[{"name":"other:80","fullPath":"/Common/other:80"}]}`)
			return
		}
		_, _ = fmt.Fprintf(w, `{"name":"test-pool"}`)
	})

	client := newDatasourceTestClient(server.URL)
	d := testPoolAttachmentResourceData(t, "", "")
	d.SetId(`{"pool":"/Common/test-pool","node":"/Common/node1:80"}`)

	_, err := resourceBigipLtmPoolAttachmentImport(context.Background(), d, client)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "cannot locate node")
}
