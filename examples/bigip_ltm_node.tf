/*
Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */
provider "bigip" {
  address  = "xxx.xxx.xxx.xxx"
  username = "xxxxx"
  password = "xxxxx"
}

resource "bigip_ltm_node" "node" {
  name    = "/Common/terraform_node1"
  address = "10.10.10.10"
}

