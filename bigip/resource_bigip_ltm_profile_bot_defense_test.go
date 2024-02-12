package bigip

import (
	"fmt"
	"testing"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var resBotDefenseName = "bigip_ltm_profile_bot_defense"

func TestAccBigipLtmProfileBotDefenseTC1(t *testing.T) {
	t.Parallel()
	var instName = "test-bot-defense-tc1"
	var TestBotDefenseName = fmt.Sprintf("/%s/%s", TestPartition, instName)
	resFullName := fmt.Sprintf("%s.%s", resBotDefenseName, instName)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckBotDefensesDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testaccbigipltmprofileBotDefenseDefaultConfig(TestPartition, TestBotDefenseName, instName),
				Check: resource.ComposeTestCheckFunc(
					testCheckBotDefenseExists(TestBotDefenseName),
					resource.TestCheckResourceAttr(resFullName, "name", TestBotDefenseName),
					resource.TestCheckResourceAttr(resFullName, "defaults_from", "/Common/bot-defense"),
				),
				Destroy: false,
			},
		},
	})
}

func testCheckBotDefenseExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*bigip.BigIP)
		p, err := client.GetBotDefenseProfile(name)
		if err != nil {
			return err
		}
		if p == nil {
			return fmt.Errorf("BotDefense %s was not created ", name)
		}

		return nil
	}
}

func testCheckBotDefensesDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(*bigip.BigIP)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bigip_ltm_profile_bot_defence" {
			continue
		}

		name := rs.Primary.ID
		BotDefense, err := client.GetBotDefenseProfile(name)
		if err != nil {
			return err
		}
		if BotDefense != nil {
			return fmt.Errorf("BotDefense %s not destroyed. ", name)
		}
	}
	return nil
}

func testaccbigipltmprofileBotDefenseDefaultConfig(partition, profileName, resourceName string) string {
	return fmt.Sprintf(`resource "bigip_ltm_profile_bot_defence" "%[3]s" {
		name = "%[2]s"
		defaults_from = "/%[1]s/bot-defense"
		description = "test-bot"
		template = "relaxed"
	}`, partition, profileName, resourceName)
}
