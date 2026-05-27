---
layout: "bigip"
page_title: "BIG-IP: bigip_gtm_server"
subcategory: "Global Traffic Manager(GTM)"
description: |-
  Provides details about bigip_gtm_server data source
---

# bigip\_gtm\_server

Use this data source (`bigip_gtm_server`) to look up an existing GTM server on the BIG-IP. This is useful when servers are shared across multiple Terraform workspaces and you need to reference them without managing them as resources.

## Example Usage

```hcl
data "bigip_gtm_server" "srv" {
  name = "my-server"
}

output "server_datacenter" {
  value = data.bigip_gtm_server.srv.datacenter
}

output "server_product" {
  value = data.bigip_gtm_server.srv.product
}

output "server_addresses" {
  value = data.bigip_gtm_server.srv.addresses
}
```

## Argument Reference

* `name` - (Required) Name of the GTM server.

## Attributes Reference

Additionally, the following attributes are exported:

* `datacenter` - Datacenter the server belongs to.

* `description` - Description of the GTM server.

* `product` - Server type (bigip, generic-host, etc.).

* `enabled` - Whether the server is enabled.

* `monitor` - Monitor assigned to the server.

* `virtual_server_discovery` - Virtual server discovery mode.

* `link_discovery` - Link discovery mode.

* `prober_preference` - Prober preference.

* `prober_fallback` - Prober fallback.

* `prober_pool` - Prober pool.

* `addresses` - IP addresses for the server. Each address contains:
  * `name` - IP address.
  * `device_name` - Device name for the address.
  * `translation` - IP translation address.

* `virtual_servers` - Virtual servers configured on the GTM server. Each virtual server contains:
  * `name` - Name of the virtual server.
  * `destination` - Destination IP address and port.
  * `enabled` - Whether the virtual server is enabled.
  * `translation_address` - Translation IP address for NAT.
  * `translation_port` - Translation port for NAT.
  * `monitor` - Health monitor for this virtual server.