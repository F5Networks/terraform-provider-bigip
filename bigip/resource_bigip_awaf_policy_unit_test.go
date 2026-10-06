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
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------
// Schema-shape assertion (no HTTP server needed)
// ---------------------------------------------------------------------

func TestResourceBigipAwafPolicySchema(t *testing.T) {
	r := resourceBigipAwafPolicy()

	if r.Schema == nil {
		t.Fatal("Expected schema to be defined")
	}
	if r.CreateContext == nil {
		t.Fatal("Expected CreateContext to be defined")
	}
	if r.ReadContext == nil {
		t.Fatal("Expected ReadContext to be defined")
	}
	if r.UpdateContext == nil {
		t.Fatal("Expected UpdateContext to be defined")
	}
	if r.DeleteContext == nil {
		t.Fatal("Expected DeleteContext to be defined")
	}
	if r.Importer == nil {
		t.Fatal("Expected Importer to be defined")
	}

	nameSchema, ok := r.Schema["name"]
	if !ok {
		t.Fatal("Expected field 'name' to exist in schema")
	}
	if !nameSchema.Required {
		t.Error("Expected field 'name' to be required")
	}
	if !nameSchema.ForceNew {
		t.Error("Expected field 'name' to be ForceNew")
	}

	partitionSchema, ok := r.Schema["partition"]
	if !ok {
		t.Fatal("Expected field 'partition' to exist in schema")
	}
	if partitionSchema.Default != "Common" {
		t.Errorf("Expected field 'partition' default to be 'Common', got %v", partitionSchema.Default)
	}

	typeSchema, ok := r.Schema["type"]
	if !ok {
		t.Fatal("Expected field 'type' to exist in schema")
	}
	if typeSchema.Default != "security" {
		t.Errorf("Expected field 'type' default to be 'security', got %v", typeSchema.Default)
	}
}

func TestResourceBigipAwafPolicyEnforcementModeValidation(t *testing.T) {
	r := resourceBigipAwafPolicy()
	s := r.Schema["enforcement_mode"]
	require.NotNil(t, s.ValidateFunc)

	for _, v := range []string{"blocking", "transparent"} {
		_, errs := s.ValidateFunc(v, "enforcement_mode")
		assert.Empty(t, errs, "expected %q to be valid", v)
	}

	_, errs := s.ValidateFunc("bogus", "enforcement_mode")
	assert.NotEmpty(t, errs, "expected 'bogus' to be invalid")
}

func TestResourceBigipAwafPolicyTypeValidation(t *testing.T) {
	r := resourceBigipAwafPolicy()
	s := r.Schema["type"]
	require.NotNil(t, s.ValidateFunc)

	for _, v := range []string{"parent", "security"} {
		_, errs := s.ValidateFunc(v, "type")
		assert.Empty(t, errs, "expected %q to be valid", v)
	}

	_, errs := s.ValidateFunc("bogus", "type")
	assert.NotEmpty(t, errs, "expected 'bogus' to be invalid")
}

func TestResourceBigipAwafPolicyImportJsonDiffSuppress(t *testing.T) {
	r := resourceBigipAwafPolicy()
	diffSuppress := r.Schema["policy_import_json"].DiffSuppressFunc
	require.NotNil(t, diffSuppress)

	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":          "test-policy",
		"template_name": "POLICY_TEMPLATE_RAPID_DEPLOYMENT",
	}, "")

	old := `{"policy":{"name":"test-policy"}}`
	same := `{"policy":{"name":"test-policy"}}`
	assert.True(t, diffSuppress("policy_import_json", old, same, d), "expected identical JSON to suppress diff")

	different := `{"policy":{"name":"other-policy"}}`
	assert.False(t, diffSuppress("policy_import_json", old, different, d), "expected different JSON to NOT suppress diff")
}

// ---------------------------------------------------------------------
// getpolicyConfig (direct function tests, no HTTP server needed)
// ---------------------------------------------------------------------

