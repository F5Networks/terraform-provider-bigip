# bigip_gtm_wideip Resource

Provides a BIG-IP GTM (Global Traffic Manager) WideIP resource. This resource allows you to configure and manage GTM WideIP objects on a BIG-IP system.

## Description

A WideIP is a DNS name that GTM resolves on behalf of an authoritative DNS server. WideIPs are the core objects in GTM that enable intelligent DNS-based load balancing and failover across multiple data centers.

GTM WideIP types correspond to different DNS record types:
- **a**: IPv4 address records
- **aaaa**: IPv6 address records
- **cname**: Canonical name records
- **mx**: Mail exchange records
- **naptr**: Naming authority pointer records
- **srv**: Service locator records

## Example Usage

### Basic WideIP

```hcl
resource "bigip_gtm_wideip" "example" {
  name      = "testwideip.local"
  type      = "a"
  partition = "Common"
  
  description = "test_wideip_a"
}
```

### WideIP with Last Resort Pool

```hcl
resource "bigip_gtm_wideip" "with_pool" {
  name             = "app.example.com"
  type             = "a"
  partition        = "Common"
  
  description      = "Application WideIP"
  last_resort_pool = "a /Common/firstpool"
  pool_lb_mode     = "round-robin"
  minimal_response = "enabled"
}
```

### Advanced WideIP Configuration

```hcl
resource "bigip_gtm_wideip" "advanced" {
  name      = "advanced.example.com"
  type      = "a"
  partition = "Common"
  
  description              = "Advanced WideIP configuration"
  enabled                  = true
  failure_rcode            = "servfail"
  failure_rcode_response   = "enabled"
  failure_rcode_ttl        = 300
  last_resort_pool         = "a /Common/backup_pool"
  minimal_response         = "disabled"
  persist_cidr_ipv4        = 24
  persist_cidr_ipv6        = 64
  persistence              = "enabled"
  pool_lb_mode             = "topology"
  ttl_persistence          = 7200
  topology_prefer_edns0_client_subnet = "enabled"
  
  load_balancing_decision_log_verbosity = ["pool-selection", "pool-member-selection"]
  aliases = ["app1.example.com", "app2.example.com"]
}
```

### IPv6 WideIP

```hcl
resource "bigip_gtm_wideip" "ipv6" {
  name      = "ipv6.example.com"
  type      = "aaaa"
  partition = "Common"
  
  description  = "IPv6 WideIP"
  enabled      = true
  pool_lb_mode = "round-robin"
}
```

## Argument Reference

The following arguments are supported:

### Required Arguments

* `name` - (Required, String) The name of the WideIP. This should be a fully qualified domain name (FQDN). Forces new resource.
* `type` - (Required, String) The type of WideIP. Valid values are: `a`, `aaaa`, `cname`, `mx`, `naptr`, `srv`. Forces new resource.

### Optional Arguments

* `partition` - (Optional, String) The partition in which the WideIP resides. Default: `Common`. Forces new resource.
* `description` - (Optional, String) User-defined description of the WideIP.
* `enabled` - (Optional, Boolean) Enable or disable the WideIP. Default: `true`.
* `disabled` - (Optional, Boolean) Disabled state of the WideIP. Default: `false`.
* `failure_rcode` - (Optional, String) Specifies the DNS RCODE (response code) used when `failure_rcode_response` is enabled. Valid values: `noerror`, `formerr`, `servfail`, `nxdomain`, `notimp`, `refused`, `yxdomain`, `yxrrset`, `nxrrset`, `notauth`, `notzone`. Default: `noerror`.
* `failure_rcode_response` - (Optional, String) Specifies whether to return a DNS RCODE response when the WideIP is unavailable. Valid values: `enabled`, `disabled`. Default: `disabled`.
* `failure_rcode_ttl` - (Optional, Integer) Specifies the negative caching TTL (in seconds) of the SOA record for the RCODE response. Default: `0`.
* `last_resort_pool` - (Optional, String) Specifies the last resort pool for the WideIP. Format: `<type> <partition>/<pool_name>` (e.g., `a /Common/firstpool`). This pool is used when all other pools are unavailable.
* `load_balancing_decision_log_verbosity` - (Optional, Set of Strings) Specifies the amount of detail logged when making load balancing decisions. Valid values: `pool-selection`, `pool-member-selection`, `pool-traversal`, `pool-member-traversal`.
* `minimal_response` - (Optional, String) Specifies whether to minimize the response to DNS queries. When enabled, returns only the resource records required to satisfy the query. Valid values: `enabled`, `disabled`. Default: `enabled`.
* `persist_cidr_ipv4` - (Optional, Integer) Specifies the number of bits in the IPv4 address to use for persistence. Valid range: 0-32. Default: `32`.
* `persist_cidr_ipv6` - (Optional, Integer) Specifies the number of bits in the IPv6 address to use for persistence. Valid range: 0-128. Default: `128`.
* `persistence` - (Optional, String) Specifies whether to use persistence for the WideIP. When enabled, GTM ensures that subsequent DNS requests from the same client IP address are directed to the same pool member. Valid values: `enabled`, `disabled`. Default: `disabled`.
* `pool_lb_mode` - (Optional, String) Specifies the load balancing method for pools in the WideIP. Valid values: `round-robin`, `ratio`, `topology`, `global-availability`, `virtual-server-capacity`, `least-connections`, `lowest-round-trip-time`, `fewest-hops`, `packet-rate`, `cpu`, `completion-rate`, `quality-of-service`, `kilobytes-per-second`, `drop-packet`, `fallback-ip`, `virtual-server-score`. Default: `round-robin`.
* `ttl_persistence` - (Optional, Integer) Specifies the time to live (TTL) in seconds for persistence records. Default: `3600`.
* `topology_prefer_edns0_client_subnet` - (Optional, String) Specifies whether to prefer EDNS0 client subnet data for topology-based load balancing. Valid values: `enabled`, `disabled`. Default: `disabled`.
* `aliases` - (Optional, Set of Strings) Specifies alternate domain names (aliases) for the WideIP. These are additional names that resolve to the same WideIP configuration. Example: `["alias1.example.com", "alias2.example.com"]`.

