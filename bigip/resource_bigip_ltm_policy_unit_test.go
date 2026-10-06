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

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------
// Test helpers
// ---------------------------------------------------------------------

// policyResourceData builds a *schema.ResourceData for the bigip_ltm_policy
// resource using schema.TestResourceDataRaw. Suitable for Read/Update/Delete
// tests that don't exercise dataToPolicy's GetRawConfig-based rule parsing.
func policyResourceData(t *testing.T, raw map[string]interface{}) *schema.ResourceData {
	t.Helper()
	return newDatasourceTestResourceData(t, resourceBigipLtmPolicy(), raw)
}

// policyResourceDataWithRawConfig builds a *schema.ResourceData whose
// GetRawConfig() returns ruleRawConfig. This is needed because
// schema.TestResourceDataRaw does not populate a usable raw config value,
// and dataToPolicy() reads rules exclusively via d.GetRawConfig() (to avoid
// stale-state contamination on TypeList index shifts). Thin wrapper around
// the shared NewTestResourceDataWithRawConfig helper
// (testing_helpers_test.go), which also backs
// clientSslResourceDataWithRawConfig in
// resource_bigip_ltm_profile_ssl_client_unit_test.go.
func policyResourceDataWithRawConfig(t *testing.T, raw map[string]interface{}, ruleRawConfig cty.Value) *schema.ResourceData {
	t.Helper()
	return NewTestResourceDataWithRawConfig(t, resourceBigipLtmPolicy(), raw, ruleRawConfig)
}

// noRulesRawConfig returns a raw-config cty.Value with a null "rule"
// attribute, exercising the "no rules configured" branch of dataToPolicy.
func noRulesRawConfig() cty.Value {
	return cty.ObjectVal(map[string]cty.Value{
		"rule": cty.NullVal(cty.List(cty.EmptyObject)),
	})
}

// ---------------------------------------------------------------------
// Schema-shape assertions
// ---------------------------------------------------------------------

func TestResourceBigipLtmPolicySchema(t *testing.T) {
	r := resourceBigipLtmPolicy()

	require.NotNil(t, r.Schema)
	require.NotNil(t, r.CreateContext)
	require.NotNil(t, r.ReadContext)
	require.NotNil(t, r.UpdateContext)
	require.NotNil(t, r.DeleteContext)
	require.NotNil(t, r.Importer)

	nameSchema, ok := r.Schema["name"]
	require.True(t, ok)
	assert.True(t, nameSchema.Required)
	assert.True(t, nameSchema.ForceNew)

	strategySchema, ok := r.Schema["strategy"]
	require.True(t, ok)
	assert.Equal(t, "/Common/first-match", strategySchema.Default)

	ruleSchema, ok := r.Schema["rule"]
	require.True(t, ok)
	assert.Equal(t, schema.TypeList, ruleSchema.Type)
}

// ---------------------------------------------------------------------
// dataToPolicy
// ---------------------------------------------------------------------

func TestDataToPolicyCommonPartitionNoRules(t *testing.T) {
	raw := map[string]interface{}{
		"name":        "/Common/mypolicy",
		"strategy":    "/Common/first-match",
		"description": "a policy",
		"controls":    []interface{}{"forwarding"},
		"requires":    []interface{}{"http"},
	}
	d := policyResourceDataWithRawConfig(t, raw, noRulesRawConfig())

	p := dataToPolicy("/Common/mypolicy", d)
	assert.Equal(t, "Drafts/mypolicy", p.Name)
	assert.Equal(t, "/Common/first-match", p.Strategy)
	assert.Equal(t, "a policy", p.Description)
	assert.ElementsMatch(t, []string{"forwarding"}, p.Controls)
	assert.ElementsMatch(t, []string{"http"}, p.Requires)
	assert.Empty(t, p.Rules)
}

func TestDataToPolicyNonCommonPartitionNoRules(t *testing.T) {
	raw := map[string]interface{}{
		"name":     "/MyPartition/mypolicy",
		"strategy": "/Common/all-match",
	}
	d := policyResourceDataWithRawConfig(t, raw, noRulesRawConfig())

	p := dataToPolicy("/MyPartition/mypolicy", d)
	assert.Equal(t, "/MyPartition/Drafts/mypolicy", p.Name)
}

