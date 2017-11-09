provider "bigip" {
  address = "10.192.74.73"
  username = "admin"
  password = "admin"
}

resource "bigip_sys_snmp_traps" "snmp_traps" {
name = "snmptraps"
community = "f5community"
host = "195.10.10.1"
description = "Setup snmp traps"
port = 111
}

