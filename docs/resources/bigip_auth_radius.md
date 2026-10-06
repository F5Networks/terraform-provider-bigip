---
layout: "bigip"
page_title: "BIG-IP: bigip_auth_radius"
subcategory: "System"
description: |-
  Provides details about bigip_auth_radius resource
---

# bigip\_auth\_radius

`bigip_auth_radius` manages the BIG-IP system's system-wide RADIUS remote authentication configuration (`auth/radius system-auth`), which references one or more [`bigip_auth_radius_server`](bigip_auth_radius_server.html) resources by name.

This is a singleton resource -- BIG-IP supports exactly one system RADIUS configuration. This resource only configures the RADIUS backend itself; it does not activate RADIUS as the active authentication source (`auth/source`).

## Example Usage

```hcl
resource "bigip_auth_radius_server" "primary" {
  name   = "radius-primary"
  server = "10.10.10.20"
  secret = var.radius_secret
}

resource "bigip_auth_radius_server" "secondary" {
  name   = "radius-secondary"
  server = "10.10.10.21"
  secret = var.radius_secret
}

resource "bigip_auth_radius" "example" {
  servers      = [bigip_auth_radius_server.primary.name, bigip_auth_radius_server.secondary.name]
  service_type = "authenticate-only"
  retries      = 3
}
```

## Argument Reference

* `servers` - (Required, type `list`) Names of the `bigip_auth_radius_server` objects to use, in priority order.

* `service_type` - (Optional, type `string`, Default `authenticate-only`) RADIUS service type requested, e.g. `authenticate-only`, `login`, `framed`, `callback-login`, `callback-framed`, `outbound`, `administrative`, `nas-prompt`, `call-check`, `callback-nas-prompt`.

* `retries` - (Optional, type `int`, Default `3`) Number of times to retry a request before failing over to the next server.

## Import

```sh
terraform import bigip_auth_radius.example /Common/system-auth
```
