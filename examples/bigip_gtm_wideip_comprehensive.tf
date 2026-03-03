# Test Configuration for GTM WideIP with All Features
# This configuration demonstrates all available fields including aliases

terraform {
  required_providers {
    bigip = {
      source  = "F5Networks/bigip"
      version = "~> 1.0"
    }
  }
}

provider "bigip" {
  address  = "10.145.71.31"
  username = "admin"
  password = "F5site02"
}

# Test 1: WideIP matching the API response exactly
resource "bigip_gtm_wideip" "api_match" {
  name      = "testwideip.local"
  type      = "a"
  partition = "Common"
  
  description = "test_wideip_a"
  enabled     = true
  
  # Failure configuration
  failure_rcode          = "noerror"
  failure_rcode_response = "disabled"
  failure_rcode_ttl      = 0
  
  # Pool configuration
  last_resort_pool = "a /Common/firstpool"
  pool_lb_mode     = "round-robin"
  
  # Response settings
  minimal_response = "disabled"
  
  # Persistence settings
  persistence       = "disabled"
  persist_cidr_ipv4 = 32
  persist_cidr_ipv6 = 128
  ttl_persistence   = 3600
  
  # Topology settings
  topology_prefer_edns0_client_subnet = "enabled"
  
  # Logging
  load_balancing_decision_log_verbosity = ["pool-selection"]
  
  # Aliases
  aliases = ["testwideip2.local"]
}

# Test 2: WideIP with multiple aliases
resource "bigip_gtm_wideip" "multi_alias" {
  name      = "app.example.com"
  type      = "a"
  partition = "Common"
  
  description = "Application WideIP with multiple aliases"
  enabled     = true
  
  pool_lb_mode     = "round-robin"
  minimal_response = "enabled"
  
  # Multiple aliases
  aliases = [
    "app-primary.example.com",
    "app-secondary.example.com",
    "app-backup.example.com",
    "www.app.example.com"
  ]
}

# Test 3: WideIP with advanced logging
resource "bigip_gtm_wideip" "advanced_logging" {
  name      = "logged.example.com"
  type      = "a"
  partition = "Common"
  
  description = "WideIP with detailed logging"
  enabled     = true
  
  pool_lb_mode = "topology"
  
  # Detailed logging
  load_balancing_decision_log_verbosity = [
    "pool-selection",
    "pool-member-selection"
  ]
  
  topology_prefer_edns0_client_subnet = "enabled"
}

# Test 4: AAAA (IPv6) WideIP with aliases
resource "bigip_gtm_wideip" "ipv6_with_aliases" {
  name      = "ipv6app.example.com"
  type      = "aaaa"
  partition = "Common"
  
  description = "IPv6 WideIP with aliases"
  enabled     = true
  
  pool_lb_mode     = "round-robin"
  minimal_response = "enabled"
  
  persist_cidr_ipv6 = 64
  
  aliases = [
    "ipv6app-alt.example.com",
    "v6.app.example.com"
  ]
}

# Test 5: CNAME WideIP with aliases
resource "bigip_gtm_wideip" "cname_with_aliases" {
  name      = "alias.example.com"
  type      = "cname"
  partition = "Common"
  
  description = "CNAME WideIP with multiple aliases"
  enabled     = true
  
  aliases = [
    "alias1.example.com",
    "alias2.example.com"
  ]
}

# Outputs
output "api_match_id" {
  value = bigip_gtm_wideip.api_match.id
}

output "api_match_aliases" {
  value = bigip_gtm_wideip.api_match.aliases
}

output "multi_alias_aliases" {
  value = bigip_gtm_wideip.multi_alias.aliases
}

output "all_wideips" {
  value = {
    api_match      = bigip_gtm_wideip.api_match.name
    multi_alias    = bigip_gtm_wideip.multi_alias.name
    advanced_log   = bigip_gtm_wideip.advanced_logging.name
    ipv6           = bigip_gtm_wideip.ipv6_with_aliases.name
    cname          = bigip_gtm_wideip.cname_with_aliases.name
  }
}
