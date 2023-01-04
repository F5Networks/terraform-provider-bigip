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

func TestAccBigipSysIappUnitInvalid(t *testing.T) {
	resourceName := "/Common/test-sys-iapp"
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipSysIappInvalid(resourceName),
				ExpectError: regexp.MustCompile(" Unsupported argument: An argument named \"invalidkey\" is not expected here"),
			},
		},
	})
}

func TestAccBigipSysIappUnitCreate(t *testing.T) {
	resourceName := "/Common/test-sys-iapp"
	//httpDefault := "/Common/http"
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
	mux.HandleFunc("/mgmt/tm/sys/application/service", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		_, _ = fmt.Fprintf(w, `{"name":"test-sys-iapp","partition":"Common","inheritedDevicegroup":"true","inheritedTrafficGroup":"true","strictUpdates":"enabled","template":"/Common/f5.http","templateModified":"no","tables":[{"columnNames":null,"name":"basic__snatpool_members","rows":null},{"columnNames":null,"name":"net__snatpool_members","rows":null},{"columnNames":null,"name":"optimizations__hosts","rows":null},{"columnNames":["name"],"name":"pool__hosts","rows":[{"row":["f5.cisco.com"]}]},{"columnNames":["addr","port","connection_limit"],"name":"pool__members","rows":[{"row":["10.0.2.167","80","0"]},{"row":["10.0.2.168","80","0"]}]},{"columnNames":null,"name":"server_pools__servers","rows":null}],"variables":[{"encrypted":"no","name":"client__http_compression","value":"/#create_new#"},{"encrypted":"no","name":"monitor__monitor","value":"/Common/http"},{"encrypted":"no","name":"net__client_mode","value":"wan"},{"encrypted":"no","name":"net__server_mode","value":"lan"},{"encrypted":"no","name":"net__v13_tcp","value":"warn"},{"encrypted":"no","name":"pool__addr","value":"10.0.1.100"},{"encrypted":"no","name":"pool__pool_to_use","value":"/#create_new#"},{"encrypted":"no","name":"pool__port","value":"80"},{"encrypted":"no","name":"ssl__mode","value":"no_ssl"},{"encrypted":"no","name":"ssl_encryption_questions__advanced","value":"no"},{"encrypted":"no","name":"ssl_encryption_questions__help","value":"hide"}]}`)
	})
	mux.HandleFunc("/mgmt/tm/sys/application/service/~Common~~Common~test-sys-iapp.app~~Common~test-sys-iapp", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"name":"test-sys-iapp","partition":"Common","inheritedDevicegroup":"true","inheritedTrafficGroup":"true","strictUpdates":"enabled","template":"/Common/f5.http","templateModified":"no","tables":[{"columnNames":null,"name":"basic__snatpool_members","rows":null},{"columnNames":null,"name":"net__snatpool_members","rows":null},{"columnNames":null,"name":"optimizations__hosts","rows":null},{"columnNames":["name"],"name":"pool__hosts","rows":[{"row":["f5.cisco.com"]}]},{"columnNames":["addr","port","connection_limit"],"name":"pool__members","rows":[{"row":["10.0.2.167","80","0"]},{"row":["10.0.2.168","80","0"]}]},{"columnNames":null,"name":"server_pools__servers","rows":null}],"variables":[{"encrypted":"no","name":"client__http_compression","value":"/#create_new#"},{"encrypted":"no","name":"monitor__monitor","value":"/Common/http"},{"encrypted":"no","name":"net__client_mode","value":"wan"},{"encrypted":"no","name":"net__server_mode","value":"lan"},{"encrypted":"no","name":"net__v13_tcp","value":"warn"},{"encrypted":"no","name":"pool__addr","value":"10.0.1.100"},{"encrypted":"no","name":"pool__pool_to_use","value":"/#create_new#"},{"encrypted":"no","name":"pool__port","value":"80"},{"encrypted":"no","name":"ssl__mode","value":"no_ssl"},{"encrypted":"no","name":"ssl_encryption_questions__advanced","value":"no"},{"encrypted":"no","name":"ssl_encryption_questions__help","value":"hide"}]}`)
	})

	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config:             testBigipSysIappCreate(resourceName, server.URL),
				ExpectNonEmptyPlan: true,
			},
			{
				Config:             testBigipSysIappModify(resourceName, server.URL),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

//	func TestAccBigipSysIappUnitReadError(t *testing.T) {
//		resourceName := "/Common/test-sys-iapp"
//		httpDefault := "/Common/http"
//		setup()
//		mux.HandleFunc("mgmt/shared/authn/login", func(w http.ResponseWriter, r *http.Request) {
//			assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
//			assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
//		})
//		mux.HandleFunc("/mgmt/tm/net/self", func(w http.ResponseWriter, r *http.Request) {
//			assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)
//			assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
//			_, _ = fmt.Fprintf(w, `{}`)
//		})
//		mux.HandleFunc("/mgmt/tm/ltm/profile/http", func(w http.ResponseWriter, r *http.Request) {
//			assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
//			_, _ = fmt.Fprintf(w, `{"name":"%s","defaultsFrom":"%s", "basicAuthRealm": "none"}`, resourceName, httpDefault)
//		})
//		mux.HandleFunc("/mgmt/tm/ltm/profile/http/~Common~test-sys-iapp", func(w http.ResponseWriter, r *http.Request) {
//			assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)
//			http.Error(w, "The requested HTTP Profile (/Common/test-sys-iapp) was not found", http.StatusNotFound)
//		})
//
//		defer teardown()
//		resource.Test(t, resource.TestCase{
//			IsUnitTest: true,
//			Providers:  testProviders,
//			Steps: []resource.TestStep{
//				{
//					Config:      testBigipSysIappCreate(resourceName, server.URL),
//					ExpectError: regexp.MustCompile("HTTP 404 :: The requested HTTP Profile \\(/Common/test-sys-iapp\\) was not found"),
//				},
//			},
//		})
//	}
func TestAccBigipSysIappUnitCreateError(t *testing.T) {
	resourceName := "/Common/test-sys-iapp"
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
	mux.HandleFunc("/mgmt/tm/sys/application/service", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		http.Error(w, "Bad Request", http.StatusBadRequest)
	})

	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipSysIappCreate(resourceName, server.URL),
				ExpectError: regexp.MustCompile("HTTP 400 :: Bad Request"),
			},
		},
	})
}

func testBigipSysIappInvalid(resourceName string) string {
	return fmt.Sprintf(`
