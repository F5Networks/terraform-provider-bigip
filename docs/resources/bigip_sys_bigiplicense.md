---
layout: "bigip"
page_title: "BIG-IP: bigip_sys_bigiplicense"
subcategory: "System"
description: |-
  Provides details about bigip_sys_bigiplicense resource
---

# bigip\_sys\_bigiplicense

`bigip_sys_bigiplicense` This resource is used to license a BIG-IP device using a registration key.

## Example Usage

```hcl
resource "bigip_sys_bigiplicense" "example" {
  command           = "install"
  registration_key  = "XXXXX-XXXXX-XXXXX-XXXXX-XXXXXXX"
}
```

## Argument Reference

* `command` - (Required) Tmsh command to execute, e.g. `install`.
* `registration_key` - (Required) A unique key F5 provides for licensing BIG-IP.

## Note

Applying a license causes the BIG-IP configuration daemon to restart, which can take several minutes. Terraform waits for this process to complete before continuing.
