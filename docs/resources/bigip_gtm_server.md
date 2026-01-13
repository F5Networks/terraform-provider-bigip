# bigip_gtm_server

Manages F5 BIG-IP GTM (Global Traffic Manager) Server resources.

A GTM server represents a BIG-IP system, a host, or a server that hosts applications. Servers are identified by their addresses and are organized within datacenters. GTM servers enable GTM to perform load balancing and provide health monitoring for distributed applications.

## Example Usage

### Basic GTM Server (BIG-IP)

```hcl
resource "bigip_gtm_datacenter" "dc1" {
  name = "datacenter1"
}

resource "bigip_gtm_server" "server1" {
  name       = "bigip_server1"
  datacenter = bigip_gtm_datacenter.dc1.name
  product    = "bigip"

  addresses {
    name = "10.1.1.1"
  }

  monitor                  = "/Common/bigip"
  virtual_server_discovery = true
  link_discovery          = "disabled"
}
```

### GTM Server with Multiple Addresses

```hcl
resource "bigip_gtm_datacenter" "dc1" {
  name = "datacenter1"
}

resource "bigip_gtm_server" "multi_address_server" {
  name       = "multi_server"
  datacenter = bigip_gtm_datacenter.dc1.name
  product    = "bigip"

  addresses {
    name        = "10.1.1.1"
    device_name = "/Common/bigip1.example.com"
    translation = "none"
  }

  addresses {
    name        = "10.1.1.2"
    device_name = "/Common/bigip2.example.com"
    translation = "none"
  }

  monitor                  = "/Common/bigip"
  virtual_server_discovery = true
}
```

### GTM Server with Address Translation

```hcl
resource "bigip_gtm_datacenter" "dc1" {
  name = "datacenter1"
}

resource "bigip_gtm_server" "nat_server" {
  name       = "nat_server"
  datacenter = bigip_gtm_datacenter.dc1.name
  product    = "bigip"

  addresses {
    name        = "10.10.10.10"
    device_name = "/Common/server.example.com"
    translation = "192.168.1.10"
  }

  monitor                  = "/Common/bigip"
  virtual_server_discovery = true
}
```

### Generic Host Server

```hcl
resource "bigip_gtm_datacenter" "dc1" {
  name = "datacenter1"
}

resource "bigip_gtm_server" "generic_host" {
  name       = "generic_server"
  datacenter = bigip_gtm_datacenter.dc1.name
  product    = "generic-host"

  addresses {
    name = "10.20.20.20"
  }

  monitor                  = "/Common/tcp"
  virtual_server_discovery = false
  link_discovery          = "disabled"
}
```

### GTM Server with Prober Settings

```hcl
resource "bigip_gtm_datacenter" "dc1" {
  name = "datacenter1"
}

resource "bigip_gtm_server" "prober_server" {
  name       = "prober_configured_server"
  datacenter = bigip_gtm_datacenter.dc1.name
  product    = "bigip"

  addresses {
    name = "10.30.30.30"
  }

  monitor                  = "/Common/bigip"
  virtual_server_discovery = true
  
  prober_preference = "inside-datacenter"
  prober_fallback   = "any-available"
  
  iq_allow_path           = true
  iq_allow_service_check  = true
  iq_allow_snmp           = true
}
```

### GTM Server with Resource Limits

