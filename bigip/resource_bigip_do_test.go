/*
Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"crypto/tls"
	"fmt"
	"github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
)

//var TEST_DEVICE_NAME = fmt.Sprintf("/%s/test-device", TEST_PARTITION)

var dir1, _ = os.Getwd()

var TEST_DO_RESOURCE = `
resource "bigip_do"  "do-example" {
     do_json = "${file("` + dir1 + `/../examples/do/example1.json")}"
     tenant_name = "do"
}
`

func TestAccBigipDo_create(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckdevicesDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TEST_DO_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckDoExists("do", true),
					resource.TestCheckResourceAttr("bigip_do.do-example", "tenant_name", "do"),
				),
			},
		},
	})
}

func testCheckDoExists(name string, exists bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client_bigip := testAccProvider.Meta().(*bigip.BigIP)
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
		client := &http.Client{Transport: tr}
		url := client_bigip.Host + "/mgmt/shared/declarative-onboarding/"
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return fmt.Errorf("[ERROR] Error while creating http request with DO json: %v", err)
		}
		req.SetBasicAuth(client_bigip.User, client_bigip.Password)
		req.Header.Set("Accept", "application/json")
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		body, err := ioutil.ReadAll(resp.Body)
		bodyString := string(body)
		if resp.Status == "204 No Content" || err != nil {
			return fmt.Errorf("[ERROR] Error while checking doresource present in bigip :%s  %v", bodyString, err)
			defer resp.Body.Close()
		}
		defer resp.Body.Close()
		return nil
	}
}
