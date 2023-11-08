---
layout: "bigip"
page_title: "BIG-IP: bigip_fast_http_app"
subcategory: "F5 Automation Tool Chain(ATC)"
description: |-
  Provides details about bigip_fast_http_app resource
---

# bigip_fast_http_app

`bigip_fast_http_app` This resource will create and manage FAST HTTP applications on BIG-IP 

[FAST documentation](https://clouddocs.f5.com/products/extensions/f5-appsvcs-templates/latest/)

## Example Usage

```hcl

resource "bigip_fast_http_app" "fast_http_app" {
  tenant      = "fasthttptenant"
  application = "fasthttpapp"
  virtual_server {
    ip   = "10.30.30.44"
    port = 443
  }
}

```

## Example Usage with service discovery

```hcl

data "bigip_fast_azure_service_discovery" "TC3" {
  resource_group  = "testazurerg"
  subscription_id = "testazuresid"
  tag_key         = "testazuretag"
  tag_value       = "testazurevalue"
}

data "bigip_fast_gce_service_discovery" "TC3" {
  tag_key   = "testgcetag"
  tag_value = "testgcevalue"
  region    = "testgceregion"
}
resource "bigip_fast_http_app" "fast_https_app" {
  tenant      = "fasthttptenant"
  application = "fasthttpapp"
  virtual_server {
    ip   = "10.30.40.44"
    port = 443
  }
  pool_members {
    addresses = ["10.11.40.120", "10.11.30.121", "10.11.30.122"]
    port      = 80
  }
  service_discovery = [data.bigip_fast_gce_service_discovery.TC3.gce_sd_json, data.bigip_fast_azure_service_discovery.TC3.azure_sd_json]
}
```

## Argument Reference

* `tenant` - (Required, `string`) Name of the FAST HTTPS application tenant.

* `application` - (Required ,`string`) Name of the FAST HTTPS application.

* `virtual_server` - (Optional,`set`) `virtual_server` block will provide `ip` and `port` options to be used for virtual server.
See [virtual server](#virtual-server) below for more details. 

* `existing_snat_pool` - (Optional,`string`) Name of an existing BIG-IP SNAT pool.

* `snat_pool_address` - (Optional,`list`) List of address to be used for FAST-Generated SNAT Pool.

* `exist_pool_name` - (Optional,`string`) Name of an existing BIG-IP pool.

* `pool_members` - (Optional,`set`) `pool_members` block takes input for FAST-Generated Pool.
See [Pool Members](#pool-members) below for more details.

* `service_discovery` - (Optional,`list`) List of different cloud service discovery config provided as string, provided `service_discovery` block to Automatically Discover Pool Members with Service Discovery on different clouds.
      
* `load_balancing_mode` - (Optional,`string`) A `load balancing method` is an algorithm that the BIG-IP system uses to select a pool member for processing a request. F5 recommends the Least Connections load balancing method
    
* `slow_ramp_time` - (Optional,`int`) Slow ramp temporarily throttles the number of connections to a new pool member. The recommended value is 300 seconds
                                            
* `existing_monitor` - (Optional,`string`) Name of an existing BIG-IP HTTPS pool monitor. Monitors are used to determine the health of the application on each server.

* `monitor` - (Optional,`set`) `monitor` block takes input for FAST-Generated Pool Monitor.
See [Pool Monitor](#pool-monitor) below for more details.

* `persistence_profile` - (Optional,`string`) Name of an existing BIG-IP persistence profile to be used.

* `persistence_type` - (Optional,`string`) Type of persistence profile to be created. Using this option will enable use of FAST generated persistence profiles.

* `fallback_persistence` - (Optional,`string`) Type of fallback persistence record to be created for each new client connection.

* `existing_waf_security_policy` - (Optional,`string`) Name of an existing WAF Security policy.

* `endpoint_ltm_policy` - (Optional,`list`) List of LTM Policies to be applied FAST HTTP Application.

* `waf_security_policy` - (Optional,`set`) `waf_security_policy` block takes input for FAST-Generated WAF Security Policy.
See [WAF Security Policy](#waf-security-policy) below for more details.

* `security_log_profiles` - (Optional,`list`) List of security log profiles to be used for FAST application

### virtual server
This IP address, combined with the port you specify below, becomes the BIG-IP virtual server address and port, which clients use to access the application

The `virtual_server` block supports the following:

* `ip` - (Optional , `string`) IP4/IPv6 address to be used for virtual server ex: `10.1.1.1`

* `port` -(Optional , `int`) Port number to used for accessing virtual server/application

### Pool Members

Using this block will `enable` for FAST-Generated Pool.

The `pool_members` block supports the following:

* `addresses` - (Optional , `list`) List of server address to be used for FAST-Generated Pool.

* `port` - (Optional , `int`) port number of serviceport to be used for FAST-Generated Pool.

* `connection_limit` - (Optional , `int`) connectionLimit value to be used for FAST-Generated Pool.

* `priority_group` - (Optional , `int`) priorityGroup value to be used for FAST-Generated Pool.

* `share_nodes` - (Optional , `bool`) shareNodes value to be used for FAST-Generated Pool.


### Pool Monitor

Using this block will `enable` for FAST-Generated Pool Monitor.

The `monitor` block supports the following:

* `monitor_auth` - (Optional , `bool`) set `true` if the servers require login credentials for web access on FAST-Generated Pool Monitor. default is `false`.

* `username` - (Optional , `string`) username for web access on FAST-Generated Pool Monitor.

* `password` - (Optional , `string`) password for web access on FAST-Generated Pool Monitor.

* `interval` - (Optional , `int`) Set the time between health checks,in seconds for FAST-Generated Pool Monitor. 

* `send_string` - (Optional , `string`) Specify data to be sent during each health check for FAST-Generated Pool Monitor.

* `response` - (Optional , `string`) The presence of this string anywhere in the HTTP response implies availability.

### WAF Security policy
Using this block will `enable` for FAST-Generated WAF Security Policy

The `waf_security_policy` block supports the following:

* `enable` - (Optional , `bool`) Setting `true` will enable FAST to create WAF Security Policy.