## Attribute Reference

In addition to the arguments listed above, the following attributes are exported:

* `id` - The unique identifier of the WideIP resource in Terraform (format: `type:/partition/name`).

## Import

GTM WideIP resources can be imported using the format `type:/partition/name`. For example:

```bash
terraform import bigip_gtm_wideip.example a:/Common/testwideip.local
```

Additional import examples:

```bash
# Import an IPv6 WideIP
terraform import bigip_gtm_wideip.ipv6_example aaaa:/Common/ipv6.example.com

# Import a CNAME WideIP
terraform import bigip_gtm_wideip.cname_example cname:/Common/alias.example.com

# Import from a non-Common partition
terraform import bigip_gtm_wideip.prod_example a:/Production/app.example.com
```

## Notes

### Pool Load Balancing Modes

The `pool_lb_mode` determines how GTM distributes DNS queries across pools:

- **round-robin**: Distributes queries equally across all available pools
- **ratio**: Distributes queries based on pool ratios
- **topology**: Uses topology records to determine the best pool
- **global-availability**: Considers pool availability and load
- **virtual-server-capacity**: Based on virtual server capacity
- **least-connections**: Selects pool with fewest active connections
- **lowest-round-trip-time**: Selects pool with lowest RTT
- **fewest-hops**: Selects pool with fewest network hops
- **packet-rate**: Based on packet transmission rate
- **cpu**: Based on CPU utilization
- **completion-rate**: Based on connection completion rate
- **quality-of-service**: Based on QoS metrics
- **kilobytes-per-second**: Based on throughput
- **drop-packet**: Drops DNS packets (used for testing)
- **fallback-ip**: Returns a fallback IP address
- **virtual-server-score**: Based on virtual server scores

### Last Resort Pool Format

The `last_resort_pool` must be specified in the format: `<type> <partition>/<pool_name>`

Examples:
- `a /Common/firstpool` - IPv4 pool
- `aaaa /Common/ipv6pool` - IPv6 pool
- `cname /Prod/alias_pool` - CNAME pool

### Persistence

When persistence is enabled:
- GTM maintains a mapping of client IP addresses to pool members
- Subsequent requests from the same client are directed to the same destination
- The `persist_cidr_ipv4` and `persist_cidr_ipv6` settings determine the subnet mask used for grouping client IPs
- Persistence records expire after `ttl_persistence` seconds

### Failure RCODE Response

When a WideIP becomes unavailable (all pools are down):
- If `failure_rcode_response` is `disabled`: GTM returns no answer (NXDOMAIN)
- If `failure_rcode_response` is `enabled`: GTM returns the specified `failure_rcode`

Common RCODE values:
- **noerror**: No error (returns empty response)
- **servfail**: Server failure
- **nxdomain**: Non-existent domain
- **refused**: Query refused

### Aliases

Aliases allow you to specify alternate domain names for the same WideIP configuration. When a DNS query is made for any of the aliases, it is handled by the same WideIP configuration.

Example:
```hcl
resource "bigip_gtm_wideip" "app" {
  name = "app.example.com"
  type = "a"
  
  aliases = [
    "app-alias1.example.com",
    "app-alias2.example.com",
    "app.backup.example.com"
  ]
}
```

In this example, DNS queries for `app.example.com`, `app-alias1.example.com`, `app-alias2.example.com`, or `app.backup.example.com` will all be handled by the same WideIP configuration.

## API Endpoints

This resource interacts with the following BIG-IP API endpoints:

- `GET /mgmt/tm/gtm/wideip/<type>/<name>` - Read WideIP configuration
- `POST /mgmt/tm/gtm/wideip/<type>` - Create WideIP
- `PUT /mgmt/tm/gtm/wideip/<type>/<name>` - Update WideIP configuration
- `DELETE /mgmt/tm/gtm/wideip/<type>/<name>` - Delete WideIP

## Related Resources

- `bigip_gtm_pool` - Manages GTM pools that can be referenced by WideIPs
- `bigip_gtm_server` - Manages GTM servers that contain virtual servers
- `bigip_gtm_datacenter` - Manages GTM data centers
- `bigip_gtm_topology` - Manages topology records for topology-based load balancing
