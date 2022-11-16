---
layout: "bigip"
page_title: "BIG-IP: bigip_sys_dns"
subcategory: "System"
description: |-
  Provides details about bigip_sys_dns resource
---

# bigip\_sys\_dns

`bigip_sys_dns` Configures DNS Name server on F5 BIG-IP

## Example Usage

```hcl
resource "bigip_sys_dns" "dns1" {
  description  = "/Common/DNS1"
  name_servers = ["1.1.1.1"]
  search       = ["f5.com"]
}
```      

## Argument Reference

* `description`- (Required,type `string` )Provide description for your DNS server

* `name_servers` - (Required,type `list` ) Specifies the name servers that the system uses to validate DNS lookups, and resolve host names.

* `number_of_dots` - (Optional,type `int` ) Configures the number of dots needed in a name before an initial absolute query will be made.

* `search` - (Optional,type `list` ) Specifies the domains that the system searches for local domain lookups, to resolve local host names.
