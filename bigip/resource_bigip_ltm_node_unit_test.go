/*
Original work from https://github.com/DealerDotCom/terraform-provider-bigip
Modifications Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file,You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/stretchr/testify/assert"
	//"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
)

func testBigipLtmNodeInvalid(resourceName string) string {
	return fmt.Sprintf(`
		resource "bigip_ltm_node" "test-node" {
			name = "%s"
			address = "10.10.10.10"
	                invalidkey = "foo"
		}
		provider "bigip" {
			address = "xxx.xxx.xxx.xxx"
			username = "xxx"
			password = "xxx"
		}
	`, resourceName)
}

func TestAccBigipLtmNodeInvalid(t *testing.T) {
	resourceName := "/Common/test-node"
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipLtmNodeInvalid(resourceName),
				ExpectError: regexp.MustCompile("Unsupported argument: An argument named \"invalidkey\" is not expected here"),
			},
		},
	})
}

func testBigipLtmNodeCreate(resourceName string, url string, address string) string {
	return fmt.Sprintf(`
		resource "bigip_ltm_node" "test-node" {
			name = "%s"
			address = "%s"
		}
		provider "bigip" {
			address = "%s"
			username = "xxxx"
			password = "xxxx"
		}
	`, resourceName, address, url)
}

func TestAccBigipLtmNodeCreate(t *testing.T) {
	resourceName := "/Common/test-node"
	address := "10.10.10.10"
	setup()
	mux.HandleFunc("/mgmt/tm/net/self", func(w http.ResponseWriter, r *http.Request) {
		//fmt.Println(r)
		assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)
		log.Println(" value of t  ")
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		fmt.Fprintf(w, `{
}`)
	})
	mux.HandleFunc("/mgmt/tm/ltm/node", func(w http.ResponseWriter, r *http.Request) {
		//fmt.Println(r)
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		//b, _ := ioutil.ReadAll(r.Body)
		//defer r.Body.Close()
		//fmt.Println(string(b))
		fmt.Fprintf(w, `{"name":"%s","address":"%s"}`, resourceName, address)
	})
	mux.HandleFunc("/mgmt/tm/ltm/node/~Common~test-node", func(w http.ResponseWriter, r *http.Request) {
		//fmt.Println(r)
		fmt.Fprintf(w, `{"name":"%s","address":"%s"}`, resourceName, address)
	})
	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config: testBigipLtmNodeCreate(resourceName, server.URL, address),
			},
		},
	})
}

var (
	// mux is the HTTP request multiplexer used with the test server.
	mux *http.ServeMux

	// server is a test HTTP server used to provide mock API responses
	server *httptest.Server
)

func setup() {
	// test server
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)
}

func teardown() {
	server.Close()
}
