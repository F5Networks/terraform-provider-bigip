/*
Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */
provider "bigip" {
  address = "10.192.74.73"
  username = "admin"
  password = "admin"
}

resource "bigip_net_selfip" "selfip1" {
	name = "/Common/internalselfIP"
	ip = "11.1.1.1/24"
	vlan = "/Common/internal"
	depends_on = ["bigip_net_vlan.vlan1"]
	}

