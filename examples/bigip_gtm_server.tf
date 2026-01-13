
# Example 1: Basic BIG-IP server
resource "bigip_gtm_server" "basic_bigip" {
  name       = "test_bigip_server"
  partition  = "Common"
  datacenter = "/Common/testdc"
  product    = "bigip"

  addresses {
    name = "10.10.10.10"
  }

  monitor = "/Common/bigip"
}

# Example 2: BIG-IP server with multiple addresses
resource "bigip_gtm_server" "multi_address" {
  name       = "multi_addr_server"
  partition  = "Common"
  datacenter = "/Common/testdc"
  product    = "bigip"

  addresses {
    name = "192.168.1.100"
  }

  addresses {
    name = "192.168.1.101"
  }

  monitor = "/Common/bigip"
}

# Example 3: BIG-IP server with NAT translation
resource "bigip_gtm_server" "nat_server" {
  name       = "nat_bigip_server"
  partition  = "Common"
  datacenter = "/Common/testdc"
  product    = "bigip"

  addresses {
    name        = "10.1.1.100"
    translation = "203.0.113.100"
  }

  monitor = "/Common/bigip"
}

# Example 4: BIG-IP server with device names
resource "bigip_gtm_server" "device_name_server" {
  name       = "device_name_server"
  partition  = "Common"
  datacenter = "/Common/testdc"
  product    = "bigip"

  addresses {
    name        = "10.2.2.100"
    device_name = "bigip1.example.com"
  }

  addresses {
    name        = "10.2.2.101"
    device_name = "bigip2.example.com"
  }

  monitor = "/Common/bigip"
}

# Example 5: Generic host server (non-BIG-IP)
resource "bigip_gtm_server" "generic_host" {
  name       = "generic_web_server"
  partition  = "Common"
  datacenter = "/Common/testdc"
  product    = "generic-host"

  addresses {
    name = "203.0.113.50"
  }

  monitor = "/Common/http"
}

# Example 6: Server with virtual server discovery enabled
resource "bigip_gtm_server" "vs_discovery" {
  name       = "vs_discovery_server"
  partition  = "Common"
  datacenter = "/Common/testdc"
  product    = "bigip"

  addresses {
    name = "10.3.3.100"
  }

  virtual_server_discovery = "enabled"
  monitor                  = "/Common/bigip"
}

# Example 7: Server with custom monitoring settings
resource "bigip_gtm_server" "custom_monitor" {
  name       = "custom_monitor_server"
  partition  = "Common"
  datacenter = "/Common/testdc"
  product    = "bigip"

  addresses {
    name = "10.4.4.100"
  }

  monitor = "/Common/bigip and /Common/tcp"
}

# Example 8: Server with prober settings
resource "bigip_gtm_server" "prober_server" {
  name       = "prober_server"
  partition  = "Common"
  datacenter = "/Common/testdc"
  product    = "bigip"

  addresses {
    name = "10.5.5.100"
  }

  prober_fallback   = "any-available"
  prober_preference = "inside-datacenter"
  monitor           = "/Common/bigip"
}

# Example 9: Server with resource limits
resource "bigip_gtm_server" "limited_server" {
  name       = "limited_server"
  partition  = "Common"
  datacenter = "/Common/testdc"
  product    = "bigip"

  addresses {
    name = "10.6.6.100"
  }

  limit_cpu_usage        = 80
  limit_cpu_usage_status = "enabled"
  limit_mem_avail        = 1024
  limit_mem_avail_status = "enabled"
  limit_max_bps          = 1000000
  limit_max_bps_status   = "enabled"

  monitor = "/Common/bigip"
}

# Example 10: Comprehensive configuration
resource "bigip_gtm_server" "comprehensive" {
  name        = "comprehensive_server"
  partition   = "Common"
  datacenter  = "/Common/testdc"
  description = "Comprehensive GTM server configuration example"
  product     = "bigip"

  addresses {
    name        = "10.7.7.100"
    device_name = "bigip-primary.example.com"
    translation = "203.0.113.200"
  }

  addresses {
    name        = "10.7.7.101"
    device_name = "bigip-secondary.example.com"
    translation = "203.0.113.201"
  }

  # Monitoring
  monitor = "/Common/bigip and /Common/tcp"

  # Prober settings
  prober_fallback   = "any-available"
  prober_preference = "inside-datacenter"

  # Resource limits
  limit_cpu_usage        = 90
  limit_cpu_usage_status = "enabled"
  limit_mem_avail        = 2048
  limit_mem_avail_status = "enabled"
  limit_max_bps          = 5000000
  limit_max_bps_status   = "enabled"
  limit_max_pps          = 10000
  limit_max_pps_status   = "enabled"

  # Connection limits
  limit_max_connections        = 50000
  limit_max_connections_status = "enabled"

  # Virtual server discovery
  virtual_server_discovery = "enabled"

  # State
  enabled = true
}

# Example 11: Disabled server
resource "bigip_gtm_server" "disabled_server" {
  name       = "disabled_server"
  partition  = "Common"
  datacenter = "/Common/testdc"
  product    = "bigip"

  addresses {
    name = "10.8.8.100"
  }

  enabled = false
  monitor = "/Common/bigip"
}

# Output examples
output "basic_server_id" {
  value       = bigip_gtm_server.basic_bigip.id
  description = "ID of the basic GTM server"
}

output "comprehensive_server_addresses" {
  value = [
    for addr in bigip_gtm_server.comprehensive.addresses :
    "${addr.name} (${addr.translation})"
  ]
  description = "Addresses of the comprehensive server"
}
