---
layout: "bigip"
page_title: "BIG-IP Provider : Index"
sidebar_current: "docs-bigip-index"
description: |-
    Provides details about provider bigip
---

# Overview

A [Terraform](terraform.io) provider for F5 BigIP. Resources are currently available for LTM.

# F5 Requirements

This provider uses the iControlREST API. Make sure that is installed and enabled on your F5 before proceeding.

# Installation

 - Download the latest [release](https://github.com/DealerDotCom/terraform-provider-bigip/releases) for your platform.
 - Rename the executable to `terraform-provider-bigip`
 - Copy somewhere on your path, or update `.terraformrc` in your home directory like so:

```
providers {
	bigip = "/path/to/terraform-provider-bigip"
}
```

# Provider Configuration

### Example

```
provider "bigip" {
  address = "${var.url}"
  username = "${var.username}"
  password = "${var.password}"
}
```

### Reference

- `address` - (Required) Address of the device
- 
- `username` - (Required) Username for authentication
- 
- `password` - (Required) Password for authentication
- 
- `token_auth` - (Optional, Default=false) Enable to use an external authentication source (LDAP, TACACS, etc)
- 
- `login_ref` - (Optional, Default="tmos") Login reference for token authentication (see BIG-IP REST docs for details)

