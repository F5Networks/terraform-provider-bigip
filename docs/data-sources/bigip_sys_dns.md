---
layout: "bigip"
page_title: "BIG-IP: bigip_sys_dns"
subcategory: "System"
description: |-
  Provides details about bigip_sys_dns data source
---

# bigip\_sys\_dns

Use this data source (`bigip_sys_dns`) to retrieve the BIG-IP system's DNS resolver configuration (name servers and search domains). `sys/dns` is a singleton object -- there is exactly one DNS configuration per BIG-IP system -- so this data source takes no arguments.

Use [`bigip_sys_dns`](../resources/bigip_sys_dns.md) (the resource) to manage this configuration.

## Example Usage

```hcl
data "bigip_sys_dns" "dns" {}

output "name_servers" {
  value = data.bigip_sys_dns.dns.name_servers
}

output "search_domains" {
  value = data.bigip_sys_dns.dns.search
}
```

## Argument Reference

This data source takes no arguments.

## Attributes Reference

* `description` - User-defined description of the system DNS configuration.

* `name_servers` - Name servers the system uses to validate DNS lookups and resolve host names.

* `search` - Domains the system searches for local domain lookups, to resolve local host names.

* `number_of_dots` - Number of dots that must appear in a host name before an initial absolute lookup is performed.
