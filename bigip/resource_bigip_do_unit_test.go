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
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------
// Schema-shape assertion (no HTTP server needed)
// ---------------------------------------------------------------------

func TestResourceBigipDoSchema(t *testing.T) {
	r := resourceBigipDo()

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

	doJsonSchema, ok := r.Schema["do_json"]
	if !ok {
		t.Fatal("Expected field 'do_json' to exist in schema")
	}
	if !doJsonSchema.Required {
		t.Error("Expected field 'do_json' to be required")
	}

	timeoutSchema, ok := r.Schema["timeout"]
	if !ok {
		t.Fatal("Expected field 'timeout' to exist in schema")
	}
	if timeoutSchema.Default != 20 {
		t.Errorf("Expected field 'timeout' default to be 20, got %v", timeoutSchema.Default)
	}
}

// ---------------------------------------------------------------------
// resourceBigipDoDelete (trivial no-op, no server needed)
// ---------------------------------------------------------------------

func TestUnitDoDelete(t *testing.T) {
	r := resourceBigipDo()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"do_json": `{"class":"Device"}`,
	}, "some-id")

	diags := resourceBigipDoDelete(context.Background(), d, nil)
	require.False(t, diags.HasError(), "unexpected delete error: %v", diags)
	assert.Equal(t, "", d.Id())
}

// ---------------------------------------------------------------------
// resourceBigipDoCreate / Read / Update
// ---------------------------------------------------------------------

const testDoJSON = `{"class":"Device","async":true}`

func doMockServer(t *testing.T, taskID string) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()
	addDoReadinessHandlers(mux)

	mux.HandleFunc("/mgmt/shared/declarative-onboarding/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprintf(w, `{"id":"%s","result":{"status":"OK"}}`, taskID)
		}
	})
	mux.HandleFunc("/mgmt/shared/declarative-onboarding/task/"+taskID, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{"id":"%s","result":{"status":"OK"},"declaration":{"class":"Device","async":true}}`, taskID)
	})

	return httptest.NewServer(mux)
}

func addDoReadinessHandlers(mux *http.ServeMux) {
	mux.HandleFunc("/mgmt/tm/net/self", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"items":[]}`)
	})
	mux.HandleFunc("/mgmt/shared/authn/login", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprint(w, `{"token":{"token":"unit-test-token"}}`)
	})
}

func TestUnitDoCreate_SyncSuccess(t *testing.T) {
	taskID := "task-sync-1"
	server := doMockServer(t, taskID)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	client.Teem = true // skip real telemetry network call

	r := resourceBigipDo()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"do_json": testDoJSON,
		"timeout": 1,
	}, "")

	diags := resourceBigipDoCreate(context.Background(), d, client)
	require.False(t, diags.HasError(), "unexpected create error: %v", diags)
	assert.Equal(t, taskID, d.Id())
}

