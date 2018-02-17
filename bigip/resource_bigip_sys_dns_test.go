package bigip

import (
	"fmt"
	"testing"

	"github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

var TEST_DNS_NAME = fmt.Sprintf("/%s/test-dns", TEST_PARTITION)

var TEST_DNS_RESOURCE = `
resource "bigip_sys_dns" "test-dns" {
   description = "` + TEST_DNS_NAME + `"
   name_servers = ["1.1.1.1"]
   numberof_dots = 2
   search = ["f5.com"]
}

`

func TestAccBigipSysdns_create(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers: testAccProviders,
		//CheckDestroy: testCheckdnssDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TEST_DNS_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckdnsExists(TEST_DNS_NAME, true),
					resource.TestCheckResourceAttr("bigip_sys_dns.test-dns", "description", TEST_DNS_NAME),
					resource.TestCheckResourceAttr("bigip_sys_dns.test-dns", "numberof_dots", "2"),
					resource.TestCheckResourceAttr("bigip_sys_dns.test-dns",
						fmt.Sprintf("name_servers.%d", schema.HashString("1.1.1.1")),
						"1.1.1.1"),
					resource.TestCheckResourceAttr("bigip_sys_dns.test-dns",
						fmt.Sprintf("search.%d", schema.HashString("f5.com")),
						"f5.com"),
				),
			},
		},
	})
}

func TestAccBigipSysdns_import(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers: testAccProviders,
		//	CheckDestroy: testCheckdnssDestroyed, ( No Delet API support)
		Steps: []resource.TestStep{
			{
				Config: TEST_DNS_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckdnsExists(TEST_DNS_NAME, true),
				),
				ResourceName:      TEST_DNS_NAME,
				ImportState:       false,
				ImportStateVerify: true,
			},
		},
	})
}

func testCheckdnsExists(description string, exists bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*bigip.BigIP)

		dns, err := client.DNSs()
		if err != nil {
			return err
		}
		if exists && dns == nil {
			return fmt.Errorf("dns %s was not created.", description)

		}
		if !exists && dns != nil {
			return fmt.Errorf("dns %s still exists.", description)

		}
		return nil
	}
}

func testCheckdnssDestroyed(s *terraform.State) error {
	/* client := testAccProvider.Meta().(*bigip.BigIP)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bigip_sys_dns" {
			continue
		}

		description := rs.Primary.ID
		dns, err := client.dnss()
		if err != nil {
			return err
		}
		if dns != nil {
			return fmt.Errorf("dns ", description, " not destroyed.")

		}
	}*/
	return nil
}