func TestGetPolicyConfig_Minimal(t *testing.T) {
	r := resourceBigipAwafPolicy()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":          "test-policy",
		"partition":     "Common",
		"template_name": "POLICY_TEMPLATE_RAPID_DEPLOYMENT",
	}, "")

	config, err := getpolicyConfig(d)
	require.NoError(t, err)
	assert.Contains(t, config, "test-policy")
	assert.Contains(t, config, "POLICY_TEMPLATE_RAPID_DEPLOYMENT")
}

func TestGetPolicyConfig_NonCommonPartition(t *testing.T) {
	r := resourceBigipAwafPolicy()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":          "test-policy",
		"partition":     "MyPartition",
		"template_name": "POLICY_TEMPLATE_RAPID_DEPLOYMENT",
	}, "")

	config, err := getpolicyConfig(d)
	require.NoError(t, err)
	assert.Contains(t, config, "/MyPartition/test-policy")
}

func TestGetPolicyConfig_WithApplicationLanguage(t *testing.T) {
	r := resourceBigipAwafPolicy()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":                 "test-policy",
		"partition":            "Common",
		"template_name":        "POLICY_TEMPLATE_RAPID_DEPLOYMENT",
		"application_language": "iso-8859-1",
		"case_insensitive":     true,
		"enable_passivemode":   true,
		"protocol_independent": true,
		"enforcement_mode":     "blocking",
		"description":          "a test policy",
	}, "")

	config, err := getpolicyConfig(d)
	require.NoError(t, err)
	assert.Contains(t, config, "iso-8859-1")
	assert.Contains(t, config, "blocking")
	assert.Contains(t, config, "a test policy")
}

func TestGetPolicyConfig_WithTemplateLink(t *testing.T) {
	r := resourceBigipAwafPolicy()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":          "test-policy",
		"partition":     "Common",
		"template_name": "POLICY_TEMPLATE_RAPID_DEPLOYMENT",
		"template_link": "https://localhost/mgmt/tm/asm/policy-templates/abc123",
	}, "")

	config, err := getpolicyConfig(d)
	require.NoError(t, err)
	assert.Contains(t, config, "abc123")
}

func TestGetPolicyConfig_WithSignaturesSettings(t *testing.T) {
	r := resourceBigipAwafPolicy()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":          "test-policy",
		"partition":     "Common",
		"template_name": "POLICY_TEMPLATE_RAPID_DEPLOYMENT",
		"signatures_settings": []interface{}{
			map[string]interface{}{
				"signature_staging":          true,
				"placesignatures_in_staging": false,
			},
		},
	}, "")

	config, err := getpolicyConfig(d)
	require.NoError(t, err)
	assert.Contains(t, config, "signatureStaging")
}

func TestGetPolicyConfig_WithPolicyBuilder(t *testing.T) {
	r := resourceBigipAwafPolicy()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":          "test-policy",
		"partition":     "Common",
		"template_name": "POLICY_TEMPLATE_RAPID_DEPLOYMENT",
		"policy_builder": []interface{}{
			map[string]interface{}{
				"learning_mode": "disabled",
			},
		},
	}, "")

	config, err := getpolicyConfig(d)
	require.NoError(t, err)
	assert.Contains(t, config, "disabled")
}

func TestGetPolicyConfig_WithGraphqlProfiles(t *testing.T) {
	r := resourceBigipAwafPolicy()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":          "test-policy",
		"partition":     "Common",
		"template_name": "POLICY_TEMPLATE_RAPID_DEPLOYMENT",
		"graphql_profiles": []interface{}{
			map[string]interface{}{
				"name":                    "test_graphql",
				"metachar_elementcheck":   true,
				"attack_signatures_check": true,
				"defense_attributes": []interface{}{
					map[string]interface{}{
						"allow_introspection_queries": true,
						"tolerate_parsing_warnings":   false,
						"maximum_batched_queries":     "10",
						"maximum_structure_depth":     "10",
						"maximum_total_length":        "100000",
						"maximum_value_length":        "10000",
					},
				},
			},
		},
	}, "")

	config, err := getpolicyConfig(d)
	require.NoError(t, err)
	assert.Contains(t, config, "test_graphql")
	assert.Contains(t, config, "graphql-profiles")
}

