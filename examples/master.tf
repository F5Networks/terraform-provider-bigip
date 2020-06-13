/*
Original work from https://github.com/DealerDotCom/terraform-provider-bigip
Modifications Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file,You can obtain one at https://mozilla.org/MPL/2.0/.
 */
provider "bigip" {
  address  = "xxx.xxx.xxx.xxx"
  username = "xxxx"
  password = "xxxx"
}

resource "bigip_sys_ntp" "ntp1" {
  description = "/Common/NTP1"
  servers     = ["time.google.com"]
  timezone    = "America/Los_Angeles"
}

resource "bigip_sys_dns" "dns1" {
  description    = "/Common/DNS1"
  name_servers   = ["8.8.8.8"]
  number_of_dots = 2
  search         = ["f5.com"]
}

resource "bigip_net_vlan" "vlan1" {
  name = "/Common/internal"
  tag  = 101
  interfaces {
    vlanport = 1.2
    tagged   = false
  }
}

resource "bigip_net_vlan" "vlan2" {
  name = "/Common/external"
  tag  = 102
  interfaces {
    vlanport = 1.1
    tagged   = false
  }
}

resource "bigip_net_selfip" "selfip1" {
  name       = "/Common/internalselfIP"
  ip         = "11.1.1.1/24"
  vlan       = "/Common/internal"
  depends_on = [bigip_net_vlan.vlan1]
}

resource "bigip_net_selfip" "selfip2" {
  name       = "/Common/externalselfIP"
  ip         = "100.1.1.1/24"
  vlan       = "/Common/external"
  depends_on = [bigip_net_vlan.vlan2]
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
}

resource "bigip_ltm_pool_attachment" "attach_node" {
  pool       = "/Common/terraform-pool"
  node       = "/Common/11.1.1.101:80"
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

