---
layout: "bigip"
page_title: "BIG-IP: bigip_ltm_request_log_profile"
subcategory: "Local Traffic Manager(LTM)"
description: |-
  Provides details about bigip_ltm_request_log_profile resource
---

# bigip\_ltm\_request\_log\_profile

`bigip_ltm_request_log_profile` Resource used for Configures request logging using the Request Logging profile

## Example Usage

```hcl
resource "bigip_ltm_request_log_profile" "request-log-profile-tc1-child" {
  name                       = "/Common/request-log-profile-tc1-child"
  defaults_from              = bigip_ltm_request_log_profile.request-log-profile-tc1.name
  request_logging            = "disabled"
  requestlog_pool            = "/Common/pool2"
  requestlog_error_pool      = "/Common/pool1"
  requestlog_protocol        = "mds-tcp"
  requestlog_error_protocol  = "mds-tcp"
  responselog_protocol       = "mds-tcp"
  responselog_error_protocol = "mds-tcp"
}

```      

## Argument Reference

* `name` (Required,type `string`) Name of the Request Logging profile,name of Profile should be full path. Full path is the combination of the `partition + profile name`,For example `/Common/request-log-profile-tc1`.

* `defaults_from` - (optional,type `string`) Specifies the profile from which this profile inherits settings. The default is the system-supplied `request-log` profile.

* `description` - (optional,type `string`) Specifies user-defined description.

* `request_logging` - (Optional,type `string`) Enables or disables request logging. The default is `disabled`, possible values are `enabled` and `disabled`.

* `requestlog_protocol` - (Optional) Specifies the protocol to be used for high-speed logging of requests. The default is `mds-udp`,possible values are `mds-udp` and `mds-tcp`.

* `requestlog_template` - (Optional) Specifies the directives and entries to be logged. More infor on requestlog_template can be found [here](https://techdocs.f5.com/en-us/bigip-15-0-0/external-monitoring-of-big-ip-systems-implementations/configuring-request-logging.html). how to use can be find [here](https://my.f5.com/manage/s/article/K00847516).

* `requestlog_error_template` - (Optional) Specifies the directives and entries to be logged for request errors.

* `requestlog_pool` - (Optional) Defines the pool to send logs to. Typically, the pool will contain one or more syslog servers. It is recommended that you create a pool specifically for logging requests. The default is `none`.

* `requestlog_error_protocol` - (Optional) Specifies the protocol to be used for high-speed logging of request errors. The default is `mds-udp`,possible values are `mds-udp` and `mds-tcp`.

* `response_logging` - (Optional,type `string`) Enables or disables response logging. The default is `disabled`, possible values are `enabled` and `disabled`.

* `responselog_protocol` - (Optional) Specifies the protocol to be used for high-speed logging of responses. The default is `mds-udp`,possible values are `mds-udp` and `mds-tcp`.

* `responselog_error_protocol` - (Optional) Specifies the protocol to be used for high-speed logging of response errors. The default is `mds-udp`,possible values are `mds-udp` and `mds-tcp`.

* `responselog_pool` - (Optional) Defines the pool to send logs to. Typically, the pool contains one or more syslog servers. It is recommended that you create a pool specifically for logging responses. The default is `none`.

* `responselog_error_pool` - (Optional) Defines the pool associated with logging response errors. The default is `none`.

* `responselog_template` - (Optional) Specifies the directives and entries to be logged. More infor on responselog_template can be found [here](https://techdocs.f5.com/en-us/bigip-15-0-0/external-monitoring-of-big-ip-systems-implementations/configuring-request-logging.html). how to use can be find [here](https://my.f5.com/manage/s/article/K00847516).

* `responselog_error_template` - (Optional) Specifies the directives and entries to be logged for request errors.

## Import

BIG-IP LTM Request Log profiles can be imported using the `name`, e.g.

```bash
terraform import bigip_ltm_request_log_profile.test-request-log /Common/test-request-log
```
