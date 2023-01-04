/*
Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/stretchr/testify/assert"
)

func TestAccBigipLtmPoolAttachUnitInvalid(t *testing.T) {
	resourceName := "/Common/test-pool-attach1"
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipLtmPoolAttachInvalid(resourceName),
				ExpectError: regexp.MustCompile(" Unsupported argument: An argument named \"invalidkey\" is not expected here"),
			},
		},
	})
}

func TestAccBigipLtmPoolAttachUnitCreate(t *testing.T) {
	resourceName := "10.10.10.10:443"
	setup()
	mux.HandleFunc("mgmt/shared/authn/login", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
	})
	mux.HandleFunc("/mgmt/tm/net/self", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		_, _ = fmt.Fprintf(w, `{}`)
	})
	mux.HandleFunc("/mgmt/tm/ltm/pool/~Common~test-pool/members", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" && r.URL.String() == "/mgmt/tm/ltm/pool/~Common~test-pool/members" {
			_, _ = fmt.Fprintf(w, `{"items":[{"name": "%s",
            "partition": "Common",
            "fullPath": "/Common/%s",
            "address": "10.10.10.10",
            "connectionLimit": 0,
            "dynamicRatio": 1,
            "ephemeral": "false",
            "fqdn": {
                "autopopulate": "disabled"
            },
            "inheritProfile": "enabled",
            "logging": "disabled",
            "monitor": "default",
            "priorityGroup": 2,
            "rateLimit": "disabled",
            "ratio": 1,
            "session": "user-enabled",
            "state": "unchecked"}]}`, resourceName, resourceName)
		}
		if r.Method == "POST" {
			_, _ = fmt.Fprintf(w, `{"name":"%s","partition":"Common","fqdn":{}}`, resourceName)
		}
	})
	mux.HandleFunc("/mgmt/tm/ltm/pool/~Common~test-pool/members/~Common~10.10.10.10:443", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"connectionLimit":2,"dynamicRatio":3,"fqdn":{},"priorityGroup":2,"rateLimit":"2","ratio":2}`)
	})
	mux.HandleFunc("/mgmt/tm/ltm/pool/~Common~test-pool", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"name": "test-pool","partition": "Common","fullPath": "/Common/test-pool"}`)
	})
	mux.HandleFunc("/mgmt/tm/ltm/pool/~Common~test-pool/members/10.10.10.10:443", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, ``)
	})

	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config: testBigipLtmPoolAttachCreate(resourceName, server.URL),
			},
			{
				Config:             testBigipLtmPoolAttachModify(resourceName, server.URL),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccBigipLtmPoolAttachUnitReadError(t *testing.T) {
	resourceName := "10.10.10.10:443"
	setup()
	mux.HandleFunc("mgmt/shared/authn/login", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
	})
	mux.HandleFunc("/mgmt/tm/net/self", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		_, _ = fmt.Fprintf(w, `{}`)
	})
	mux.HandleFunc("/mgmt/tm/ltm/pool/~Common~test-pool", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"name": "test-pool","partition": "Common","fullPath": "/Common/test-pool"}`)
	})
	mux.HandleFunc("/mgmt/tm/ltm/pool/~Common~test-pool/members", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" && r.URL.String() == "/mgmt/tm/ltm/pool/~Common~test-pool/members" {
			_, _ = fmt.Fprintf(w, `{"items":[]}`)
		}
		if r.Method == "POST" {
			_, _ = fmt.Fprintf(w, `{"name":"%s","partition":"Common","fqdn":{}}`, resourceName)
		}
	})
	mux.HandleFunc("/mgmt/tm/ltm/pool/~Common~test-pool/members/~Common~10.10.10.10:443", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"connectionLimit":2,"dynamicRatio":3,"fqdn":{},"priorityGroup":2,"rateLimit":"2","ratio":2}`)
	})
	mux.HandleFunc("/mgmt/tm/ltm/pool/~Common~test-pool/members/10.10.10.10:443", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, ``)
	})

	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipLtmPoolAttachCreate(resourceName, server.URL),
				ExpectError: regexp.MustCompile("Not able to attached Node :10.10.10.10:443 to pool /Common/test-pool"),
			},
		},
	})
}

func TestAccBigipLtmPoolAttachUnitCreateError(t *testing.T) {
	resourceName := "10.10.10.10:443"
	setup()
	mux.HandleFunc("mgmt/shared/authn/login", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
	})
	mux.HandleFunc("/mgmt/tm/net/self", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		_, _ = fmt.Fprintf(w, `{}`)
	})
	mux.HandleFunc("/mgmt/tm/ltm/pool/~Common~test-pool/members", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" && r.URL.String() == "/mgmt/tm/ltm/pool/~Common~test-pool/members" {
			_, _ = fmt.Fprintf(w, `{"items":[]}`)
		}
		if r.Method == "POST" {
			http.Error(w, "Bad Request", http.StatusBadRequest)
		}
	})

	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipLtmPoolAttachCreate(resourceName, server.URL),
				ExpectError: regexp.MustCompile("Failure adding node 10.10.10.10:443 to pool /Common/test-pool: HTTP 400 :: Bad Request"),
			},
		},
	})
}

func TestAccBigipLtmPoolAttachUnitImport(t *testing.T) {
	resourceName := "10.10.10.10:443"
	setup()
	mux.HandleFunc("mgmt/shared/authn/login", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
	})
	mux.HandleFunc("/mgmt/tm/net/self", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		_, _ = fmt.Fprintf(w, `{}`)
	})
	mux.HandleFunc("/mgmt/tm/ltm/pool/~Common~test-pool/members", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" && r.URL.String() == "/mgmt/tm/ltm/pool/~Common~test-pool/members" {
			_, _ = fmt.Fprintf(w, `{"items":[{"name": "%s",
            "partition": "Common",
            "fullPath": "/Common/%s",
            "address": "10.10.10.10",
            "connectionLimit": 0,
            "dynamicRatio": 1,
            "ephemeral": "false",
            "fqdn": {
                "autopopulate": "disabled"
            },
            "inheritProfile": "enabled",
            "logging": "disabled",
            "monitor": "default",
            "priorityGroup": 2,
            "rateLimit": "disabled",
            "ratio": 1,
            "session": "user-enabled",
            "state": "unchecked"}]}`, resourceName, resourceName)
		}
		if r.Method == "POST" {
			_, _ = fmt.Fprintf(w, `{"name":"%s","partition":"Common","fqdn":{}}`, resourceName)
		}
	})
	mux.HandleFunc("/mgmt/tm/ltm/pool/~Common~test-pool/members/~Common~10.10.10.10:443", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"connectionLimit":2,"dynamicRatio":3,"fqdn":{},"priorityGroup":2,"rateLimit":"2","ratio":2}`)
	})
	mux.HandleFunc("/mgmt/tm/ltm/pool/~Common~test-pool", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"name": "test-pool","partition": "Common","fullPath": "/Common/test-pool"}`)
	})
	mux.HandleFunc("/mgmt/tm/ltm/pool/~Common~test-pool/members/10.10.10.10:443", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, ``)
	})

	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config: testBigipLtmPoolAttachCreate(resourceName, server.URL),
			},
			{
				Config:       testBigipLtmPoolAttachImport(server.URL),
				ResourceName: "bigip_ltm_pool_attachment.test-pool-attach1",
				ImportState:  true,
				//ImportStateVerify: true,
				ImportStateId: `{"pool": "/Common/test-pool", "node": "/Common/10.10.10.10:443"}`,
			},
		},
	})
}

func TestAccBigipLtmPoolAttachUnitCreateFqdn(t *testing.T) {
	resourceName := "www.f5.com:443"
	setup()
	mux.HandleFunc("mgmt/shared/authn/login", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
	})
	mux.HandleFunc("/mgmt/tm/net/self", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		_, _ = fmt.Fprintf(w, `{}`)
	})
	mux.HandleFunc("/mgmt/tm/ltm/pool/~Common~test-pool/members", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" && r.URL.String() == "/mgmt/tm/ltm/pool/~Common~test-pool/members" {
			_, _ = fmt.Fprintf(w, `{"items":[{
    "name": "%s",
    "partition": "Common",
    "fullPath": "/Common/%s",
    "address": "any6",
    "connectionLimit": 0,
    "dynamicRatio": 1,
    "ephemeral": "false",
    "fqdn": {
        "autopopulate": "enabled",
        "tmName": "www.f5.com"
    },
    "inheritProfile": "enabled",
    "logging": "disabled",
    "monitor": "default",
    "priorityGroup": 0,
    "rateLimit": "disabled",
    "ratio": 1,
    "session": "user-enabled",
    "state": "fqdn-up"
}]}`, resourceName, resourceName)
		}
		if r.Method == "POST" {
			_, _ = fmt.Fprintf(w, `{"name":"%s","partition":"Common","fqdn":{"autopopulate": "enabled","tmName": "www.f5.com"}}`, resourceName)
		}
	})
	mux.HandleFunc("/mgmt/tm/ltm/pool/~Common~test-pool/members/~Common~www.f5.com:443", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"connectionLimit":2,"dynamicRatio":3,"fqdn":{"autopopulate": "enabled","tmName": "www.f5.com"}},"priorityGroup":2,"rateLimit":"2","ratio":2}`)
	})
	mux.HandleFunc("/mgmt/tm/ltm/pool/~Common~test-pool", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"name": "test-pool","partition": "Common","fullPath": "/Common/test-pool"}`)
	})
	mux.HandleFunc("/mgmt/tm/ltm/pool/~Common~test-pool/members/www.f5.com:443", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, ``)
	})

	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config: testBigipLtmPoolAttachCreateFqdn(resourceName, server.URL),
			},
		},
	})
}

func TestAccBigipLtmPoolAttachUnitCreateFqdnTC2(t *testing.T) {
	resourceName := "/Common/www.f5.com:443"
	setup()
	mux.HandleFunc("mgmt/shared/authn/login", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
	})
	mux.HandleFunc("/mgmt/tm/net/self", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		_, _ = fmt.Fprintf(w, `{}`)
	})
	mux.HandleFunc("/mgmt/tm/ltm/node/~Common~www.f5.com", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{
    "name": "www.f5.com",
    "partition": "Common",
    "fullPath": "/Common/www.f5.com",
    "address": "any6",
    "connectionLimit": 0,
    "dynamicRatio": 1,
    "ephemeral": "false",
    "fqdn": {
        "addressFamily": "all",
        "autopopulate": "enabled",
        "downInterval": 5,
        "interval": "3600",
        "tmName": "www.f5.com"
    },
    "logging": "disabled",
    "monitor": "default",
    "rateLimit": "disabled",
    "ratio": 1,
    "session": "user-enabled",
    "state": "fqdn-up"
}`)
	})
	mux.HandleFunc("/mgmt/tm/ltm/pool/~Common~test-pool/members", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" && r.URL.String() == "/mgmt/tm/ltm/pool/~Common~test-pool/members" {
			_, _ = fmt.Fprintf(w, `{"items":[{
    "name": "%s",
    "partition": "Common",
    "fullPath": "%s",
    "address": "any6",
    "connectionLimit": 0,
    "dynamicRatio": 1,
    "ephemeral": "false",
    "fqdn": {
        "autopopulate": "enabled",
        "tmName": "www.f5.com"
    },
    "inheritProfile": "enabled",
    "logging": "disabled",
    "monitor": "default",
    "priorityGroup": 0,
    "rateLimit": "disabled",
    "ratio": 1,
    "session": "user-enabled",
    "state": "fqdn-up"
}]}`, resourceName, resourceName)
		}
		if r.Method == "POST" {
			_, _ = fmt.Fprintf(w, `{"name":"%s","partition":"Common","fqdn":{"autopopulate": "enabled","tmName": "www.f5.com"}}`, resourceName)
		}
	})
	mux.HandleFunc("/mgmt/tm/ltm/pool/~Common~test-pool/members/~Common~www.f5.com:443", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"connectionLimit":2,"dynamicRatio":3,"fqdn":{"autopopulate": "enabled","tmName": "www.f5.com"}},"priorityGroup":2,"rateLimit":"2","ratio":2}`)
	})
	mux.HandleFunc("/mgmt/tm/ltm/pool/~Common~test-pool", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"name": "test-pool","partition": "Common","fullPath": "/Common/test-pool"}`)
	})
	mux.HandleFunc("/mgmt/tm/ltm/pool/~Common~test-pool/members/www.f5.com:443", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, ``)
	})

	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config: testBigipLtmPoolAttachCreateFqdn(resourceName, server.URL),
			},
		},
	})
}

func testBigipLtmPoolAttachInvalid(resourceName string) string {
	return fmt.Sprintf(`
resource "bigip_ltm_pool_attachment" "test-pool-attach1" {
   pool = "/Common/test-pool"
   node = "%s"
   invalidkey = "foo"
}

provider "bigip" {
  address  = "xxx.xxx.xxx.xxx"
  username = "xxx"
  password = "xxx"
}`, resourceName)
}

func testBigipLtmPoolAttachCreate(resourceName, url string) string {
	return fmt.Sprintf(`
resource "bigip_ltm_pool_attachment" "test-pool-attach1" {
   pool = "/Common/test-pool"
   node = "%s"
   ratio                 = 2
   connection_limit      = 2
   connection_rate_limit = 2
   priority_group        = 2
   dynamic_ratio         = 3
}
provider "bigip" {
  address  = "%s"
  username = ""
  password = ""
  login_ref = ""
}`, resourceName, url)
}

func testBigipLtmPoolAttachModify(resourceName, url string) string {
	return fmt.Sprintf(`
resource "bigip_ltm_pool_attachment" "test-pool-attach1" {
   pool = "/Common/test-pool"
   node = "%s"
   ratio                 = 2
   connection_limit      = 2
   connection_rate_limit = 2
   priority_group        = 3
   dynamic_ratio         = 3
}
provider "bigip" {
  address  = "%s"
  username = ""
  password = ""
  login_ref = ""
}`, resourceName, url)
}
func testBigipLtmPoolAttachImport(url string) string {
	return fmt.Sprintf(`
resource "bigip_ltm_pool_attachment" "test-pool-attach1" {
    pool = "/Partition/Name"
    node = "/Partition/Name"
}
provider "bigip" {
  address  = "%s"
  username = ""
  password = ""
  login_ref = ""
}`, url)
}
func testBigipLtmPoolAttachCreateFqdn(resourceName, url string) string {
	return fmt.Sprintf(`
resource "bigip_ltm_pool_attachment" "test-pool-attach1" {
   pool = "/Common/test-pool"
   node = "%s"
   ratio                 = 1
   connection_limit      = 0
   priority_group        = 0
   dynamic_ratio         = 1
}
provider "bigip" {
  address  = "%s"
  username = ""
  password = ""
  login_ref = ""
}`, resourceName, url)
}
