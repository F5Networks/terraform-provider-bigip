package bigip

import (
	"fmt"
	"log"
	"testing"

	"github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

var TEST_IAPP_NAME = "/" + TEST_PARTITION + "/test-iapp"


var TEST_IAPP_RESOURCE = `
	resource "bigip_sys_iapp" "test-iapp" {
		name = "test-iapp"
		jsonfile = <<EOF
		{
			"name":"test-iapp",
			"partition": "Common",
			"inheritedDevicegroup": "true",
			"inheritedTrafficGroup": "true",
			"strictUpdates": "enabled",
			"template": "/Common/appsvcs_integration_v2.0_001_sap",
			"execute-action": "definition",
			"lists": [
			  {
				"name": "vs__BundledItems",
				"encrypted": "no",
				"value": [ ]
			  }
			],
			"tables": [
			  {
				"name": "feature__easyL4FirewallBlacklist",
				"columnNames": [ "CIDRRange" ],
				"rows": [ ]
			  },
			  {
				"name": "feature__easyL4FirewallSourceList",
				"columnNames": [ "CIDRRange" ],
				"rows": [
				  { "row": [ "0.0.0.0/0" ] }
				]
			  },
			  {
				"name": "l7policy__rulesAction",
				"columnNames": [ "Group", "Target", "Parameter" ],
				"rows": [ ]
			  },
			  {
				"name": "l7policy__rulesMatch",
				"columnNames": [ "Group", "Operand", "Negate", "Condition", "Value", "CaseSensitive", "Missing" ],
				"rows": [ ]
			  },
			  {
				"name": "monitor__Monitors",
				"columnNames": [ "Index", "Name", "Type", "Options" ],
				"rows": [
				  { "row": [ "0", "/Common/http", "none", "none" ] }
				]
			  },
			  {
				"name": "pool__Members",
				"columnNames": [ "Index", "IPAddress", "Port", "ConnectionLimit", "Ratio", "PriorityGroup", "State", "AdvOptions" ],
				"rows": [
				  { "row": [ "0", "192.168.32.32", "80", "0", "1", "0", "enabled", "none" ] }
				]
			  },
			  {
				"name": "pool__Pools",
				"columnNames": [ "Index", "Name", "Description", "LbMethod", "Monitor", "AdvOptions" ],
				"rows": [
				  { "row": [ "0", "", "", "round-robin", "0", "none" ] }
				]
			  },
			  {
				"name": "vs__Listeners",
				"columnNames": [ "Listener", "Destination" ],
				"rows": [ ]
			  }
			],
			"variables": [
			  {
				"name": "extensions__Field1",
				"encrypted": "no"
			  },
			  {
				"name": "extensions__Field2",
				"encrypted": "no"
			  },
			  {
				"name": "extensions__Field3",
				"encrypted": "no"
			  },
			  {
				"name": "feature__easyL4Firewall",
				"encrypted": "no",
				"value": "disabled"
			  },
			  {
				"name": "feature__insertXForwardedFor",
				"encrypted": "no",
				"value": "disabled"
			  },
			  {
				"name": "feature__redirectToHTTPS",
				"encrypted": "no",
				"value": "disabled"
			  },
			  {
				"name": "feature__securityEnableHSTS",
				"encrypted": "no",
				"value": "disabled"
			  },
			  {
				"name": "feature__sslEasyCipher",
				"encrypted": "no",
				"value": "disabled"
			  },
			  {
				"name": "feature__statsHTTP",
				"encrypted": "no",
				"value": "auto"
			  },
			  {
				"name": "feature__statsTLS",
				"encrypted": "no",
				"value": "auto"
			  },
			  {
				"name": "iapp__apmDeployMode",
				"encrypted": "no",
				"value": "preserve-bypass"
			  },
			  {
				"name": "iapp__appStats",
				"encrypted": "no",
				"value": "enabled"
			  },
			  {
				"name": "iapp__asmDeployMode",
				"encrypted": "no",
				"value": "preserve-bypass"
			  },
			  {
				"name": "iapp__logLevel",
				"encrypted": "no",
				"value": "7"
			  },
			  {
				"name": "iapp__mode",
				"encrypted": "no",
				"value": "auto"
			  },
			  {
				"name": "iapp__routeDomain",
				"encrypted": "no",
				"value": "auto"
			  },
			  {
				"name": "iapp__strictUpdates",
				"encrypted": "no",
				"value": "enabled"
			  },
			  {
				"name": "l7policy__defaultASM",
				"encrypted": "no",
				"value": "bypass"
			  },
			  {
				"name": "l7policy__defaultL7DOS",
				"encrypted": "no",
				"value": "bypass"
			  },
			  {
				"name": "l7policy__strategy",
				"encrypted": "no",
				"value": "/Common/first-match"
			  },
			  {
				"name": "pool__DefaultPoolIndex",
				"encrypted": "no",
				"value": "0"
			  },
			  {
				"name": "pool__MemberDefaultPort",
				"encrypted": "no",
				"value": "80"
			  },
			  {
				"name": "pool__addr",
				"encrypted": "no",
				"value": "255.255.255.254"
			  },
			  {
				"name": "pool__mask",
				"encrypted": "no",
				"value": "255.255.255.255"
			  },
			  {
				"name": "pool__port",
				"encrypted": "no",
				"value": "443"
			  },
			  {
				"name": "vs__AdvOptions",
				"encrypted": "no",
				"value": ""
			  },
			  {
				"name": "vs__AdvPolicies",
				"encrypted": "no",
				"value": ""
			  },
			  {
				"name": "vs__AdvProfiles",
				"encrypted": "no",
				"value": ""
			  },
			  {
				"name": "vs__ConnectionLimit",
				"encrypted": "no",
				"value": "0"
			  },
			  {
				"name": "vs__Description",
				"encrypted": "no",
				"value": ""
			  },
			  {
				"name": "vs__IpProtocol",
				"encrypted": "no",
				"value": "tcp"
			  },
			  {
				"name": "vs__Irules",
				"encrypted": "no",
				"value": ""
			  },
			  {
				"name": "vs__Name",
				"encrypted": "no",
				"value": ""
			  },
			  {
				"name": "vs__OptionConnectionMirroring",
				"encrypted": "no",
				"value": "disabled"
			  },
			  {
				"name": "vs__OptionSourcePort",
				"encrypted": "no",
				"value": "preserve"
			  },
			  {
				"name": "vs__ProfileAccess",
				"encrypted": "no",
				"value": ""
			  },
			  {
				"name": "vs__ProfileAnalytics",
				"encrypted": "no",
				"value": ""
			  },
			  {
				"name": "vs__ProfileClientProtocol",
				"encrypted": "no",
				"value": ""
			  },
			  {
				"name": "vs__ProfileClientSSL",
				"encrypted": "no",
				"value": ""
			  },
			  {
				"name": "vs__ProfileClientSSLAdvOptions",
				"encrypted": "no",
				"value": ""
			  },
			  {
				"name": "vs__ProfileClientSSLCert",
				"encrypted": "no",
				"value": ""
			  },
			  {
				"name": "vs__ProfileClientSSLChain",
				"encrypted": "no",
				"value": ""
			  },
			  {
				"name": "vs__ProfileClientSSLCipherString",
				"encrypted": "no",
				"value": ""
			  },
			  {
				"name": "vs__ProfileClientSSLKey",
				"encrypted": "no",
				"value": ""
			  },
			  {
				"name": "vs__ProfileCompression",
				"encrypted": "no",
				"value": ""
			  },
			  {
				"name": "vs__ProfileConnectivity",
				"encrypted": "no",
				"value": ""
			  },
			  {
				"name": "vs__ProfileDefaultPersist",
				"encrypted": "no",
				"value": ""
			  },
			  {
				"name": "vs__ProfileFallbackPersist",
				"encrypted": "no",
				"value": ""
			  },
			  {
				"name": "vs__ProfileHTTP",
				"encrypted": "no",
				"value": ""
			  },
			  {
				"name": "vs__ProfileOneConnect",
				"encrypted": "no",
				"value": ""
			  },
			  {
				"name": "vs__ProfilePerRequest",
				"encrypted": "no",
				"value": ""
			  },
			  {
				"name": "vs__ProfileRequestLogging",
				"encrypted": "no",
				"value": ""
			  },
			  {
				"name": "vs__ProfileSecurityDoS",
				"encrypted": "no",
				"value": ""
			  },
			  {
				"name": "vs__ProfileSecurityIPBlacklist",
				"encrypted": "no",
				"value": "none"
			  },
			  {
				"name": "vs__ProfileSecurityLogProfiles",
				"encrypted": "no",
				"value": ""
			  },
			  {
				"name": "vs__ProfileServerProtocol",
				"encrypted": "no",
				"value": ""
			  },
			  {
				"name": "vs__ProfileServerSSL",
				"encrypted": "no",
				"value": ""
			  },
			  {
				"name": "vs__RouteAdv",
				"encrypted": "no",
				"value": "disabled"
			  },
			  {
				"name": "vs__SNATConfig",
				"encrypted": "no",
				"value": "automap"
			  },
			  {
				"name": "vs__SourceAddress",
				"encrypted": "no",
				"value": "0.0.0.0/0"
			  },
			  {
				"name": "vs__VirtualAddrAdvOptions",
				"encrypted": "no",
				"value": ""
			  }
			]
		  }

EOF
	}`

func TestAccBigipSysIapp_create(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckIappDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TEST_IAPP_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckIappExists(TEST_IAPP_NAME),
				),
			},
		},
	})
}

