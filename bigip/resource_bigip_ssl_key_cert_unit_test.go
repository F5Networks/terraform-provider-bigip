/*
Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------
// Schema-shape assertion (no HTTP server needed)
// ---------------------------------------------------------------------

func TestResourceBigipSSLKeyCertSchema(t *testing.T) {
	r := resourceBigipSSLKeyCert()

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

	keyNameSchema, ok := r.Schema["key_name"]
	if !ok {
		t.Fatal("Expected field 'key_name' to exist in schema")
	}
	if !keyNameSchema.Required {
		t.Error("Expected field 'key_name' to be required")
	}
}

// ---------------------------------------------------------------------
// resourceBigipSSLKeyCertCreate / Update mutex-release-on-error behavior
//
// These target the same class of bug fixed in resource_bigip_awaf_policy.go:
// mutex is a package-level sync.Mutex shared between this file and
// resource_bigip_awaf_policy.go. Before the defer mutex.Unlock() fix, an
// error between StartTransaction and CommitTransaction (or, for the AWAF
// resource, between ImportAwafJson and the final Unlock) would leave the
// mutex locked forever, hanging every subsequent Create/Update call across
// *both* resource types -- surfacing as a test-suite timeout, not a failed
// assertion. These tests make that release explicit and independently
// verifiable via mutex.TryLock().
// ---------------------------------------------------------------------

func TestUnitSSLKeyCertCreate_AddKeyError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/shared/file-transfer/uploads/testkey", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"remainingByteCount":0,"usedChunks":{},"totalByteCount":10,"localFilePath":"","temporaryFilePath":"","generation":0,"lastUpdateMicros":0}`)
	})
	mux.HandleFunc("/mgmt/tm/transaction", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"transId":1,"state":"STARTED"}`)
	})
	mux.HandleFunc("/mgmt/tm/sys/file/ssl-key", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"add key failed"}`)
	})
	mux.HandleFunc("/mgmt/tm/transaction/1", func(w http.ResponseWriter, r *http.Request) {
		// The rollback PATCH ({"state":"ROLLED_BACK"}) lands here too;
		// respond OK so it doesn't itself produce a spurious log error.
		_, _ = fmt.Fprint(w, `{"transId":1,"state":"ROLLED_BACK"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipSSLKeyCert()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"key_name":     "testkey",
		"key_content":  "fake-key-content",
		"cert_name":    "testcert",
		"cert_content": "fake-cert-content",
		"partition":    "Common",
	}, "")

	diags := resourceBigipSSLKeyCertCreate(context.Background(), d, client)
	require.True(t, diags.HasError())

	// Explicitly verify the package-level mutex was released on this error
	// path, rather than relying on the test merely not hanging.
	require.True(t, mutex.TryLock(), "mutex must be released on the error path")
	mutex.Unlock()
}

func TestUnitSSLKeyCertUpdate_ModifyKeyError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/shared/file-transfer/uploads/testkey", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"remainingByteCount":0,"usedChunks":{},"totalByteCount":10,"localFilePath":"","temporaryFilePath":"","generation":0,"lastUpdateMicros":0}`)
	})
	mux.HandleFunc("/mgmt/tm/transaction", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"transId":2,"state":"STARTED"}`)
	})
	mux.HandleFunc("/mgmt/tm/sys/file/ssl-key/~Common~testkey", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"modify key failed"}`)
	})
	mux.HandleFunc("/mgmt/tm/transaction/2", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"transId":2,"state":"ROLLED_BACK"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipSSLKeyCert()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"key_name":     "testkey",
		"key_content":  "fake-key-content",
		"cert_name":    "testcert",
		"cert_content": "fake-cert-content",
		"partition":    "Common",
	}, "testkey_testcert")

	diags := resourceBigipSSLKeyCertUpdate(context.Background(), d, client)
	require.True(t, diags.HasError())

	require.True(t, mutex.TryLock(), "mutex must be released on the error path")
	mutex.Unlock()
}

func TestUnitSSLKeyCertCreate_StartTransactionError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/shared/file-transfer/uploads/testkey", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"remainingByteCount":0,"usedChunks":{},"totalByteCount":10,"localFilePath":"","temporaryFilePath":"","generation":0,"lastUpdateMicros":0}`)
	})
	mux.HandleFunc("/mgmt/tm/transaction", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"start transaction failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipSSLKeyCert()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"key_name":     "testkey",
		"key_content":  "fake-key-content",
		"cert_name":    "testcert",
		"cert_content": "fake-cert-content",
		"partition":    "Common",
	}, "")

	// StartTransaction fails before any Add/Modify/Upload call, so there is
	// nothing to roll back -- this exercises the "no rollback needed" early
	// error branch, still under the scoped-lock closure.
	diags := resourceBigipSSLKeyCertCreate(context.Background(), d, client)
	require.True(t, diags.HasError())

	require.True(t, mutex.TryLock(), "mutex must be released on the error path")
	mutex.Unlock()
}

func TestUnitSSLKeyCertCreate_UploadKeyError(t *testing.T) {
	// UploadKey is called before mutex.Lock(), so an error here never
	// touches the mutex at all -- included for completeness alongside the
	// other Create error-path tests in this file.
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/shared/file-transfer/uploads/testkey", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"upload failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipSSLKeyCert()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"key_name":     "testkey",
		"key_content":  "fake-key-content",
		"cert_name":    "testcert",
		"cert_content": "fake-cert-content",
		"partition":    "Common",
	}, "")

	diags := resourceBigipSSLKeyCertCreate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

// ---------------------------------------------------------------------
// rollbackTransaction
// ---------------------------------------------------------------------

func TestRollbackTransaction_Success(t *testing.T) {
	var gotBody string
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/transaction/42", func(w http.ResponseWriter, r *http.Request) {
		AssertRequestMethod(t, r, http.MethodPatch)
		body, _ := io.ReadAll(r.Body)
		gotBody = string(body)
		_, _ = fmt.Fprint(w, `{"transId":42,"state":"ROLLED_BACK"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	rollbackTransaction(client, 42)

	assert.Contains(t, gotBody, `"state":"ROLLED_BACK"`)
}

func TestRollbackTransaction_RequestErrorIsLoggedNotPanicked(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/transaction/99", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"rollback failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	// Must not panic or return a value (rollbackTransaction has no return);
	// the request failure is logged and swallowed by design.
	rollbackTransaction(client, 99)
}

// ---------------------------------------------------------------------
// Happy-path CRUD lifecycle
//
// Successful Create/Read/Update/Delete, plus related Read/Delete error
// paths and the fqdn helper. Complements the mutex-release-on-error tests
// above with the resource's standard (non-error) lifecycle.
// ---------------------------------------------------------------------

func TestUnitSSLKeyCertCreateReadUpdateDelete(t *testing.T) {
	keyName := "testkey"
	certName := "testcert"
	keyMangled := "/mgmt/tm/sys/file/ssl-key/" + MangleFullPath("/Common/"+keyName)
	certMangled := "/mgmt/tm/sys/file/ssl-cert/" + MangleFullPath("/Common/"+certName)

	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/shared/file-transfer/uploads/"+keyName, func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"remainingByteCount":0,"usedChunks":{},"totalByteCount":10,"localFilePath":"","temporaryFilePath":"","generation":0,"lastUpdateMicros":0}`)
	})
	mux.HandleFunc("/mgmt/shared/file-transfer/uploads/"+certName, func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"remainingByteCount":0,"usedChunks":{},"totalByteCount":10,"localFilePath":"","temporaryFilePath":"","generation":0,"lastUpdateMicros":0}`)
	})
	mux.HandleFunc("/mgmt/tm/transaction", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"transId":7,"state":"STARTED"}`)
	})
	mux.HandleFunc("/mgmt/tm/transaction/7", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"transId":7,"state":"VALIDATING"}`)
	})
	mux.HandleFunc("/mgmt/tm/sys/file/ssl-key", func(w http.ResponseWriter, r *http.Request) {
		AssertRequestMethod(t, r, http.MethodPost)
		_, _ = fmt.Fprintf(w, `{"name":"%s","partition":"Common","fullPath":"/Common/%s"}`, keyName, keyName)
	})
	mux.HandleFunc("/mgmt/tm/sys/file/ssl-cert", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"name":"%s","partition":"Common","fullPath":"/Common/%s"}`, certName, certName)
	})
	keyHandler := func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"name":"%s","partition":"Common","fullPath":"/Common/%s"}`, keyName, keyName)
	}
	mux.HandleFunc(keyMangled, keyHandler)
	certHandler := func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodDelete:
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprint(w, `{}`)
		default:
			_, _ = fmt.Fprintf(w, `{"name":"%s","partition":"Common","fullPath":"/Common/%s","issuerCert":"/Common/issuer","certValidationOptions":["certificate"]}`, certName, certName)
		}
	}
	mux.HandleFunc(certMangled, certHandler)
	mux.HandleFunc(certMangled+"/", certHandler)
	mux.HandleFunc(keyMangled+"/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodDelete {
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprint(w, `{}`)
		}
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipSSLKeyCert()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"key_name":             keyName,
		"key_content":          "fake-key-content",
		"cert_name":            certName,
		"cert_content":         "fake-cert-content",
		"partition":            "Common",
		"cert_monitoring_type": "certificate",
		"issuer_cert":          "/Common/issuer",
		"cert_ocsp":            "/Common/my-ocsp",
	}, "")

	ctx := context.Background()

	diags := resourceBigipSSLKeyCertCreate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected create error: %v", diags)
	require.Equal(t, keyName+"_"+certName, d.Id())

	diags = resourceBigipSSLKeyCertRead(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected read error: %v", diags)
	require.Equal(t, "certificate", d.Get("cert_monitoring_type").(string))

	diags = resourceBigipSSLKeyCertUpdate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected update error: %v", diags)

	diags = resourceBigipSSLKeyCertDelete(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected delete error: %v", diags)
	require.Equal(t, "", d.Id())
}

func TestUnitSSLKeyCertRead_KeyNotFound(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipSSLKeyCert()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"key_name":     "testkey",
		"key_content":  "fake-key-content",
		"cert_name":    "testcert",
		"cert_content": "fake-cert-content",
		"partition":    "Common",
	}, "testkey_testcert")

	diags := resourceBigipSSLKeyCertRead(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitSSLKeyCertDelete_LogsErrorsButSucceeds(t *testing.T) {
	// DeleteKey/DeleteCertificate errors are logged (not returned) inside
	// Delete, so the transaction still commits and Delete itself succeeds.
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/transaction", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"transId":9,"state":"STARTED"}`)
	})
	mux.HandleFunc("/mgmt/tm/transaction/9", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"transId":9,"state":"VALIDATING"}`)
	})
	// No handlers registered for the key/cert delete endpoints themselves,
	// so those requests 404 and DeleteKey/DeleteCertificate return errors
	// that Delete logs and swallows.
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipSSLKeyCert()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"key_name":     "testkey",
		"key_content":  "fake-key-content",
		"cert_name":    "testcert",
		"cert_content": "fake-cert-content",
		"partition":    "Common",
	}, "testkey_testcert")

	diags := resourceBigipSSLKeyCertDelete(context.Background(), d, client)
	require.False(t, diags.HasError(), "unexpected delete error: %v", diags)
	require.Equal(t, "", d.Id())
}

func TestUnitSSLKeyCertDelete_StartTransactionError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/transaction", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"start transaction failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipSSLKeyCert()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"key_name":     "testkey",
		"key_content":  "fake-key-content",
		"cert_name":    "testcert",
		"cert_content": "fake-cert-content",
		"partition":    "Common",
	}, "testkey_testcert")

	diags := resourceBigipSSLKeyCertDelete(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitSSLKeyCertDelete_CommitTransactionError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/transaction", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"transId":11,"state":"STARTED"}`)
	})
	mux.HandleFunc("/mgmt/tm/transaction/11", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"commit failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipSSLKeyCert()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"key_name":     "testkey",
		"key_content":  "fake-key-content",
		"cert_name":    "testcert",
		"cert_content": "fake-cert-content",
		"partition":    "Common",
	}, "testkey_testcert")

	diags := resourceBigipSSLKeyCertDelete(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitSSLKeyCertUpdate_CommitTransactionError(t *testing.T) {
	keyName := "testkey"
	certName := "testcert"
	keyMangled := "/mgmt/tm/sys/file/ssl-key/" + MangleFullPath("/Common/"+keyName)

	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/shared/file-transfer/uploads/"+keyName, func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"remainingByteCount":0,"usedChunks":{},"totalByteCount":10,"localFilePath":"","temporaryFilePath":"","generation":0,"lastUpdateMicros":0}`)
	})
	mux.HandleFunc("/mgmt/shared/file-transfer/uploads/"+certName, func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"remainingByteCount":0,"usedChunks":{},"totalByteCount":10,"localFilePath":"","temporaryFilePath":"","generation":0,"lastUpdateMicros":0}`)
	})
	mux.HandleFunc("/mgmt/tm/transaction", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"transId":13,"state":"STARTED"}`)
	})
	mux.HandleFunc("/mgmt/tm/transaction/13", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"commit failed"}`)
	})
	mux.HandleFunc(keyMangled, func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"name":"%s","partition":"Common","fullPath":"/Common/%s"}`, keyName, keyName)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipSSLKeyCert()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"key_name":     keyName,
		"key_content":  "fake-key-content",
		"cert_name":    certName,
		"cert_content": "fake-cert-content",
		"partition":    "Common",
	}, keyName+"_"+certName)

	diags := resourceBigipSSLKeyCertUpdate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitSSLKeyCertCreate_CommitTransactionError(t *testing.T) {
	keyName := "testkey"
	certName := "testcert"

	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/shared/file-transfer/uploads/"+keyName, func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"remainingByteCount":0,"usedChunks":{},"totalByteCount":10,"localFilePath":"","temporaryFilePath":"","generation":0,"lastUpdateMicros":0}`)
	})
	mux.HandleFunc("/mgmt/shared/file-transfer/uploads/"+certName, func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"remainingByteCount":0,"usedChunks":{},"totalByteCount":10,"localFilePath":"","temporaryFilePath":"","generation":0,"lastUpdateMicros":0}`)
	})
	mux.HandleFunc("/mgmt/tm/transaction", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"transId":15,"state":"STARTED"}`)
	})
	mux.HandleFunc("/mgmt/tm/transaction/15", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"commit failed"}`)
	})
	mux.HandleFunc("/mgmt/tm/sys/file/ssl-key", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"name":"%s","partition":"Common","fullPath":"/Common/%s"}`, keyName, keyName)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipSSLKeyCert()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"key_name":     keyName,
		"key_content":  "fake-key-content",
		"cert_name":    certName,
		"cert_content": "fake-cert-content",
		"partition":    "Common",
	}, "")

	diags := resourceBigipSSLKeyCertCreate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

// ---------------------------------------------------------------------
// fqdn helper
// ---------------------------------------------------------------------

func TestFqdn(t *testing.T) {
	require.Equal(t, "/Common/testkey", fqdn("Common", "testkey"))
	require.Equal(t, "/Common/testkey", fqdn("Common", "/Common/testkey"))
	require.Equal(t, "testkey", fqdn("", "testkey"))
	require.Equal(t, "/Common/testkey", fqdn("", "/Common/testkey"))
}
