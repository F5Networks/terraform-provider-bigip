provider "bigip" {
  address = "10.192.74.73"
  username = "admin"
  password = "admin"
}

resource "bigip_net_selfip" "selfip1" {
	name = "/Common/internalselfIP"
	ip = "11.1.1.1/24"
	vlan = "/Common/internal"
	depends_on = ["bigip_net_vlan.vlan1"]
	}

