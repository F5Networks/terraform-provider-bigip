/*
Original work from https://github.com/DealerDotCom/terraform-provider-bigip
Modifications Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file,You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccBigipLtmWafPolicyTestCases(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testMonitorsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: loadFixtureString("../examples/awaf/awaftest_issue822.tf"),
				Check:  resource.ComposeTestCheckFunc(
				// testCheckMonitorExists("/Common/test_monitor_tc1"),
				// testCheckMonitorExists("/Common/test_monitor_tc2"),
				// testCheckMonitorExists("/Common/test_monitor_tc3"),
				// testCheckMonitorExists("/Common/test_monitor_tc4"),
				// testCheckMonitorExists("/Common/test_monitor_tc5"),
				),
			},
		},
	})
}

// func testCheckMonitorExists(name string) resource.TestCheckFunc {
// 	return func(s *terraform.State) error {
// 		client := testAccProvider.Meta().(*bigip.BigIP)

// 		monitors, err := client.Monitors()
// 		if err != nil {
// 			return err
// 		}

// 		for _, m := range monitors {
// 			if m.FullPath == name {
// 				return nil
// 			}
// 		}
// 		return fmt.Errorf("Monitor %s was not created ", name)
// 	}
// }

// func testMonitorsDestroyed(s *terraform.State) error {
// 	client := testAccProvider.Meta().(*bigip.BigIP)

// 	monitors, err := client.Monitors()
// 	if err != nil {
// 		return err
// 	}

// 	for _, rs := range s.RootModule().Resources {
// 		if rs.Type != "bigip_ltm_monitor" {
// 			continue
// 		}

// 		name := rs.Primary.ID
// 		for _, m := range monitors {
// 			if m.FullPath == name {
// 				return fmt.Errorf("Monitor %s not destroyed ", name)
// 			}
// 		}
// 	}
// 	return nil
// }
