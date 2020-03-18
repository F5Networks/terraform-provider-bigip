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

#Loading from a file is the preferred method
resource "bigip_ltm_irule" "rule" {
  name  = "/Common/terraform_irule"
  irule = file("myirule.tcl")
}

resource "bigip_ltm_monitor" "monitor" {
  name     = "/Common/terraform_monitor"
  parent   = "/Common/http"
  send     = "GET /some/path\r\n"
  timeout  = "999"
  interval = "999"
}

resource "bigip_ltm_pool" "pool" {
  name                = "/Common/terraform-pool"
  load_balancing_mode = "round-robin"

  #nodes = ["11.1.1.101:80", "11.1.1.102:80"]
  monitors   = ["/Common/terraform_monitor"]
  allow_snat = true
}

resource "bigip_ltm_virtual_server" "http" {
  pool                       = "/Common/terraform-pool"
  name                       = "/Common/terraform_vs_http"
  destination                = "100.1.1.100"
  port                       = 80
  ip_protocol                = "tcp"
  profiles                   = ["/Common/sanjay-http"]
  irules                     = ["/Common/terraform_irule"]
  source_address_translation = "automap"
  depends_on                 = [bigip_ltm_pool.pool]
}

