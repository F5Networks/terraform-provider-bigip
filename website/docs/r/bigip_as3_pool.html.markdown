---
layout: "bigip"
page_title: "BIG-IP: bigip_as3_pool"
sidebar_current: "docs-bigip-datasource-pool-x"
description: |-
   Provides details about bigip_as3_pool datasource
---
 
# bigip\_as3\_pool
 
`bigip_as3_pool` Manages a Pool class, which  contain your servers as well as health monitors and load balancing methods and more.
 
In the following example, our pool is web_pool, it’s using the default HTTP health monitor, and includes two servers on port 80.
 
## Example Usage
 
 
```hcl
data "bigip_as3_pool" "mydataas3pool" {
  name = "web_pool"
  loadbalancing_mode = "round-robin"
  servicedown_action = "none"
  pool_members {
    connection_limit = 10
    rate_limit=10
    dynamic_ratio=100
    service_port=8080
    ratio=90
    priority_group=30
    sharenodes=true
    adminstate="enable"
    server_addresses=["192.168.30.1","192.168.25.1"]
  }
  minimummembers_active=1
  reselect_tries=0
  slowramp_time=10
  minimum_monitors=1
  monitors=["http"]
}
}
```
 
## Argument Reference
 
* `name` - (Required) Name of the pool
 
* `loadbalancing_mode` - (Optional) Specifies method used for automatic balancing and distributing traffic across real physical servers
 
* `pool_members` - (Required) Set of Pool members
 
* `label` - (Optional) Optional friendly name for this object
 
* `remark` - (Optional) Arbitrary (brief) text pertaining to this object
 
* `servicedown_action` - (Optional) Specifies connection handling when member is non-responsive
 
* `minimummembers_active`- (Optional) Pool is down when fewer than this number of members are up
 
* `reselect_tries` - (Optional) Maximum number of attempts to find a responsive member for a connection
 
* `slowramp_time` - (Optional) AS3 slowly the connection rate to a newly-active member slowly during this interval (seconds)
 
* `minimum_monitors` - (Optional) Member is down when fewer than minimum monitors report it healthy. Specify ‘all’ to require all monitors to be up
 
* `monitors` - (Required) List of health monitors (each by name or AS3 pointer)
 
* Below attributes needs to be configured under pool_members option.
 
* `connection_limit` - (Optional) Maximum concurrent connections to member
 
* `rate_limit` - (Optional) Value zero prevents use of member
 
* `dynamic_ratio` - (Optional) Specifies a range of numbers that you want the system to use in conjunction with the ratio load balancing method
 
* `service_port` - (Required) Service L4 port (optional port-discovery may override)
 
* `ratio` - (Optional) Specifies the weight of the pool member for load balancing purposes
 
* `priority_group` - (Optional) Specifies the priority group within the pool for this pool member
 
* `sharenodes` - (Optional) If enabled, nodes are created in /Common instead of the tenant’s partition
 
* `adminstate` - (Optional) Setting adminState to enable will create the node in an operational state. Set to disable to disallow new connections but allow existing connections to drain. Set to offline to force immediate termination of all connections.
 
* `server_addresses` - (Required) Static IP addresses of servers (nodes)
 
 