func TestUnitDoCreate_AsyncSuccess(t *testing.T) {
	// The first poll returns an unexpected status (exercising the poll
	// loop's "default" branch), then succeeds on the second. This test
	// takes just over 1 real second due to the loop's unconditional
	// time.Sleep(1 * time.Second) between iterations.
	taskID := "task-async-1"
	var pollCount int

	mux := http.NewServeMux()
	addDoReadinessHandlers(mux)
	mux.HandleFunc("/mgmt/shared/declarative-onboarding/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
		_, _ = fmt.Fprintf(w, `{"id":"%s"}`, taskID)
	})
	mux.HandleFunc("/mgmt/shared/declarative-onboarding/task/"+taskID, func(w http.ResponseWriter, r *http.Request) {
		pollCount++
		if pollCount == 1 {
			w.WriteHeader(http.StatusServiceUnavailable)
			_, _ = fmt.Fprint(w, `{}`)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{"id":"%s","result":{"status":"OK"},"declaration":{"class":"Device"}}`, taskID)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	client.Teem = true

	r := resourceBigipDo()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"do_json": testDoJSON,
		"timeout": 1,
	}, "")

	diags := resourceBigipDoCreate(context.Background(), d, client)
	require.False(t, diags.HasError(), "unexpected create error: %v", diags)
	assert.GreaterOrEqual(t, pollCount, 2)
	assert.Equal(t, taskID, d.Id())
}

func TestUnitDoCreate_AsyncStillRunningThenTimeout(t *testing.T) {
	taskID := "task-running-1"

	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/shared/declarative-onboarding/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
		_, _ = fmt.Fprintf(w, `{"id":"%s"}`, taskID)
	})
	mux.HandleFunc("/mgmt/shared/declarative-onboarding/task/"+taskID, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
		_, _ = fmt.Fprint(w, `{"result":{"status":"RUNNING"}}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	client.Teem = true

	r := resourceBigipDo()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"do_json": testDoJSON,
		// timeout is in minutes and gets multiplied by 60 -- the smallest
		// value the schema's TypeInt allows here that still exercises the
		// full poll-until-timeout loop quickly is 0, which becomes a
		// 0-second timeout so the "for" loop body runs at most once or
		// twice before the outer time.Since check fails.
		"timeout": 0,
	}, "")

	diags := resourceBigipDoCreate(context.Background(), d, client)
	require.True(t, diags.HasError(), "expected a timeout error since status never reaches RUNNING->done")
	assert.Equal(t, "", d.Id())
}

func TestUnitDoCreate_AsyncNonRunningStatus(t *testing.T) {
	taskID := "task-failed-1"

	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/shared/declarative-onboarding/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
		_, _ = fmt.Fprintf(w, `{"id":"%s"}`, taskID)
	})
	mux.HandleFunc("/mgmt/shared/declarative-onboarding/task/"+taskID, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
		_, _ = fmt.Fprint(w, `{"result":{"status":"FAILED","message":"boom"}}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	client.Teem = true

	r := resourceBigipDo()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"do_json": testDoJSON,
		"timeout": 1,
	}, "")

	diags := resourceBigipDoCreate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitDoCreate_PostError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/shared/declarative-onboarding/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"post failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	client.Teem = true // skip real telemetry network call
	r := resourceBigipDo()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"do_json": testDoJSON,
		"timeout": 1,
	}, "")

	diags := resourceBigipDoCreate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitDoRead_Success(t *testing.T) {
	taskID := "task-read-1"
	server := doMockServer(t, taskID)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipDo()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"do_json": testDoJSON,
	}, taskID)

	diags := resourceBigipDoRead(context.Background(), d, client)
	require.False(t, diags.HasError(), "unexpected read error: %v", diags)
	assert.Contains(t, d.Get("do_json").(string), "Device")
}

func TestUnitDoRead_WithConnectBigIP(t *testing.T) {
	// Exercises the connectBigIP branch in Read: when bigip_address/user/
	// password are all set, Read reconnects instead of using the
	// provider-configured client. Unlike Create, Read never sends TEEM
	// telemetry, so this is safe to run without a real network dependency.
	taskID := "task-read-connect-1"

	mux := http.NewServeMux()
	addDoReadinessHandlers(mux)
	mux.HandleFunc("/mgmt/shared/declarative-onboarding/task/"+taskID, func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"id":"%s","declaration":{"class":"Device"}}`, taskID)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	address, port := splitTestServerURL(t, server.URL)

	outerClient := NewUnitTestClient("http://127.0.0.1:0")

	r := resourceBigipDo()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"do_json":        testDoJSON,
		"bigip_address":  address,
		"bigip_port":     port,
		"bigip_user":     "admin",
		"bigip_password": "admin",
	}, taskID)

	diags := resourceBigipDoRead(context.Background(), d, outerClient)
	require.False(t, diags.HasError(), "unexpected read error: %v", diags)
	assert.Contains(t, d.Get("do_json").(string), "Device")
}

