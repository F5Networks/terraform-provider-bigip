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

func testCheckSysSyslogRemoteServer(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*bigip.BigIP)
		var cfg sysSyslogConfig
		found, err := restGet(client, "sys/syslog", &cfg)
		if err != nil {
			return err
		}
		if !found {
			return fmt.Errorf("sys/syslog config not found")
		}
		for _, rs := range cfg.RemoteServers {
			if stripPartitionPrefix(rs.Name) == name {
				return nil
			}
		}
		return fmt.Errorf("remote syslog server %q not found in %+v", name, cfg.RemoteServers)
	}
}

func TestAccBigipSysSyslog_create(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAcctPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: `
resource "bigip_sys_syslog" "test-syslog" {
  auth_priv_from = "notice"
  auth_priv_to   = "emerg"
  console_log    = "enabled"
  remote_servers {
    name        = "test-syslog-remote"
    host        = "192.0.2.10"
    remote_port = 514
  }
}
`,
				Check: resource.ComposeTestCheckFunc(
					testCheckSysSyslogRemoteServer("test-syslog-remote"),
					resource.TestCheckResourceAttr("bigip_sys_syslog.test-syslog", "auth_priv_from", "notice"),
					resource.TestCheckResourceAttr("bigip_sys_syslog.test-syslog", "remote_servers.0.host", "192.0.2.10"),
				),
			},
		},
	})
}

func TestAccBigipSysSyslog_import(t *testing.T) {
	resourceName := "bigip_sys_syslog.test-syslog"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAcctPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: `
resource "bigip_sys_syslog" "test-syslog" {
  auth_priv_from = "notice"
  auth_priv_to   = "emerg"
}
`,
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateId:     "syslog",
				ImportStateVerify: true,
			},
		},
	})
}
