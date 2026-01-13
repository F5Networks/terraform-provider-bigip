# bigip_gtm_pool Resource

Provides a BIG-IP GTM (Global Traffic Manager) Pool resource. This resource allows you to configure and manage GTM Pool objects on a BIG-IP system.

## Description

A GTM pool is a collection of virtual servers or other pool members that can be distributed across multiple data centers. GTM pools are used by WideIPs to intelligently distribute DNS traffic based on various load balancing algorithms and health monitoring.

GTM Pool types correspond to different DNS record types:
- **a**: IPv4 address pools
- **aaaa**: IPv6 address pools
- **cname**: Canonical name pools
- **mx**: Mail exchange pools
- **naptr**: Naming authority pointer pools
- **srv**: Service locator pools

## Example Usage

### Basic Pool

```hcl
resource "bigip_gtm_pool" "example" {
  name      = "my_pool"
  type      = "a"
  partition = "Common"
  
  load_balancing_mode = "round-robin"
  monitor             = "/Common/https"
}
```

### Pool with Members

```hcl
resource "bigip_gtm_pool" "with_members" {
  name      = "app_pool"
  type      = "a"
  partition = "Common"
  
  load_balancing_mode = "round-robin"
  monitor             = "/Common/https"
  ttl                 = 30
  
  members {
    name         = "server1:/Common/vs_app"
    enabled      = true
    ratio        = 1
    member_order = 0
  }
  
  members {
    name         = "server2:/Common/vs_app"
    enabled      = true
    ratio        = 1
    member_order = 1
  }
}
```

### Advanced Pool Configuration

```hcl
resource "bigip_gtm_pool" "advanced" {
  name      = "advanced_pool"
  type      = "a"
  partition = "Common"
  
  # Load balancing configuration
  load_balancing_mode   = "round-robin"
  alternate_mode        = "topology"
  fallback_mode         = "return-to-dns"
  fallback_ip           = "192.0.2.1"
  
  # Response configuration
  max_answers_returned  = 2
  ttl                   = 60
  
  # Monitoring
  monitor               = "/Common/https"
  verify_member_availability = "enabled"
  
  # QoS weights
  qos_hit_ratio         = 10
  qos_hops              = 5
  qos_kilobytes_second  = 5
  qos_lcs               = 50
  qos_packet_rate       = 5
  qos_rtt               = 100
  
  # Connection limits
  limit_max_connections        = 5000
  limit_max_connections_status = "enabled"
  limit_max_bps                = 100000000
  limit_max_bps_status         = "enabled"
  
  # Minimum members requirement
  min_members_up_mode  = "at-least"
  min_members_up_value = 2
  
  members {
    name                         = "server1:/Common/vs_app"
    enabled                      = true
    ratio                        = 2
    member_order                 = 0
    monitor                      = "default"
    limit_max_connections        = 2000
    limit_max_connections_status = "enabled"
  }
  
  members {
    name         = "server2:/Common/vs_app"
    enabled      = true
    ratio        = 1
    member_order = 1
  }
}
```

## Argument Reference

The following arguments are supported:

### Required Arguments

* `name` - (Required, String) The name of the GTM pool. Forces new resource.
* `type` - (Required, String) The type of GTM pool. Valid values are: `a`, `aaaa`, `cname`, `mx`, `naptr`, `srv`. Forces new resource.

### Optional Arguments

#### General Settings

* `partition` - (Optional, String) The partition in which the pool resides. Default: `Common`. Forces new resource.
* `enabled` - (Optional, Boolean) Enable or disable the pool. Default: `true`.
* `disabled` - (Optional, Boolean) Disabled state of the pool. Default: `false`.
* `monitor` - (Optional, String) Specifies the health monitor for the pool. Example: `/Common/https`.

#### Load Balancing Settings

