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

resource "bigip_ltm_monitor" "monitor" {
  name       = "/Common/terraform_monitor"
  parent     = "/Common/http"
  send       = "GET /some/path\r\n"
  timeout    = "999"
  interval   = "999"
  depends_on = [bigip_sys_provision.provision-afm]
}