func TestUnitDoRead_ConnectBigIPError(t *testing.T) {
	// connectBigIP retries transport-level failures (see its own comment);
	// force a single attempt so this deliberately-unreachable address
	// fails fast instead of exhausting the default retry/backoff budget.
	t.Setenv("API_RETRIES", "1")
	outerClient := NewUnitTestClient("http://127.0.0.1:0")

	r := resourceBigipDo()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"do_json":        testDoJSON,
		"bigip_address":  "127.0.0.1",
		"bigip_port":     "0",
		"bigip_user":     "admin",
		"bigip_password": "admin",
	}, "some-id")

	diags := resourceBigipDoRead(context.Background(), d, outerClient)
	require.True(t, diags.HasError())
}

func TestUnitDoRead_Error(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/shared/declarative-onboarding/task/missing-task", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = fmt.Fprint(w, `{"code":404,"message":"not found"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipDo()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"do_json": testDoJSON,
	}, "missing-task")

	diags := resourceBigipDoRead(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitDoUpdate_SyncSuccess(t *testing.T) {
	taskID := "task-update-1"
	server := doMockServer(t, taskID)
	defer server.Close()

	client := NewUnitTestClient(server.URL)

	r := resourceBigipDo()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"do_json": testDoJSON,
		"timeout": 1,
	}, taskID)

	diags := resourceBigipDoUpdate(context.Background(), d, client)
	require.False(t, diags.HasError(), "unexpected update error: %v", diags)
}