* `load_balancing_mode` - (Optional, String) Specifies the preferred load balancing mode for the pool. Valid values: `round-robin`, `ratio`, `topology`, `global-availability`, `virtual-server-capacity`, `least-connections`, `lowest-round-trip-time`, `fewest-hops`, `packet-rate`, `cpu`, `completion-rate`, `quality-of-service`, `kilobytes-per-second`, `drop-packet`, `fallback-ip`, `virtual-server-score`, `dynamic-ratio`. Default: `round-robin`.
* `alternate_mode` - (Optional, String) Specifies the load balancing mode to use if the preferred and fallback modes are unsuccessful. Default: `round-robin`.
* `fallback_mode` - (Optional, String) Specifies the load balancing mode that the system uses if the pool's preferred and alternate modes are unsuccessful. Default: `return-to-dns`.
* `fallback_ip` - (Optional, String) Specifies the IPv4 or IPv6 address of the server to which the system directs requests when it cannot use one of its pools. Default: `any`.
* `dynamic_ratio` - (Optional, String) Enables or disables the dynamic ratio load balancing algorithm. Valid values: `enabled`, `disabled`. Default: `disabled`.
* `manual_resume` - (Optional, String) Specifies whether manual resume is enabled. Valid values: `enabled`, `disabled`. Default: `disabled`.

#### Response Settings

* `max_answers_returned` - (Optional, Integer) Specifies the maximum number of available virtual servers that the system lists in a response. Default: `1`.
* `ttl` - (Optional, Integer) Specifies the time to live (TTL) in seconds for the pool. Default: `30`.
* `verify_member_availability` - (Optional, String) Specifies whether the system verifies the availability of pool members before sending traffic. Valid values: `enabled`, `disabled`. Default: `enabled`.

#### QoS Settings

QoS (Quality of Service) weights determine how the system ranks virtual servers when using QoS load balancing modes:

* `qos_hit_ratio` - (Optional, Integer) Specifies the weight for QoS hit ratio. Default: `5`.
* `qos_hops` - (Optional, Integer) Specifies the weight for QoS hops. Default: `0`.
* `qos_kilobytes_second` - (Optional, Integer) Specifies the weight for QoS kilobytes per second. Default: `3`.
* `qos_lcs` - (Optional, Integer) Specifies the weight for QoS link capacity. Default: `30`.
* `qos_packet_rate` - (Optional, Integer) Specifies the weight for QoS packet rate. Default: `1`.
* `qos_rtt` - (Optional, Integer) Specifies the weight for QoS round trip time. Default: `50`.
* `qos_topology` - (Optional, Integer) Specifies the weight for QoS topology. Default: `0`.
* `qos_vs_capacity` - (Optional, Integer) Specifies the weight for QoS virtual server capacity. Default: `0`.
* `qos_vs_score` - (Optional, Integer) Specifies the weight for QoS virtual server score. Default: `0`.

#### Connection Limits

* `limit_max_bps` - (Optional, Integer) Specifies the maximum allowable data throughput rate in bits per second. Default: `0`.
* `limit_max_bps_status` - (Optional, String) Enables or disables the limit_max_bps option. Valid values: `enabled`, `disabled`. Default: `disabled`.
* `limit_max_connections` - (Optional, Integer) Specifies the maximum number of concurrent connections. Default: `0`.
* `limit_max_connections_status` - (Optional, String) Enables or disables the limit_max_connections option. Valid values: `enabled`, `disabled`. Default: `disabled`.
* `limit_max_pps` - (Optional, Integer) Specifies the maximum allowable data transfer rate in packets per second. Default: `0`.
* `limit_max_pps_status` - (Optional, String) Enables or disables the limit_max_pps option. Valid values: `enabled`, `disabled`. Default: `disabled`.

#### Minimum Members Settings

* `min_members_up_mode` - (Optional, String) Specifies whether the minimum number of members must be up for the pool to be active. Valid values: `off`, `at-least`, `percent`. Default: `off`.
* `min_members_up_value` - (Optional, Integer) Specifies the minimum number (or percentage) of pool members that must be up. Default: `0`.

#### Members Block

* `members` - (Optional, Set) A set of pool members. Each member supports:
  * `name` - (Required, String) Name of the pool member in the format `<server_name>:<virtual_server_name>`. Example: `server1:/Common/vs_app`.
  * `enabled` - (Optional, Boolean) Enable or disable the pool member. Default: `true`.
  * `disabled` - (Optional, Boolean) Disabled state of the pool member. Default: `false`.
  * `ratio` - (Optional, Integer) Specifies the weight of the pool member for load balancing. Default: `1`.
  * `member_order` - (Optional, Integer) Specifies the order in which the member will be used. Default: `0`.
  * `monitor` - (Optional, String) Specifies the health monitor for this pool member. Default: `default`.
  * `limit_max_bps` - (Optional, Integer) Specifies the maximum allowable data throughput rate for this member. Default: `0`.
  * `limit_max_bps_status` - (Optional, String) Enables or disables the limit_max_bps option for this member. Default: `disabled`.
  * `limit_max_connections` - (Optional, Integer) Specifies the maximum number of concurrent connections for this member. Default: `0`.
  * `limit_max_connections_status` - (Optional, String) Enables or disables the limit_max_connections option for this member. Default: `disabled`.
  * `limit_max_pps` - (Optional, Integer) Specifies the maximum allowable data transfer rate in packets per second for this member. Default: `0`.
  * `limit_max_pps_status` - (Optional, String) Enables or disables the limit_max_pps option for this member. Default: `disabled`.

