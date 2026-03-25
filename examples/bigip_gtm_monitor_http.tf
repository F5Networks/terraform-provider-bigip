# Basic HTTP monitor with defaults
resource "bigip_gtm_monitor_http" "basic" {
  name = "/Common/my_http_monitor"
}

# HTTP monitor with custom settings
resource "bigip_gtm_monitor_http" "custom" {
  name                 = "/Common/my_custom_http_monitor"
  defaults_from        = "/Common/http"
  destination          = "*:80"
  interval             = 10
  timeout              = 60
  probe_timeout        = 3
  ignore_down_response = "disabled"
  transparent          = "disabled"
  reverse              = "disabled"
  send                 = "GET /health\\r\\n"
  receive              = "200 OK"
}
