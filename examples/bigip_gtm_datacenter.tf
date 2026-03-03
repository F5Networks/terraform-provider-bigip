resource "bigip_gtm_datacenter" "example" {
    contact           = "test@f5.com"
    description       = "testdc"
    enabled           = true
    location          = "had"
    name              = "testdc"
    partition         = "Common"
    prober_fallback   = "any-available"
    prober_preference = "inside-datacenter"
}

# resource "bigip_gtm_datacenter" "example2" {
#   name      = "example_datacenter"
#   partition = "Common"

#   # Location and contact information
#   location    = "Seattle, WA"
#   contact     = "admin@example.com"
#   description = "Primary datacenter for west coast operations"

#   # Enable the datacenter
#   enabled = true

#   # Prober settings
#   prober_preference = "inside-datacenter"
#   prober_fallback   = "any-available"
# }

# # Example with custom prober settings
# resource "bigip_gtm_datacenter" "custom_prober" {
#   name      = "custom_datacenter"
#   partition = "Common"

#   location    = "New York, NY"
#   contact     = "ops@example.com"
#   description = "East coast datacenter with custom prober configuration"

#   enabled = true

#   # Use outside-datacenter preference instead of pool (which requires prober_pool)
#   prober_preference = "outside-datacenter"
#   prober_fallback   = "any-available"
# }

# # Example of a disabled datacenter
# resource "bigip_gtm_datacenter" "maintenance" {
#   name      = "maintenance_dc"
#   partition = "Common"

#   location    = "Chicago, IL"
#   contact     = "maintenance@example.com"
#   description = "Datacenter currently under maintenance"

#   # Disabled during maintenance
#   enabled = false

#   prober_preference = "inside-datacenter"
#   prober_fallback   = "any-available"
# }