## Attribute Reference

In addition to the arguments listed above, the following attributes are exported:

* `id` - The unique identifier of the pool resource in Terraform (format: `/partition/name:type`).

## Import

GTM Pool resources can be imported using the format `/<partition>/<name>:<type>`. For example:

```bash
terraform import bigip_gtm_pool.example /Common/my_pool:a
```

## Notes

### Pool Member Name Format

Pool members must be specified in the format: `<server_name>:<virtual_server_name>`

Examples:
- `server1:/Common/vs_app` - References virtual server `vs_app` on server `server1`
- `dc1_server:/Prod/app_vs` - References virtual server `app_vs` in partition `Prod` on server `dc1_server`

The server and virtual server must already exist in the GTM configuration.

### Load Balancing Modes

The `load_balancing_mode` determines how GTM distributes DNS queries across pool members:

- **round-robin**: Distributes queries equally across all available members
- **ratio**: Distributes queries based on member ratios
- **topology**: Uses topology records to determine the best member
- **global-availability**: Considers member availability and load
- **virtual-server-capacity**: Based on virtual server capacity
- **least-connections**: Selects member with fewest active connections
- **lowest-round-trip-time**: Selects member with lowest RTT
- **fewest-hops**: Selects member with fewest network hops
- **packet-rate**: Based on packet transmission rate
- **cpu**: Based on CPU utilization
- **completion-rate**: Based on connection completion rate
- **quality-of-service**: Based on QoS metrics
- **kilobytes-per-second**: Based on throughput
- **dynamic-ratio**: Dynamically adjusts member ratios
- **drop-packet**: Drops DNS packets (used for testing)
- **fallback-ip**: Returns a fallback IP address
- **virtual-server-score**: Based on virtual server scores

### QoS Weights

QoS (Quality of Service) weights are used when the load balancing mode is set to `quality-of-service`. Higher weights give more importance to specific metrics:

- **qos_hit_ratio**: Cache hit ratio
- **qos_hops**: Number of router hops
- **qos_kilobytes_second**: Data throughput
- **qos_lcs**: Link capacity score
- **qos_packet_rate**: Packet transmission rate
- **qos_rtt**: Round trip time
- **qos_topology**: Topology distance
- **qos_vs_capacity**: Virtual server capacity
- **qos_vs_score**: Virtual server score

### Connection Limits

Connection limits can be set at both the pool level and individual member level:
- Pool-level limits apply to the entire pool
- Member-level limits apply to individual members
- Both limits must have their corresponding `_status` field set to `enabled` to take effect

### Minimum Members

The `min_members_up_mode` and `min_members_up_value` work together:
- **off**: No minimum requirement
- **at-least**: At least `min_members_up_value` members must be up
- **percent**: At least `min_members_up_value` percent of members must be up

Example: If you have 5 members and set `min_members_up_mode = "at-least"` and `min_members_up_value = 2`, the pool will be marked down if fewer than 2 members are available.

## API Endpoints

This resource interacts with the following BIG-IP API endpoints:

- `GET /mgmt/tm/gtm/pool/<type>/<name>?expandSubcollections=true` - Read pool configuration
- `POST /mgmt/tm/gtm/pool/<type>` - Create pool
- `PUT /mgmt/tm/gtm/pool/<type>/<name>` - Update pool configuration
- `DELETE /mgmt/tm/gtm/pool/<type>/<name>` - Delete pool

## Related Resources

- `bigip_gtm_wideip` - Manages GTM WideIPs that reference pools
- `bigip_gtm_server` - Manages GTM servers that contain virtual servers
- `bigip_gtm_datacenter` - Manages GTM data centers
- `bigip_gtm_monitor` - Manages GTM health monitors