func TestDataToPolicyWithRuleActionCondition(t *testing.T) {
	raw := map[string]interface{}{
		"name":     "/Common/mypolicy",
		"strategy": "/Common/first-match",
	}

	forwardAction := cty.ObjectVal(map[string]cty.Value{
		"forward": cty.True,
		"pool":    cty.StringVal("/Common/mypool"),
	})
	condition := cty.ObjectVal(map[string]cty.Value{
		"http_uri": cty.True,
		"contains": cty.True,
		"values":   cty.ListVal([]cty.Value{cty.StringVal("/some/path")}),
	})
	rule1 := cty.ObjectVal(map[string]cty.Value{
		"name":        cty.StringVal("rule1"),
		"description": cty.StringVal("first rule"),
		"action":      cty.ListVal([]cty.Value{forwardAction}),
		"condition":   cty.ListVal([]cty.Value{condition}),
	})
	root := cty.ObjectVal(map[string]cty.Value{
		"rule": cty.ListVal([]cty.Value{rule1}),
	})

	d := policyResourceDataWithRawConfig(t, raw, root)
	p := dataToPolicy("/Common/mypolicy", d)

	require.Len(t, p.Rules, 1)
	assert.Equal(t, "rule1", p.Rules[0].Name)
	assert.Equal(t, "first rule", p.Rules[0].Description)
	require.Len(t, p.Rules[0].Actions, 1)
	assert.True(t, p.Rules[0].Actions[0].Forward)
	assert.Equal(t, "/Common/mypool", p.Rules[0].Actions[0].Pool)
	require.Len(t, p.Rules[0].Conditions, 1)
	assert.True(t, p.Rules[0].Conditions[0].HttpUri)
}

func TestDataToPolicyWithDisableAction(t *testing.T) {
	raw := map[string]interface{}{
		"name": "/Common/mypolicy",
	}

	disableAction := cty.ObjectVal(map[string]cty.Value{
		"disable": cty.True,
		"policy":  cty.StringVal("/Common/other-policy"),
		"select":  cty.True,
		"forward": cty.True,
		"pool":    cty.StringVal("/Common/mypool"),
	})
	rule := cty.ObjectVal(map[string]cty.Value{
		"name":        cty.StringVal("rule2"),
		"description": cty.NullVal(cty.String),
		"action":      cty.ListVal([]cty.Value{disableAction}),
		"condition":   cty.NullVal(cty.List(cty.EmptyObject)),
	})
	root := cty.ObjectVal(map[string]cty.Value{
		"rule": cty.ListVal([]cty.Value{rule}),
	})

	d := policyResourceDataWithRawConfig(t, raw, root)
	p := dataToPolicy("/Common/mypolicy", d)

	require.Len(t, p.Rules, 1)
	require.Len(t, p.Rules[0].Actions, 1)
	// Disable branch should have cleared Policy/Select/Forward/Pool.
	a := p.Rules[0].Actions[0]
	assert.True(t, a.Disable)
	assert.Empty(t, a.Policy)
	assert.False(t, a.Select)
	assert.False(t, a.Forward)
	assert.Empty(t, a.Pool)
}

// ---------------------------------------------------------------------
// policyToData
// ---------------------------------------------------------------------

func TestPolicyToDataSuccess(t *testing.T) {
	d := policyResourceData(t, map[string]interface{}{
		"name": "/Common/mypolicy",
	})

	p := &bigip.Policy{
		FullPath: "/Common/mypolicy",
		Strategy: "/Common/first-match",
		Controls: []string{"forwarding"},
		Requires: []string{"http"},
		Rules: []bigip.PolicyRule{
			{Name: "rule2", Ordinal: 1},
			{Name: "rule1", Ordinal: 0, Actions: []bigip.PolicyRuleAction{{Forward: true, Pool: "/Common/mypool"}}},
		},
	}

	diags := policyToData(p, d)
	require.False(t, diags.HasError())
	assert.Equal(t, "first-match", d.Get("strategy"))
	assert.Equal(t, "/Common/mypolicy", d.Get("name"))
	rules := d.Get("rule").([]interface{})
	require.Len(t, rules, 2)
	// Verify sorted by ordinal: rule1 (ordinal 0) then rule2 (ordinal 1)
	rule0 := rules[0].(map[string]interface{})
	assert.Equal(t, "rule1", rule0["name"])
}

