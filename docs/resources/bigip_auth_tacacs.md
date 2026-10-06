---
layout: "bigip"
page_title: "BIG-IP: bigip_auth_tacacs"
subcategory: "System"
description: |-
  Provides details about bigip_auth_tacacs resource
---

# bigip\_auth\_tacacs

`bigip_auth_tacacs` manages the BIG-IP system's TACACS+ remote authentication configuration (`auth/tacacs system-auth`).

This is a singleton resource -- BIG-IP supports exactly one system TACACS+ configuration. This resource only configures the TACACS+ backend itself; it does not activate TACACS+ as the active authentication source (`auth/source`).

~> **Note** BIG-IP never returns `secret` on read. Terraform will not detect out-of-band changes to it.

## Example Usage

```hcl
resource "bigip_auth_tacacs" "example" {
  servers    = ["10.10.10.30", "10.10.10.31"]
  secret     = var.tacacs_secret
  service    = "ppp"
  protocol   = "ip"
  encryption = "enabled"
  accounting = "send-to-first-server"
}
```

## Argument Reference

* `servers` - (Required, type `list`) TACACS+ server IP addresses or hostnames, in priority order.

* `secret` - (Optional, type `string`, Sensitive) Shared secret used to authenticate to the TACACS+ servers.

* `service` - (Optional, type `string`, Default `ppp`) TACACS+ service requested, e.g. `ppp`, `shell`, `slip`, `system`.

* `protocol` - (Optional, type `string`, Default `ip`) TACACS+ protocol name associated with `service`, e.g. `ip`, `lcp`, `vpdn`.

* `encryption` - (Optional, type `string`, Default `enabled`) Whether the system encrypts TACACS+ traffic.

* `accounting` - (Optional, type `string`, Default `send-to-first-server`) Whether to send accounting messages to the first available server (`send-to-first-server`) or all available servers (`send-to-all-servers`).

* `authentication_via` - (Optional, type `string`) Whether to authenticate via PAP or CHAP.

## Import

```sh
terraform import bigip_auth_tacacs.example /Common/system-auth
```
