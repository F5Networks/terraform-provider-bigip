---
layout: "bigip"
page_title: "BIG-IP: bigip_sys_ntp"
subcategory: "System"
description: |-
  Provides details about bigip_sys_ntp resource
---

# bigip\_sys\_ntp

`bigip_sys_ntp` resource is helpful when configuring NTP server on the BIG-IP.

## Example Usage

```hcl
resource "bigip_sys_ntp" "ntp1" {
  description = "/Common/NTP1"
  servers     = ["time.facebook.com"]
  timezone    = "America/Los_Angeles"
}
```      

## Argument Reference

* `description` - (Required,type `string`) User defined description.

* `servers` - (Required,type `list`) Specifies the time servers that the system uses to update the system time.

* `timezone` - (Optional,type `string`) Specifies the time zone that you want to use for the system time.