func TestGetPolicyConfig_WithHostNames(t *testing.T) {
	r := resourceBigipAwafPolicy()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":          "test-policy",
		"partition":     "Common",
		"template_name": "POLICY_TEMPLATE_RAPID_DEPLOYMENT",
		"host_names": []interface{}{
			map[string]interface{}{
				"name": "www.example.com",
			},
		},
	}, "")

	config, err := getpolicyConfig(d)
	require.NoError(t, err)
	assert.Contains(t, config, "www.example.com")
	assert.Contains(t, config, "host-names")
}

func TestGetPolicyConfig_WithFileTypes(t *testing.T) {
	r := resourceBigipAwafPolicy()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":          "test-policy",
		"partition":     "Common",
		"template_name": "POLICY_TEMPLATE_RAPID_DEPLOYMENT",
		"file_types": []interface{}{
			map[string]interface{}{
				"name":    "testfiletype",
				"type":    "explicit",
				"allowed": true,
			},
		},
	}, "")

	config, err := getpolicyConfig(d)
	require.NoError(t, err)
	assert.Contains(t, config, "testfiletype")
	assert.Contains(t, config, "filetypes")
}

func TestGetPolicyConfig_WithIpExceptions(t *testing.T) {
	r := resourceBigipAwafPolicy()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":          "test-policy",
		"partition":     "Common",
		"template_name": "POLICY_TEMPLATE_RAPID_DEPLOYMENT",
		"ip_exceptions": []interface{}{
			map[string]interface{}{
				"ip_address":              "100.10.10.0",
				"ip_mask":                 "255.255.255.0",
				"description":             "trusted range",
				"block_requests":          "never",
				"trustedby_policybuilder": true,
				"ignore_anomalies":        true,
				"ignore_ipreputation":     true,
			},
		},
	}, "")

	config, err := getpolicyConfig(d)
	require.NoError(t, err)
	assert.Contains(t, config, "100.10.10.0")
	assert.Contains(t, config, "whitelist-ips")
}

func TestGetPolicyConfig_WithServerTechnologies(t *testing.T) {
	r := resourceBigipAwafPolicy()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":                "test-policy",
		"partition":           "Common",
		"template_name":       "POLICY_TEMPLATE_RAPID_DEPLOYMENT",
		"server_technologies": []interface{}{"MySQL", "Unix/Linux"},
	}, "")

	config, err := getpolicyConfig(d)
	require.NoError(t, err)
	assert.Contains(t, config, "MySQL")
	assert.Contains(t, config, "server-technologies")
}

func TestGetPolicyConfig_WithUrlsParametersSignatures(t *testing.T) {
	r := resourceBigipAwafPolicy()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":          "test-policy",
		"partition":     "Common",
		"template_name": "POLICY_TEMPLATE_RAPID_DEPLOYMENT",
		"urls": []interface{}{
			`{"name":"/foo","protocol":"https","type":"explicit"}`,
		},
		"parameters": []interface{}{
			`{"name":"param1","type":"explicit"}`,
		},
		"signature_sets": []interface{}{
			`{"name":"sigset1","block":true}`,
		},
		"signatures": []interface{}{
			`{"name":"sig1","enabled":true}`,
		},
		"open_api_files": []interface{}{
			"https://example.com/openapi.json",
		},
	}, "")

	config, err := getpolicyConfig(d)
	require.NoError(t, err)
	assert.Contains(t, config, "/foo")
	assert.Contains(t, config, "param1")
	assert.Contains(t, config, "sigset1")
	assert.Contains(t, config, "sig1")
	assert.Contains(t, config, "openapi.json")
}

func TestGetPolicyConfig_WithModifications(t *testing.T) {
	r := resourceBigipAwafPolicy()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":          "test-policy",
		"partition":     "Common",
		"template_name": "POLICY_TEMPLATE_RAPID_DEPLOYMENT",
		"modifications": []interface{}{
			`{"suggestions":[{"action":"add","entity":"url"}]}`,
		},
	}, "")

	config, err := getpolicyConfig(d)
	require.NoError(t, err)
	assert.Contains(t, config, "modifications")
	assert.Contains(t, config, "add")
}

