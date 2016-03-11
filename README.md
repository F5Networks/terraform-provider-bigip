# Overview

A [Terraform](terraform.io) provider for F5 BigIP. Resources are currently available for LTM. [![Build Status](https://travis-ci.org/DealerDotCom/terraform-provider-bigip.svg?branch=master)](https://travis-ci.org/DealerDotCom/terraform-provider-bigip)

# Installation

 - Download the latest [release](https://github.com/DealerDotCom/terraform-provider-bigip/releases) for your platform.
 - Rename the executable to `terraform-provider-bigip`
 - Copy somewhere on your path or update `.terraformrc` in your home directory like so:
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

`address` - (Required) Address of the device

`username` - (Required) Username for authentication

`password` - (Required) Password for authentication

`loginReference` - (Optional) Login reference for token auth, needed
if you are using non-local accounts (LDAP, Radius, etc).


# Resources

## bigip_ltm_monitor

Configures a custom monitor for use by health checks.

### Example
```
resource "bigip_ltm_monitor" "monitor" {
  name = "terraform_monitor"
  parent = "http"
  send = "GET /some/path\r\n"
  timeout = "999"
  interval = "999"
}
```

### Reference

`name` - (Required) Name of the monitor

`parent` - (Required) Existing LTM monitor to inherit from

`partition` - (Required) LTM partition to create the resource in. Default = Common.

`interval` - (Optional) Check interval in seconds

`timeout` - (Optional) Timeout in seconds

`send` - (Optional) Request string to send

`receive` - (Optional) Expected response string

`receive_disable` - (Optional)

`reverse` - (Optional)

`transparent` - (Optional)

`manual_resume` - (Optional)

`ip_dscp` - (Optional)

`time_until_up` - (Optional)

## bigip_ltm_node

Manages a node configuration

### Example

```
resource "bigip_ltm_node" "node" {
  name = "terraform_node1"
  address = "10.10.10.10"
}
```

### Reference

`name` - (Required) Name of the node

`address` - (Required) IP or hostname of the node

## bigip_ltm_pool

### Example

```
resource "bigip_ltm_pool" "pool" {
  name = "terraform-pool"
  load_balancing_mode = "round-robin"
  nodes = ["${bigip_ltm_node.node.name}:80"]
  monitors = ["${bigip_ltm_monitor.monitor.name}","${bigip_ltm_monitor.monitor2.name}"]
  allow_snat = false
}
```

### Reference

`name` - (Required) Name of the pool

`partition` - (Required) LTM partition to create the resource in. Default = Common.

`nodes` - (Optional) Nodes to add to the pool. Format node_name:port. e.g. `node01:443`

`monitors` - (Optional) List of monitor names to associate with the pool

`allow_nat` - (Optional)

`allow_snat` - (Optional)

`load_balancing_mode` - (Optional, Default = round-robin)

## bigip_ltm_virtual_server

Configures a Virtual Server

### Example

```
resource "bigip_ltm_virtual_server" "vs" {
  name = "terraform_vs_http"
  destination = "10.12.12.12"
  port = 80
  pool = "${bigip_ltm_pool.pool.name}"
}
```

### Reference

`name` - (Required) Name of the virtual server

`partition` - (Required, Default=Common) LTM partition to create the resource in.

`port` - (Required) Listen port for the virtual server

`destination` - (Required) Destination IP

`pool` - (Optional) Default pool name

`mask` - (Optional) Mask can either be in CIDR notation or decimal, i.e.: `24` or `255.255.255.0`. A CIDR mask of `0` is the same as `0.0.0.0`

`source_address_translation` - (Optional) Can be either omitted for `none` or the values `automap` or `snat` 

## bigip_ltm_irule

Creates iRules

### Example

```
# Loading from a file is the preferred method
resource "bigip_ltm_irule" "rule" {
  name = "terraform_irule"
  irule = "${file("myirule.tcl")}"
}

resource "bigip_ltm_irule" "rule2" {
  name = "terraform_irule2"
  irule = <<EOF
when CLIENT_ACCEPTED {
     log local0. "test"
   }
EOF
}
```

### Reference

`name` - (Required) Name of the iRule

`irule` - (Required) Body of the iRule


## bigip_ltm_virtual_address

Configures a Virtual Address

### Example 

```
resource "bigip_ltm_virtual_address" "" {
  
}
```

