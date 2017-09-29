provider "bigip" {
 address = "10.192.74.73"
 username = "admin"
 password = "admin"
}


module  "sjvlan1" {
  source = "./vlanmodule"
  name = "/Common/intvlan"
  tag = 101
  vlanport = "1.1"
  tagged = true
 }

module "sjvlan2"  {
  source = "./vlanmodule"
  name = "/Common/extvlan"
  tag = 102
  vlanport = "1.2"
  tagged = true
 }

