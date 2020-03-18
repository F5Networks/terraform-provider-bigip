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

resource "bigip_ltm_profile_http2" "nyhttp2" {
  name                              = "/Common/NewYork_http2"
  defaults_from                     = "/Common/http2"
  concurrent_streams_per_connection = 10
  connection_idle_timeout           = 30
  activation_modes                  = ["alpn", "npn"]
}