func TestAccBigipSysIapp_import(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckIappDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TEST_IAPP_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckIappExists(TEST_IAPP_NAME),
				),
				ResourceName:      TEST_IAPP_NAME,
				ImportState:       false,
				ImportStateVerify: true,
			},
		},
	})
}

func testCheckIappExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*bigip.BigIP)

		jsonfile, err := client.Iapp(name)
		log.Println(" I am here in Exists !!!!!!!!!!!!", name)
		if err != nil {
			return fmt.Errorf("Error while fetching iapp: %v", err)

		}
		body := s.RootModule().Resources["bigip_sys_iapp.test-iapp"].Primary.Attributes["name"]
		if jsonfile.Name == body {
			return fmt.Errorf("jsonfile  body does not match. Expecting %s got %s.", body, jsonfile.Name)
		}

		jsonfile_name := fmt.Sprintf("/%s/%s", jsonfile.Partition, jsonfile.Name)
		if jsonfile_name == name {
			return fmt.Errorf("Jsonfile name does not match. Expecting %s got %s.", name, jsonfile_name)
		}
		return nil
	}
}

func testCheckIappDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(*bigip.BigIP)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bigip_sys_iapp" {
			continue
		}

		name := rs.Primary.ID
		log.Println(" I am in Destroy function currently +++++++++++++++++++++++++++ ", name)

		// Join three strings into one.
		jsonfile, err := client.Iapp(name)

		if err != nil {
			return nil
		}

		if jsonfile == nil {
			return fmt.Errorf("Iapp %s not destroyed.", name)
		}
	}
	return nil
}
