/*
Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

resource "bigip_ltm_monitor" "pa_tc1" {
  name     = "/Common/test_monitor_pa_tc1"
  parent   = "/Common/http"
  send     = "GET /some/path\r\n"
  timeout  = "999"
  interval = "998"
}

resource "bigip_ltm_pool" "pa_tc1" {
  name                = "/Common/test_pool_pa_tc1"
  load_balancing_mode = "round-robin"
  monitors            = [bigip_ltm_monitor.pa_tc1.name]
  allow_snat          = "yes"
  allow_nat           = "yes"
}

resource "bigip_ltm_node" "pa_tc1" {
  name    = "/Common/test_node_pa_tc1"
  address = "192.168.30.2"
}

resource "bigip_ltm_pool_attachment" "pa_tc1" {
  pool = bigip_ltm_pool.pa_tc1.name
  node = "${bigip_ltm_node.pa_tc1.name}:80"
}

resource "bigip_ltm_pool_attachment" "pa_tc2" {
  pool = bigip_ltm_pool.pa_tc1.name
  node = "192.168.30.1:80"
}

resource "bigip_ltm_pool_attachment" "pa_tc3" {
  pool = bigip_ltm_pool.pa_tc1.name
  node = "10.10.10.11:80"
}

resource "bigip_ltm_pool_attachment" "pa_tc4" {
  pool = bigip_ltm_pool.pa_tc1.name
  node = "test2.com:80"
}

resource "bigip_ltm_node" "pa_tc5" {
  name    = "/Common/test_node_pa_tc5"
  address = "192.168.30.5"
}

resource "bigip_ltm_pool_attachment" "pa_tc5" {
  pool = bigip_ltm_pool.pa_tc1.name
  node = format("%s:80", bigip_ltm_node.pa_tc5.name)
}

resource "bigip_ltm_node" "pa_tc6" {
  name    = "/Common/test3.com"
  address = "test3.com"
}
resource "bigip_ltm_pool_attachment" "pa_tc6" {
  pool = bigip_ltm_pool.pa_tc1.name
  node = format("%s:80", bigip_ltm_node.pa_tc6.name)
}

resource "bigip_ltm_pool_attachment" "pa_tc7" {
  pool = bigip_ltm_pool.pa_tc1.name
  node = "1.1.11.2:80"
  //node                  = ""
  ratio                 = 2
  connection_limit      = 2
  connection_rate_limit = 2
  priority_group        = 2
  dynamic_ratio         = 3
}

resource "bigip_ltm_pool_attachment" "pa_tc8" {
  pool = bigip_ltm_pool.pa_tc1.name
  node = "1.1.12.2:80"
}

resource "bigip_ltm_pool" "pa_tc9" {
  name = "/Common/test_pool_pa_tc9"
}

resource "bigip_ltm_pool_attachment" "pa_tc9" {
  pool              = bigip_ltm_pool.pa_tc9.name
  node              = "facebook.com:80"
  fqdn_autopopulate = "disabled"
}

resource "bigip_command" "pa_tc10" {
  commands = ["create auth partition TEST3", "create net route-domain /TEST3/testdomain id 50"]
  when     = "apply"
}

locals {
  partition_name = "TEST3"
}

resource "bigip_ltm_pool" "pa_tc10" {
  name       = "/${local.partition_name}/test_pool_pa_tc10"
  depends_on = [bigip_command.pa_tc10]
}

resource "bigip_ltm_pool_attachment" "pa_tc10" {
  pool             = bigip_ltm_pool.pa_tc10.name
  node             = "2.3.2.2%50:8080"
  connection_limit = 11
}

resource "bigip_ltm_pool_attachment" "pa_tc11" {
  pool             = bigip_ltm_pool.pa_tc10.name
  node             = "10.1.100.1:8080"
  connection_limit = 11
}

//resource "bigip_ltm_pool" "pool" {
//  name                = "/Common/Axiom_Environment_APP1_Pool"
//  load_balancing_mode = "round-robin"
//  //minimum_active_members = 1
//  monitors = [bigip_ltm_monitor.monitor.name]
//}