---
layout: "bigip"
page_title: "BIG-IP: bigip_auth_ldap"
subcategory: "System"
description: |-
  Provides details about bigip_auth_ldap resource
---

# bigip\_auth\_ldap

`bigip_auth_ldap` manages the BIG-IP system's LDAP remote authentication configuration (`auth/ldap system-auth`).

This is a singleton resource -- BIG-IP supports exactly one system LDAP configuration. This resource only configures the LDAP backend itself; it does not activate LDAP as the active authentication source (`auth/source`).

~> **Note** BIG-IP never returns `bind_pw` on read. Terraform will not detect out-of-band changes to it; manage rotation separately (e.g. `lifecycle { ignore_changes = [bind_pw] }` if it's rotated outside Terraform).

## Example Usage

```hcl
resource "bigip_auth_ldap" "example" {
  servers        = ["ldap.example.com"]
  port           = 636
  bind_dn        = "cn=admin,dc=example,dc=com"
  bind_pw        = var.ldap_bind_password
  search_base_dn = "dc=example,dc=com"
  login_attribute = "uid"
  ssl            = "enabled"
  ssl_check_peer = "enabled"
}
```

## Argument Reference

* `servers` - (Required, type `list`) LDAP server IP addresses or hostnames.

* `port` - (Optional, type `int`, Default `389`) Port used to communicate with the LDAP servers.

* `bind_dn` - (Optional, type `string`) Distinguished name used to bind to the LDAP server.

* `bind_pw` - (Optional, type `string`, Sensitive) Password for `bind_dn`.

* `bind_timeout` - (Optional, type `int`, Default `30`) Seconds to wait for a bind response.

* `search_base_dn` - (Optional, type `string`) Base distinguished name for user searches.

* `login_attribute` - (Optional, type `string`) LDAP attribute containing the user login name.

* `user_template` - (Optional, type `string`) Template used to construct a user's distinguished name, e.g. `"uid=%s,ou=people,dc=example,dc=com"`.

* `filter` - (Optional, type `string`) LDAP search filter used when searching for a user.

* `group_dn` - (Optional, type `string`) Distinguished name of the group used for group-based role assignment.

* `group_member_attribute` - (Optional, type `string`) LDAP attribute used to determine group membership.

* `ssl` - (Optional, type `string`, Default `disabled`) Whether to use SSL/TLS: `enabled`, `disabled`, or `start-tls`.

* `ssl_ca_cert_file` - (Optional, type `string`) Full path on the BIG-IP system of the CA certificate used to verify the LDAP server's certificate.

* `ssl_check_peer` - (Optional, type `string`, Default `disabled`) Whether to verify the LDAP server's certificate.

* `version` - (Optional, type `int`, Default `3`) LDAP protocol version.

## Import

```sh
terraform import bigip_auth_ldap.example /Common/system-auth
```
