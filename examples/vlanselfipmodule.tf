/*
Original work from https://github.com/DealerDotCom/terraform-provider-bigip
Modifications Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file,You can obtain one at https://mozilla.org/MPL/2.0/.
 */
provider "bigip" {
  address  = "xxx.xxx.xxx.xxx"
  username = "xxxx"
  password = "xxxx"
}

module "sjvlan1" {
  source   = "./vlanmodule"
  name     = "/Common/intvlan"
  tag      = 101
  vlanport = "1.1"
  tagged   = true
}

resource "bigip_net_selfip" "selfip" {
  name       = "/Common/InternalselfIP"
  ip         = "100.1.1.1/24"
  vlan       = "/Common/intvlan"
  depends_on = [module.sjvlan1]
}

