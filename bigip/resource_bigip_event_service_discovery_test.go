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

func TestAccBigipEventServiceDiscoveryUnitInvalid(t *testing.T) {
	resourceName := "examples/simple_http"
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccBigipEventServiceDiscoveryUnitInvalid(resourceName),
				ExpectError: regexp.MustCompile(" Unsupported argument: An argument named \"invalidkey\" is not expected here"),
			},
		},
	})
}

func TestAccBigipEventServiceDiscoveryUnitCreate(t *testing.T) {
	resourceName := "examples/simple_http"
	setup()
	mux.HandleFunc("mgmt/shared/authn/login", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
	})
	mux.HandleFunc("/mgmt/shared/service-discovery/task/~Sample_event_sd~My_app~My_pool/nodes", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			_, _ = fmt.Fprintf(w, `{"result": { "providerOptions": { "nodeList" : [{"id":"newNode1","ip":"192.168.2.3","port":8080},{"id":"newNode2","ip":"192.168.2.4","port":8080}]}}}`)
		}
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{}`)
	})

	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccBigipEventServiceDiscoveryUnitCreate(resourceName, server.URL),
			},
			{
				Config:             testAccBigipEventServiceDiscoveryUnitModify(resourceName, server.URL),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccBigipEventServiceDiscoveryUnitInvalid(resourceName string) string {
	return fmt.Sprintf(`
resource "bigip_event_service_discovery" "test" {
  taskid = "~Sample_event_sd~My_app~My_pool"
  invalidkey = "foo"
  node {
    id   = "newNode1"
    ip   = "192.168.2.3"
    port = 8080
  }
  node {
    id   = "newNode2"
    ip   = "192.168.2.4"
    port = 8080
  }
}
provider "bigip" {
  address  = "xxx.xxx.xxx.xxx"
  username = "xxx"
  password = "xxx"
}`)
}

func testAccBigipEventServiceDiscoveryUnitCreate(resourceName, url string) string {
	return fmt.Sprintf(`
resource "bigip_event_service_discovery" "test" {
  taskid = "~Sample_event_sd~My_app~My_pool"
  node {
    id   = "newNode1"
    ip   = "192.168.2.3"
    port = 8080
  }
  node {
    id   = "newNode2"
    ip   = "192.168.2.4"
    port = 8080
  }
}
provider "bigip" {
  address  = "%s"
  username = ""
  password = ""
  login_ref = ""
}`, url)
}

func testAccBigipEventServiceDiscoveryUnitModify(resourceName, url string) string {
	return fmt.Sprintf(`
resource "bigip_event_service_discovery" "test" {
  taskid = "~Sample_event_sd~My_app~My_pool"
  node {
    id   = "newNode1"
    ip   = "192.168.3.3"
    port = 8080
  }
  node {
    id   = "newNode2"
    ip   = "192.168.2.4"
    port = 8080
  }
}
provider "bigip" {
  address  = "%s"
  username = ""
  password = ""
  login_ref = ""
}`, url)
}