```hcl
resource "bigip_gtm_datacenter" "dc1" {
  name = "datacenter1"
}

resource "bigip_gtm_server" "limited_server" {
  name       = "resource_limited_server"
  datacenter = bigip_gtm_datacenter.dc1.name
  product    = "bigip"

  addresses {
    name = "10.40.40.40"
  }

  monitor                  = "/Common/bigip"
  virtual_server_discovery = true

  limit_max_connections = 10000
  limit_max_bps         = 1000000
  limit_max_pps         = 50000
  limit_cpu_usage       = 80
  limit_mem_avail       = 1024
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the GTM server. Must be unique within the partition.

* `datacenter` - (Required) The datacenter where this server resides. Must be a valid datacenter name or full path (e.g., `/Common/datacenter1`).

* `partition` - (Optional) Partition or tenant the server belongs to. Default is `Common`.

* `product` - (Optional) Type of server. Valid values are:
  - `bigip` (default) - F5 BIG-IP system
  - `generic-host` - Generic host server
  - `redundant-bigip` - Redundant BIG-IP pair
  - `single-bigip` - Single BIG-IP system

* `addresses` - (Optional) List of IP addresses for the server. Each address block supports:
  - `name` - (Required) IP address
  - `device_name` - (Optional) Device name associated with the address
  - `translation` - (Optional) IP translation address. Default is `none`

* `monitor` - (Optional) Monitor assigned to check server health (e.g., `/Common/bigip`, `/Common/tcp`).

* `virtual_server_discovery` - (Optional) Enable or disable virtual server discovery. Default is `true`. When enabled, GTM automatically discovers virtual servers on BIG-IP systems.

* `link_discovery` - (Optional) Link discovery mode. Valid values:
  - `disabled` (default) - No link discovery
  - `enabled` - Enable link discovery
  - `enabled-no-delete` - Enable link discovery but don't delete existing links

* `prober_preference` - (Optional) Preferred type of prober. Valid values:
  - `inherit` (default) - Inherit from datacenter
  - `inside-datacenter` - Prefer probers inside the datacenter
  - `outside-datacenter` - Prefer probers outside the datacenter
  - `pool` - Use specific prober pool

* `prober_fallback` - (Optional) Fallback prober selection. Valid values:
  - `inherit` (default) - Inherit from datacenter
  - `any-available` - Use any available prober
  - `inside-datacenter` - Use probers inside datacenter
  - `outside-datacenter` - Use probers outside datacenter
  - `pool` - Use specific prober pool

* `expose_route_domains` - (Optional) Allow GTM server to expose route domains. Default is `false`.

* `iq_allow_path` - (Optional) Enable iQuery path probing. Default is `true`.

* `iq_allow_service_check` - (Optional) Enable iQuery service checking. Default is `true`.

* `iq_allow_snmp` - (Optional) Enable iQuery SNMP. Default is `true`.

* `limit_cpu_usage` - (Optional) Maximum CPU usage allowed (percent). 0 means no limit. Default is `0`.

* `limit_cpu_usage_status` - (Optional/Computed) CPU usage limit status.

* `limit_max_bps` - (Optional) Maximum bits per second. 0 means no limit. Default is `0`.

* `limit_max_bps_status` - (Optional/Computed) Maximum bps limit status.

* `limit_max_connections` - (Optional) Maximum concurrent connections. 0 means no limit. Default is `0`.

* `limit_max_connections_status` - (Optional/Computed) Maximum connections limit status.

* `limit_max_pps` - (Optional) Maximum packets per second. 0 means no limit. Default is `0`.

* `limit_max_pps_status` - (Optional/Computed) Maximum pps limit status.

* `limit_mem_avail` - (Optional) Available memory limit (MB). 0 means no limit. Default is `0`.

* `limit_mem_avail_status` - (Optional/Computed) Available memory limit status.

## Attribute Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `id` - The full path of the GTM server resource (e.g., `/Common/server1`).

## Import

GTM servers can be imported using the server name or full path:

```bash
terraform import bigip_gtm_server.example /Common/server1
```

or

```bash
terraform import bigip_gtm_server.example server1
```

## Notes

* When creating a GTM server of type `bigip`, ensure that the BIG-IP device is accessible and properly configured for GTM communication.

* Virtual server discovery requires proper iQuery communication between GTM systems.

* Address translation is useful when servers are behind NAT.

* Multiple addresses can be specified for servers with multiple network interfaces or for redundancy.

* Resource limits help prevent a single server from consuming all available capacity in load balancing decisions.

* Prober settings control how GTM monitors server health from different network locations.
