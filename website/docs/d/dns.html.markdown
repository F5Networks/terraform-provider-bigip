---
layout: "bigip"
page_title: "BIG-IP: bigip_sys_dns"
sidebar_current: "docs-bigip-datasource-dns-x"
description: |-
    Provides details about bigip dns
---

# bigip\_dns

`bigip_sys_dns` provides details about a specific bigip

This resource is helpful when configuring DNS server on the BIG-IP. 
## Example Usage


```hcl
provider "bigip" {
  address = "10.192.74.73"
  username = "admin"
  password = "admin"
}


resource "bigip_sys_dns" "dns1" {
   description = "/Common/DNS1"
   name_servers = ["1.1.1.1"]
   numberof_dots = 2
   search = ["f5.com"]
}

```      

## Argument Reference

* `bigip_sys_dns` - Is the resource is used to configure dns name server on the BIG-IP.

* `/Common/DNS` - Is the description of the DNS server in the main or common partition of BIG-IP.

* `1.1.1.1` - Is the name_Server, you can configure a set of name servers. You can also provide DNS names. 