func TestUnitDoUpdate_WithConnectBigIP(t *testing.T) {
	// Exercises the connectBigIP branch in Update. Unlike Create, Update
	// never sends TEEM telemetry, so this is safe to run without a real
	// network dependency.
	taskID := "task-update-connect-1"

	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/net/self", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"items":[]}`)
	})
	mux.HandleFunc("/mgmt/shared/authn/login", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprint(w, `{"token":{"token":"unit-test-token"}}`)
	})
	mux.HandleFunc("/mgmt/shared/declarative-onboarding/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{"id":"%s","result":{"status":"OK"}}`, taskID)
	})
	mux.HandleFunc("/mgmt/shared/declarative-onboarding/task/"+taskID, func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"id":"%s","result":{"status":"OK"},"declaration":{"class":"Device"}}`, taskID)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	address, port := splitTestServerURL(t, server.URL)

	outerClient := NewUnitTestClient("http://127.0.0.1:0")

	r := resourceBigipDo()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"do_json":        testDoJSON,
		"timeout":        1,
		"bigip_address":  address,
		"bigip_port":     port,
		"bigip_user":     "admin",
		"bigip_password": "admin",
	}, "")

	diags := resourceBigipDoUpdate(context.Background(), d, outerClient)
	require.False(t, diags.HasError(), "unexpected update error: %v", diags)
	assert.Equal(t, taskID, d.Id())
}

func TestUnitDoUpdate_AsyncUnexpectedStatusThenSuccess(t *testing.T) {
	// Exercises the poll loop's "default" branch (an HTTP status other than
	// 200/202) on the first poll, then succeeds on the second. This test
	// takes just over 1 real second due to the loop's unconditional
	// time.Sleep(1 * time.Second) between iterations.
	taskID := "task-update-unexpected-1"
	var pollCount int

	mux := http.NewServeMux()
	addDoReadinessHandlers(mux)
	mux.HandleFunc("/mgmt/shared/declarative-onboarding/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
		_, _ = fmt.Fprintf(w, `{"id":"%s"}`, taskID)
	})
	mux.HandleFunc("/mgmt/shared/declarative-onboarding/task/"+taskID, func(w http.ResponseWriter, r *http.Request) {
		pollCount++
		if pollCount == 1 {
			w.WriteHeader(http.StatusServiceUnavailable)
			_, _ = fmt.Fprint(w, `{}`)
			return
		}
		_, _ = fmt.Fprintf(w, `{"id":"%s","result":{"status":"OK"},"declaration":{"class":"Device"}}`, taskID)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)

	r := resourceBigipDo()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"do_json": testDoJSON,
		"timeout": 1,
	}, "")

	diags := resourceBigipDoUpdate(context.Background(), d, client)
	require.False(t, diags.HasError(), "unexpected update error: %v", diags)
	assert.Equal(t, taskID, d.Id())
	assert.GreaterOrEqual(t, pollCount, 2)
}

func TestUnitDoUpdate_ConnectBigIPError(t *testing.T) {
	// See the comment in TestUnitDoRead_ConnectBigIPError.
	t.Setenv("API_RETRIES", "1")
	outerClient := NewUnitTestClient("http://127.0.0.1:0")

	r := resourceBigipDo()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"do_json":        testDoJSON,
		"bigip_address":  "127.0.0.1",
		"bigip_port":     "0",
		"bigip_user":     "admin",
		"bigip_password": "admin",
	}, "")

	diags := resourceBigipDoUpdate(context.Background(), d, outerClient)
	require.True(t, diags.HasError())
}

func TestUnitDoUpdate_PostError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/shared/declarative-onboarding/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"update post failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipDo()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"do_json": testDoJSON,
		"timeout": 1,
	}, "some-id")

	diags := resourceBigipDoUpdate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitDoUpdate_AsyncNonRunningStatus(t *testing.T) {
	taskID := "task-update-failed-1"

	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/shared/declarative-onboarding/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
		_, _ = fmt.Fprintf(w, `{"id":"%s"}`, taskID)
	})
	mux.HandleFunc("/mgmt/shared/declarative-onboarding/task/"+taskID, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
		_, _ = fmt.Fprint(w, `{"result":{"status":"FAILED","message":"boom"}}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)

	r := resourceBigipDo()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"do_json": testDoJSON,
		"timeout": 1,
	}, "")

	diags := resourceBigipDoUpdate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitDoUpdate_AsyncSuccess(t *testing.T) {
	taskID := "task-update-async-1"

	mux := http.NewServeMux()
	addDoReadinessHandlers(mux)
	mux.HandleFunc("/mgmt/shared/declarative-onboarding/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
		_, _ = fmt.Fprintf(w, `{"id":"%s"}`, taskID)
	})
	mux.HandleFunc("/mgmt/shared/declarative-onboarding/task/"+taskID, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{"id":"%s","result":{"status":"OK"},"declaration":{"class":"Device"}}`, taskID)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)

	r := resourceBigipDo()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"do_json": testDoJSON,
		"timeout": 1,
	}, "")

	diags := resourceBigipDoUpdate(context.Background(), d, client)
	require.False(t, diags.HasError(), "unexpected update error: %v", diags)
	assert.Equal(t, taskID, d.Id())
}

func TestUnitDoUpdate_AsyncStillRunningThenTimeout(t *testing.T) {
	taskID := "task-update-running-1"

	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/shared/declarative-onboarding/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
		_, _ = fmt.Fprintf(w, `{"id":"%s"}`, taskID)
	})
	mux.HandleFunc("/mgmt/shared/declarative-onboarding/task/"+taskID, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
		_, _ = fmt.Fprint(w, `{"result":{"status":"RUNNING"}}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)

	r := resourceBigipDo()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"do_json": testDoJSON,
		"timeout": 0,
	}, "")

	diags := resourceBigipDoUpdate(context.Background(), d, client)
	require.True(t, diags.HasError(), "expected a timeout error since status never leaves RUNNING")
	assert.Equal(t, "", d.Id())
}

