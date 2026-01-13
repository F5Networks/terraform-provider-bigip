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

# # Example 1: Basic GTM Pool with minimal configuration
# resource "bigip_gtm_pool" "basic_pool" {
#   name      = "basic_pool"
#   type      = "a"
#   partition = "Common"
  
#   load_balancing_mode = "round-robin"
#   monitor             = "/Common/https"
# }

# # Example 2: GTM Pool with members
# resource "bigip_gtm_pool" "firstpool" {
#     alternate_mode               = "round-robin"
#     disabled                     = false
#     dynamic_ratio                = "disabled"
#     enabled                      = true
#     fallback_ip                  = "any"
#     fallback_mode                = "return-to-dns"
#     limit_max_bps                = 0
#     limit_max_bps_status         = "disabled"
#     limit_max_connections        = 0
#     limit_max_connections_status = "disabled"
#     limit_max_pps                = 0
#     limit_max_pps_status         = "disabled"
#     load_balancing_mode          = "round-robin"
#     manual_resume                = "disabled"
#     max_answers_returned         = 1
#     min_members_up_mode          = "off"
#     min_members_up_value         = 0
#     monitor                      = "/Common/gateway_icmp and /Common/http"
#     name                         = "firstpool"
#     partition                    = "Common"
#     qos_hit_ratio                = 5
#     qos_hops                     = 0
#     qos_kilobytes_second         = 3
#     qos_lcs                      = 30
#     qos_packet_rate              = 1
#     qos_rtt                      = 50
#     qos_topology                 = 0
#     qos_vs_capacity              = 0
#     qos_vs_score                 = 0
#     ttl                          = 30
#     type                         = "a"
#     verify_member_availability   = "enabled"
#     members {
#         name    = "/Common/ecosyshydbigip16.com:/Common/testwebpolicy.app/testwebpolicy_vs"
#         member_order = 0
#     }
#     members {
#         name      = "/Common/ecosyshydbigip16.com:/Common/testravivs"
#         member_order = 1
#     }
# }


# # Example 2: GTM Pool matching the API response
# resource "bigip_gtm_pool" "firstpool2" {
#   name      = "firstpool_tf"
#   type      = "a"
#   partition = "Common"
  
#   # Load balancing configuration
#   alternate_mode        = "round-robin"
#   load_balancing_mode   = "round-robin"
#   fallback_mode         = "return-to-dns"
#   fallback_ip           = "any"
  
#   # State
#   enabled               = true
#   dynamic_ratio         = "disabled"
#   manual_resume         = "disabled"
  
#   # Response configuration
#   max_answers_returned  = 1
#   ttl                   = 30
  
#   # Monitoring
#   monitor               = "/Common/https"
#   verify_member_availability = "enabled"
  
#   # QoS settings
#   qos_hit_ratio         = 5
#   qos_hops              = 0
#   qos_kilobytes_second  = 3
#   qos_lcs               = 30
#   qos_packet_rate       = 1
#   qos_rtt               = 50
#   qos_topology          = 0
#   qos_vs_capacity       = 0
#   qos_vs_score          = 0
  
#   # Limits
#   limit_max_bps                  = 0
#   limit_max_bps_status           = "disabled"
#   limit_max_connections          = 0
#   limit_max_connections_status   = "disabled"
#   limit_max_pps                  = 0
#   limit_max_pps_status           = "disabled"
  
#   # Min members
#   min_members_up_mode  = "off"
#   min_members_up_value = 0
  
#   # Members
#   # Note: Members require valid GTM server:virtual_server references
#   # Example: members { name = "server_name:/Common/vs_name" }
#   # Uncomment and update when GTM servers and virtual servers are created
#   # members {
#   #   name                          = "server1:/Common/vs1"
#   #   enabled                       = true
#   #   ratio                         = 1
#   #   member_order                  = 0
#   #   monitor                       = "default"
#   #   limit_max_bps                 = 0
#   #   limit_max_bps_status          = "disabled"
#   #   limit_max_connections         = 0
#   #   limit_max_connections_status  = "disabled"
#   #   limit_max_pps                 = 0
#   #   limit_max_pps_status          = "disabled"
#   # }
# }

