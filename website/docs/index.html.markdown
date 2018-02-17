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
To create your own binary follow the below steps
# create working dir
export GOPATH=$HOME/workspace
mkdir -p $GOPATH/src/github.com/f5devcentral && cd $GOPATH

# get source code
go get github.com/f5devcentral/terraform-provider-bigip

# build and move bin into plugin folder
cd src/github.com/f5devcentral/terraform-provider-bigip/
go build
PLUGIN_DIR=$HOME/.terraform.d/plugins/linux_amd64/
mkdir -p $PLUGIN_DIR && cp terraform-provider-bigip $PLUGIN_DIR/

cat > $HOME/.terraformrc <<- EOM
providers {
  bigip = "$HOME/.terraform.d/plugins/linux_amd64/terraform-provider-bigip"
}
EOM

...

# init
$ terraform init

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
- `username` - (Required) Username for authentication
- `password` - (Required) Password for authentication
- `token_auth` - (Optional, Default=false) Enable to use an external authentication source (LDAP, TACACS, etc)
- `login_ref` - (Optional, Default="tmos") Login reference for token authentication (see BIG-IP REST docs for details)
