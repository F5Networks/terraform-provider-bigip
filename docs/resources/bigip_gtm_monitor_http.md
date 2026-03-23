# bigip_gtm_monitor_http Resource

Provides a BIG-IP GTM (Global Traffic Manager) HTTP Monitor resource. This resource allows you to configure and manage GTM HTTP health monitors on a BIG-IP system.

## Description

A GTM HTTP monitor verifies the availability and performance of HTTP services across your GTM infrastructure. The monitor sends HTTP requests to target resources and evaluates the responses to determine health status. HTTP monitors are commonly used to check web server availability and verify that expected content is being served.

## Example Usage

### Basic HTTP Monitor

```hcl
resource "bigip_gtm_monitor_http" "example" {
  name = "/Common/my_http_monitor"
}
```

### Full HTTP Monitor Configuration

```hcl
resource "bigip_gtm_monitor_http" "advanced" {
  name                 = "/Common/my_http_monitor"
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
```

## Argument Reference

The following arguments are supported:

### Required Arguments

* `name` - (Required, String) The full path name of the GTM HTTP monitor (e.g., `/Common/my_http_monitor`). Forces new resource.

### Optional Arguments

* `defaults_from` - (Optional, String) Specifies the parent monitor from which this monitor inherits settings. Default: `/Common/http`.
* `destination` - (Optional, String) Specifies the IP address and service port of the resource being monitored. Format: `ip:port`. Default: `*:*`.
* `interval` - (Optional, Integer) Specifies, in seconds, the frequency at which the system issues the monitor check. Default: `30`.
* `timeout` - (Optional, Integer) Specifies the number of seconds the target has in which to respond to the monitor request. Default: `120`.
* `probe_timeout` - (Optional, Integer) Specifies the number of seconds after which the system times out the probe request. Default: `5`.
* `ignore_down_response` - (Optional, String) Specifies whether the monitor ignores a down response from the system it is monitoring. Valid values: `enabled`, `disabled`. Default: `disabled`.
* `transparent` - (Optional, String) Specifies whether the monitor operates in transparent mode. Valid values: `enabled`, `disabled`. Default: `disabled`.
* `reverse` - (Optional, String) Instructs the system to mark the target resource down when the test is successful. Valid values: `enabled`, `disabled`. Default: `disabled`.
* `send` - (Optional, String) Specifies the text string that the monitor sends to the target object. Default: `GET /`.
* `receive` - (Optional, String) Specifies the text string that the monitor looks for in the returned resource.

## Attribute Reference

In addition to the arguments listed above, the following attributes are exported:

* `id` - The full path name of the GTM HTTP monitor.

## Import

GTM HTTP Monitor resources can be imported using the full path name:

```bash
terraform import bigip_gtm_monitor_http.example /Common/my_http_monitor
```
