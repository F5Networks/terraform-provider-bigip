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
  name     = "/Common/terraform_monitor"
  parent   = "/Common/http"
  send     = "GET /some/path\r\n"
  timeout  = "999"
  interval = "999"
}

resource "bigip_ltm_pool" "pool" {
  name                = "/Common/terraform-pool"
  load_balancing_mode = "round-robin"
  monitors            = ["/Common/terraform_monitor"]
  allow_snat          = "yes"
  allow_nat           = "yes"
}

resource "bigip_ltm_pool_attachment" "attach_node1" {
  pool       = "/Common/terraform-pool"
  node       = "/Common/11.1.1.101:80"
  depends_on = [bigip_ltm_pool.pool]
}

resource "bigip_ltm_pool_attachment" "attach_node2" {
  pool       = "/Common/terraform-pool"
  node       = "/Common/11.1.1.102:80"
  depends_on = [bigip_ltm_pool.pool]
}

resource "bigip_ltm_virtual_server" "http" {
  pool                       = "/Common/terraform-pool"
  name                       = "/Common/terraform_vs_http"
  destination                = "100.1.1.100"
  port                       = 80
  source_address_translation = "automap"
  depends_on                 = [bigip_ltm_pool.pool]
}

