package bigip

import (
	"fmt"
	"testing"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var resSaasBotDefenseName = "bigip_saas_bot_defense_profile"

func TestAccBigipSaasBotDefenseProfileTC1(t *testing.T) {
	t.Parallel()
	var instName = "test-saas-bot-defense-tc1"
	var TestBotDefenseName = fmt.Sprintf("/%s/%s", TestPartition, instName)
	resFullName := fmt.Sprintf("%s.%s", resSaasBotDefenseName, instName)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckSaasBotDefensesDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testaccbigipSaasBotDefenseDefaultConfig(TestPartition, TestBotDefenseName, instName),
				Check: resource.ComposeTestCheckFunc(
					testCheckSaasBotDefenseExists(TestBotDefenseName),
					resource.TestCheckResourceAttr(resFullName, "name", TestBotDefenseName),
					resource.TestCheckResourceAttr(resFullName, "defaults_from", "/Common/bd"),
				),
				Destroy: false,
			},
		},
	})
}

func testCheckSaasBotDefenseExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*bigip.BigIP)
		p, err := client.GetSaasBotDefenseProfile(name)
		if err != nil {
			return err
		}
		if p == nil {
			return fmt.Errorf("BotDefense %s was not created ", name)
		}

		return nil
	}
}

func testCheckSaasBotDefensesDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(*bigip.BigIP)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bigip_saas_bot_defense_profile" {
			continue
		}

		name := rs.Primary.ID
		BotDefense, err := client.GetSaasBotDefenseProfile(name)
		if err != nil {
			return nil
		}
		if BotDefense != nil {
			return fmt.Errorf("BotDefense %s not destroyed. ", name)
		}
	}
	return nil
}

func testaccbigipSaasBotDefenseDefaultConfig(partition, profileName, resourceName string) string {
	return fmt.Sprintf(`resource "bigip_saas_bot_defense_profile" "%[3]s" {
		name                  = "%[2]s"
		application_id        = "89fb0bfcb4bf4c578fad9adb37ce3b19"
		tenant_id             = "a-aavN9vaYOV"
		api_key               = "49840d1dd6fa4c4d86c88762eb398eee"
		shape_protection_pool = "/%[1]s/cs1.pool"
		ssl_profile           = "/%[1]s/cloud-service-default-ssl"
		protected_endpoints {
			name     = "pe1"
			host     = "abc.com"
			endpoint = "/login"
			post     = "enabled"
		}
}`, partition, profileName, resourceName)
}
