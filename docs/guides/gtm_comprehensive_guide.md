# GTM (Global Traffic Manager) Comprehensive Usage Guide

This guide provides a comprehensive overview of using F5 BIG-IP GTM resources with the Terraform provider. It covers all available GTM resources, their relationships, and provides end-to-end examples for common use cases.

## Table of Contents

1. [Overview](#overview)
2. [GTM Architecture](#gtm-architecture)
3. [Resource Hierarchy](#resource-hierarchy)
4. [Available Resources](#available-resources)
5. [Quick Start](#quick-start)
6. [Detailed Resource Usage](#detailed-resource-usage)
7. [End-to-End Examples](#end-to-end-examples)
8. [Best Practices](#best-practices)
9. [Troubleshooting](#troubleshooting)
10. [API Reference](#api-reference)

---

## Overview

F5 BIG-IP GTM (Global Traffic Manager) is a DNS-based traffic management solution that provides intelligent DNS resolution for geographically distributed applications. It enables:

- **Geographic Load Balancing**: Route users to the nearest or most appropriate data center
- **Disaster Recovery**: Automatic failover when primary sites become unavailable
- **Application Availability**: Distribute traffic based on application health and performance
- **DNS Resolution**: Act as an authoritative DNS server for your domains

## GTM Architecture

```
                    ┌─────────────────┐
                    │    WideIP       │
                    │ (DNS Name)      │
                    └────────┬────────┘
                             │
                    ┌────────▼────────┐
                    │     Pools       │
                    │ (Load Balance)  │
                    └────────┬────────┘
                             │
              ┌──────────────┼──────────────┐
              │              │              │
     ┌────────▼────┐  ┌──────▼─────┐  ┌─────▼───────┐
     │   Server    │  │   Server   │  │   Server    │
     │ (Big-IP/    │  │ (Generic)  │  │ (Big-IP)    │
     │  Generic)   │  │            │  │             │
     └──────┬──────┘  └─────┬──────┘  └──────┬──────┘
            │               │                │
     ┌──────▼──────┐  ┌─────▼──────┐  ┌──────▼──────┐
     │   Virtual   │  │   Virtual  │  │   Virtual   │
     │   Servers   │  │   Servers  │  │   Servers   │
     └─────────────┘  └────────────┘  └─────────────┘
            │               │                │
     ┌──────▼──────┐  ┌─────▼──────┐  ┌──────▼──────┐
     │ Datacenter  │  │ Datacenter │  │ Datacenter  │
     │   (West)    │  │  (Central) │  │   (East)    │
     └─────────────┘  └────────────┘  └─────────────┘
```

## Resource Hierarchy

GTM resources must be created in a specific order due to dependencies:

1. **Datacenters** - Physical/logical locations (no dependencies)
2. **Servers** - Systems hosting applications (depend on Datacenters)
3. **Virtual Servers** - Services on servers (part of Server resource)
4. **Pools** - Groups of virtual servers (depend on Servers/Virtual Servers)
5. **WideIPs** - DNS names (depend on Pools)

## Available Resources

| Resource | Description | Status |
|----------|-------------|--------|
| `bigip_gtm_datacenter` | Logical representation of a physical location | ✅ Available |
| `bigip_gtm_server` | BIG-IP system or generic host server | ✅ Available |
| `bigip_gtm_pool` | Collection of virtual servers for load balancing | ✅ Available |
| `bigip_gtm_wideip` | DNS name that GTM resolves | ✅ Available |
| `bigip_gtm_monitor_http` | HTTP health monitor | ✅ Available |
| `bigip_gtm_monitor_https` | HTTPS health monitor | ✅ Available |
| `bigip_gtm_monitor_tcp` | TCP health monitor | ✅ Available |
| `bigip_gtm_monitor_postgresql` | PostgreSQL health monitor | ✅ Available |
| `bigip_gtm_monitor_bigip` | BIG-IP health monitor | ✅ Available |
| `bigip_gtm_prober_pool` | Pool of probers for health checking | ❌ Not Yet Implemented |
| `bigip_gtm_topology` | Topology records for geographic routing | ❌ Not Yet Implemented |

---

## Quick Start

### Minimal Configuration

This example creates a basic GTM setup with a datacenter, server, pool, and wideip:

```hcl
terraform {
  required_providers {
    bigip = {
      source  = "f5networks/bigip"
      version = ">= 1.22.0"
    }
  }
}

provider "bigip" {
  address  = "10.0.0.1"
  username = "admin"
  password = "password"
}

# 1. Create Datacenter
resource "bigip_gtm_datacenter" "dc1" {
  name      = "datacenter1"
  partition = "Common"
  enabled   = true
}

# 2. Create Server with Virtual Servers
resource "bigip_gtm_server" "server1" {
  name       = "app_server"
  partition  = "Common"
  datacenter = bigip_gtm_datacenter.dc1.id
  product    = "generic-host"

  addresses {
    name = "192.168.1.100"
  }

  virtual_server_discovery = "disabled"

  virtual_servers {
    name        = "vs_web"
    destination = "192.168.1.100:80"
    enabled     = true
  }

  enabled = true
}

# 3. Create Pool
resource "bigip_gtm_pool" "pool1" {
  name                = "web_pool"
  type                = "a"
  partition           = "Common"
  load_balancing_mode = "round-robin"

  members {
    name    = "/Common/app_server:vs_web"
    enabled = true
    ratio   = 1
  }

  depends_on = [bigip_gtm_server.server1]
}

# 4. Create WideIP
resource "bigip_gtm_wideip" "wideip1" {
  name             = "www.example.com"
  type             = "a"
  partition        = "Common"
  pool_lb_mode     = "round-robin"
  last_resort_pool = "a /Common/web_pool"
  enabled          = true

  depends_on = [bigip_gtm_pool.pool1]
}
```

---

## Detailed Resource Usage

### 1. GTM Datacenter (`bigip_gtm_datacenter`)

A datacenter represents a physical location where servers reside.

#### Basic Example

```hcl
resource "bigip_gtm_datacenter" "west_coast" {
  name              = "west_coast_dc"
  partition         = "Common"
  enabled           = true
  location          = "Seattle, WA"
  contact           = "ops@example.com"
  description       = "Primary West Coast datacenter"
  prober_preference = "inside-datacenter"
  prober_fallback   = "any-available"
}
```

#### Key Arguments

| Argument | Type | Default | Description |
|----------|------|---------|-------------|
| `name` | string | required | Name of the datacenter |
| `partition` | string | "Common" | Partition |
| `enabled` | bool | true | Enable/disable the datacenter |
| `prober_preference` | string | "inside-datacenter" | Preferred prober location |
| `prober_fallback` | string | "any-available" | Fallback prober selection |

---

### 2. GTM Server (`bigip_gtm_server`)

A server represents a BIG-IP system or generic host.

#### Server Types

| Product Type | Description | Virtual Server Discovery |
|--------------|-------------|-------------------------|
| `bigip` | F5 BIG-IP system | Automatic via iQuery |
| `generic-host` | Non-F5 server | Manual definition required |
| `redundant-bigip` | HA BIG-IP pair | Automatic via iQuery |
| `single-bigip` | Standalone BIG-IP | Automatic via iQuery |

#### BIG-IP Server Example

```hcl
resource "bigip_gtm_server" "bigip_ltm" {
  name       = "bigip_ltm_server"
  partition  = "Common"
  datacenter = bigip_gtm_datacenter.west_coast.id
  product    = "bigip"

  addresses {
    name        = "10.145.71.31"
    device_name = "/Common/bigip1.example.com"
  }

  virtual_server_discovery = "enabled"
  monitor                  = "/Common/bigip"
  enabled                  = true

  # Prober settings
  prober_preference = "inside-datacenter"
  prober_fallback   = "any-available"

  # iQuery settings
  iq_allow_path          = true
  iq_allow_service_check = true
  iq_allow_snmp          = true
}
```

#### Generic Host Server with Virtual Servers

```hcl
resource "bigip_gtm_server" "generic_app" {
  name       = "generic_app_server"
  partition  = "Common"
  datacenter = bigip_gtm_datacenter.west_coast.id
  product    = "generic-host"

  addresses {
    name = "192.168.10.100"
  }

  virtual_server_discovery = "disabled"
  link_discovery           = "disabled"

  # Manually define virtual servers
  virtual_servers {
    name        = "vs_http"
    destination = "192.168.10.100:80"
    enabled     = true
  }

  virtual_servers {
    name        = "vs_https"
    destination = "192.168.10.100:443"
    enabled     = true
  }

  virtual_servers {
    name        = "vs_api"
    destination = "192.168.10.100:8080"
    enabled     = true
  }

  enabled = true
}
```

#### Server with Resource Limits

```hcl
resource "bigip_gtm_server" "limited_server" {
  name       = "limited_server"
  datacenter = bigip_gtm_datacenter.west_coast.id
  product    = "bigip"

  addresses {
    name = "10.10.10.10"
  }

  virtual_server_discovery = "enabled"
  monitor                  = "/Common/bigip"

  # Resource limits
  limit_max_connections        = 10000
  limit_max_connections_status = "enabled"
  limit_max_bps                = 1000000000
  limit_max_bps_status         = "enabled"
  limit_max_pps                = 100000
  limit_max_pps_status         = "enabled"
  limit_cpu_usage              = 80
  limit_cpu_usage_status       = "enabled"
  limit_mem_avail              = 2048
  limit_mem_avail_status       = "enabled"

  enabled = true
}
```

---

### 3. GTM Pool (`bigip_gtm_pool`)

A pool is a collection of virtual servers used for load balancing.

#### Pool Types

| Type | Description | DNS Record |
|------|-------------|------------|
| `a` | IPv4 address pool | A record |
| `aaaa` | IPv6 address pool | AAAA record |
| `cname` | Canonical name pool | CNAME record |
| `mx` | Mail exchange pool | MX record |
| `naptr` | Naming authority pointer | NAPTR record |
| `srv` | Service locator pool | SRV record |

#### A Record Pool Example

```hcl
resource "bigip_gtm_pool" "web_pool" {
  name      = "web_pool"
  type      = "a"
  partition = "Common"

  # Load balancing configuration
  load_balancing_mode = "round-robin"
  alternate_mode      = "global-availability"
  fallback_mode       = "return-to-dns"
  fallback_ip         = "192.168.1.1"

  # Response settings
  ttl                  = 300
  max_answers_returned = 1

  # Monitoring
  verify_member_availability = "enabled"

  # Minimum members
  min_members_up_mode  = "at-least"
  min_members_up_value = 1

  # Pool members
  members {
    name         = "/Common/generic_app_server:vs_http"
    enabled      = true
    ratio        = 2
    member_order = 0
    monitor      = "default"
  }

  members {
    name         = "/Common/generic_app_server:vs_https"
    enabled      = true
    ratio        = 1
    member_order = 1
    monitor      = "default"
  }
}
```

#### Pool with QoS Settings

```hcl
resource "bigip_gtm_pool" "qos_pool" {
  name                = "qos_pool"
  type                = "a"
  partition           = "Common"
  load_balancing_mode = "quality-of-service"

  # QoS weights
  qos_hit_ratio        = 10
  qos_hops             = 5
  qos_kilobytes_second = 5
  qos_lcs              = 50
  qos_packet_rate      = 5
  qos_rtt              = 100
  qos_topology         = 10
  qos_vs_capacity      = 10
  qos_vs_score         = 10

  ttl = 60
}
```

#### Load Balancing Modes

| Mode | Description |
|------|-------------|
| `round-robin` | Equal distribution across members |
| `ratio` | Distribution based on member ratios |
| `topology` | Based on topology records |
| `global-availability` | Based on member availability |
| `virtual-server-capacity` | Based on VS capacity |
| `least-connections` | Fewest active connections |
| `lowest-round-trip-time` | Lowest RTT |
| `fewest-hops` | Fewest network hops |
| `quality-of-service` | Based on QoS metrics |

---

### 4. GTM WideIP (`bigip_gtm_wideip`)

A WideIP is the DNS name that GTM resolves.

#### Basic WideIP

```hcl
resource "bigip_gtm_wideip" "www" {
  name         = "www.example.com"
  type         = "a"
  partition    = "Common"
  description  = "Main website WideIP"
  pool_lb_mode = "round-robin"
  enabled      = true
}
```

#### WideIP with Last Resort Pool

```hcl
resource "bigip_gtm_wideip" "app" {
  name             = "app.example.com"
  type             = "a"
  partition        = "Common"
  description      = "Application WideIP"
  pool_lb_mode     = "round-robin"
  last_resort_pool = "a /Common/web_pool"
  minimal_response = "enabled"
  enabled          = true
}
```

#### WideIP with Persistence

```hcl
resource "bigip_gtm_wideip" "persistent" {
  name              = "session.example.com"
  type              = "a"
  partition         = "Common"
  pool_lb_mode      = "round-robin"
  persistence       = "enabled"
  persist_cidr_ipv4 = 24
  persist_cidr_ipv6 = 64
  ttl_persistence   = 7200
  enabled           = true
}
```

#### WideIP with Failure RCODE

```hcl
resource "bigip_gtm_wideip" "failover" {
  name                   = "failover.example.com"
  type                   = "a"
  partition              = "Common"
  pool_lb_mode           = "global-availability"
  failure_rcode          = "servfail"
  failure_rcode_response = "enabled"
  failure_rcode_ttl      = 300
  enabled                = true
}
```

#### WideIP with Aliases

```hcl
resource "bigip_gtm_wideip" "with_aliases" {
  name         = "primary.example.com"
  type         = "a"
  partition    = "Common"
  pool_lb_mode = "round-robin"
  enabled      = true

  aliases = [
    "www.example.com",
    "web.example.com",
    "app.example.com"
  ]
}
```

#### WideIP with Topology-based Load Balancing

```hcl
resource "bigip_gtm_wideip" "topology" {
  name                                = "geo.example.com"
  type                                = "a"
  partition                           = "Common"
  pool_lb_mode                        = "topology"
  topology_prefer_edns0_client_subnet = "enabled"
  enabled                             = true

  load_balancing_decision_log_verbosity = [
    "pool-selection",
    "pool-member-selection"
  ]
}
```

---

## End-to-End Examples

### Example 1: Multi-Datacenter Web Application

This example demonstrates a complete GTM setup for a web application distributed across two datacenters:

```hcl
# ===========================================
# Provider Configuration
# ===========================================
terraform {
  required_providers {
    bigip = {
      source  = "f5networks/bigip"
      version = ">= 1.22.0"
    }
  }
}

provider "bigip" {
  address  = "10.145.71.31"
  username = "admin"
  password = "F5site02"
}

# ===========================================
# Variables
# ===========================================
variable "environment" {
  default = "production"
}

# ===========================================
# Datacenters
# ===========================================
resource "bigip_gtm_datacenter" "west" {
  name              = "${var.environment}_dc_west"
  partition         = "Common"
  enabled           = true
  location          = "Seattle, WA"
  description       = "West Coast Datacenter"
  prober_preference = "inside-datacenter"
  prober_fallback   = "any-available"
}

resource "bigip_gtm_datacenter" "east" {
  name              = "${var.environment}_dc_east"
  partition         = "Common"
  enabled           = true
  location          = "New York, NY"
  description       = "East Coast Datacenter"
  prober_preference = "inside-datacenter"
  prober_fallback   = "any-available"
}

# ===========================================
# Servers
# ===========================================
resource "bigip_gtm_server" "web_west" {
  name       = "${var.environment}_web_west"
  partition  = "Common"
  datacenter = bigip_gtm_datacenter.west.id
  product    = "generic-host"

  addresses {
    name = "192.168.1.100"
  }

  virtual_server_discovery = "disabled"
  link_discovery           = "disabled"

  virtual_servers {
    name        = "vs_http"
    destination = "192.168.1.100:80"
    enabled     = true
  }

  virtual_servers {
    name        = "vs_https"
    destination = "192.168.1.100:443"
    enabled     = true
  }

  enabled = true
}

resource "bigip_gtm_server" "web_east" {
  name       = "${var.environment}_web_east"
  partition  = "Common"
  datacenter = bigip_gtm_datacenter.east.id
  product    = "generic-host"

  addresses {
    name = "192.168.2.100"
  }

  virtual_server_discovery = "disabled"
  link_discovery           = "disabled"

  virtual_servers {
    name        = "vs_http"
    destination = "192.168.2.100:80"
    enabled     = true
  }

  virtual_servers {
    name        = "vs_https"
    destination = "192.168.2.100:443"
    enabled     = true
  }

  enabled = true
}

# ===========================================
# Pools
# ===========================================
resource "bigip_gtm_pool" "www_pool" {
  name      = "${var.environment}_www_pool"
  type      = "a"
  partition = "Common"

  load_balancing_mode        = "round-robin"
  alternate_mode             = "global-availability"
  fallback_mode              = "return-to-dns"
  fallback_ip                = "192.168.1.1"
  ttl                        = 300
  max_answers_returned       = 1
  verify_member_availability = "enabled"

  members {
    name         = "/Common/${var.environment}_web_west:vs_https"
    enabled      = true
    ratio        = 1
    member_order = 0
  }

  members {
    name         = "/Common/${var.environment}_web_east:vs_https"
    enabled      = true
    ratio        = 1
    member_order = 1
  }

  depends_on = [
    bigip_gtm_server.web_west,
    bigip_gtm_server.web_east
  ]
}

# ===========================================
# WideIP
# ===========================================
resource "bigip_gtm_wideip" "www" {
  name              = "www.myapp.com"
  type              = "a"
  partition         = "Common"
  description       = "Main website"
  pool_lb_mode      = "round-robin"
  last_resort_pool  = "a /Common/${var.environment}_www_pool"
  persistence       = "enabled"
  persist_cidr_ipv4 = 24
  ttl_persistence   = 3600
  minimal_response  = "enabled"
  enabled           = true

  aliases = [
    "myapp.com",
    "app.myapp.com"
  ]

  depends_on = [bigip_gtm_pool.www_pool]
}

# ===========================================
# Outputs
# ===========================================
output "datacenter_west_id" {
  value = bigip_gtm_datacenter.west.id
}

output "datacenter_east_id" {
  value = bigip_gtm_datacenter.east.id
}

output "wideip_name" {
  value = bigip_gtm_wideip.www.name
}
```

### Example 2: API Gateway with Multiple Endpoints

```hcl
# API Gateway GTM Configuration

# Datacenter
resource "bigip_gtm_datacenter" "api_dc" {
  name    = "api_datacenter"
  enabled = true
}

# API Server
resource "bigip_gtm_server" "api" {
  name       = "api_server"
  datacenter = bigip_gtm_datacenter.api_dc.id
  product    = "generic-host"

  addresses {
    name = "10.0.1.100"
  }

  virtual_server_discovery = "disabled"

  # Different API endpoints
  virtual_servers {
    name        = "vs_api_v1"
    destination = "10.0.1.100:8081"
    enabled     = true
  }

  virtual_servers {
    name        = "vs_api_v2"
    destination = "10.0.1.100:8082"
    enabled     = true
  }

  virtual_servers {
    name        = "vs_graphql"
    destination = "10.0.1.100:4000"
    enabled     = true
  }

  enabled = true
}

# Pool for API v1
resource "bigip_gtm_pool" "api_v1" {
  name                = "api_v1_pool"
  type                = "a"
  load_balancing_mode = "round-robin"
  ttl                 = 60

  members {
    name    = "/Common/api_server:vs_api_v1"
    enabled = true
  }

  depends_on = [bigip_gtm_server.api]
}

# Pool for API v2
resource "bigip_gtm_pool" "api_v2" {
  name                = "api_v2_pool"
  type                = "a"
  load_balancing_mode = "round-robin"
  ttl                 = 60

  members {
    name    = "/Common/api_server:vs_api_v2"
    enabled = true
  }

  depends_on = [bigip_gtm_server.api]
}

# WideIP for API v1
resource "bigip_gtm_wideip" "api_v1" {
  name             = "api-v1.example.com"
  type             = "a"
  pool_lb_mode     = "round-robin"
  last_resort_pool = "a /Common/api_v1_pool"
  enabled          = true

  depends_on = [bigip_gtm_pool.api_v1]
}

# WideIP for API v2
resource "bigip_gtm_wideip" "api_v2" {
  name             = "api-v2.example.com"
  type             = "a"
  pool_lb_mode     = "round-robin"
  last_resort_pool = "a /Common/api_v2_pool"
  enabled          = true

  depends_on = [bigip_gtm_pool.api_v2]
}
```

---

## Best Practices

### 1. Resource Naming

Use consistent naming conventions:

```hcl
locals {
  env_prefix = "${var.environment}_"
  app_name   = "myapp"
}

resource "bigip_gtm_datacenter" "dc" {
  name = "${local.env_prefix}${var.region}_dc"
}
```

### 2. Use Dependencies

Always specify `depends_on` for proper creation order:

```hcl
resource "bigip_gtm_pool" "pool" {
  # ... configuration ...
  depends_on = [bigip_gtm_server.server]
}

resource "bigip_gtm_wideip" "wideip" {
  # ... configuration ...
  depends_on = [bigip_gtm_pool.pool]
}
```

### 3. Enable Health Monitoring

Always configure proper monitoring:

```hcl
resource "bigip_gtm_pool" "pool" {
  verify_member_availability = "enabled"
  # ...
}
```

### 4. Set Appropriate TTLs

Use appropriate TTL values for your use case:

| Use Case | Recommended TTL |
|----------|----------------|
| Static content | 3600 (1 hour) |
| Dynamic applications | 300 (5 minutes) |
| Real-time failover | 30-60 seconds |
| Testing | 10-30 seconds |

### 5. Use Persistence Carefully

Enable persistence only when needed:

```hcl
# For stateful applications
resource "bigip_gtm_wideip" "stateful" {
  persistence       = "enabled"
  persist_cidr_ipv4 = 24 # /24 subnet
  ttl_persistence   = 7200
}

# For stateless applications - no persistence needed
resource "bigip_gtm_wideip" "stateless" {
  persistence = "disabled"
}
```

---

## Troubleshooting

### Common Issues

#### 1. Pool Member Not Found

**Error**: `Pool member not found: server:vs_name`

**Solution**: Ensure the server and virtual server exist:
```bash
curl -sk -u admin:password \
  https://bigip/mgmt/tm/gtm/server/~Common~server_name/virtual-servers
```

#### 2. Datacenter Reference Error

**Error**: `Datacenter not found`

**Solution**: Use the full path or ID:

```text
# Option 1: Use resource reference (recommended)
datacenter = bigip_gtm_datacenter.dc.id

# Option 2: Use full path string
datacenter = "/Common/datacenter_name"
```

#### 3. WideIP Pool Reference

**Error**: `Invalid last resort pool format`

**Solution**: Use correct format: `<type> <partition>/<pool_name>`

```text
# Correct format (includes type prefix)
last_resort_pool = "a /Common/pool_name"

# Wrong format (missing type prefix)
last_resort_pool = "/Common/pool_name"
```

### Verification Commands

```bash
# List all datacenters
curl -sk -u admin:password https://bigip/mgmt/tm/gtm/datacenter

# List all servers
curl -sk -u admin:password https://bigip/mgmt/tm/gtm/server

# List virtual servers on a server
curl -sk -u admin:password \
  "https://bigip/mgmt/tm/gtm/server/~Common~server_name/virtual-servers"

# List pools (A record type)
curl -sk -u admin:password https://bigip/mgmt/tm/gtm/pool/a

# List pool members
curl -sk -u admin:password \
  "https://bigip/mgmt/tm/gtm/pool/a/~Common~pool_name/members"

# List wideips (A record type)
curl -sk -u admin:password https://bigip/mgmt/tm/gtm/wideip/a
```

---

## API Reference

### GTM API Endpoints

| Resource | Endpoint |
|----------|----------|
| Datacenter | `/mgmt/tm/gtm/datacenter` |
| Server | `/mgmt/tm/gtm/server` |
| Virtual Server | `/mgmt/tm/gtm/server/{name}/virtual-servers` |
| Pool (A) | `/mgmt/tm/gtm/pool/a` |
| Pool (AAAA) | `/mgmt/tm/gtm/pool/aaaa` |
| Pool (CNAME) | `/mgmt/tm/gtm/pool/cname` |
| WideIP (A) | `/mgmt/tm/gtm/wideip/a` |
| WideIP (AAAA) | `/mgmt/tm/gtm/wideip/aaaa` |
| WideIP (CNAME) | `/mgmt/tm/gtm/wideip/cname` |

### Import Commands

```bash
# Import datacenter
terraform import bigip_gtm_datacenter.dc /Common/datacenter_name

# Import server
terraform import bigip_gtm_server.server /Common/server_name

# Import pool
terraform import bigip_gtm_pool.pool /Common/pool_name:a

# Import wideip
terraform import bigip_gtm_wideip.wideip a:/Common/www.example.com
```

---

## Not Yet Implemented Resources

The following GTM resources are not yet available in the Terraform provider but can be configured via direct API calls:

### GTM Monitors

```bash
# Create HTTP monitor
curl -sk -u admin:password -H "Content-Type: application/json" \
  -X POST https://bigip/mgmt/tm/gtm/monitor/http \
  -d '{"name":"my_http_monitor","interval":30,"timeout":120}'
```

### GTM Prober Pool

```bash
# Create prober pool
curl -sk -u admin:password -H "Content-Type: application/json" \
  -X POST https://bigip/mgmt/tm/gtm/prober-pool \
  -d '{"name":"my_prober_pool","enabled":true}'
```

### GTM Topology Records

```bash
# Create topology record
curl -sk -u admin:password -H "Content-Type: application/json" \
  -X POST https://bigip/mgmt/tm/gtm/topology \
  -d '{"name":"ldns: subnet 10.0.0.0/8 server: datacenter dc1"}'
```

---

## Version History

| Provider Version | Features |
|-----------------|----------|
| 1.22.0+ | GTM Server virtual_servers block support |
| 1.21.0+ | GTM Pool members with full attributes |
| 1.20.0+ | GTM basic resources (datacenter, server, pool, wideip) |

---

## Additional Resources

- [F5 GTM Documentation](https://techdocs.f5.com/en-us/bigip-16-1-0/big-ip-dns-services-implementations.html)
- [Terraform Provider Documentation](https://registry.terraform.io/providers/F5Networks/bigip/latest/docs)
- [F5 iControl REST API Reference](https://clouddocs.f5.com/api/icontrol-rest/)
