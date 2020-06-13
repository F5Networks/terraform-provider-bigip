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

module "sjvlan1" {
  source   = "./vlanmodule"
  name     = "/Common/intvlan"
  tag      = 101
  vlanport = "1.1"
  tagged   = true
}

module "sjvlan2" {
  source   = "./vlanmodule"
  name     = "/Common/extvlan"
  tag      = 102
  vlanport = "1.2"
  tagged   = true
}