func TestGetPolicyConfig_WithPolicyImportJson(t *testing.T) {
	r := resourceBigipAwafPolicy()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":               "test-policy",
		"partition":          "Common",
		"template_name":      "POLICY_TEMPLATE_RAPID_DEPLOYMENT",
		"policy_import_json": `{"policy":{"name":"old-name","fullPath":"/Common/old-name"}}`,
	}, "")

	config, err := getpolicyConfig(d)
	require.NoError(t, err)
	assert.Contains(t, config, "test-policy")
	assert.Contains(t, config, "/Common/test-policy")
}

func TestGetPolicyConfig_WithPolicyImportJsonAndExistingLists(t *testing.T) {
	r := resourceBigipAwafPolicy()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":          "test-policy",
		"partition":     "Common",
		"template_name": "POLICY_TEMPLATE_RAPID_DEPLOYMENT",
		"urls": []interface{}{
			`{"name":"/new-url"}`,
		},
		"parameters": []interface{}{
			`{"name":"newparam"}`,
		},
		"signature_sets": []interface{}{
			`{"name":"newsigset"}`,
		},
		"file_types": []interface{}{
			map[string]interface{}{"name": "newft", "type": "explicit", "allowed": true},
		},
		"ip_exceptions": []interface{}{
			map[string]interface{}{"ip_address": "1.2.3.4", "ip_mask": "255.255.255.255"},
		},
		"host_names": []interface{}{
			map[string]interface{}{"name": "new.example.com"},
		},
		"server_technologies": []interface{}{"NewTech"},
		"open_api_files":      []interface{}{"https://example.com/new-openapi.json"},
		"graphql_profiles": []interface{}{
			map[string]interface{}{
				"name": "newgraphql",
				"defense_attributes": []interface{}{
					map[string]interface{}{},
				},
			},
		},
		"modifications": []interface{}{
			`{"suggestions":[{"action":"add"}]}`,
		},
		"policy_import_json": `{"policy":{
			"name":"old-name",
			"fullPath":"/Common/old-name",
			"description":"existing description",
			"urls":[{"name":"/existing-url"}],
			"parameters":[{"name":"existingparam"}],
			"signature-sets":[{"name":"existingsigset"}],
			"filetypes":[{"name":"existingft"}],
			"whitelist-ips":[{"ipAddress":"5.6.7.8"}],
			"host-names":[{"name":"existing.example.com"}],
			"server-technologies":[{"serverTechnologyName":"ExistingTech"}],
			"open-api-files":[{"link":"https://example.com/existing-openapi.json"}],
			"graphql-profiles":[{"name":"existinggraphql"}]
		}}`,
	}, "")

	config, err := getpolicyConfig(d)
	require.NoError(t, err)
	assert.Contains(t, config, "existing-url")
	assert.Contains(t, config, "new-url")
	assert.Contains(t, config, "existingparam")
	assert.Contains(t, config, "newparam")
	assert.Contains(t, config, "existingsigset")
	assert.Contains(t, config, "newsigset")
	assert.Contains(t, config, "existingft")
	assert.Contains(t, config, "newft")
	assert.Contains(t, config, "5.6.7.8")
	assert.Contains(t, config, "existing.example.com")
	assert.Contains(t, config, "new.example.com")
	assert.Contains(t, config, "ExistingTech")
	assert.Contains(t, config, "NewTech")
	assert.Contains(t, config, "existing-openapi.json")
	assert.Contains(t, config, "new-openapi.json")
	assert.Contains(t, config, "existinggraphql")
	assert.Contains(t, config, "newgraphql")
}

