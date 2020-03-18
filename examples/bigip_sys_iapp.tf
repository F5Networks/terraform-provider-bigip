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

resource "bigip_sys_iapp" "waf_asm" {
  name     = "policywaf"
  jsonfile = file("policywaf.json")
}

resource "bigip_sys_iapp" "pool_deployed" {
  name     = "sap-dmzpool-rp1-80"
  jsonfile = file("sap-dmzpool-rp1-80.json")
}

