/*
Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

resource "bigip_ltm_monitor" "vs_tc1" {
  name     = "/Common/test_monitor_vs_tc1"
  parent   = "/Common/http"
  send     = "GET /some/path\r\n"
  timeout  = "999"
  interval = "990"
}

resource "bigip_ltm_pool" "vs_tc1" {
  name                = "/Common/test_pool_vs_tc1"
  load_balancing_mode = "round-robin"
  monitors            = [bigip_ltm_monitor.vs_tc1.name]
  allow_snat          = "yes"
  allow_nat           = "yes"
}

resource "bigip_ltm_pool_attachment" "vs_tc1" {
  pool = bigip_ltm_pool.vs_tc1.name
  node = "11.1.1.101:80"
}

resource "bigip_ltm_pool_attachment" "vs_tc1a" {
  pool = bigip_ltm_pool.vs_tc1.name
  node = "11.1.1.102:80"
}

resource "bigip_ltm_virtual_server" "vs_tc1" {
  pool                       = bigip_ltm_pool.vs_tc1.name
  name                       = "/Common/test_vs_tc1"
  destination                = "100.1.1.100"
  port                       = 80
  source_address_translation = "automap"
}

resource "bigip_ltm_virtual_server" "vs_tc2" {
  name                       = "/Common/test_vs_tc2"
  destination                = "10.0.0.2"
  mask                       = "255.255.255.254"
  ip_protocol                = "any"
  description                = "Virtual server"
  port                       = 80
  source_address_translation = "automap"
  profiles = [
    "/Common/fastL4",
  ]
  source_port = "preserve"
}

resource "bigip_ltm_virtual_server" "vs_tc3" {
  name                       = "/Common/test_vs_tc3"
  destination                = "10.0.1.2"
  mask                       = "31"
  ip_protocol                = "any"
  description                = "Virtual server"
  port                       = 80
  source_address_translation = "automap"
  profiles = [
    "/Common/fastL4",
  ]
  source_port = "preserve"
}

resource "bigip_ltm_profile_http" "vs_tc4" {
  name              = "/Common/test_profile_vs_tc4"
  defaults_from     = "/Common/http"
  response_chunking = "rechunk"
}

# A Virtual server with separate client and server profiles
resource "bigip_ltm_virtual_server" "vs_tc4" {
  name                       = "/Common/test_vs_tc4"
  destination                = "10.255.255.254"
  description                = "VirtualServer-test"
  port                       = 443
  profiles                   = ["/Common/tcp", bigip_ltm_profile_http.vs_tc4.name]
  client_profiles            = ["/Common/clientssl"]
  server_profiles            = ["/Common/serverssl"]
  security_log_profiles      = ["/Common/global-network"]
  source_address_translation = "automap"
}

resource "bigip_ltm_virtual_server" "vs_tc5" {
  name        = "/Common/test_vs_tc5"
  destination = "10.12.12.12"
  port        = 80
  //  snatpool                   = "snat-pool-1234"
  //  snatpool = "/Common/testsnatpool"
  source_address_translation = "none"
  //  source_address_translation = "snat"
  //  pool = "terraform-pool-8443"
}

resource "bigip_ltm_virtual_server" "vs_tc6" {
  name                       = "/Common/test_vs_tc6"
  destination                = "10.13.12.12"
  port                       = 80
  source_address_translation = "none"
  vlans_enabled              = true
  vlans                      = ["/Common/external"]
}