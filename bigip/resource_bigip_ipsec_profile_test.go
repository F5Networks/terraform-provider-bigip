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

func TestAccBigipNetIPsecProfile_create(t *testing.T) {
	t.Parallel()
	resName = "bigip_ipsec_profile"
	var instName = "test-ipsec-policy"
	var instFullName = fmt.Sprintf("/%s/%s", TestPartition, instName)
	resFullName := fmt.Sprintf("%s.%s", resName, instName)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckIPSecProfileDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testaccNetIpsecProfileDefaultCreate(instName),
				Check: resource.ComposeTestCheckFunc(
					testCheckIPSecProfileExists(instFullName, true),
					resource.TestCheckResourceAttr(resFullName, "name", instFullName),
				),
			},
		},
	})
}

func testCheckIPSecProfileExists(name string, exists bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*bigip.BigIP)
		p, err := client.GetIPSecPolicy(name)
		if err != nil {
			return err
		}
		if exists && p == nil {
			return fmt.Errorf(" IPSec Policy %s was not created.", name)
		}
		if !exists && p == nil {
			return fmt.Errorf(" IPSec Policy %s still exists.", name)
		}
		return nil
	}
}

func testCheckIPSecProfileDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(*bigip.BigIP)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != resName {
			continue
		}
		name := rs.Primary.ID
		ipsec, err := client.GetIPSecProfile(name)
		if err != nil {
			return err
		}
		if ipsec.Name != "" {
			return fmt.Errorf(" IPSec Policy %s not destroyed.", name)
		}
	}
	return nil
}
func testaccNetIpsecProfileDefaultCreate(instName string) string {
	return fmt.Sprintf(`
resource "%[1]s" "%[2]s" {
  name = "/Common/%[2]s"
}
		`, resName, instName)
}
