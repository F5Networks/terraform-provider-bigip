# bigip_gtm_monitor_bigip Resource

Provides a BIG-IP GTM (Global Traffic Manager) BIG-IP Monitor resource. This resource allows you to configure and manage GTM BIG-IP health monitors on a BIG-IP system.

## Description

A GTM BIG-IP monitor is designed to monitor BIG-IP systems themselves within your GTM infrastructure. Unlike protocol-specific monitors (HTTP, TCP, etc.), the BIG-IP monitor collects health and performance data directly from BIG-IP devices, including metrics such as CPU usage, memory utilization, and throughput. This monitor type does not use `send`/`receive` strings or `probe_timeout` and instead focuses on aggregated system-level metrics.

## Example Usage

### Basic BIG-IP Monitor

```hcl
resource "bigip_gtm_monitor_bigip" "example" {
  name = "/Common/my_bigip_monitor"
}
```

### Full BIG-IP Monitor Configuration

```hcl
resource "bigip_gtm_monitor_bigip" "advanced" {
  name                 = "/Common/my_bigip_monitor"
  defaults_from        = "/Common/bigip"
  destination          = "*:*"
  interval             = 10
  timeout              = 30
  ignore_down_response = "disabled"
  aggregation_type     = "average-nodes"
}
```

## Argument Reference

The following arguments are supported:

### Required Arguments

* `name` - (Required, String) The full path name of the GTM BIG-IP monitor (e.g., `/Common/my_bigip_monitor`). Forces new resource.

### Optional Arguments

* `defaults_from` - (Optional, String) Specifies the parent monitor from which this monitor inherits settings. Default: `/Common/bigip`.
* `destination` - (Optional, String) Specifies the IP address and service port of the resource being monitored. Format: `ip:port`. Default: `*:*`. Note: when a specific IP address is provided, a specific port is also required (BIG-IP API constraint).
* `interval` - (Optional, Integer) Specifies, in seconds, the frequency at which the system issues the monitor check. Default: `30`.
* `timeout` - (Optional, Integer) Specifies the number of seconds the target has in which to respond to the monitor request. Default: `90`.
* `ignore_down_response` - (Optional, String) Specifies whether the monitor ignores a down response from the system it is monitoring. Valid values: `enabled`, `disabled`. Default: `disabled`.
* `aggregation_type` - (Optional, String) Specifies how the system combines the monitor information it collects for a pool of monitored resources. Default: `none`.

## Attribute Reference

In addition to the arguments listed above, the following attributes are exported:

* `id` - The full path name of the GTM BIG-IP monitor.

## Import

GTM BIG-IP Monitor resources can be imported using the full path name:

```bash
terraform import bigip_gtm_monitor_bigip.example /Common/my_bigip_monitor
```
