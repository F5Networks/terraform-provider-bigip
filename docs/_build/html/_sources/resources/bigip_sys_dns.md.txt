---
layout: "bigip"
page_title: "BIG-IP: bigip_sys_dns"
sidebar_current: "docs-bigip-resource-dns-x"
description: |-
    Provides details about bigip_sys_dns resource
---

# bigip\_sys\_dns

`bigip_sys_dns` Configures DNS server on F5 BIG-IP




## Example Usage


```hcl
resource "bigip_sys_dns" "dns1" {
  description    = "/Common/DNS1"
  name_servers   = ["1.1.1.1"]
  number_of_dots = 2
  search         = ["f5.com"]
}
```      

## Argument Reference


* `description`- Provide description for your DNS server

* `name_servers` - Name or IP address of the DNS server

* `number_of_dots` - Configures the number of dots needed in a name before an initial absolute query will be made.

* `search` - Specify what domains you want to search
