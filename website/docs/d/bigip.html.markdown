---
layout: "bigip"
page_title: "BIG-IP Provider : bigip"
sidebar_current: "docs-bigip-datasource-bigip-x"
description: |-
    Provides details about provider  bigip
---

# bigip

`bigip` provides details about terraform BIG-IP Provider

This is provider for F5 BIG-IP and is  helpful when configuring the BIG-IP

 ## Example Usage


```hcl
provider "bigip" {
  address = "${var.url}"
  username = "${var.username}"
  password = "${var.password}"
}

```      

## Argument Reference

* `bigip` - Is terraform F5 BIG-IP provider used to configure  BIG-IP.

* `address` - (Required) Address of the device

* `username` - (Required) Username for authentication

* `password` - (Required) Password for authentication

* `token_auth` - (Optional, Default=false) Enable to use an external authentication source (LDAP, TACACS, etc)

login_ref - (Optional, Default="tmos") Login reference for token authentication (see BIG-IP REST docs for details)
