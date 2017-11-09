provider "bigip" {
address = "10.192.74.73"
username = "admin"
password = "admin"

}

resource "bigip_net_vlan" "vlan1" {
	name = "/Common/internal"
	tag = 101
	interfaces = {
		name = 1.1,
		tagged = false
	}	

}

resource "bigip_net_vlan" "vlan2" {
        name = "/Common/external"
        tag = 102
        interfaces = {
                name = 1.2,
                tagged = false
        }

}

resource "bigip_net_selfip" "selfip1" {
	name = "/Common/internalselfIP"
	ip = " 11.1.1.1/24"
	vlan = "/Common/internal"
	depends_on = ["bigip_net_vlan.vlan1"]
	}

resource "bigip_net_selfip" "selfip2" {
        name = "/Common/externalselfIP"
        ip = " 100.1.1.1/24"
        vlan = "/Common/external"
        depends_on = ["bigip_net_vlan.vlan2"]
        }



