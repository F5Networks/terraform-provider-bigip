# Basic BIG-IP monitor with defaults
resource "bigip_gtm_monitor_bigip" "basic" {
  name = "/Common/my_bigip_monitor"
}

# BIG-IP monitor with custom settings
resource "bigip_gtm_monitor_bigip" "custom" {
  name                 = "/Common/my_custom_bigip_monitor"
  defaults_from        = "/Common/bigip"
  destination          = "*:*"
  interval             = 10
  timeout              = 30
  ignore_down_response = "disabled"
  aggregation_type     = "average-nodes"
}
