---
layout: "bigip"
page_title: "BIG-IP: bigip_auth_user"
subcategory: "System"
description: |-
  Provides details about bigip_auth_user resource
---

# bigip\_auth\_user

`bigip_auth_user` manages a BIG-IP local user account and its partition/role assignments (`auth/user`).

~> **Note** BIG-IP never returns `password` on read (only an opaque `encryptedPassword`, which this resource does not track). Terraform will not detect out-of-band password changes; set it once and manage rotation separately, e.g. `lifecycle { ignore_changes = [password] }`.

## Example Usage

```hcl
resource "bigip_auth_user" "operator" {
  name        = "operator1"
  description = "Read-only operator account"
  password    = var.operator_password
  shell       = "tmsh"

  partition_access {
    partition = "Common"
    role      = "guest"
  }
}
```

## Argument Reference

* `name` - (Required, type `string`, ForceNew) Login name of the user account.

* `description` - (Optional, type `string`, Computed) User-defined description, commonly used for the account's display/full name. If not set, BIG-IP defaults this to the account's name.

* `password` - (Optional, type `string`, Sensitive) Password for the account.

* `shell` - (Optional, type `string`, Default `none`) Login shell assigned to the user, e.g. `bash`, `tmsh`, or `none`. Note BIG-IP restricts `bash` to accounts with the `admin` role.

* `partition_access` - (Required, type `list`, min 1 item) One or more partition/role assignments. Each block supports:
  * `partition` - (Required, type `string`) Partition name the role applies to, or `all-partitions` for every partition.
  * `role` - (Required, type `string`) Role granted, e.g. `admin`, `resource-admin`, `manager`, `guest`, `operator`, `application-editor`, `auditor`, `no-access`.

## Import

```sh
terraform import bigip_auth_user.operator operator1
```
