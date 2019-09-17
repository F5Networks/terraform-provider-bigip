/*
Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */
provider "bigip" {
  address = "10.192.74.61"
  username = "admin"
  password = "admin"
}

resource "bigip_net_route" "route2" {
  name = "sanjay-route2"
  network = "10.10.10.0/24"
  gw      = "1.1.1.2"
}

