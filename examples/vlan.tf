provider "bigip" {
  address = "10.192.74.73"
  username = "admin"
  password = "admin"
}
resource "bigip_ltm_vlan" "vlan1" {
	name = "/Common/Internal"
	tag = 101
	interfaces = {
                vlanport = 1.2,
		tagged = false
	}
}
