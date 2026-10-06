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
	"testing"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// policyDatasourceResourceData builds a *schema.ResourceData for the
// bigip_ltm_policy data source using schema.TestResourceDataRaw.
func policyDatasourceResourceData(t *testing.T, raw map[string]interface{}) *schema.ResourceData {
	t.Helper()
	return newDatasourceTestResourceData(t, dataSourceBigipLtmPolicy(), raw)
}

// ---------------------------------------------------------------------
// Schema-shape assertions
// ---------------------------------------------------------------------

func TestDataSourceBigipLtmPolicySchema(t *testing.T) {
	r := dataSourceBigipLtmPolicy()

	require.NotNil(t, r.Schema)
	require.NotNil(t, r.ReadContext)

	nameSchema, ok := r.Schema["name"]
	require.True(t, ok)
	assert.True(t, nameSchema.Required)

	strategySchema, ok := r.Schema["strategy"]
	require.True(t, ok)
	assert.Equal(t, "/Common/first-match", strategySchema.Default)

	ruleSchema, ok := r.Schema["rule"]
	require.True(t, ok)
	assert.Equal(t, schema.TypeList, ruleSchema.Type)
}

// ---------------------------------------------------------------------
// parsePolicyPath
// ---------------------------------------------------------------------

func TestParsePolicyPathTwoFields(t *testing.T) {
	fields, diags := parsePolicyPath("/Common/mypolicy")
	require.Nil(t, diags)
	assert.Equal(t, []string{"Common", "mypolicy"}, fields)
}

func TestParsePolicyPathThreeFields(t *testing.T) {
	fields, diags := parsePolicyPath("/Common/Drafts/mypolicy")
	require.Nil(t, diags)
	assert.Equal(t, []string{"Common", "Drafts", "mypolicy"}, fields)
}

func TestParsePolicyPathInvalid(t *testing.T) {
	fields, diags := parsePolicyPath("mypolicy-no-slashes")
	require.NotNil(t, diags)
	assert.Nil(t, fields)
}

// ---------------------------------------------------------------------
// DatapolicyToData
// ---------------------------------------------------------------------

func TestDatapolicyToDataSuccess(t *testing.T) {
	d := policyDatasourceResourceData(t, map[string]interface{}{"name": "/Common/mypolicy"})

	p := &bigip.Policy{
		FullPath: "/Common/mypolicy",
		Strategy: "/Common/first-match",
		Controls: []string{"forwarding"},
		Requires: []string{"http"},
		Rules: []bigip.PolicyRule{
			{Name: "rule2", Ordinal: 1},
			{Name: "rule1", Ordinal: 0, Actions: []bigip.PolicyRuleAction{{Forward: true, Pool: "/Common/mypool"}},
				Conditions: []bigip.PolicyRuleCondition{{HttpUri: true}}},
		},
	}

	diags := DatapolicyToData(p, d)
	require.False(t, diags.HasError())
	assert.Equal(t, "first-match", d.Get("strategy"))
	assert.Equal(t, "/Common/mypolicy", d.Get("name"))
	rules := d.Get("rule").([]interface{})
	require.Len(t, rules, 2)
	rule0 := rules[0].(map[string]interface{})
	assert.Equal(t, "rule1", rule0["name"])
}

func TestDatapolicyToDataStrategyRegexFail(t *testing.T) {
	d := policyDatasourceResourceData(t, map[string]interface{}{"name": "/Common/mypolicy"})
	p := &bigip.Policy{
		FullPath: "/Common/mypolicy",
		Strategy: "badstrategy-no-slashes",
	}
	diags := DatapolicyToData(p, d)
	require.True(t, diags.HasError())
}

func TestDatapolicyToDataNoRules(t *testing.T) {
	d := policyDatasourceResourceData(t, map[string]interface{}{"name": "/Common/mypolicy"})
	p := &bigip.Policy{FullPath: "/Common/mypolicy"}
	diags := DatapolicyToData(p, d)
	require.False(t, diags.HasError())
}

// ---------------------------------------------------------------------
// DataflattenPolicyRules / Actions / Conditions / DatainterfaceToResourceData
// ---------------------------------------------------------------------

func TestDataflattenPolicyRules(t *testing.T) {
	rules := []bigip.PolicyRule{
		{
			Name:       "rule1",
			Actions:    []bigip.PolicyRuleAction{{Forward: true, Pool: "/Common/pool1"}},
			Conditions: []bigip.PolicyRuleCondition{{HttpUri: true}},
		},
		{
			Name: "rule2",
		},
	}
	att := DataflattenPolicyRules(rules)
	require.Len(t, att, 2)
	obj0 := att[0].(map[string]any)
	assert.Equal(t, "rule1", obj0["name"])
	assert.NotNil(t, obj0["action"])
	assert.NotNil(t, obj0["condition"])

	obj1 := att[1].(map[string]any)
	assert.Equal(t, "rule2", obj1["name"])
	_, hasAction := obj1["action"]
	assert.False(t, hasAction)
}

