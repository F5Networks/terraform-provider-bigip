provider "bigip" {
  address = "10.192.74.73"
  username = "admin"
  password = "admin"
}

resource "bigip_snmp" "snmp" {
  sys_contact = " NetOPsAdmin s.shitole@f5.com" 
  sys_location = "SeattleHQ"
  allowedaddresses = ["202.10.10.2"]
}

