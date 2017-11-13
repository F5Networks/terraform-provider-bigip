
resource "bigip_net_vlan" "vlan" {
	name = "${var.name}"
	tag = "${var.tag}"
	interfaces = {
		vlanport = "${var.vlanport}",
		tagged = "${var.tagged}"
	}	

}


