---
layout: "bigip"
page_title: "BIG-IP: bigip_auth_radius_server"
subcategory: "System"
description: |-
  Provides details about bigip_auth_radius_server resource
---

# bigip\_auth\_radius\_server

`bigip_auth_radius_server` manages a single RADIUS server (`auth/radius-server`) that BIG-IP can use for remote authentication. Reference one or more of these by name from a [`bigip_auth_radius`](bigip_auth_radius.html) resource's `servers` list.

~> **Note** BIG-IP never returns `secret` on read. Terraform will not detect out-of-band changes to it.

## Example Usage

```hcl
resource "bigip_auth_radius_server" "primary" {
  name    = "radius-primary"
  server  = "10.10.10.20"
  secret  = var.radius_secret
  port    = 1812
  timeout = 3
}
```

## Argument Reference

* `name` - (Required, type `string`, ForceNew) Name of the RADIUS server configuration object.

* `server` - (Required, type `string`) IP address or hostname of the RADIUS server.

* `secret` - (Required, type `string`, Sensitive) Shared secret used to authenticate to the RADIUS server.

* `port` - (Optional, type `int`, Default `1812`) Port used to communicate with the RADIUS server.

* `timeout` - (Optional, type `int`, Default `3`) Seconds to wait for a response from the RADIUS server.

## Import

```sh
terraform import bigip_auth_radius_server.primary radius-primary
```
