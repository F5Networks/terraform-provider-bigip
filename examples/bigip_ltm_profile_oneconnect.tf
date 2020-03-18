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

resource "bigip_ltm_profile_oneconnect" "oneconnect-sanjose" {
  name                  = "sanjose"
  partition             = "Common"
  defaults_from         = "/Common/oneconnect"
  idle_timeout_override = "disabled"
  max_age               = 3600
  max_reuse             = 1000
  max_size              = 1000
  share_pools           = "disabled"
  source_mask           = "255.255.255.255"
}

