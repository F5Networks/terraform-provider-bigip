/*
Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */
provider "bigip" {
  address = "10.192.74.68"
  username = "admin"
  password = "admin"
}


resource "bigip_ltm_node" "node" {
  name = "/Common/terraform_node1"
  address = "10.10.10.10"
}