func TestUnitDoCreate_AsyncMalformedTaskJSON(t *testing.T) {
	// The task-poll GET succeeds with a 200 but a body that isn't valid
	// JSON, so json.Unmarshal in the "case taskResp.StatusCode == 200"
	// branch errors -- this surfaces as a hard error from Create, not a
	// timeout.
	taskID := "task-malformed-1"

	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/shared/declarative-onboarding/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
		_, _ = fmt.Fprintf(w, `{"id":"%s"}`, taskID)
	})
	mux.HandleFunc("/mgmt/shared/declarative-onboarding/task/"+taskID, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprint(w, `not-json`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	client.Teem = true

	r := resourceBigipDo()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"do_json": testDoJSON,
		"timeout": 1,
	}, "")

	diags := resourceBigipDoCreate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitDoCreate_FinalTimeoutMalformedJSON(t *testing.T) {
	// timeout=0 means timeoutSec is 0, so the "for time.Since(start).Seconds()
	// < 0" poll loop body never executes at all (its condition is false
	// immediately) -- Create falls straight through to the standalone
	// "!doSuccess" final GET made after the loop. If that final GET's body
	// isn't valid JSON, the resulting error comes from that final-check
	// json.Unmarshal rather than the generic timeout message.
	taskID := "task-final-malformed-1"

	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/shared/declarative-onboarding/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
		_, _ = fmt.Fprintf(w, `{"id":"%s"}`, taskID)
	})

	var pollCount int
	mux.HandleFunc("/mgmt/shared/declarative-onboarding/task/"+taskID, func(w http.ResponseWriter, r *http.Request) {
		pollCount++
		w.WriteHeader(http.StatusAccepted)
		_, _ = fmt.Fprint(w, `not-json`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	client.Teem = true

	r := resourceBigipDo()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"do_json": testDoJSON,
		"timeout": 0,
	}, "")

	diags := resourceBigipDoCreate(context.Background(), d, client)
	require.True(t, diags.HasError())
	assert.GreaterOrEqual(t, pollCount, 1)
}

func TestUnitDoCreate_WithTokenAuth(t *testing.T) {
	// Async (202) so Create actually enters the task-polling loop, exercising
	// the "clientBigip.Token != \"\"" branch (X-F5-Auth-Token header) on both
	// the initial POST and the poll GET, instead of the default Basic Auth
	// branch used by the other tests. A 200 response here would set doSuccess
	// immediately and skip polling entirely, leaving the poll GET's auth
	// branch untested.
	taskID := "task-token-1"
	var sawPostAuthHeader, sawPollAuthHeader string

	mux := http.NewServeMux()
	addDoReadinessHandlers(mux)
	mux.HandleFunc("/mgmt/shared/declarative-onboarding/", func(w http.ResponseWriter, r *http.Request) {
		sawPostAuthHeader = r.Header.Get("X-F5-Auth-Token")
		w.WriteHeader(http.StatusAccepted)
		_, _ = fmt.Fprintf(w, `{"id":"%s"}`, taskID)
	})
	mux.HandleFunc("/mgmt/shared/declarative-onboarding/task/"+taskID, func(w http.ResponseWriter, r *http.Request) {
		sawPollAuthHeader = r.Header.Get("X-F5-Auth-Token")
		_, _ = fmt.Fprintf(w, `{"id":"%s","result":{"status":"OK"},"declaration":{"class":"Device"}}`, taskID)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	client.Teem = true
	client.Token = "test-token-value"

	r := resourceBigipDo()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"do_json": testDoJSON,
		"timeout": 1,
	}, "")

	diags := resourceBigipDoCreate(context.Background(), d, client)
	require.False(t, diags.HasError(), "unexpected create error: %v", diags)
	assert.Equal(t, "test-token-value", sawPostAuthHeader)
	assert.Equal(t, "test-token-value", sawPollAuthHeader)
}

