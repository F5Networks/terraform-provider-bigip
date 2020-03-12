/*
Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */
provider "bigip" {
  address  = "xxx.xxx.xxx.xxx"
  username = "xxxx"
  password = "xxxxx"
}

resource "bigip_fast" "fasttest-example1" {
  name             = "examples/simple_http"
  tenant_name      = "fasttest10"
  application_name = "fasttest10app"
  virtual_port     = 80
  virtual_address  = "20.10.20.1"
  server_port      = 8080
  server_addresses = ["20.10.40.1", "20.10.50.1"]
}