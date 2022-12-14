---
layout: "bigip"
page_title: "BIG-IP: bigip_vcmp_guest"
subcategory: "Network"
description: |-
  Provides details about bigip_vcmp_guest resource
---

# bigip\_vcmp\_guest

`bigip_vcmp_guest` Manages a vCMP guest configuration

Resource does not wait for vCMP guest to reach the desired state, it only ensures that a desired configuration is set on the target device.


## Example Usage


```hcl
resource "bigip_vcmp_guest" "vcmp-test" {
  name                = "tf_guest"
  initial_image       = "12.1.2.iso"
  mgmt_network        = "bridged"
  mgmt_address        = "10.1.1.1/24"
  mgmt_route          = "none"
  state               = "provisioned"
  cores_per_slot      = 2
  number_of_slots     = 1
  min_number_of_slots = 1
}


```      

## Argument Reference

* `name` - (Required) Name of the vCMP guest

* `initial_image` - (Optional, `string`) Specifies the base software release ISO image file for installing the TMOS hypervisor instance.

* `initial_hotfix` - (Optional, `string`) Specifies the hotfix ISO image file which is applied on top of the base image.

* `vlans` - (Optional, `list`) Specifies the list of VLANs the vCMP guest uses to communicate with other guests, the host, and with the external network. The naming format must be the combination of the partition + name. For example /Common/my-vlan

* `mgmt_network` - (Optional, `string`) Specifies the method by which the management address is used in the vCMP guest. options : [`bridged`,`isolated`,`host-only`].

* `mgmt_address` - (Optional, `string`) Specifies the IP address and subnet or subnet mask you use to access the guest when you want to manage a module running within the guest.

* `mgmt_route` - (Optional, `string`) Specifies the gateway address for the `mgmt_address`. Can be set to `none` to remove the value from the configuration.

* `state` - (Optional, `string`) Specifies the state of the vCMP guest on the system. options : [`configured`,`provisioned`,`deployed`].

* `cores_per_slot` - (Optional, `int`) Specifies the number of cores the system allocates to the guest.

* `number_of_slots` - (Optional, `int`) Specifies the number of slots for the system to use when creating the guest.

* `min_number_of_slots` - (Optional, `int`) Specifies the minimum number of slots the guest must be assigned to in order to deploy.

* `delete_virtual_disk` - (Optional, `bool`) Indicates if virtual disk associated with vCMP guest should be removed during remove operation.  The default is `true`

