provider "bigip" {
  address = "10.192.74.68"
  username = "admin"
  password = "admin"
}

resource "bigip_sys_iapp" "waf_asm" {
  name = "policywaf"
  jsonfile = "${file("policywaf.json")}"
}

resource "bigip_sys_iapp" "pool_deployed" {
  name = "sap-dmzpool-rp1-80"
  jsonfile = "${file("sap-dmzpool-rp1-80.json")}"
}

