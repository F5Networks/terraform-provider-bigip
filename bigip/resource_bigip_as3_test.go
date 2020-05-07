/*
Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"testing"

	"github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

//var TEST_DEVICE_NAME = fmt.Sprintf("/%s/test-device", TEST_PARTITION)

var dir, err = os.Getwd()

var TEST_AS3_RESOURCE = `
resource "bigip_as3"  "as3-example" {
     as3_json = "${file("` + dir + `/../examples/as3/example1.json")}"
    // tenant_name = "as3"
}
`

var TEST_AS3_RESOURCE_INVALID_JSON = `
resource "bigip_as3"  "as3-example" {
     as3_json = "${file("` + dir + `/../examples/as3/invalid.json")}"
    // tenant_name = "as3"
}
`

func TestAccBigipAs3_create(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckdevicesDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TEST_AS3_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckAs3Exists("as3", true),
					//					resource.TestCheckResourceAttr("bigip_as3.as3-example", "tenant_name", "as3"),
				),
			},
		},
	})
}

func testCheckAs3Exists(name string, exists bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client_bigip := testAccProvider.Meta().(*bigip.BigIP)
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
		client := &http.Client{Transport: tr}
		url := client_bigip.Host + "/mgmt/shared/appsvcs/declare"
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return fmt.Errorf("[ERROR] Error while creating http request with AS3 json: %v", err)
		}
		req.SetBasicAuth(client_bigip.User, client_bigip.Password)
		req.Header.Set("Accept", "application/json")
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		body, err := ioutil.ReadAll(resp.Body)
		bodyString := string(body)
		if resp.Status == "204 No Content" || err != nil {
			return fmt.Errorf("[ERROR] Error while checking as3resource present in bigip :%s  %v", bodyString, err)
			defer resp.Body.Close()
		}
		defer resp.Body.Close()
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
				Config:      TEST_AS3_RESOURCE_INVALID_JSON,
				ExpectError: regexp.MustCompile(`"as3_json" contains an invalid JSON:.*`),
			},
		},
	})
}
