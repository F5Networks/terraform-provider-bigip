/*
Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */
provider "bigip" {
  address  = "xxx.xxx.xxx.xxx"
  username = "xxxx"
  password = "xxxx"
}

resource "bigip_ltm_policy" "test-policy" {
  name           = "my_policy"
  strategy       = "first-match"
  requires       = ["http"]
  published_copy = "Drafts/my_policy"
  controls       = ["forwarding"]
  rule {
    name = "rule6"

    action {
      tm_name    = "20"
      forward    = true
      connection = false
      pool       = "/Common/mypool"
    }
  }
  depends_on = [bigip_ltm_pool.mypool]
}

resource "bigip_ltm_pool" "mypool" {
  name                = "/Common/mypool"
  monitors            = ["/Common/http"]
  allow_nat           = "yes"
  allow_snat          = "yes"
  load_balancing_mode = "round-robin"
}

