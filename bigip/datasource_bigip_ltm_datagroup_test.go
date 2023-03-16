/*
Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/

package bigip

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccBigipLtmDataGroup_basic(t *testing.T) {
	t.Parallel()
	resName := "bigip_ltm_datagroup.DGTEST"
	dsName := "data.bigip_ltm_datagroup.DGTEST"
	var dataGroupName = "test-rg"
	var dataGroupFullName = fmt.Sprintf("/%s/%s", TestPartition, dataGroupName)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAcctPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckDataGroupDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDatagroupConfigBasic(dataGroupName),
				Check: resource.ComposeTestCheckFunc(
					testCheckDataGroupExists(dataGroupFullName),
					resource.TestCheckResourceAttr(resName, "name", dataGroupFullName),
				),
			},
			{
				Config: testAccCheckDatagroupConfigBasic(dataGroupName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dsName, "name", regexp.MustCompile(dataGroupName)),
				),
			},
		},
	})
}

func testAccCheckDatagroupConfigBasic(dataGroupName string) string {
	return fmt.Sprintf(`
resource "bigip_ltm_datagroup" "DGTEST" {
  name = "/%s/%s"
  type = "string"
  record {
    name = "test-name1"
    data = "test-data1"
  }
  record {
    name = "test-name2"
  }
}
data "bigip_ltm_datagroup" "DGTEST" {
  name       = "%s"
  partition  = "Common"
  depends_on = [bigip_ltm_datagroup.DGTEST]
}
`, "Common", dataGroupName, dataGroupName)
}
