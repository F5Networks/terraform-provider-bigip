# bigip_gtm_monitor_https Resource

Provides a BIG-IP GTM (Global Traffic Manager) HTTPS Monitor resource. This resource allows you to configure and manage GTM HTTPS health monitors on a BIG-IP system.

## Description

A GTM HTTPS monitor verifies the availability and performance of HTTPS (SSL/TLS) services across your GTM infrastructure. The monitor establishes secure connections to target resources and evaluates the responses to determine health status. HTTPS monitors support client certificate authentication and configurable cipher lists, making them suitable for monitoring secure web services.

## Example Usage

### Basic HTTPS Monitor

```hcl
resource "bigip_gtm_monitor_https" "example" {
  name = "/Common/my_https_monitor"
}
```

### HTTPS Monitor with Client Certificate

```hcl
resource "bigip_gtm_monitor_https" "with_cert" {
  name          = "/Common/my_https_monitor"
  defaults_from = "/Common/https"
  destination   = "*:443"
  interval      = 10
  timeout       = 60
  probe_timeout = 3
  send          = "GET /health\\r\\n"
  receive       = "200 OK"
  cert          = "/Common/my_client_cert"
  key           = "/Common/my_client_key"
  compatibility = "enabled"
}
```

## Argument Reference

The following arguments are supported:

### Required Arguments

* `name` - (Required, String) The full path name of the GTM HTTPS monitor (e.g., `/Common/my_https_monitor`). Forces new resource.

### Optional Arguments

#### General Settings

* `defaults_from` - (Optional, String) Specifies the parent monitor from which this monitor inherits settings. Default: `/Common/https`.
* `destination` - (Optional, String) Specifies the IP address and service port of the resource being monitored. Format: `ip:port`. Default: `*:*`.
* `interval` - (Optional, Integer) Specifies, in seconds, the frequency at which the system issues the monitor check. Default: `30`.
* `timeout` - (Optional, Integer) Specifies the number of seconds the target has in which to respond to the monitor request. Default: `120`.
* `probe_timeout` - (Optional, Integer) Specifies the number of seconds after which the system times out the probe request. Default: `5`.
* `ignore_down_response` - (Optional, String) Specifies whether the monitor ignores a down response from the system it is monitoring. Valid values: `enabled`, `disabled`. Default: `disabled`.
* `transparent` - (Optional, String) Specifies whether the monitor operates in transparent mode. Valid values: `enabled`, `disabled`. Default: `disabled`.
* `reverse` - (Optional, String) Instructs the system to mark the target resource down when the test is successful. Valid values: `enabled`, `disabled`. Default: `disabled`.
* `send` - (Optional, String) Specifies the text string that the monitor sends to the target object. Default: `GET /\r\n`.
* `receive` - (Optional, String) Specifies the text string that the monitor looks for in the returned resource.

#### SSL Settings

* `cert` - (Optional, String) Specifies a fully-qualified path for a client certificate that the monitor sends to the target SSL server.
* `key` - (Optional, String) Specifies the key for the client certificate that the monitor sends to the target SSL server.
* `cipherlist` - (Optional, String) Specifies the list of ciphers for this monitor. Default: `DEFAULT:+SHA:+3DES:+kEDH`.
* `compatibility` - (Optional, String) Specifies the SSL version compatibility. Valid values: `enabled`, `disabled`. Default: `enabled`.

## Attribute Reference

In addition to the arguments listed above, the following attributes are exported:

* `id` - The full path name of the GTM HTTPS monitor.

## Import

GTM HTTPS Monitor resources can be imported using the full path name:

```bash
terraform import bigip_gtm_monitor_https.example /Common/my_https_monitor
```
