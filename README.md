# Overview

A [Terraform](terraform.io) provider for F5 BigIP. Resources are currently available for LTM.

[![Build Status](https://travis-ci.org/DealerDotCom/terraform-provider-bigip.svg?branch=master)](https://travis-ci.org/DealerDotCom/terraform-provider-bigip)


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

`auth_token` - (Optional, Default=false) Enable to use an external authentication source (LDAP, TACACS, etc)

`login_ref` - (Optional, Default="tmos") Login reference for token authentication (see BIG-IP REST docs for details)

# Resources

For resources should be named with their "full path". The full path is the combination of the partition + name of the resource.
For example `/Common/my-pool`.

## bigip_ltm_monitor

Configures a custom monitor for use by health checks.

### Example
```
resource "bigip_ltm_monitor" "monitor" {
  name = "/Common/terraform_monitor"
  parent = "http"
  send = "GET /some/path\r\n"
  timeout = "999"
  interval = "999"
}
```

### Reference

`name` - (Required) Name of the monitor

`parent` - (Required) Existing LTM monitor to inherit from

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
  name = "/Common/terraform_node1"
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
  name = "/Common/terraform-pool"
  load_balancing_mode = "round-robin"
  nodes = ["${bigip_ltm_node.node.name}:80"]
  monitors = ["${bigip_ltm_monitor.monitor.name}","${bigip_ltm_monitor.monitor2.name}"]
  allow_snat = false
}
```

### Reference

`name` - (Required) Name of the pool

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
  name = "/Common/terraform_vs_http"
  destination = "10.12.12.12"
  port = 80
  pool = "${bigip_ltm_pool.pool.name}"
}
```

### Reference

`name` - (Required) Name of the virtual server

`port` - (Required) Listen port for the virtual server

`destination` - (Required) Destination IP

`pool` - (Optional) Default pool name

`mask` - (Optional) Mask can either be in CIDR notation or decimal, i.e.: `24` or `255.255.255.0`. A CIDR mask of `0` is the same as `0.0.0.0`

`source_address_translation` - (Optional) Can be either omitted for `none` or the values `automap` or `snat`

`ip_protocol` - (Optional) Specify the IP protocol to use with the the virtual server (all, tcp, or udp are valid)

## bigip_ltm_irule

Creates iRules

### Example

```
# Loading from a file is the preferred method
resource "bigip_ltm_irule" "rule" {
  name = "/Common/terraform_irule"
  irule = "${file("myirule.tcl")}"
}

resource "bigip_ltm_irule" "rule2" {
  name = "/Common/terraform_irule2"
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

Configures a Virtual Address. NOTE: create/delete are not implemented
since the virtual addresses should be created/deleted automatically
with the corresponding virtual server.

### Example 

```
resource "bigip_ltm_virtual_address" "vs_va" {

    name = "/Common/${bigip_ltm_virtual_server.vs.destination}"
    advertize_route = true
}
```

### Reference

`name` - (Required) Name of the virtual address

`description` - (Optional) Description of the virtual address

`advertize_route` - (Optional) Enabled dynamic routing of the address

`conn_limit` - (Optional, Default=0) Max number of connections for virtual address

`enabled` - (Optional, Default=true) Enable or disable the virtual address

`arp` - (Optional, Default=true) Enable or disable ARP for the virtual address

`auto_delete` - (Optional, Default=true) Automatically delete the virtual address with the virtual server 

`icmp_echo` - (Optional, Default=true) Enable/Disable ICMP response to the virtual address

`traffic_group` - (Optional, Default=/Common/traffic-group-1) Specify the partition and traffic group

## bigip_ltm_policy

Configure [local traffic policies](https://support.f5.com/kb/en-us/solutions/public/15000/000/sol15085.html).
This is a fairly low level resource that does little to make actually using policies any simpler. A solid
understanding of how policies and their associated rules, actions and conditions
are managed through iControlREST is recommended.

### Example 

```
resource "bigip_ltm_policy" "policy" {
  name = "/Common/my_policy"
  strategy = "/Common/first-match"
  requires = ["http"]
  controls = ["forwarding"]
  rule {
    name = "/Common/rule1"

    condition {
      httpUri = true
      startsWith = true
      values = ["/foo"]
    }

    condition {
      httpMethod = true
      values = ["GET"]
    }

    action {
      forward = true
      pool = "/Common/my_pool"
    }
  }
}
```

### Reference

`name` - (Required) Name of the policy

`strategy` - (Required) Strategy selection when more than one rule matches.

`requires` - (Required) Defines the types of conditions that you can use when configuring a rule.

`controls` - (Required) Defines the types of actions that you can use when configuring a rule.

`rule` - defines a single rule to add to the policy. Multiple rules can be defined for a single policy.
 
**Rules**
 
 Actions and Conditions support all fields available via the iControlREST API. You can see all of the 
 available fields in the [iControlREST API documentation](https://devcentral.f5.com/d/icontrol-rest-api-reference-version-120).
 Each field in the actions and conditions objects is available. Pro tip: Create your policy via the GUI first then use
 the REST API to figure out how to configure the terraform resource.
 
 `name` (Required) - Name of the rule
 
 `action` - Defines a single action. Multiple actions can exist per rule.
 
 `condition` - Defines a single condition. Multiple conditions can exist per rule.


# Building

Create the distributable packages like so:

```
make && make dist
```

# Testing

Running the acceptance test suite requires an F5 to test against. Set `BIGIP_HOST`, `BIGIP_USER`
and `BIGIP_PASSWORD` to a device to run the tests against. By default tests will use the `Common` 
partition for creating objects. You can change the partition by setting `BIGIP_TEST_PARTITION`.

```
BIGIP_HOST=f5.mycompany.com BIGIP_USER=foo BIGIP_PASSWORD=secret make testacc
```


Read [here](https://github.com/hashicorp/terraform/blob/master/.github/CONTRIBUTING.md#running-an-acceptance-test) for
more information about acceptance testing in Terraform.