resource "bigip_sys_iapp" "test-sys-iapp" {
  name       = "%s"
  invalidkey = "foo"
}
provider "bigip" {
  address  = "xxx.xxx.xxx.xxx"
  username = "xxx"
  password = "xxx"
}`, resourceName)
}

func testBigipSysIappCreate(resourceName, url string) string {
	return fmt.Sprintf(`
resource "bigip_sys_iapp" "test-sys-iapp" {
  name    = "%s"
  jsonfile = <<EOF
		{
			"inheritedDevicegroup": "true",
			"inheritedTrafficGroup": "true",
			"name": "test-sys-iapp",
			"partition": "Common",
			"strictUpdates": "enabled",
			"tables": [
					{
							"name": "basic__snatpool_members"
					},
					{
							"name": "net__snatpool_members"
					},
					{
							"name": "optimizations__hosts"
					},
					{
							"columnNames": [
									"name"
							],
							"name": "pool__hosts",
							"rows": [
									{
											"row": [
													"f5.cisco.com"
											]
									}
							]
					},
					{
							"columnNames": [
									"addr",
									"port",
									"connection_limit"
							],
							"name": "pool__members",
							"rows": [
									{
											"row": [
													"10.0.2.167",
													"80",
													"0"
											]
									},
									{
											"row": [
													"10.0.2.168",
													"80",
													"0"
											]
									}
							]
					},
					{
							"name": "server_pools__servers"
					}
			],
			"template": "/Common/f5.http",
			"templateModified": "no",
			"templateReference": {
					"link": "https://localhost/mgmt/tm/sys/application/template/~Common~f5.http?ver=13.0.0"
			},
		 
			"variables": [
					{
							"encrypted": "no",
							"name": "client__http_compression",
							"value": "/#create_new#"
					},
					{
							"encrypted": "no",
							"name": "monitor__monitor",
							"value": "/Common/http"
					},
					{
							"encrypted": "no",
							"name": "net__client_mode",
							"value": "wan"
					},
					{
							"encrypted": "no",
							"name": "net__server_mode",
							"value": "lan"
					},
					{
							"encrypted": "no",
							"name": "net__v13_tcp",
							"value": "warn"
					},
					{
							"encrypted": "no",
							"name": "pool__addr",
							"value": "10.0.1.100"
					},
					{
							"encrypted": "no",
							"name": "pool__pool_to_use",
							"value": "/#create_new#"
					},
					{
							"encrypted": "no",
							"name": "pool__port",
							"value": "80"
					},
					{
							"encrypted": "no",
							"name": "ssl__mode",
							"value": "no_ssl"
					},
					{
							"encrypted": "no",
							"name": "ssl_encryption_questions__advanced",
							"value": "no"
					},
					{
							"encrypted": "no",
							"name": "ssl_encryption_questions__help",
							"value": "hide"
					}
			]
	}
EOF
}
provider "bigip" {
  address  = "%s"
  username = ""
  password = ""
  login_ref = ""
}`, resourceName, url)
}

func testBigipSysIappModify(resourceName, url string) string {
	return fmt.Sprintf(`
