# Example based on the API response provided
# This creates a WideIP that matches the configuration shown in the API

resource "bigip_gtm_wideip" "testwideip" {
  name      = "testwideip.local"
  type      = "a"
  partition = "Common"
  
  description = "test_wideip_a"
  
  # State
  enabled  = true
  disabled = false
  
  # Failure handling
  failure_rcode          = "noerror"
  failure_rcode_response = "disabled"
  failure_rcode_ttl      = 0
  
  # Pool configuration
  last_resort_pool = "a /Common/firstpool"
  pool_lb_mode     = "round-robin"
  
  # Response configuration
  minimal_response = "disabled"
  
  # Persistence
  persistence       = "disabled"
  persist_cidr_ipv4 = 32
  persist_cidr_ipv6 = 128
  ttl_persistence   = 3600
  
  # Topology
  topology_prefer_edns0_client_subnet = "enabled"
  
  # Logging
  load_balancing_decision_log_verbosity = ["pool-selection"]
  
  # Aliases
  aliases = ["testwideip2.local"]
}

# Output the WideIP details
output "wideip_full_path" {
  description = "The full path of the WideIP"
  value       = "/Common/${bigip_gtm_wideip.testwideip.name}"
}

output "wideip_id" {
  description = "The Terraform ID of the WideIP"
  value       = bigip_gtm_wideip.testwideip.id
}
