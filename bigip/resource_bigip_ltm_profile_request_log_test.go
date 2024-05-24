/*
Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"fmt"
	"strings"
	"testing"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var resRequestLogName = "bigip_ltm_request_log_profile"

func TestAccBigipLtmProfileRequestLogTC1(t *testing.T) {
	t.Parallel()
	var instName = "request-log-profile-tc1"
	var testPartition = "Common"
	var testRequestLogProfileName = fmt.Sprintf("/%s/%s", testPartition, instName)
	resFullName := fmt.Sprintf("%s.%s", resRequestLogName, instName)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckProfileRequestLogDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccBigipLtmProfileRequestLogTC1Config(testPartition, testRequestLogProfileName, instName),
				Check: resource.ComposeTestCheckFunc(
					testCheckRequestLogExists(testRequestLogProfileName),
					resource.TestCheckResourceAttr(resFullName, "name", testRequestLogProfileName),
					resource.TestCheckResourceAttr(resFullName, "defaults_from", "/Common/request-log"),
				),
			},
		},
	})
}

func TestAccBigipLtmProfileRequestLogTC2(t *testing.T) {
	t.Parallel()
	var instName = "request-log-profile-tc2"
	var testPartition = "Common"
	var testRequestLogProfileName = fmt.Sprintf("/%s/%s", testPartition, instName)
	resFullName := fmt.Sprintf("%s.%s", resRequestLogName, instName)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckProfileRequestLogDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccBigipLtmProfileRequestLogTC1Config(testPartition, testRequestLogProfileName, instName),
				Check: resource.ComposeTestCheckFunc(
					testCheckRequestLogExists(testRequestLogProfileName),
					resource.TestCheckResourceAttr(resFullName, "name", testRequestLogProfileName),
					resource.TestCheckResourceAttr(resFullName, "defaults_from", "/Common/request-log"),
				),
			},
			{
				Config: testAccBigipLtmProfileRequestLogTC2Config(testPartition, testRequestLogProfileName, instName),
				Check: resource.ComposeTestCheckFunc(
					testCheckRequestLogExists(testRequestLogProfileName),
					resource.TestCheckResourceAttr(resFullName, "name", testRequestLogProfileName),
					resource.TestCheckResourceAttr(resFullName, "defaults_from", "/Common/request-log"),
					resource.TestCheckResourceAttr(resFullName, "request_logging", "disabled"),
					resource.TestCheckResourceAttr(resFullName, "requestlog_protocol", "mds-tcp"),
					resource.TestCheckResourceAttr(resFullName, "requestlog_error_protocol", "mds-tcp"),
				),
			},
			{
				Config: testAccBigipLtmProfileRequestLogTC2UpdateConfig(testPartition, testRequestLogProfileName, instName),
				Check: resource.ComposeTestCheckFunc(
					testCheckRequestLogExists(testRequestLogProfileName),
					resource.TestCheckResourceAttr(resFullName, "name", testRequestLogProfileName),
					resource.TestCheckResourceAttr(resFullName, "defaults_from", "/Common/request-log"),
					resource.TestCheckResourceAttr(resFullName, "request_logging", "enabled"),
					resource.TestCheckResourceAttr(resFullName, "requestlog_protocol", "mds-udp"),
					resource.TestCheckResourceAttr(resFullName, "requestlog_error_protocol", "mds-udp"),
				),
			},
		},
	})
}

func TestAccBigipLtmProfileRequestLogTC3(t *testing.T) {
	t.Parallel()
	var instName = "request-log-profile-tc3"
	var testPartition = "Common"
	var testRequestLogProfileName = fmt.Sprintf("/%s/%s", testPartition, instName)
	resFullName := fmt.Sprintf("%s.%s", resRequestLogName, instName)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckProfileRequestLogDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccBigipLtmProfileRequestLogTC3Config(testPartition, testRequestLogProfileName, instName),
				Check: resource.ComposeTestCheckFunc(
					testCheckRequestLogExists(testRequestLogProfileName),
					resource.TestCheckResourceAttr(resFullName, "name", testRequestLogProfileName),
					resource.TestCheckResourceAttr(resFullName, "requestlog_template", "<134> ${TIME_MSECS} ${TIME_OFFSET} bigip_host=${BIGIP_HOSTNAME} type=request x-client-cert-subject=\"$X-Client-Cert-Subject\""),
					resource.TestCheckResourceAttr(resFullName, "responselog_template", "<134> ${TIME_MSECS} ${TIME_OFFSET} bigip_host=${BIGIP_HOSTNAME} type=response x-client-cert-subject=\"$X-Client-Cert-Subject\""),
				),
			},
		},
	})
}

func testCheckRequestLogExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*bigip.BigIP)
		p, err := client.GetRequestLogProfile(name)
		if err != nil {
			return err
		}
		if p == nil {
			return fmt.Errorf("http %s was not created", name)
		}

		return nil
	}
}

func testCheckProfileRequestLogDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(*bigip.BigIP)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != resRequestLogName {
			continue
		}
		name := rs.Primary.ID
		http, err := client.GetRequestLogProfile(name)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				return nil
			}
			return err
		}
		if http != nil {
			return fmt.Errorf("http %s not destroyed ", name)
		}
	}
	return nil
}

func testAccBigipLtmProfileRequestLogTC1Config(partition, profileName, resourceName string) string {
	return fmt.Sprintf(`
resource "%[4]s" "%[3]s" {
  name          = "%[2]s"
  defaults_from = "/%[1]s/request-log"
}
`, partition, profileName, resourceName, resRequestLogName)
}

func testAccBigipLtmProfileRequestLogTC2Config(partition, profileName, resourceName string) string {
	return fmt.Sprintf(`resource "%[4]s" "%[3]s" {
	name                       = "%[2]s"
	defaults_from              = "/%[1]s/request-log"
	request_logging            = "disabled"
	requestlog_protocol        = "mds-tcp"
	requestlog_error_protocol  = "mds-tcp"
	responselog_protocol       = "mds-tcp"
	responselog_error_protocol = "mds-tcp"
}`, partition, profileName, resourceName, resRequestLogName)
}

func testAccBigipLtmProfileRequestLogTC3Config(partition, profileName, resourceName string) string {
	return fmt.Sprintf(`resource "%[4]s" "%[3]s" {
	name           	= "%[2]s"
	defaults_from   = "/%[1]s/request-log"
	request_logging    = "enabled"
	requestlog_template 	= "<134> $${TIME_MSECS} $${TIME_OFFSET} bigip_host=$${BIGIP_HOSTNAME} type=request x-client-cert-subject=\"$X-Client-Cert-Subject\""
	response_logging    	= "enabled"
	responselog_template 	= "<134> $${TIME_MSECS} $${TIME_OFFSET} bigip_host=$${BIGIP_HOSTNAME} type=response x-client-cert-subject=\"$X-Client-Cert-Subject\""
}`, partition, profileName, resourceName, resRequestLogName)
}

func testAccBigipLtmProfileRequestLogTC2UpdateConfig(partition, profileName, resourceName string) string {
	return fmt.Sprintf(`resource "%[4]s" "%[3]s" {
	name                       = "%[2]s"
	defaults_from              = "/%[1]s/request-log"
	request_logging            = "enabled"
	requestlog_protocol        = "mds-udp"
	requestlog_error_protocol  = "mds-udp"
	responselog_protocol       = "mds-udp"
	responselog_error_protocol = "mds-udp"
}`, partition, profileName, resourceName, resRequestLogName)
}
