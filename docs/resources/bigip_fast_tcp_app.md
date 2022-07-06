---
layout: "bigip"
page_title: "BIG-IP: bigip_fast_tcp_app"
subcategory: "F5 Automation Tool Chain(ATC)"
description: |-
  Provides details about bigip_fast_tcp_app resource
---

# bigip_fast_tcp_app

`bigip_fast_tcp_app` This resource will create and manage FAST TCP applications on BIG-IP from provided JSON declaration. 


## Example Usage


```hcl

resource "bigip_fast_tcp_app" "fast-tcp-app" {
  application = "tcp_app_2"
  tenant      = "tcp_app_tenant"

  virtual_server = {
    ip   = "11.12.16.30"
    port = 443
  }

  fast_create_pool_members {
    addresses        = ["10.11.34.65", "56.43.23.76"]
    port             = 443
    priority_group   = 1
    connection_limit = 4
    share_nodes      = true
  }
}

```

## Argument Reference

* `application` - (Required) Name of the FAST TCP application.

* `tenant` - (Required) Name of the FAST TCP application tenant.
  
* `virtual_server` - (Optional,`set`) `virtual_server` block will provide `ip` and `port` options to be used for virtual server.
See [virtual server](#virtual-server) below for more details. 

* `existing_snat_pool` - (Optional,`string`) Name of an existing BIG-IP SNAT pool.

* `fast_create_snat_pool_address` - (Optional,`list`) List of address to be used for FAST-Generated SNAT Pool.

* `exist_pool_name` - (Optional,`string`) Name of an existing BIG-IP pool.

* `fast_create_pool_members` - (Optional,`set`) `fast_create_pool_members` block takes input for FAST-Generated Pool.
See [Pool Members](#pool-members) below for more details.

* `load_balancing_mode` - (Optional,`string`) A `load balancing method` is an algorithm that the BIG-IP system uses to select a pool member for processing a request. F5 recommends the Least Connections load balancing method

* `slow_ramp_time` - (Optional,`int`) Slow ramp temporarily throttles the number of connections to a new pool member. The recommended value is 300 seconds

* `existing_monitor` - (Optional,`string`) Name of an existing BIG-IP HTTPS pool monitor. Monitors are used to determine the health of the application on each server.

* `fast_create_monitor` - (Optional,`set`) `fast_create_monitor` block takes input for FAST-Generated Pool Monitor.
See [Pool Monitor](#pool-monitor) below for more details.


### virtual server
This IP address, combined with the port you specify below, becomes the BIG-IP virtual server address and port, which clients use to access the application

The `virtual_server` block supports the following:

* `ip` - (Optional , `string`) IP4/IPv6 address to be used for virtual server ex: `10.1.1.1`

* `port` -(Optional , `int`) Port number to used for accessing virtual server/application


### Pool Members

Using this block will `enable` for FAST-Generated Pool.

The `fast_create_pool_members` block supports the following:

* `addresses` - (Optional , `list`) List of server address to be used for FAST-Generated Pool.

* `port` - (Optional , `int`) port number of serviceport to be used for FAST-Generated Pool.

* `connection_limit` - (Optional , `int`) connectionLimit value to be used for FAST-Generated Pool.

* `priority_group` - (Optional , `int`) priorityGroup value to be used for FAST-Generated Pool.

* `share_nodes` - (Optional , `bool`) shareNodes value to be used for FAST-Generated Pool.


### Pool Monitor

Using this block will `enable` for FAST-Generated Pool Monitor. Unlike FAST HTTP and HTTPS apps, TCP only has the `interval` option here.

The `fast_create_monitor` block supports the following:

* `interval` - (Optional , `int`) Set the time between health checks,in seconds for FAST-Generated Pool Monitor. 
