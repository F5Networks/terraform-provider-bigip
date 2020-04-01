---
layout: "bigip"
page_title: "BIG-IP: bigip_command"
sidebar_current: "docs-bigip-resource-device-x"
description: |-
    Provides details about bigip device
---

# bigip_command

`bigip_command` Run TMSH commands on F5 devices

This resource is helpful to send TMSH command to an BIG-IP node and returns the results read from the device
## Example Usage


```hcl
resource "bigip_command" "test-command" {
  commands   = ["show sys version"]
}
```   
## Argument Reference

* `commands` - (Required) The commands to send to the remote BIG-IP device over the configured provider. The resulting output from the command is returned and added to `command_result` 
* `name` - (Optional) Optional Name to to Uniquely identify state.    


## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `command_result` - The resulting output from the `commands` executed
