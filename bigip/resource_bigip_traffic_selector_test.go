/*
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/

package bigip

import (
	"fmt"
	"testing"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccBigipnetIPsecTrafficselector_Defaultcreate(t *testing.T) {
	t.Parallel()
	resName = "bigip_traffic_selector"
	var instName = "test-selector"
	var instFullName = fmt.Sprintf("/%s/%s", TestPartition, instName)
	resFullName := fmt.Sprintf("%s.%s", resName, instName)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckIPSectsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testaccNetIpsecTsDefaultcreate(instName),
				Check: resource.ComposeTestCheckFunc(
					testCheckIPSectsExists(instFullName, true),
					resource.TestCheckResourceAttr(resFullName, "name", instFullName),
				),
			},
		},
	})
}

func testCheckIPSectsExists(name string, exists bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*bigip.BigIP)
		p, err := client.GetTrafficselctor(name)
		if err != nil {
			return err
		}
		if exists && p == nil {
			return fmt.Errorf(" IPSec traffic-selector %s was not created.", name)
		}
		if !exists && p == nil {
			return fmt.Errorf(" IPSec traffic-selector %s still exists.", name)
		}
		return nil
	}
}

func testCheckIPSectsDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(*bigip.BigIP)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != resName {
			continue
		}
		name := rs.Primary.ID
		ipsecTs, err := client.GetTrafficselctor(name)
		if err != nil {
			return err
		}
		if ipsecTs.Name != "" {
			return fmt.Errorf(" IPSec traffic-selector %s not destroyed.", name)
		}
	}
	return nil
}
func testaccNetIpsecTsDefaultcreate(instName string) string {
	return fmt.Sprintf(`
resource "%[1]s" "%[2]s" {
  name                = "/Common/%[2]s"
  destination_address = "3.10.11.2/32"
  source_address      = "2.10.11.12/32"
}
		`, resName, instName)
}