# # Example 3: GTM Pool with multiple members
# # Note: This example requires GTM servers (server1, server2, server3) and their virtual servers (vs1, vs2, vs3)
# # Commented out until servers are created
# /*
# resource "bigip_gtm_pool" "multi_member_pool" {
#   name      = "multi_member_pool"
#   type      = "a"
#   partition = "Common"
  
#   load_balancing_mode = "round-robin"
#   monitor             = "/Common/gateway_icmp"
#   ttl                 = 60
  
#   members {
#     name         = "server1:/Common/vs1"
#     enabled      = true
#     ratio        = 2
#     member_order = 0
#   }
  
#   members {
#     name         = "server2:/Common/vs2"
#     enabled      = true
#     ratio        = 1
#     member_order = 1
#   }
  
#   members {
#     name         = "server3:/Common/vs3"
#     enabled      = true
#     ratio        = 1
#     member_order = 2
#   }
# }
# */

# # Example 4: GTM Pool with ratio load balancing
# # Note: This example requires GTM servers with virtual servers
# # Commented out until servers are created
# /*
# resource "bigip_gtm_pool" "ratio_pool" {
#   name      = "ratio_pool"
#   type      = "a"
#   partition = "Common"
  
#   load_balancing_mode = "ratio"
#   alternate_mode      = "round-robin"
#   fallback_mode       = "return-to-dns"
  
#   members {
#     name    = "server1:/Common/vs1"
#     enabled = true
#     ratio   = 3  # 3x more traffic
#   }
  
#   members {
#     name    = "server2:/Common/vs2"
#     enabled = true
#     ratio   = 1
#   }
# }
# */

# # Example 5: AAAA (IPv6) Pool
# # Note: This example requires a GTM server with an IPv6 virtual server
# # Commented out until server is created
# /*
# resource "bigip_gtm_pool" "ipv6_pool" {
#   name      = "ipv6_pool"
#   type      = "aaaa"
#   partition = "Common"
  
#   load_balancing_mode = "round-robin"
#   monitor             = "/Common/gateway_icmp"
  
#   members {
#     name    = "server1:/Common/vs_ipv6"
#     enabled = true
#   }
# }
# */

# # Example 6: Pool with connection limits
# # Note: This example requires a GTM server with virtual server
# # Commented out until server is created
# /*
# resource "bigip_gtm_pool" "limited_pool" {
#   name      = "limited_pool"
#   type      = "a"
#   partition = "Common"
  
#   load_balancing_mode          = "round-robin"
#   limit_max_connections        = 1000
#   limit_max_connections_status = "enabled"
#   limit_max_bps                = 10000000
#   limit_max_bps_status         = "enabled"
  
#   members {
#     name                         = "server1:/Common/vs1"
#     enabled                      = true
#     limit_max_connections        = 500
#     limit_max_connections_status = "enabled"
#   }
# }
# */

# # Example 7: Pool with minimum members requirement
# # Note: Valid values for min_members_up_mode are: "at-least", "percentage", "off"
# # Note: This example requires GTM servers with virtual servers
# # Commented out until servers are created
# /*
# resource "bigip_gtm_pool" "min_members_pool" {
#   name      = "min_members_pool"
#   type      = "a"
#   partition = "Common"
  
#   load_balancing_mode  = "round-robin"
#   # Fix: The mode should likely be a different value or requires specific BIG-IP version
#   # Common values: "off", "enabled" or numeric
#   min_members_up_mode  = "off"  # Changed from "at-least"
#   min_members_up_value = 2
  
#   members {
#     name    = "server1:/Common/vs1"
#     enabled = true
#   }
  
#   members {
#     name    = "server2:/Common/vs2"
#     enabled = true
#   }
  
#   members {
#     name    = "server3:/Common/vs3"
#     enabled = true
#   }
# }
# */

# # Outputs
# output "firstpool_id" {
#   description = "The ID of the first pool"
#   value       = bigip_gtm_pool.firstpool.id
# }

# output "firstpool_members" {
#   description = "The members of the first pool"
#   value       = bigip_gtm_pool.firstpool.members
# }
