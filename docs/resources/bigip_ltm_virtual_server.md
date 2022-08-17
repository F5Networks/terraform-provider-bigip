---
layout: "bigip"
page_title: "BIG-IP: bigip_ltm_virtual_server"
subcategory: "Local Traffic Manager(LTM)"
description: |-
  Provides details about bigip_ltm_virtual_server resource
---

# bigip\_ltm\_virtual\_server

`bigip_ltm_virtual_server` Configures Virtual Server

For resources should be named with their `full path`. The full path is the combination of the `partition + name` of the resource (example: `/Common/test-virtualserver` ) or `partition + directory + name` of the resource (example: `/Common/test/test-virtualserver` ).
When including directory in `fullpath` we have to make sure it is created in the given partition before using it.



## Example Usage


```hcl
resource "bigip_ltm_virtual_server" "http" {
  name        = "/Common/terraform_vs_http"
  destination = "10.12.12.12"
  port        = 80
  pool        = "/Common/the-default-pool"
}

# A Virtual server with SSL enabled
resource "bigip_ltm_virtual_server" "https" {
  name                       = "/Common/terraform_vs_https"
  destination                = var.vip_ip
  description                = "VirtualServer-test"
  port                       = 443
  pool                       = var.pool
  profiles                   = ["/Common/tcp", "/Common/my-awesome-ssl-cert", "/Common/http"]
  source_address_translation = "automap"
  translate_address          = "enabled"
  translate_port             = "enabled"
}

# A Virtual server with separate client and server profiles
resource "bigip_ltm_virtual_server" "https" {
  name                       = "/Common/terraform_vs_https"
  destination                = "10.255.255.254"
  description                = "VirtualServer-test"
  port                       = 443
  client_profiles            = ["/Common/clientssl"]
  server_profiles            = ["/Common/serverssl"]
  security_log_profiles      = ["/Common/global-network"]
  source_address_translation = "automap"
}

```      

## Argument Reference


* `name`- (Required) Name of the virtual server

* `port` - (Required) Listen port for the virtual server

* `destination` - (Required) Destination IP

* `description` - (Optional) Description of Virtual server

* `pool` - (Optional) Default pool name

* `mask` - (Optional) Mask can either be in CIDR notation or decimal, i.e.: 24 or 255.255.255.0. A CIDR mask of 0 is the same as 0.0.0.0

* `source_address_translation` - (Optional) Can be either omitted for none or the values automap or snat

* `translate_address` - Enables or disables address translation for the virtual server. Turn address translation off for a virtual server if you want to use the virtual server to load balance connections to any address. This option is useful when the system is load balancing devices that have the same IP address.

* `translate_port` - Enables or disables port translation. Turn port translation off for a virtual server if you want to use the virtual server to load balance connections to any service

* `ip_protocol`- (Optional) Specify the IP protocol to use with the the virtual server (all, tcp, or udp are valid)

* `profiles` - (Optional) List of profiles associated both client and server contexts on the virtual server. This includes protocol, ssl, http, etc.

* `client_profiles` - (Optional) List of client context profiles associated on the virtual server. Not mutually exclusive with profiles and server_profiles

* `server_profiles` - (Optional) List of server context profiles associated on the virtual server. Not mutually exclusive with profiles and client_profiles

* `source` -  (Optional) Specifies an IP address or network from which the virtual server will accept traffic.

* `irules` - (Optional) The iRules list you want run on this virtual server. iRules help automate the intercepting, processing, and routing of application traffic.

* `snatpool` - (Optional) Specifies the name of an existing SNAT pool that you want the virtual server to use to implement selective and intelligent SNATs. DEPRECATED - see Virtual Server Property Groups source-address-translation

* `vlans` - (Optional) The virtual server is enabled/disabled on this set of VLANs,enable/disabled will be desided by attribute `vlan_enabled`

* `vlans_enabled` - (Optional Bool) Enables the virtual server on the VLANs specified by the `vlans` option.
By default it is `false` i.e vlanDisabled on specified vlans, if we want enable virtual server on VLANs specified by `vlans`, mark this attribute to `true`.

* `persistence_profiles` - (Optional) List of persistence profiles associated with the Virtual Server.

* `fallback_persistence_profile` - (Optional) Specifies a fallback persistence profile for the Virtual Server to use when the default persistence profile is not available.

* `security_log_profiles` - (Optional) Specifies the log profile applied to the virtual server.

## Importing
An existing virtual-server can be imported into this resource by supplying virtual-server Name in `full path` as `id`.
An example is below:
```sh
$ terraform import bigip_ltm_virtual_server.http /Common/terraform_vs_http
```
