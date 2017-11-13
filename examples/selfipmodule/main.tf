variable "name" {}
variable "ip" {}
variable "vlan" {}

resource "bigip_net_selfip" "selfip" {
        name = "${var.name}"
        ip = "${var.ip}"
        vlan = "${var.vlan}"
        }
