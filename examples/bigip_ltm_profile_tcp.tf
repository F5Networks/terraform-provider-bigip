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

resource "bigip_ltm_profile_tcp" "sanjose-tcp-lan-profile" {
  name               = "sanjose-tcp-lan-profile"
  idle_timeout       = 200
  close_wait_timeout = 5
  finwait_2timeout   = 5
  finwait_timeout    = 300
  keepalive_interval = 1700
  deferred_accept    = "enabled"
  fast_open          = "enabled"
}

