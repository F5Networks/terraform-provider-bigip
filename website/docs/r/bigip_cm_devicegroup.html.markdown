---
layout: "bigip"
page_title: "BIG-IP: bigip_cm_devicegroup"
sidebar_current: "docs-bigip-resource-devicegroup-x"
description: |-
    Provides details about bigip device_name
---

# bigip_cm_devicegroup

`bigip_cm_devicegroup` A device group is a collection of BIG-IP devices that are configured to securely synchronize their BIG-IP configuration data, and fail over when needed.


## Example Usage


```hcl
resource "bigip_cm_devicegroup" "my_new_devicegroup" {
  name              = "sanjose_devicegroup"
  auto_sync         = "enabled"
  full_load_on_sync = "true"
  type              = "sync-only"
  device {
    name = "bigip1.cisco.com"
  }
  device {
    name = "bigip200.f5.com"
  }
}
```      

## Argument Reference

* `bigip_cm_devicegroup` - Is the resource  used to configure new device group on the BIG-IP.

* `name` - Is the name of the device Group

* `auto_sync` - Specifies if the device-group will automatically sync configuration data to its members

* `type` - Specifies if the device-group will be used for failover or resource syncing

* `device` - Name of the device to be included in device group, this need to be configured before using devicegroup resource
