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

resource "bigip_sys_snmp_traps" "snmp_traps" {
  name        = "snmptraps"
  community   = "f5community"
  host        = "195.10.10.1"
  description = "Setup snmp traps"
  port        = 111
}

