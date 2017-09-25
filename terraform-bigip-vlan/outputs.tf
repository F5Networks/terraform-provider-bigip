output "name" {
 description = " This is the Name of the Vlan"
 value = "${var.name}"
}

output "tag" {
 description = " This is the Vlan tag used "
 value = "${var.tag}"
}

output "tagged" {
 description = " This is boolean to enable tagging "
 value = "${var.tagged}"
}

output "vlanport" {
 description = " This is the port number used by Vlan "
 value = "${var.vlanport}"
}