func TestDataflattenPolicyRuleActionsAndConditions(t *testing.T) {
	actions := []bigip.PolicyRuleAction{{Forward: true, Pool: "/Common/pool1"}}
	att := DataflattenPolicyRuleActions(actions)
	require.Len(t, att, 1)
	m := att[0].(map[string]any)
	assert.Equal(t, true, m["forward"])
	assert.Equal(t, "/Common/pool1", m["pool"])

	conditions := []bigip.PolicyRuleCondition{{HttpUri: true, Contains: true}}
	catt := DataflattenPolicyRuleConditions(conditions)
	require.Len(t, catt, 1)
	cm := catt[0].(map[string]any)
	assert.Equal(t, true, cm["http_uri"])
	assert.Equal(t, true, cm["contains"])
}

func TestDatainterfaceToResourceData(t *testing.T) {
	a := bigip.PolicyRuleAction{
		Name:    "0",
		Forward: true,
		Pool:    "/Common/mypool",
	}
	m := DatainterfaceToResourceData(a)
	_, hasName := m["name"]
	assert.False(t, hasName)
	assert.Equal(t, true, m["forward"])
	assert.Equal(t, "/Common/mypool", m["pool"])
	_, hasCache := m["cache"]
	assert.False(t, hasCache)
}

// ---------------------------------------------------------------------
// dataSourceBigipLtmPolicyRead
// ---------------------------------------------------------------------

func TestDataSourceBigipLtmPolicyReadTwoFieldSuccess(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/policy/~Common~mypolicy", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"name":"mypolicy","fullPath":"/Common/mypolicy","strategy":"/Common/first-match"}`)
	})
	mux.HandleFunc("/mgmt/tm/ltm/policy/~Common~mypolicy/rules", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"items":[{"name":"rule1"}]}`)
	})
	mux.HandleFunc("/mgmt/tm/ltm/policy/~Common~mypolicy/rules/rule1/actions", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"items":[{"forward":true,"pool":"/Common/mypool"}]}`)
	})
	mux.HandleFunc("/mgmt/tm/ltm/policy/~Common~mypolicy/rules/rule1/conditions", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"items":[]}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newDatasourceTestClient(server.URL)
	client.Teem = true

	d := policyDatasourceResourceData(t, map[string]interface{}{"name": "/Common/mypolicy"})

	diags := dataSourceBigipLtmPolicyRead(context.Background(), d, client)
	require.False(t, diags.HasError(), "unexpected error: %v", diags)
	assert.Equal(t, "/Common/mypolicy", d.Id())
	assert.Equal(t, "/Common/mypolicy", d.Get("name"))
	rules := d.Get("rule").([]interface{})
	require.Len(t, rules, 1)
}

func TestDataSourceBigipLtmPolicyReadThreeFieldSuccess(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/policy/~Common~Drafts~mypolicy", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"name":"mypolicy","fullPath":"/Common/Drafts/mypolicy"}`)
	})
	mux.HandleFunc("/mgmt/tm/ltm/policy/~Common~Drafts~mypolicy/rules", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"items":[]}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newDatasourceTestClient(server.URL)
	client.Teem = true

	d := policyDatasourceResourceData(t, map[string]interface{}{"name": "/Common/Drafts/mypolicy"})

	diags := dataSourceBigipLtmPolicyRead(context.Background(), d, client)
	require.False(t, diags.HasError(), "unexpected error: %v", diags)
	assert.Equal(t, "/Common/Drafts/mypolicy", d.Id())
}

func TestDataSourceBigipLtmPolicyReadInvalidName(t *testing.T) {
	client := newDatasourceTestClient("http://127.0.0.1:0")
	client.Teem = true
	d := policyDatasourceResourceData(t, map[string]interface{}{"name": "mypolicy-no-slashes"})

	diags := dataSourceBigipLtmPolicyRead(context.Background(), d, client)
	require.True(t, diags.HasError())
	assert.Contains(t, diags[0].Summary, "failed to parse policy path")
}

func TestDataSourceBigipLtmPolicyReadFetchError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/policy/~Common~mypolicy", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "boom", http.StatusInternalServerError)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newDatasourceTestClient(server.URL)
	client.Teem = true

	d := policyDatasourceResourceData(t, map[string]interface{}{"name": "/Common/mypolicy"})

	diags := dataSourceBigipLtmPolicyRead(context.Background(), d, client)
	require.True(t, diags.HasError())
	assert.Empty(t, d.Id())
}