func TestGetPolicyConfig_WithPolicyImportJsonNoTemplateChangeNeeded(t *testing.T) {
	// Exercises the branch where fullPath already matches policyWaf.Name so the
	// name/fullPath rewrite is skipped, and template.Name is empty so the
	// template rewrite is skipped too.
	r := resourceBigipAwafPolicy()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":               "test-policy",
		"partition":          "Common",
		"policy_import_json": `{"policy":{"name":"test-policy","fullPath":"test-policy"}}`,
	}, "")

	config, err := getpolicyConfig(d)
	require.NoError(t, err)
	assert.Contains(t, config, "test-policy")
}

// ---------------------------------------------------------------------
// resourceBigipAwafPolicyCreate / Read / Update / Delete
// (direct CRUD-function tests using NewUnitTestClient/NewTestResourceData
// and a locally-owned httptest.Server)
// ---------------------------------------------------------------------

func awafPolicyMockServer(t *testing.T, policyID, policyName, partition string) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()

	mux.HandleFunc("/mgmt/tm/sys/provision/asm", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"name":"asm","level":"nominal"}`)
	})
	mux.HandleFunc("/mgmt/tm/asm/file-transfer/uploads/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"remainingByteCount":0,"usedChunks":{},"totalByteCount":10,"localFilePath":"","temporaryFilePath":"","generation":0,"lastUpdateMicros":0}`)
	})
	mux.HandleFunc("/mgmt/tm/asm/tasks/import-policy", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"status":"COMPLETED","id":"import-task-1"}`)
	})
	mux.HandleFunc("/mgmt/tm/asm/tasks/import-policy/import-task-1", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"status":"COMPLETED","id":"import-task-1"}`)
	})
	mux.HandleFunc("/mgmt/tm/asm/policies/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"items":[{"name":"%s","partition":"%s","id":"%s"}]}`, policyName, partition, policyID)
	})
	mux.HandleFunc("/mgmt/tm/asm/tasks/apply-policy", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"status":"COMPLETED","id":"apply-task-1"}`)
	})
	mux.HandleFunc("/mgmt/tm/asm/tasks/apply-policy/apply-task-1", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"status":"COMPLETED","id":"apply-task-1"}`)
	})
	mux.HandleFunc("/mgmt/tm/asm/policies/"+policyID, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodDelete {
			w.WriteHeader(http.StatusOK)
			return
		}
		_, _ = fmt.Fprintf(w, `{"name":"%s","partition":"%s","id":"%s","fullPath":"/%s/%s"}`, policyName, partition, policyID, partition, policyName)
	})
	mux.HandleFunc("/mgmt/tm/asm/tasks/export-policy", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"status":"COMPLETED","id":"export-task-1"}`)
	})
	mux.HandleFunc("/mgmt/tm/asm/tasks/export-policy/export-task-1", func(w http.ResponseWriter, r *http.Request) {
		file := fmt.Sprintf(`{"policy":{"name":"%s","fullPath":"/%s/%s","type":"security","template":{"name":"POLICY_TEMPLATE_RAPID_DEPLOYMENT"}}}`, policyName, partition, policyName)
		_, _ = fmt.Fprintf(w, `{"status":"COMPLETED","id":"export-task-1","result":{"file":%q}}`, file)
	})

	return httptest.NewServer(mux)
}

func TestUnitAwafPolicyCreateReadUpdateDelete(t *testing.T) {
	// Shrink the real 10-second post-import settle delay in
	// resourceBigipAwafPolicyCreate down to effectively nothing for this
	// test, restoring the original value afterward so other tests (and any
	// tests added later) aren't affected.
	originalDelay := awafPolicySettleDelay
	awafPolicySettleDelay = time.Millisecond
	t.Cleanup(func() { awafPolicySettleDelay = originalDelay })

	policyID := "abc-123"
	policyName := "test-policy"
	partition := "Common"

	server := awafPolicyMockServer(t, policyID, policyName, partition)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	client.Teem = true // skip real telemetry network call

	r := resourceBigipAwafPolicy()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":          policyName,
		"partition":     partition,
		"template_name": "POLICY_TEMPLATE_RAPID_DEPLOYMENT",
	}, "")

	ctx := context.Background()

	diags := resourceBigipAwafPolicyCreate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected create error: %v", diags)
	assert.Equal(t, policyID, d.Id())

	diags = resourceBigipAwafPolicyRead(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected read error: %v", diags)
	assert.Equal(t, policyName, d.Get("name"))
	assert.Equal(t, partition, d.Get("partition"))

	diags = resourceBigipAwafPolicyUpdate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected update error: %v", diags)

	diags = resourceBigipAwafPolicyDelete(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected delete error: %v", diags)
	assert.Equal(t, "", d.Id())
}

func TestUnitAwafPolicyCreate_NotProvisioned(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/sys/provision/asm", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"name":"asm","level":"none"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipAwafPolicy()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":          "test-policy",
		"partition":     "Common",
		"template_name": "POLICY_TEMPLATE_RAPID_DEPLOYMENT",
	}, "")

	diags := resourceBigipAwafPolicyCreate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitAwafPolicyCreate_ProvisionCheckError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/sys/provision/asm", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"internal error"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipAwafPolicy()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":          "test-policy",
		"partition":     "Common",
		"template_name": "POLICY_TEMPLATE_RAPID_DEPLOYMENT",
	}, "")

	diags := resourceBigipAwafPolicyCreate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitAwafPolicyCreate_ImportError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/sys/provision/asm", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"name":"asm","level":"nominal"}`)
	})
	mux.HandleFunc("/mgmt/tm/asm/file-transfer/uploads/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"remainingByteCount":0,"usedChunks":{},"totalByteCount":10,"localFilePath":"","temporaryFilePath":"","generation":0,"lastUpdateMicros":0}`)
	})
	mux.HandleFunc("/mgmt/tm/asm/tasks/import-policy", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"import failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipAwafPolicy()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":          "test-policy",
		"partition":     "Common",
		"template_name": "POLICY_TEMPLATE_RAPID_DEPLOYMENT",
	}, "")

	diags := resourceBigipAwafPolicyCreate(context.Background(), d, client)
	require.True(t, diags.HasError())

	// Explicitly verify the package-level mutex was released on this error
	// path. Before the defer mutex.Unlock() fix, this would not fail an
	// assertion -- it would hang the entire test binary (mutex.Lock() below
	// blocks forever), surfacing as a timeout rather than a test failure.
	require.True(t, mutex.TryLock(), "mutex must be released on the error path")
	mutex.Unlock()
}

