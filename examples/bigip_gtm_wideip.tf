terraform {
  required_providers {
    bigip = {
      source  = "F5Networks/bigip"
    #   version = "~> 1.0"
    }
  }
}

provider "bigip" {
  address  = "10.145.71.31"
  username = "admin"
  password = "F5site02"
}

# Example 1: Basic GTM WideIP of type 'a' with minimal configuration
resource "bigip_gtm_wideip" "example_basic" {
  name      = "testwideip1.local"
  type      = "a"
  partition = "Common"
  
  description = "test_wideip_a"
}

# Example 2: GTM WideIP with full configuration
resource "bigip_gtm_wideip" "example_full" {
  name      = "fullwideip.example.com"
  type      = "a"
  partition = "Common"
  
  description              = "Full configuration WideIP example"
  enabled                  = true
  failure_rcode            = "noerror"
  failure_rcode_response   = "disabled"
  failure_rcode_ttl        = 0
  last_resort_pool         = "a /Common/firstpool"
  minimal_response         = "disabled"
  persist_cidr_ipv4        = 32
  persist_cidr_ipv6        = 128
  persistence              = "disabled"
  pool_lb_mode             = "round-robin"
  ttl_persistence          = 3600
  topology_prefer_edns0_client_subnet = "enabled"
  
  # Optional: Enable load balancing decision logging
  load_balancing_decision_log_verbosity = ["pool-selection"]
  
  # Optional: Add aliases for the WideIP
  aliases = ["fullwideip-alias1.local", "fullwideip-alias2.example.com"]
}

# Example 3: GTM WideIP of type 'aaaa' (IPv6)
resource "bigip_gtm_wideip" "example_ipv6" {
  name      = "ipv6wideip.example.com"
  type      = "aaaa"
  partition = "Common"
  
  description       = "IPv6 WideIP"
  enabled           = true
  pool_lb_mode      = "topology"
  minimal_response  = "disabled"
}

# Example 4: GTM WideIP of type 'cname'
resource "bigip_gtm_wideip" "example_cname" {
  name      = "alias.example.com"
  type      = "cname"
  partition = "Common"
  
  description  = "CNAME WideIP"
  enabled      = true
}

# Example 5: Disabled WideIP
resource "bigip_gtm_wideip" "example_disabled" {
  name      = "disabled.example.com"
  type      = "a"
  partition = "Common"
  
  description = "Disabled WideIP"
  enabled     = false
  disabled    = true
}

# Example 6: WideIP with custom persistence settings
resource "bigip_gtm_wideip" "example_persistence" {
  name      = "persistent.example.com"
  type      = "a"
  partition = "Common"
  
  description       = "WideIP with persistence enabled"
  enabled           = true
  persistence       = "enabled"
  persist_cidr_ipv4 = 24
  persist_cidr_ipv6 = 64
  ttl_persistence   = 7200
}

# Output examples
output "wideip_id" {
  description = "The ID of the WideIP"
  value       = bigip_gtm_wideip.example_basic.id
}

output "wideip_name" {
  description = "The name of the WideIP"
  value       = bigip_gtm_wideip.example_basic.name
}