func TestUnitDoCreate_WithBasicAuth(t *testing.T) {
	// Mirrors TestUnitDoCreate_WithTokenAuth but with clientBigip.Token left
	// empty, exercising the "else { req.SetBasicAuth(...) }" branch on both
	// the initial POST and the async task-polling GET.
	taskID := "task-basic-auth-1"
	var sawPostUser, sawPostPass string
	var sawPollUser, sawPollPass string
	var sawPostOK, sawPollOK bool

	mux := http.NewServeMux()
	addDoReadinessHandlers(mux)
	mux.HandleFunc("/mgmt/shared/declarative-onboarding/", func(w http.ResponseWriter, r *http.Request) {
		sawPostUser, sawPostPass, sawPostOK = r.BasicAuth()
		w.WriteHeader(http.StatusAccepted)
		_, _ = fmt.Fprintf(w, `{"id":"%s"}`, taskID)
	})
	mux.HandleFunc("/mgmt/shared/declarative-onboarding/task/"+taskID, func(w http.ResponseWriter, r *http.Request) {
		sawPollUser, sawPollPass, sawPollOK = r.BasicAuth()
		_, _ = fmt.Fprintf(w, `{"id":"%s","result":{"status":"OK"},"declaration":{"class":"Device"}}`, taskID)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	client.Teem = true
	client.User = "admin"
	client.Password = "admin"

	r := resourceBigipDo()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"do_json": testDoJSON,
		"timeout": 1,
	}, "")

	diags := resourceBigipDoCreate(context.Background(), d, client)
	require.False(t, diags.HasError(), "unexpected create error: %v", diags)
	assert.True(t, sawPostOK, "expected Basic Auth on the initial POST")
	assert.True(t, sawPollOK, "expected Basic Auth on the task-polling GET")
	assert.Equal(t, "admin", sawPostUser)
	assert.Equal(t, "admin", sawPostPass)
	assert.Equal(t, "admin", sawPollUser)
	assert.Equal(t, "admin", sawPollPass)
}

func TestUnitDoUpdate_AsyncMalformedTaskJSON(t *testing.T) {
	taskID := "task-update-malformed-1"

	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/shared/declarative-onboarding/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
		_, _ = fmt.Fprintf(w, `{"id":"%s"}`, taskID)
	})
	mux.HandleFunc("/mgmt/shared/declarative-onboarding/task/"+taskID, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprint(w, `not-json`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)

	r := resourceBigipDo()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"do_json": testDoJSON,
		"timeout": 1,
	}, "")

	diags := resourceBigipDoUpdate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitDoUpdate_AsyncNonRunningMalformedJSON(t *testing.T) {
	// timeout=1 (60 real seconds) rather than 0 so the "for
	// time.Since(start).Seconds() < timeoutSec" loop body actually executes
	// at least once; the mock's malformed JSON causes that first iteration
	// to return immediately via the case-202 json.Unmarshal error, well
	// before the 60s ceiling would matter.
	taskID := "task-update-202-malformed-1"

	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/shared/declarative-onboarding/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
		_, _ = fmt.Fprintf(w, `{"id":"%s"}`, taskID)
	})
	mux.HandleFunc("/mgmt/shared/declarative-onboarding/task/"+taskID, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
		_, _ = fmt.Fprint(w, `not-json`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)

	r := resourceBigipDo()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"do_json": testDoJSON,
		"timeout": 1,
	}, "")

	diags := resourceBigipDoUpdate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitDoUpdate_WithTokenAuth(t *testing.T) {
	// Async (202) so the update also exercises the poll loop's own
	// "clientBigip.Token != \"\"" header branch, not just the initial POST's.
	taskID := "task-update-token-1"
	var sawPostAuthHeader, sawPollAuthHeader string

	mux := http.NewServeMux()
	addDoReadinessHandlers(mux)
	mux.HandleFunc("/mgmt/shared/declarative-onboarding/", func(w http.ResponseWriter, r *http.Request) {
		sawPostAuthHeader = r.Header.Get("X-F5-Auth-Token")
		w.WriteHeader(http.StatusAccepted)
		_, _ = fmt.Fprintf(w, `{"id":"%s"}`, taskID)
	})
	mux.HandleFunc("/mgmt/shared/declarative-onboarding/task/"+taskID, func(w http.ResponseWriter, r *http.Request) {
		sawPollAuthHeader = r.Header.Get("X-F5-Auth-Token")
		_, _ = fmt.Fprintf(w, `{"id":"%s","result":{"status":"OK"},"declaration":{"class":"Device"}}`, taskID)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	client.Token = "test-token-value"

	r := resourceBigipDo()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"do_json": testDoJSON,
		"timeout": 1,
	}, "")

	diags := resourceBigipDoUpdate(context.Background(), d, client)
	require.False(t, diags.HasError(), "unexpected update error: %v", diags)
	assert.Equal(t, "test-token-value", sawPostAuthHeader)
	assert.Equal(t, "test-token-value", sawPollAuthHeader)
}

