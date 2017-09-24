variable "name" {}
variable  tag  {}
variable "vlanport"    {}
variable "tagged"    {}

resource "bigip_ltm_vlan" "vlan" {
	name = "${var.name}"
	tag = "${var.tag}"
	interfaces = {
		vlanport = "${var.vlanport}",
		tagged = "${var.tagged}"
	}	

}