func TestUnitAwafPolicyRead_Error(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/asm/policies/bad-id", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = fmt.Fprint(w, `{"code":404,"message":"not found"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipAwafPolicy()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":          "test-policy",
		"partition":     "Common",
		"template_name": "POLICY_TEMPLATE_RAPID_DEPLOYMENT",
	}, "bad-id")

	diags := resourceBigipAwafPolicyRead(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitAwafPolicyUpdate_Error(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/asm/file-transfer/uploads/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"remainingByteCount":0,"usedChunks":{},"totalByteCount":10,"localFilePath":"","temporaryFilePath":"","generation":0,"lastUpdateMicros":0}`)
	})
	mux.HandleFunc("/mgmt/tm/asm/tasks/import-policy", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"import failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipAwafPolicy()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":          "test-policy",
		"partition":     "Common",
		"template_name": "POLICY_TEMPLATE_RAPID_DEPLOYMENT",
	}, "policy-id-1")

	diags := resourceBigipAwafPolicyUpdate(context.Background(), d, client)
	require.True(t, diags.HasError())

	// See TestUnitAwafPolicyCreate_ImportError: verify the mutex was
	// actually released, not just that the test didn't hang.
	require.True(t, mutex.TryLock(), "mutex must be released on the error path")
	mutex.Unlock()
}

func TestUnitAwafPolicyDelete_Error(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/asm/policies/bad-id", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"delete failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipAwafPolicy()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":          "test-policy",
		"partition":     "Common",
		"template_name": "POLICY_TEMPLATE_RAPID_DEPLOYMENT",
	}, "bad-id")

	diags := resourceBigipAwafPolicyDelete(context.Background(), d, client)
	require.True(t, diags.HasError())
}
