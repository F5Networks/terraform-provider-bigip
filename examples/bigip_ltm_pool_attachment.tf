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
  interval = "998"
}

resource "bigip_ltm_pool" "pool" {
  name                = "/Common/terraform-pool"
  load_balancing_mode = "round-robin"
  monitors            = ["/Common/terraform_monitor"]
  allow_snat          = "yes"
  allow_nat           = "yes"
  depends_on          = [bigip_ltm_monitor.monitor]
}

resource "bigip_ltm_node" "node" {
  name    = "/Common/terraform_node1"
  address = "192.168.30.2"
}

resource "bigip_ltm_pool_attachment" "attach_node" {
  pool       = "/Common/terraform-pool"
  node       = "${bigip_ltm_node.node.name}:80"
  depends_on = [bigip_ltm_pool.pool]
}


