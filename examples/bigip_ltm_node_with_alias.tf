/*
Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */
provider "bigip" {
  address  = "xxx.xxx.xxx.xxx"
  alias    = "east"
  username = "xxxx"
  password = "xxxx"
}

provider "bigip" {
  alias    = "west"
  address  = "xxx.xxx.xxx.xxx"
  username = "xxxx"
  password = "xxxx"
}

resource "bigip_ltm_node" "node_west" {
  name     = "/Common/terraform_node1"
  provider = bigip.west
  address  = "1.1.1.1"
  state    = "user-up"
}

resource "bigip_ltm_node" "node_east" {
  name     = "/Common/terraform_node1"
  provider = bigip.east
  address  = "1.1.1.1"
  state    = "user-down"
}