func TestPolicyToDataStrategyRegexFail(t *testing.T) {
	d := policyResourceData(t, map[string]interface{}{
		"name": "/Common/mypolicy",
	})
	p := &bigip.Policy{
		FullPath: "/Common/mypolicy",
		Strategy: "badstrategy-no-slashes",
	}
	diags := policyToData(p, d)
	require.True(t, diags.HasError())
}

func TestPolicyToDataDescriptionGetOk(t *testing.T) {
	d := policyResourceData(t, map[string]interface{}{
		"name":        "/Common/mypolicy",
		"description": "existing description",
	})
	p := &bigip.Policy{FullPath: "/Common/mypolicy"}
	diags := policyToData(p, d)
	require.False(t, diags.HasError())
	assert.Equal(t, "existing description", p.Description)
}

func TestPolicyToDataNoRules(t *testing.T) {
	d := policyResourceData(t, map[string]interface{}{
		"name": "/Common/mypolicy",
	})
	p := &bigip.Policy{FullPath: "/Common/mypolicy"}
	diags := policyToData(p, d)
	require.False(t, diags.HasError())
}

// ---------------------------------------------------------------------
// flatten* / interfaceToResourceData
// ---------------------------------------------------------------------

func TestFlattenPolicyRules(t *testing.T) {
	rules := []bigip.PolicyRule{
		{
			Name:        "rule1",
			Description: "desc1",
			Actions:     []bigip.PolicyRuleAction{{Forward: true, Pool: "/Common/pool1"}},
			Conditions:  []bigip.PolicyRuleCondition{{HttpUri: true}},
		},
		{
			Name: "rule2",
		},
	}
	att := flattenPolicyRules(rules)
	require.Len(t, att, 2)
	obj0 := att[0].(map[string]interface{})
	assert.Equal(t, "rule1", obj0["name"])
	assert.Equal(t, "desc1", obj0["description"])
	assert.NotNil(t, obj0["action"])
	assert.NotNil(t, obj0["condition"])

	obj1 := att[1].(map[string]interface{})
	assert.Equal(t, "rule2", obj1["name"])
	_, hasDesc := obj1["description"]
	assert.False(t, hasDesc)
	_, hasAction := obj1["action"]
	assert.False(t, hasAction)
}

func TestFlattenPolicyRuleActionsAndConditions(t *testing.T) {
	actions := []bigip.PolicyRuleAction{{Forward: true, Pool: "/Common/pool1"}}
	att := flattenPolicyRuleActions(actions)
	require.Len(t, att, 1)
	m := att[0].(map[string]interface{})
	assert.Equal(t, true, m["forward"])
	assert.Equal(t, "/Common/pool1", m["pool"])

	conditions := []bigip.PolicyRuleCondition{{HttpUri: true, Contains: true}}
	catt := flattenPolicyRuleConditions(conditions)
	require.Len(t, catt, 1)
	cm := catt[0].(map[string]interface{})
	assert.Equal(t, true, cm["http_uri"])
	assert.Equal(t, true, cm["contains"])
}

func TestInterfaceToResourceData(t *testing.T) {
	a := bigip.PolicyRuleAction{
		Name:    "0",
		Forward: true,
		Pool:    "/Common/mypool",
	}
	m := interfaceToResourceData(a)
	// "name" field is always excluded.
	_, hasName := m["name"]
	assert.False(t, hasName)
	assert.Equal(t, true, m["forward"])
	assert.Equal(t, "/Common/mypool", m["pool"])
	// Zero-value fields should be omitted.
	_, hasCache := m["cache"]
	assert.False(t, hasCache)
}

// ---------------------------------------------------------------------
// CRUD - Create
// ---------------------------------------------------------------------

