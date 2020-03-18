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

resource "bigip_ltm_snat" "snat_list" {
  name        = "NewSnatList"
  translation = "136.1.1.1"
  origins {
    name = "2.2.2.2"
  }
  origins {
    name = "3.3.3.3"
  }
}

