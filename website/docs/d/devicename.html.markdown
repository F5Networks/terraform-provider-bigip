---
layout: "bigip"
page_title: "BIG-IP: bigip_device_name"
sidebar_current: "docs-bigip-datasource-device_name-x"
description: |-
    Provides details about bigip device_name 
---

# bigip\_device_name

`bigip_device_name` provides details about a specific bigip

This resource is helpful when configuring the BIG-IP device_name in cluster or in HA mode. 
## Example Usage


```hcl
provider "bigip" {
  address = "10.192.74.73"
  username = "admin"
  password = "admin"
}


resource "bigip_device_name" "my_new_device_name"

{
            name = "bigip1i00.f5.com"
            command = "mv"
            target = "bigipNew.f5.com"
        }
```      

## Argument Reference

* `bigip_device_name` - Is the resource is used to configure new device name on the BIG-IP this is often reqquired while setting up BIG-IP in HA pairs.

* `bigip1i00.f5.com` - Is the existing name of the BIG-IP device

* `mv` - Is the tmsh move instruction which is attribute of command 

* `target` - Is the new name of the BIG-IP
