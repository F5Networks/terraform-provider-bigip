provider "bigip" {
  address = "10.192.74.73"
  username = "admin"
  password = "admin"
}

resource "bigip_snat" "snat3" {
  // this is using snatpool translation is not required
  name = "snat3"
  origins = ["6.1.6.6"]
  mirror = "false"
  snatpool = "/Common/sanjaysnatpool"
}

resource "bigip_snat" "snat_list" {
 name = "NewSnatList"
 translation = "136.1.1.1"
 origins = ["2.2.2.2", "3.3.3.3"]
}

