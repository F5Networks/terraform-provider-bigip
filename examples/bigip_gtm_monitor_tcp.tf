# Basic TCP monitor with defaults
resource "bigip_gtm_monitor_tcp" "basic" {
  name = "/Common/my_tcp_monitor"
}

# TCP monitor with send/receive strings
resource "bigip_gtm_monitor_tcp" "custom" {
  name                 = "/Common/my_custom_tcp_monitor"
  defaults_from        = "/Common/tcp"
  destination          = "*:3306"
  interval             = 10
  timeout              = 60
  probe_timeout        = 3
  ignore_down_response = "disabled"
  transparent          = "disabled"
  reverse              = "disabled"
  send                 = "ping"
  receive              = "pong"
}
