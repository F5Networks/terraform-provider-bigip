variable "name" {}
variable "ip" {}
variable "vlan" {}

resource "bigip_ltm_selfip" "selfip" {
        name = "${var.name}"
        ip = "${var.ip}"
        vlan = "${var.vlan}"
        }