func TestResourceBigipLtmPolicyCreateSuccess(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/policy", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{}`)
	})
	mux.HandleFunc("/mgmt/tm/ltm/policy/~Common~mypolicy", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"name":"mypolicy","fullPath":"/Common/mypolicy","strategy":"/Common/first-match"}`)
	})
	mux.HandleFunc("/mgmt/tm/ltm/policy/~Common~mypolicy/rules", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"items":[]}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newDatasourceTestClient(server.URL)
	client.Teem = true

	d := policyResourceDataWithRawConfig(t, map[string]interface{}{
		"name":     "/Common/mypolicy",
		"strategy": "/Common/first-match",
	}, noRulesRawConfig())

	diags := resourceBigipLtmPolicyCreate(context.Background(), d, client)
	require.False(t, diags.HasError(), "unexpected error: %v", diags)
	assert.Equal(t, "/Common/mypolicy", d.Id())
}

func TestResourceBigipLtmPolicyCreateBadName(t *testing.T) {
	client := newDatasourceTestClient("http://127.0.0.1:0")
	client.Teem = true
	d := policyResourceDataWithRawConfig(t, map[string]interface{}{
		"name": "mypolicy-no-slashes",
	}, noRulesRawConfig())

	diags := resourceBigipLtmPolicyCreate(context.Background(), d, client)
	require.True(t, diags.HasError())
	assert.Contains(t, diags[0].Summary, "failed to match the regex")
}

func TestResourceBigipLtmPolicyCreateCreatePolicyError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/policy", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "internal error", http.StatusInternalServerError)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newDatasourceTestClient(server.URL)
	client.Teem = true

	d := policyResourceDataWithRawConfig(t, map[string]interface{}{
		"name": "/Common/mypolicy",
	}, noRulesRawConfig())

	diags := resourceBigipLtmPolicyCreate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestResourceBigipLtmPolicyCreatePublishError(t *testing.T) {
	calls := 0
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/policy", func(w http.ResponseWriter, r *http.Request) {
		calls++
		if calls == 1 {
			// CreatePolicy succeeds
			w.Header().Set("Content-Type", "application/json")
			_, _ = fmt.Fprint(w, `{}`)
			return
		}
		// PublishPolicy fails
		http.Error(w, "publish failed", http.StatusInternalServerError)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newDatasourceTestClient(server.URL)
	client.Teem = true

	d := policyResourceDataWithRawConfig(t, map[string]interface{}{
		"name": "/Common/mypolicy",
	}, noRulesRawConfig())

	diags := resourceBigipLtmPolicyCreate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

// ---------------------------------------------------------------------
// CRUD - Read
// ---------------------------------------------------------------------

func TestResourceBigipLtmPolicyReadSuccess(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/policy/~Common~mypolicy", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"name":"mypolicy","fullPath":"/Common/mypolicy","strategy":"/Common/first-match"}`)
	})
	mux.HandleFunc("/mgmt/tm/ltm/policy/~Common~mypolicy/rules", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"items":[]}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newDatasourceTestClient(server.URL)
	client.Teem = true

	d := policyResourceData(t, map[string]interface{}{"name": "/Common/mypolicy"})
	d.SetId("/Common/mypolicy")

	diags := resourceBigipLtmPolicyRead(context.Background(), d, client)
	require.False(t, diags.HasError(), "unexpected error: %v", diags)
	assert.Equal(t, "/Common/mypolicy", d.Get("name"))
}

