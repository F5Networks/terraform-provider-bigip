# Basic HTTPS monitor with defaults
resource "bigip_gtm_monitor_https" "basic" {
  name = "/Common/my_https_monitor"
}

# HTTPS monitor with client certificate and custom settings
resource "bigip_gtm_monitor_https" "with_cert" {
  name                 = "/Common/my_secure_https_monitor"
  defaults_from        = "/Common/https"
  destination          = "*:443"
  interval             = 10
  timeout              = 60
  probe_timeout        = 3
  ignore_down_response = "disabled"
  transparent          = "disabled"
  reverse              = "disabled"
  send                 = "GET /health\\r\\n"
  receive              = "200 OK"
  cert                 = "/Common/my_client_cert"
  key                  = "/Common/my_client_key"
  compatibility        = "enabled"
}