resource "bigip_sys_iapp" "test-sys-iapp" {
  name    = "%s"
  jsonfile = <<EOF
		{
			"inheritedDevicegroup": "true",
			"inheritedTrafficGroup": "true",
			"name": "test-sys-iapp",
			"partition": "Common",
			"strictUpdates": "enabled",
			"tables": [
					{
							"name": "basic__snatpool_members"
					},
					{
							"name": "net__snatpool_members"
					},
					{
							"name": "optimizations__hosts"
					},
					{
							"columnNames": [
									"name"
							],
							"name": "pool__hosts",
							"rows": [
									{
											"row": [
													"f5.cisco.com"
											]
									}
							]
					},
					{
							"columnNames": [
									"addr",
									"port",
									"connection_limit"
							],
							"name": "pool__members",
							"rows": [
									{
											"row": [
													"10.0.2.167",
													"80",
													"0"
											]
									},
									{
											"row": [
													"10.0.2.168",
													"80",
													"0"
											]
									}
							]
					},
					{
							"name": "server_pools__servers"
					}
			],
			"template": "/Common/f5.http",
			"templateModified": "no",
			"templateReference": {
					"link": "https://localhost/mgmt/tm/sys/application/template/~Common~f5.http?ver=13.0.0"
			},
		 
			"variables": [
					{
							"encrypted": "no",
							"name": "client__http_compression",
							"value": "/#create_new#"
					},
					{
							"encrypted": "no",
							"name": "monitor__monitor",
							"value": "/Common/http"
					},
					{
							"encrypted": "no",
							"name": "net__client_mode",
							"value": "wan"
					},
					{
							"encrypted": "no",
							"name": "net__server_mode",
							"value": "lan"
					},
					{
							"encrypted": "no",
							"name": "net__v13_tcp",
							"value": "warn"
					},
					{
							"encrypted": "no",
							"name": "pool__addr",
							"value": "10.0.1.100"
					},
					{
							"encrypted": "no",
							"name": "pool__pool_to_use",
							"value": "/#create_new#"
					},
					{
							"encrypted": "no",
							"name": "pool__port",
							"value": "80"
					},
					{
							"encrypted": "no",
							"name": "ssl__mode",
							"value": "no_ssl"
					},
					{
							"encrypted": "no",
							"name": "ssl_encryption_questions__advanced",
							"value": "no"
					},
					{
							"encrypted": "no",
							"name": "ssl_encryption_questions__help",
							"value": "hide"
					}
			]
	}
EOF
}
provider "bigip" {
  address  = "%s"
  username = ""
  password = ""
  login_ref = ""
}`, resourceName, url)
}
