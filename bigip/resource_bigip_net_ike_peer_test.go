/*
Copyright 2019 F5 Networks Inc.
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

var TEST_IKE_PEER_NAME = "/Common/testpeer"

var TEST_IKE_PEER_RESOURCE = `
resource "bigip_net_ike_peer"  "test_ike_peer" {
    name = "` + TEST_IKE_PEER_NAME + `"
    dpd_delay                      = 30
    generate_policy                = "off"
    lifetime                       = 1440
    mode                           = "main"
    my_cert_file                   = "/Common/default.crt"
    my_cert_key_file               = "/Common/default.key"
    my_id_type                     = "address"
    nat_traversal                  = "off"
    passive                        = "false"
    peers_cert_type                = "none"
    peers_id_type                  = "address"
    phase1_auth_method             = "rsa-signature"
    phase1_encrypt_algorithm       = "3des"
    phase1_hash_algorithm          = "sha256"
    phase1_perfect_forward_secrecy = "modp1024"
    prf                            = "sha256"
    proxy_support                  = "enabled"
    remote_address                 = "1.5.3.4"
    replay_window_size             = 64
    state                          = "enabled"
    verify_cert                    = "false"
}
`

func TestAccBigipNetIkePeerCreate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: resource.ComposeTestCheckFunc(testCheckBigipNetIkePeerDestroyed),
		Steps: []resource.TestStep{
			{
				Config: TEST_IKE_PEER_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testBigipNetIkePeerExists(TEST_IKE_PEER_NAME, true),
					resource.TestCheckResourceAttr("bigip_net_ike_peer.test_ike_peer", "name", TEST_IKE_PEER_NAME),
					resource.TestCheckResourceAttr("bigip_net_ike_peer.test_ike_peer", "dpd_delay", "30"),
					resource.TestCheckResourceAttr("bigip_net_ike_peer.test_ike_peer", "generate_policy", "off"),
					resource.TestCheckResourceAttr("bigip_net_ike_peer.test_ike_peer", "lifetime", "1440"),
					resource.TestCheckResourceAttr("bigip_net_ike_peer.test_ike_peer", "mode", "main"),
					resource.TestCheckResourceAttr("bigip_net_ike_peer.test_ike_peer", "my_cert_file", "/Common/default.crt"),
					resource.TestCheckResourceAttr("bigip_net_ike_peer.test_ike_peer", "my_cert_key_file", "/Common/default.key"),
					resource.TestCheckResourceAttr("bigip_net_ike_peer.test_ike_peer", "my_id_type", "address"),
					resource.TestCheckResourceAttr("bigip_net_ike_peer.test_ike_peer", "nat_traversal", "off"),
					resource.TestCheckResourceAttr("bigip_net_ike_peer.test_ike_peer", "passive", "false"),
					resource.TestCheckResourceAttr("bigip_net_ike_peer.test_ike_peer", "peers_cert_type", "none"),
					resource.TestCheckResourceAttr("bigip_net_ike_peer.test_ike_peer", "peers_id_type", "address"),
					resource.TestCheckResourceAttr("bigip_net_ike_peer.test_ike_peer", "phase1_auth_method", "rsa-signature"),
					resource.TestCheckResourceAttr("bigip_net_ike_peer.test_ike_peer", "phase1_encrypt_algorithm", "3des"),
					resource.TestCheckResourceAttr("bigip_net_ike_peer.test_ike_peer", "phase1_hash_algorithm", "sha256"),
					resource.TestCheckResourceAttr("bigip_net_ike_peer.test_ike_peer", "phase1_perfect_forward_secrecy", "modp1024"),
					resource.TestCheckResourceAttr("bigip_net_ike_peer.test_ike_peer", "prf", "sha256"),
					resource.TestCheckResourceAttr("bigip_net_ike_peer.test_ike_peer", "proxy_support", "enabled"),
					resource.TestCheckResourceAttr("bigip_net_ike_peer.test_ike_peer", "remote_address", "1.5.3.4"),
					resource.TestCheckResourceAttr("bigip_net_ike_peer.test_ike_peer", "replay_window_size", "64"),
					resource.TestCheckResourceAttr("bigip_net_ike_peer.test_ike_peer", "state", "enabled"),
					resource.TestCheckResourceAttr("bigip_net_ike_peer.test_ike_peer", "verify_cert", "false"),
				),
			},
		},
	})

}
func TestAccBigipNetIkePeerImport(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckBigipNetIkePeerDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TEST_IKE_PEER_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testBigipNetIkePeerExists(TEST_IKE_PEER_NAME, true),
				),
				ResourceName:      TEST_IKE_PEER_NAME,
				ImportState:       false,
				ImportStateVerify: true,
			},
		},
	})
}

func testBigipNetIkePeerExists(name string, exists bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*bigip.BigIP)

		pp, err := client.GetIkePeer(name)
		if err != nil {
			return err
		}
		if exists && pp == nil {
			return fmt.Errorf("IkePeer %s does not exist.", name)
		}
		if !exists && pp != nil {
			return fmt.Errorf("IkePeer %s exists.", name)
		}
		return nil
	}
}

func testCheckBigipNetIkePeerDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(*bigip.BigIP)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bigip_net_ike_peer" {
			continue
		}

		name := rs.Primary.ID
		pp, err := client.GetIkePeer(name)
		if err != nil {
			return err
		}

		if pp != nil {
			return fmt.Errorf("IkePeer %s not destroyed.", name)
		}
	}
	return nil
}
