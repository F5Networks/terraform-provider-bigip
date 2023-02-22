/*
Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */
resource "bigip_ltm_node" "TC1" {
  name    = "/Common/test_node_tc1"
  address = "10.10.10.10"
}
resource "bigip_ltm_node" "TC2" {
  name    = "/Common/test_node_tc2"
  address = "192.168.30.1"
}
resource "bigip_ltm_node" "TC3" {
  name    = "/Common/test_node_tc3"
  address = "192.168.30.2"
}
data "bigip_ltm_node" "TC2" {
  name      = split("/", bigip_ltm_node.TC2.name)[2]
  partition = split("/", bigip_ltm_node.TC2.name)[1]
}
resource "bigip_ltm_node" "TC4" {
  name             = "/Common/test_node_tc4"
  address          = "f5.com"
  connection_limit = "0"
  dynamic_ratio    = "1"
  monitor          = "default"
  rate_limit       = "disabled"
  fqdn { interval = "3000" }
  state = "user-up"
  ratio = "19"
}

resource "bigip_ltm_pool" "TC6" {
  name = "/Common/test_pool_tc6"
}

//resource "bigip_ltm_node" "TC5" {
//  for_each = toset(["3.3.3.3"])
//  name     = format("/%s/%s", "Common", each.value)
//  address  = each.value
//}
//resource "bigip_ltm_pool_attachment" "TC6" {
//  for_each   = toset([bigip_ltm_node.TC2.name, bigip_ltm_node.TC3.name])
//  pool       = bigip_ltm_pool.TC6.name
//  node       = "${each.key}:80"
//  depends_on = [bigip_ltm_node.TC2, bigip_ltm_node.TC3]
//}
resource "bigip_ltm_pool" "TC7" {
  name = "/Common/test_pool_tc7"
}
resource "bigip_ltm_pool_attachment" "TC7" {
  pool = bigip_ltm_pool.TC7.name
  node = format("%s:%s", bigip_ltm_node.TC2.name, 90)
}
resource "bigip_ltm_node" "TC8" {
  name    = "/Common/test_node_tc8"
  address = "test1.com"
  fqdn {
    address_family = "ipv4"
    interval       = "3000"
  }
}