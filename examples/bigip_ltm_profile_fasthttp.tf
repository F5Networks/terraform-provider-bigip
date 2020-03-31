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

resource "bigip_ltm_profile_fasthttp" "sjfasthttpprofile" {
  name                         = "sjfasthttpprofile"
  defaults_from                = "/Common/fasthttp"
  idle_timeout                 = 300
  connpoolidle_timeoutoverride = 0
  connpool_maxreuse            = 2
  connpool_maxsize             = 2048
  connpool_minsize             = 0
  connpool_replenish           = "enabled"
  connpool_step                = 4
  forcehttp_10response         = "disabled"
  maxheader_size               = 32768
}

