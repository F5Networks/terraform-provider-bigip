/*
Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

resource "bigip_ltm_monitor" "TC1" {
  name     = "/Common/test_monitor_tc1"
  parent   = "/Common/http"
  send     = "GET /some/path\r\n"
  timeout  = "999"
  interval = "990"
}

resource "bigip_ltm_monitor" "TC2" {
  name        = "/Common/test_monitor_tc2"
  parent      = "/Common/http"
  timeout     = "999"
  interval    = "990"
  destination = "1.2.3.4:1234"
}

resource "bigip_ltm_monitor" "TC3" {
  name     = "/Common/test_monitor_tc3"
  parent   = "/Common/http"
  send     = "GET /\r\n"
  reverse  = "disabled"
  interval = 5
}

//resource "bigip_ltm_monitor" "TC3_CHILD" {
//  name = "/Common/test_monitor_tc3_child"
//  defaults_from = bigip_ltm_monitor.TC3.name
////  parent      = bigip_ltm_monitor.TC3.name
//  send        = "GET /\r\n"
//  receive     = ""
//  destination = "*:8008"
//  reverse     = "disabled"
//}

resource "bigip_ltm_monitor" "TC4" {
  name     = "/Common/test_monitor_tc4"
  parent   = "/Common/ldap"
  username = "testuser"
}

resource "bigip_ltm_monitor" "TC5" {
  name   = "/Common/test_monitor_tc5"
  parent = "/Common/mysql"
}

resource "bigip_ltm_monitor" "TC6" {
  name     = "/Common/HC4_HTTP_OUR_CUSTOM_MONITOR_PARENT"
  parent   = "/Common/http"
  send     = "GET /some/path\r\n"
  timeout  = "999"
  interval = "990"
}
resource "bigip_ltm_monitor" "TC7" {
  name          = "/Common/HC4_HTTP_OUR_CUSTOM_MONITOR_CHILD"
  parent        = "/Common/http"
  custom_parent = bigip_ltm_monitor.TC6.name
  send          = "GET /some/path\r\n"
  timeout       = "999"
  interval      = "980"
}