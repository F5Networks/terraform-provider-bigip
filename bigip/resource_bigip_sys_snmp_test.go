package bigip

import (
	"fmt"
	"testing"

	"github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

var TEST_SNMP_NAME = fmt.Sprintf("/%s/test-snmp", TEST_PARTITION)

var TEST_SNMP_RESOURCE = `
resource "bigip_sys_snmp" "test-snmp" {
  sys_contact = "NetOPsAdmin s.shitole@f5.com"
  sys_location = "SeattleHQ"
  allowedaddresses = ["202.10.10.2"]
}
`

func TestBigipSyssnmp_create(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers: testAccProviders,
		//CheckDestroy: testChecksnmpsDestroyed, (delete API not supported )
		Steps: []resource.TestStep{
			{
				Config: TEST_SNMP_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testChecksnmpExists(TEST_SNMP_NAME, true),
					resource.TestCheckResourceAttr("bigip_sys_snmp.test-snmp", "sys_contact", "NetOPsAdmin s.shitole@f5.com"),
					resource.TestCheckResourceAttr("bigip_sys_snmp.test-snmp", "sys_location", "SeattleHQ"),
					resource.TestCheckResourceAttr("bigip_sys_snmp.test-snmp",
						fmt.Sprintf("allowedaddresses.%d", schema.HashString("202.10.10.2")),
						"202.10.10.2"),
				),
			},
		},
	})
}

func TestBigipSyssnmp_import(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers: testAccProviders,
		//CheckDestroy: testChecksnmpsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TEST_SNMP_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testChecksnmpExists(TEST_SNMP_NAME, true),
				),
				ResourceName:      TEST_SNMP_NAME,
				ImportState:       false,
				ImportStateVerify: true,
			},
		},
	})
}

func testChecksnmpExists(name string, exists bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*bigip.BigIP)
		p, err := client.SNMPs()
		if err != nil {
			return err
		}
		if exists && p == nil {
			return fmt.Errorf("snmp %s was not created.", name)
		}
		if !exists && p != nil {
			return fmt.Errorf("snmp %s still exists.", name)
		}
		return nil
	}
}

/*func testChecksnmpsDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(*bigip.BigIP)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bigip_sys_snmp" {
			continue
		}

		name := rs.Primary.ID
		snmp, err := client.snmps(name)
		if err != nil {
			return err
		}
		if snmp == nil {
			return fmt.Errorf("snmp ", name, " not destroyed.")
		}
	}
	return nil
}
*/
