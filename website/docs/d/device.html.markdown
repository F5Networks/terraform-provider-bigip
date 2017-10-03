---
layout: "bigip"
page_title: "BIG-IP: bigip_device"
sidebar_current: "docs-bigip-datasource-device-x"
description: |-
    Provides details about bigip device 
---

# bigip\_device

`bigip_device` provides details about a specific bigip

This resource is helpful when configuring the BIG-IP device in cluster or in HA mode. 
## Example Usage


```hcl
provider "bigip" {
  address = "10.192.74.73"
  username = "admin"
  password = "admin"
}


resource "bigip_device" "my_new_device"

        {
            name = "bigip300.f5.com"
            configsync_ip = "2.2.2.2"
            mirror_ip = "10.10.10.10"
            mirror_secondary_ip = "11.11.11.11"
        }
        ```
