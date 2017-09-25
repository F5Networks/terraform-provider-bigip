provider "bigip" {
address = "10.192.74.73"
username = "admin"
password = "admin"

}
module  "vlan" {
  source = "./terraform-bigip-vlan"
  name = "/Common/intvlan"
  tag = 101
  vlanport = "1.1"
  tagged = true
 }
 output "vlanport" {
 value = "${module.vlan.vlanport}"
 }

