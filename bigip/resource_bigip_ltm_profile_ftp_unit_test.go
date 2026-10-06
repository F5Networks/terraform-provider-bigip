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

	"github.com/stretchr/testify/require"
)

// versionMockHandler returns an httptest handler for /mgmt/tm/cli/version
// that reports the given "active" version description, used to drive the
// resourceBigipLtmProfileFtp version-branch logic (regex ^(12)|(13).*).
func versionMockHandler(description string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"entries":{"https://localhost/mgmt/tm/cli/version/0":{"nestedStats":{"entries":{"active":{"description":"%s"}}}}}}`, description)
	}
}

// ---------------------------------------------------------------------
// Schema-shape assertion (no HTTP server needed)
// ---------------------------------------------------------------------

func TestResourceBigipLtmProfileFtpSchema(t *testing.T) {
	r := resourceBigipLtmProfileFtp()

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

	nameSchema, ok := r.Schema["name"]
	if !ok {
		t.Fatal("Expected field 'name' to exist in schema")
	}
	if !nameSchema.Required {
		t.Error("Expected field 'name' to be required")
	}
}

// ---------------------------------------------------------------------
// Direct CRUD-function unit tests (style 2)
// ---------------------------------------------------------------------

func TestUnitLtmProfileFtpCreateReadUpdateDelete_ModernVersion(t *testing.T) {
	name := "test-ftp"

	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/cli/version", versionMockHandler("  Version     15.1.0"))
	mux.HandleFunc("/mgmt/tm/ltm/profile/ftp", func(w http.ResponseWriter, r *http.Request) {
		AssertRequestMethod(t, r, http.MethodPost)
		_, _ = fmt.Fprintf(w, `{"name":"%s"}`, name)
	})
	mux.HandleFunc("/mgmt/tm/ltm/profile/ftp/"+name, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			_, _ = fmt.Fprintf(w, `{"name":"%s","defaultsFrom":"/Common/ftp","allowFtps":"enabled","appService":"my-app","description":"my ftp","inheritParentProfile":"enabled","logProfile":"/Common/log","inheritVlanList":"enabled","logPublisher":"/Common/publisher","port":21,"security":"disabled","translateExtended":"enabled","ftpsMode":"allow","enforceTlsSesionReuse":"enabled","allowActiveMode":"enabled"}`, name)
		case http.MethodPatch:
			_, _ = fmt.Fprintf(w, `{"name":"%s"}`, name)
		case http.MethodDelete:
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprint(w, `{}`)
		}
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipLtmProfileFtp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":                     name,
		"defaults_from":            "/Common/ftp",
		"allow_ftps":               "enabled",
		"app_service":              "my-app",
		"description":              "my ftp",
		"inherit_parent_profile":   "enabled",
		"log_profile":              "/Common/log",
		"inherit_vlan_list":        "enabled",
		"log_publisher":            "/Common/publisher",
		"port":                     21,
		"security":                 "disabled",
		"translate_extended":       "enabled",
		"ftps_mode":                "allow",
		"enforce_tlssession_reuse": "enabled",
		"allow_active_mode":        "enabled",
	}, "")

	ctx := context.Background()

	diags := resourceBigipLtmProfileFtpCreate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected create error: %v", diags)
	require.Equal(t, name, d.Id())

	diags = resourceBigipLtmProfileFtpRead(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected read error: %v", diags)
	require.Equal(t, "allow", d.Get("ftps_mode").(string))

	diags = resourceBigipLtmProfileFtpUpdate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected update error: %v", diags)

	diags = resourceBigipLtmProfileFtpDelete(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected delete error: %v", diags)
	require.Equal(t, "", d.Id())
}

func TestUnitLtmProfileFtpCreateReadUpdate_LegacyVersion(t *testing.T) {
	// BIG-IP 12.x/13.x: the ftps_mode/enforce_tlssession_reuse/allow_active_mode
	// fields are omitted from the request/state entirely (legacy branch).
	name := "test-ftp"

	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/cli/version", versionMockHandler("  Version     13.1.0"))
	mux.HandleFunc("/mgmt/tm/ltm/profile/ftp", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"name":"%s"}`, name)
	})
	mux.HandleFunc("/mgmt/tm/ltm/profile/ftp/"+name, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			_, _ = fmt.Fprintf(w, `{"name":"%s","defaultsFrom":"/Common/ftp","port":21}`, name)
		case http.MethodPatch:
			_, _ = fmt.Fprintf(w, `{"name":"%s"}`, name)
		}
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipLtmProfileFtp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": name,
		"port": 21,
	}, "")

	ctx := context.Background()

	diags := resourceBigipLtmProfileFtpCreate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected create error: %v", diags)

	diags = resourceBigipLtmProfileFtpRead(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected read error: %v", diags)

	diags = resourceBigipLtmProfileFtpUpdate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected update error: %v", diags)
}

func TestUnitLtmProfileFtpCreate_VersionError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/cli/version", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"version failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipLtmProfileFtp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": "test-ftp",
	}, "")

	diags := resourceBigipLtmProfileFtpCreate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitLtmProfileFtpCreate_Error(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/cli/version", versionMockHandler("  Version     15.1.0"))
	mux.HandleFunc("/mgmt/tm/ltm/profile/ftp", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"create failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipLtmProfileFtp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": "test-ftp",
	}, "")

	diags := resourceBigipLtmProfileFtpCreate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitLtmProfileFtpRead_VersionError(t *testing.T) {
	name := "test-ftp"
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/cli/version", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"version failed"}`)
	})
	mux.HandleFunc("/mgmt/tm/ltm/profile/ftp/"+name, func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"name":"%s"}`, name)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipLtmProfileFtp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": name,
	}, name)

	diags := resourceBigipLtmProfileFtpRead(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitLtmProfileFtpRead_Error(t *testing.T) {
	name := "test-ftp"
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/profile/ftp/"+name, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"read failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipLtmProfileFtp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": name,
	}, name)

	diags := resourceBigipLtmProfileFtpRead(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitLtmProfileFtpUpdate_VersionError(t *testing.T) {
	name := "test-ftp"
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/cli/version", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"version failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipLtmProfileFtp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": name,
	}, name)

	diags := resourceBigipLtmProfileFtpUpdate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitLtmProfileFtpUpdate_Error(t *testing.T) {
	name := "test-ftp"
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/cli/version", versionMockHandler("  Version     15.1.0"))
	mux.HandleFunc("/mgmt/tm/ltm/profile/ftp/"+name, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"update failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipLtmProfileFtp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": name,
	}, name)

	diags := resourceBigipLtmProfileFtpUpdate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitLtmProfileFtpUpdate_LegacyVersionError(t *testing.T) {
	name := "test-ftp"
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/cli/version", versionMockHandler("  Version     12.1.0"))
	mux.HandleFunc("/mgmt/tm/ltm/profile/ftp/"+name, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"update failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipLtmProfileFtp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": name,
	}, name)

	diags := resourceBigipLtmProfileFtpUpdate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitLtmProfileFtpDelete_Error(t *testing.T) {
	name := "test-ftp"
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/profile/ftp/"+name, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"delete failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipLtmProfileFtp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": name,
	}, name)

	diags := resourceBigipLtmProfileFtpDelete(context.Background(), d, client)
	require.True(t, diags.HasError())
}
