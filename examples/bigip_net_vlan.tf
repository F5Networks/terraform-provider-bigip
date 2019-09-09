provider "bigip" {
  address = "10.192.74.73"
  username = "admin"
  password = "admin"
}




resource "bigip_net_vlan" "vlan1" {
	name = "/Common/internal"
	tag = 101
	interfaces = {
		vlanport = 1.2,
		tagged = false
	}	
}