func TestResourceBigipLtmPolicyReadBadName(t *testing.T) {
	client := newDatasourceTestClient("http://127.0.0.1:0")
	client.Teem = true
	d := policyResourceData(t, map[string]interface{}{"name": "x"})
	d.SetId("mypolicy-no-slashes")

	diags := resourceBigipLtmPolicyRead(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestResourceBigipLtmPolicyReadError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/policy/~Common~mypolicy", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "boom", http.StatusInternalServerError)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newDatasourceTestClient(server.URL)
	client.Teem = true

	d := policyResourceData(t, map[string]interface{}{"name": "x"})
	d.SetId("/Common/mypolicy")

	diags := resourceBigipLtmPolicyRead(context.Background(), d, client)
	require.True(t, diags.HasError())
}

// ---------------------------------------------------------------------
// CRUD - Update
// ---------------------------------------------------------------------

func TestResourceBigipLtmPolicyUpdateSuccessDraftExists(t *testing.T) {
	mux := http.NewServeMux()
	// CheckDraftPolicy GET -> draft exists
	mux.HandleFunc("/mgmt/tm/ltm/policy/~Common~Drafts~2Fmypolicy", func(w http.ResponseWriter, r *http.Request) {
		// URL-encoded slash variant; not expected to be hit directly since
		// go-bigip encodes internally, handled by catch-all below instead.
	})
	// Catch-all handler for the ~Common~... family of paths (draft check,
	// update PATCH, and the read-back GET/rules GET all live under this
	// prefix with URL-encoded path segments).
	mux.HandleFunc("/mgmt/tm/ltm/policy/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.Method == http.MethodGet && strings.HasSuffix(r.URL.Path, "/rules"):
			_, _ = fmt.Fprint(w, `{"items":[]}`)
		case r.Method == http.MethodGet && strings.Contains(r.URL.Path, "Drafts"):
			// CheckDraftPolicy: report the draft already exists.
			_, _ = fmt.Fprint(w, `{"name":"mypolicy","fullPath":"/Common/Drafts/mypolicy"}`)
		case r.Method == http.MethodGet:
			_, _ = fmt.Fprint(w, `{"name":"mypolicy","fullPath":"/Common/mypolicy","strategy":"/Common/first-match"}`)
		case r.Method == http.MethodPatch:
			_, _ = fmt.Fprint(w, `{}`)
		default:
			http.Error(w, "unexpected method", http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/mgmt/tm/ltm/policy", func(w http.ResponseWriter, r *http.Request) {
		// PublishPolicy POST
		assert.Equal(t, http.MethodPost, r.Method)
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newDatasourceTestClient(server.URL)
	client.Teem = true

	d := policyResourceDataWithRawConfig(t, map[string]interface{}{
		"name":     "/Common/mypolicy",
		"strategy": "/Common/first-match",
	}, noRulesRawConfig())
	d.SetId("/Common/mypolicy")

	diags := resourceBigipLtmPolicyUpdate(context.Background(), d, client)
	require.False(t, diags.HasError(), "unexpected error: %v", diags)
}

func TestResourceBigipLtmPolicyUpdateCreateDraftThenSuccess(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/policy/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.Method == http.MethodGet && strings.HasSuffix(r.URL.Path, "/rules"):
			_, _ = fmt.Fprint(w, `{"items":[]}`)
		case r.Method == http.MethodGet && strings.Contains(r.URL.Path, "Drafts"):
			// CheckDraftPolicy: report the draft doesn't exist.
			w.WriteHeader(http.StatusNotFound)
			_, _ = fmt.Fprint(w, `{"code":404}`)
		case r.Method == http.MethodGet:
			_, _ = fmt.Fprint(w, `{"name":"mypolicy","fullPath":"/Common/mypolicy","strategy":"/Common/first-match"}`)
		case r.Method == http.MethodPatch:
			// Both CreatePolicyDraft and UpdatePolicy PATCH here.
			_, _ = fmt.Fprint(w, `{}`)
		default:
			http.Error(w, "unexpected method", http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/mgmt/tm/ltm/policy", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newDatasourceTestClient(server.URL)
	client.Teem = true

	d := policyResourceDataWithRawConfig(t, map[string]interface{}{
		"name":     "/Common/mypolicy",
		"strategy": "/Common/first-match",
	}, noRulesRawConfig())
	d.SetId("/Common/mypolicy")

	diags := resourceBigipLtmPolicyUpdate(context.Background(), d, client)
	require.False(t, diags.HasError(), "unexpected error: %v", diags)
}

func TestResourceBigipLtmPolicyUpdateBadName(t *testing.T) {
	client := newDatasourceTestClient("http://127.0.0.1:0")
	client.Teem = true
	d := policyResourceDataWithRawConfig(t, map[string]interface{}{
		"name": "x",
	}, noRulesRawConfig())
	d.SetId("mypolicy-no-slashes")

	diags := resourceBigipLtmPolicyUpdate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestResourceBigipLtmPolicyUpdateCreateDraftError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/policy/", func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && strings.Contains(r.URL.Path, "Drafts"):
			w.WriteHeader(http.StatusNotFound)
			_, _ = fmt.Fprint(w, `{"code":404}`)
		case r.Method == http.MethodPatch:
			http.Error(w, "create draft failed", http.StatusInternalServerError)
		default:
			http.Error(w, "unexpected", http.StatusMethodNotAllowed)
		}
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newDatasourceTestClient(server.URL)
	client.Teem = true

	d := policyResourceDataWithRawConfig(t, map[string]interface{}{
		"name": "/Common/mypolicy",
	}, noRulesRawConfig())
	d.SetId("/Common/mypolicy")

	diags := resourceBigipLtmPolicyUpdate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestResourceBigipLtmPolicyUpdateUpdatePolicyError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/policy/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.Method == http.MethodGet && strings.Contains(r.URL.Path, "Drafts"):
			_, _ = fmt.Fprint(w, `{"name":"mypolicy","fullPath":"/Common/Drafts/mypolicy"}`)
		case r.Method == http.MethodPatch:
			http.Error(w, "update failed", http.StatusInternalServerError)
		default:
			http.Error(w, "unexpected", http.StatusMethodNotAllowed)
		}
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newDatasourceTestClient(server.URL)
	client.Teem = true

	d := policyResourceDataWithRawConfig(t, map[string]interface{}{
		"name": "/Common/mypolicy",
	}, noRulesRawConfig())
	d.SetId("/Common/mypolicy")

	diags := resourceBigipLtmPolicyUpdate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestResourceBigipLtmPolicyUpdatePublishPolicyError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/policy/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.Method == http.MethodGet && strings.Contains(r.URL.Path, "Drafts"):
			_, _ = fmt.Fprint(w, `{"name":"mypolicy","fullPath":"/Common/Drafts/mypolicy"}`)
		case r.Method == http.MethodPatch:
			_, _ = fmt.Fprint(w, `{}`)
		default:
			http.Error(w, "unexpected", http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/mgmt/tm/ltm/policy", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "publish failed", http.StatusInternalServerError)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newDatasourceTestClient(server.URL)
	client.Teem = true

	d := policyResourceDataWithRawConfig(t, map[string]interface{}{
		"name": "/Common/mypolicy",
	}, noRulesRawConfig())
	d.SetId("/Common/mypolicy")

	diags := resourceBigipLtmPolicyUpdate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

// ---------------------------------------------------------------------
// CRUD - Delete
// ---------------------------------------------------------------------

func TestResourceBigipLtmPolicyDeleteSuccess(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/policy/", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		w.WriteHeader(http.StatusOK)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newDatasourceTestClient(server.URL)
	client.Teem = true

	d := policyResourceData(t, map[string]interface{}{"name": "x"})
	d.SetId("/Common/mypolicy")

	diags := resourceBigipLtmPolicyDelete(context.Background(), d, client)
	require.False(t, diags.HasError(), "unexpected error: %v", diags)
	assert.Empty(t, d.Id())
}

func TestResourceBigipLtmPolicyDeleteBadName(t *testing.T) {
	client := newDatasourceTestClient("http://127.0.0.1:0")
	client.Teem = true
	d := policyResourceData(t, map[string]interface{}{"name": "x"})
	d.SetId("mypolicy-no-slashes")

	diags := resourceBigipLtmPolicyDelete(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestResourceBigipLtmPolicyDeleteError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/policy/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "delete failed", http.StatusInternalServerError)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newDatasourceTestClient(server.URL)
	client.Teem = true

	d := policyResourceData(t, map[string]interface{}{"name": "x"})
	d.SetId("/Common/mypolicy")

	diags := resourceBigipLtmPolicyDelete(context.Background(), d, client)
	require.True(t, diags.HasError())
}
