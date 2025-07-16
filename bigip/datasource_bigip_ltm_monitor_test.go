package bigip

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceBigipLtmMonitor_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAcctPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceBigipLtmMonitorConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.bigip_ltm_monitor.test", "name"),
					resource.TestCheckResourceAttrSet("data.bigip_ltm_monitor.test", "partition"),
					resource.TestCheckResourceAttrSet("data.bigip_ltm_monitor.test", "interval"),
					resource.TestCheckResourceAttrSet("data.bigip_ltm_monitor.test", "timeout"),
				),
			},
		},
	})
}

func testAccDataSourceBigipLtmMonitorConfig() string {
	return `
	data "bigip_ltm_monitor" "test" {
		name      = "http"
		partition = "Common"
	}
	`
}