func TestUnitDoRead_WithTokenAuth(t *testing.T) {
	taskID := "task-read-token-1"
	var sawAuthHeader string

	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/shared/declarative-onboarding/task/"+taskID, func(w http.ResponseWriter, r *http.Request) {
		sawAuthHeader = r.Header.Get("X-F5-Auth-Token")
		_, _ = fmt.Fprintf(w, `{"id":"%s","declaration":{"class":"Device"}}`, taskID)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	client.Token = "test-token-value"

	r := resourceBigipDo()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"do_json": testDoJSON,
	}, taskID)

	diags := resourceBigipDoRead(context.Background(), d, client)
	require.False(t, diags.HasError(), "unexpected read error: %v", diags)
	assert.Equal(t, "test-token-value", sawAuthHeader)
}

// ---------------------------------------------------------------------
// connectBigIP
// ---------------------------------------------------------------------

// splitTestServerURL breaks an httptest.Server URL (e.g.
// "http://127.0.0.1:54321") into a bare "http://127.0.0.1" address and its
// port, matching how connectBigIP/bigip.NewSession assemble
// "<address>:<port>" -- passing the full server.URL as bigip_address would
// double up the port (":54321:443").
func splitTestServerURL(t *testing.T, rawURL string) (address, port string) {
	t.Helper()
	u, err := url.Parse(rawURL)
	require.NoError(t, err)
	return u.Scheme + "://" + u.Hostname(), u.Port()
}

func TestConnectBigIP_Success(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/net/self", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"items":[]}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	address, port := splitTestServerURL(t, server.URL)

	r := resourceBigipDo()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"do_json":        testDoJSON,
		"bigip_address":  address,
		"bigip_port":     port,
		"bigip_user":     "admin",
		"bigip_password": "admin",
	}, "")

	client, err := connectBigIP(d)
	require.NoError(t, err)
	assert.NotNil(t, client)
}

func TestConnectBigIP_ValidateConnectionError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/net/self", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"internal error"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	address, port := splitTestServerURL(t, server.URL)

	r := resourceBigipDo()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"do_json":        testDoJSON,
		"bigip_address":  address,
		"bigip_port":     port,
		"bigip_user":     "admin",
		"bigip_password": "admin",
	}, "")

	_, err := connectBigIP(d)
	require.Error(t, err)
}

// Note: a test exercising the successful connectBigIP-then-post path in
// Create was intentionally omitted. connectBigIP always returns a fresh
// client with Teem defaulting to false (there's no way to pre-seed it,
// unlike the provider-configured client which this resource's other tests
// set Teem=true on directly), so any Create call that reconnects via
// connectBigIP and then succeeds would fall through to the real TEEM
// telemetry call against product.apis.f5.com -- a real network dependency
// this test suite otherwise avoids. TestConnectBigIP_Success covers
// connectBigIP's own success path directly instead.

func TestUnitDoCreate_ConnectBigIPError(t *testing.T) {
	// See the comment in TestUnitDoRead_ConnectBigIPError.
	t.Setenv("API_RETRIES", "1")
	outerClient := NewUnitTestClient("http://127.0.0.1:0")

	r := resourceBigipDo()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"do_json":        testDoJSON,
		"bigip_address":  "127.0.0.1",
		"bigip_port":     "0",
		"bigip_user":     "admin",
		"bigip_password": "admin",
	}, "")

	diags := resourceBigipDoCreate(context.Background(), d, outerClient)
	require.True(t, diags.HasError())
}
