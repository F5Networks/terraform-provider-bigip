/*
Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"testing"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var dir, _ = os.Getwd()

var TestAs3Resource = `
resource "bigip_as3"  "as3-example" {
     as3_json = "${file("` + dir + `/../examples/as3/example1.json")}"
}
`
var TestAs3Resource1 = `
resource "bigip_as3"  "as3-multitenant-example" {
     as3_json = "${file("` + dir + `/../examples/as3/as3_example1.json")}"
}
`
var TestAs3Resource2 = `
resource "bigip_as3"  "as3-partialsuccess-example" {
     as3_json = "${file("` + dir + `/../examples/as3/as3_example2.json")}"
}
`
var TestAs3Resource3 = `
resource "bigip_as3"  "as3-tenantadd-example" {
     as3_json = "${file("` + dir + `/../examples/as3/as3_example3.json")}"
}
`
var TestAs3Resource4 = `
resource "bigip_as3"  "as3-tenantfilter-example" {
     as3_json = "${file("` + dir + `/../examples/as3/as3_example1.json")}"
     tenant_filter = "Sample_01"
}
`
var TestAs3ResourceInvalidJson = `
resource "bigip_as3"  "as3-example" {
     as3_json = "${file("` + dir + `/../examples/as3/invalid.json")}"
}
`
var TestAs3Resource5 = `
resource "bigip_as3"  "as3-example" {
     as3_json = "${file("` + dir + `/../examples/as3/example3.json")}"
}
`

var TestAs3Resourcegithub592 = `
resource "bigip_as3"  "as3-example" {
     as3_json = "${file("` + dir + `/../examples/as3/github592.json")}"
}
`

var TestAs3Resourcegithub600 = `
resource "bigip_as3"  "as3-example" {
     as3_json = "${file("` + dir + `/../examples/as3/github600.json")}"
}
`

var TestAs3Resourcegithub601a = `
resource "bigip_as3"  "as3-example" {
	as3_json = "${file("` + dir + `/../examples/as3/github601_a.json")}"
}
`

var TestAs3Resourcegithub601b = `
resource "bigip_as3"  "as3-example" {
	as3_json = "${file("` + dir + `/../examples/as3/github601_b.json")}"
}
`

func TestAccBigipAs3_create_SingleTenant(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckAs3Destroy,
		Steps: []resource.TestStep{
			{
				Config: TestAs3Resource,
				Check: resource.ComposeTestCheckFunc(
					testCheckAs3Exists("Sample_new", true),
				),
			},
		},
	})
}

func TestAccBigipAs3_create_MultiTenants(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckAs3Destroy,
		Steps: []resource.TestStep{
			{
				Config: TestAs3Resource1,
				Check: resource.ComposeTestCheckFunc(
					testCheckAs3Exists("Sample_01,Sample_02", true),
				),
			},
		},
	})
}
func TestAccBigipAs3_create_PartialSuccess(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckAs3Destroy,
		Steps: []resource.TestStep{
			{
				Config: TestAs3Resource2,
				Check: resource.ComposeTestCheckFunc(
					testCheckAs3Exists("Sample_03", true),
					testCheckAs3Exists("Sample_04", false),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}
func TestAccBigipAs3_addTenantFilter(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckAs3Destroy,
		Steps: []resource.TestStep{
			{
				Config: TestAs3Resource4,
				Check: resource.ComposeTestCheckFunc(
					testCheckAs3Exists("Sample_01", true),
					testCheckAs3Exists("Sample_02", false),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccBigipAs3_update_addTenant(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckAs3Destroy,
		Steps: []resource.TestStep{
			{
				Config: TestAs3Resource1,
				Check: resource.ComposeTestCheckFunc(
					testCheckAs3Exists("Sample_01,Sample_02", true),
				),
			},
			{
				Config: TestAs3Resource3,
				Check: resource.ComposeTestCheckFunc(
					testCheckAs3Exists("Sample_01,Sample_02,Sample_03", true),
				),
			},
		},
	})
}
func TestAccBigipAs3_update_deleteTenant(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckAs3Destroy,
		Steps: []resource.TestStep{
			{
				Config: TestAs3Resource3,
				Check: resource.ComposeTestCheckFunc(
					testCheckAs3Exists("Sample_01,Sample_02,Sample_03", true),
				),
			},
			{
				Config: TestAs3Resource1,
				Check: resource.ComposeTestCheckFunc(
					testCheckAs3Exists("Sample_01,Sample_02", true),
					testCheckAs3Exists("Sample_03", false),
				),
			},
		},
	})
}

func TestAccBigipAs3_update_config(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckAs3Destroy,
		Steps: []resource.TestStep{
			{
				Config: TestAs3Resource,
				Check: resource.ComposeTestCheckFunc(
					testCheckAs3Exists("Sample_new", true),
				),
			},
			{
				Config: TestAs3Resource5,
				Check: resource.ComposeTestCheckFunc(
					testCheckAs3Exists("Sample_new", true),
				),
			},
		},
	})
}

func TestAccBigipAs3Issue592(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckAs3Destroy,
		Steps: []resource.TestStep{
			{
				Config: TestAs3Resourcegithub592,
				Check: resource.ComposeTestCheckFunc(
					testCheckAs3Exists("A1", true),
				),
			},
		},
	})
}

func TestAccBigipAs3Issue600(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckAs3Destroy,
		Steps: []resource.TestStep{
			{
				Config: TestAs3Resourcegithub592,
				Check: resource.ComposeTestCheckFunc(
					testCheckAs3Exists("A1", true),
				),
			},
			{
				Config: TestAs3Resourcegithub600,
				Check: resource.ComposeTestCheckFunc(
					testCheckAs3Exists("A1", true),
				),
			},
		},
	})
}

// func TestAccBigipAs3Issue601(t *testing.T) {
//	resource.Test(t, resource.TestCase{
//		PreCheck: func() {
//			testAcctPreCheck(t)
//		},
//		Providers:    testAccProviders,
//		CheckDestroy: testCheckAs3Destroy,
//		Steps: []resource.TestStep{
//			{
//				Config: TestAs3Resourcegithub601a,
//				Check: resource.ComposeTestCheckFunc(
//					testCheckAs3Exists("Sample_01", true),
//				),
//			},
//			{
//				Config: TestAs3Resourcegithub601b,
//				Check: resource.ComposeTestCheckFunc(
//					testCheckAs3Exists("Sample_02", true),
//					testCheckAs3Exists("Sample_01", false),
//				),
//			},
//		},
//	})
// }

func TestAccBigipAs3_import_SingleTenant(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckAs3Destroy,
		Steps: []resource.TestStep{
			{
				Config: TestAs3Resource,
				Check: resource.ComposeTestCheckFunc(
					testCheckAs3Exists("Sample_new", true),
				),
				ResourceName:      "as3-example",
				ImportState:       false,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccBigipAs3_import_MultiTenants(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckAs3Destroy,
		Steps: []resource.TestStep{
			{
				Config: TestAs3Resource1,
				Check: resource.ComposeTestCheckFunc(
					testCheckAs3Exists("Sample_01,Sample_02", true),
				),
				ResourceName:      "as3-multitenant-example",
				ImportState:       false,
				ImportStateVerify: true,
			},
		},
	})
}

func testCheckAs3Exists(name string, exists bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		clientBigip := testAccProvider.Meta().(*bigip.BigIP)
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
		client := &http.Client{Transport: tr}
		url := clientBigip.Host + "/mgmt/shared/appsvcs/declare/" + name
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return fmt.Errorf("[ERROR] Error while creating http request with AS3 json: %v", err)
		}
		req.SetBasicAuth(clientBigip.User, clientBigip.Password)
		req.Header.Set("Accept", "application/json")
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			return err
		}

		defer func() {
			if err := resp.Body.Close(); err != nil {
				log.Printf("[DEBUG] Could not close the request to %s", url)
			}
		}()
		var body bytes.Buffer
		_, err = io.Copy(&body, resp.Body)
		// body, err := ioutil.ReadAll(resp.Body)
		bodyString := body.String()
		if (resp.Status == "204 No Content" || err != nil || resp.StatusCode == 404) && exists {
			return fmt.Errorf("[ERROR] Error while checking as3resource present in bigip :%s  %v", bodyString, err)
		}

		return nil
	}
}

func TestAccBigipAs3_badJSON(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckdevicesDestroyed,
		Steps: []resource.TestStep{
			{
				Config:      TestAs3ResourceInvalidJson,
				ExpectError: regexp.MustCompile(`"as3_json" contains an invalid JSON:.*`),
			},
		},
	})
}
func testCheckAs3Destroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*bigip.BigIP)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bigip_as3" {
			continue
		}

		name := rs.Primary.ID
		err, failedTenants := client.DeleteAs3Bigip(name)
		if err != nil || failedTenants != "" {
			return err
		}
	}
	return nil
}
