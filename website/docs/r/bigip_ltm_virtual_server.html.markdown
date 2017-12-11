---
layout: "bigip"
page_title: "BIG-IP: bigip_ltm_virtual_server"
sidebar_current: "docs-bigip-resource-virtual_server-x"
description: |-
    Provides details about bigip_ltm_virtual_server resource
---

# bigip\_ltm\_virtual\_server

`bigip_ltm_virtual_server` Configures Virtual Server

For resources should be named with their "full path". The full path is the combination of the partition + name of the resource. For example /Common/my-pool.


## Example Usage


```hcl
resource "bigip_ltm_virtual_server" "http" {
  name = "/Common/terraform_vs_http"
  destination = "10.12.12.12"
  port = 80
  pool = "/Common/the-default-pool"
}

# A Virtual server with SSL enabled
resource "bigip_ltm_virtual_server" "https" {
  name = "/Common/terraform_vs_https"
  destination = "${var.vip_ip}"
  port = 443
  pool = "${var.pool}"
  profiles = ["/Common/tcp","/Common/my-awesome-ssl-cert","/Common/http"]
  source_address_translation = "automap"
}

# A Virtual server with separate client and server profiles
resource "bigip_ltm_virtual_server" "https" {
  name = "/Common/terraform_vs_https"
  destination = "${var.vip_ip}"
  port = 443
  pool = "${var.pool}"
  client_profiles = ["/Common/tcp"]
  server_profiles = ["/Common/tcp-lan-optimized"]
  source_address_translation = "automap"
}


```      

## Argument Reference


* `name`- (Required) Name of the virtual server

* `port` - (Required) Listen port for the virtual server

* `destination` - (Required) Destination IP

* `pool` - (Optional) Default pool name

* `mask` - (Optional) Mask can either be in CIDR notation or decimal, i.e.: 24 or 255.255.255.0. A CIDR mask of 0 is the same as 0.0.0.0

* `source_address_translation` - (Optional) Can be either omitted for none or the values automap or snat

* `ip_protocol`- (Optional) Specify the IP protocol to use with the the virtual server (all, tcp, or udp are valid)

* `profiles` - (Optional) List of profiles associated both client and server contexts on the virtual server. This includes protocol, ssl, http, etc.

* `client_profiles` - (Optional) List of client context profiles associated on the virtual server. Not mutually exclusive with profiles and server_profiles

* `server_profiles` - (Optional) List of server context profiles associated on the virtual server. Not mutually exclusive with profiles and client_profiles
