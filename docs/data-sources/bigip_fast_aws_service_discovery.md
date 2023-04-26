---
layout: "bigip"
page_title: "BIG-IP: bigip_fast_aws_service_discovery"
subcategory: "F5 Automation Tool Chain(ATC)"
description: |-
  Provides details for AWS Service discovery config
---

# bigip\_fast\_aws\_service\_discovery


Use this data source (`bigip_fast_aws_service_discovery`) to get the AWS Service discovery config to be used for `http`/`https` app deployment in FAST.

## Example Usage

```hcl
data "bigip_fast_aws_service_discovery" "TC2" {
  tag_key   = "testawstagkey"
  tag_value = "testawstagvalue"
}
```      

## Argument Reference

* `tag_key` - (`Required`,type `string`) The tag key associated with the node to add to this pool.
* `tag_value` - (`Required`,type `string`) The tag value associated with the node to add to this pool.
* `port` - (`optional`,type `int`)Port to be used for AWS service discovery,default `80`.
* `address_realm` - (`optional`,type `string`)Specifies whether to look for public or private IP addresses,default `private`.
* `undetectable_action` - (`optional`,type `string`)Action to take when node cannot be detected,default `remove`.
* `credential_update` - (`optional`,type `bool`) Specifies whether you are updating your credentials,default `false`.
* `aws_region` - (`optional`,type `string`) AWS region in which ADC is running,default Empty string.
* `aws_access_key` - (`optional`,type `string`)Information for discovering AWS nodes that are not in the same region as your BIG-IP (also requires the `aws_secret_access_key` field)
* `aws_secret_access_key` - (`optional`,type `string`)Information for discovering AWS nodes that are not in the same region as your BIG-IP (also requires the `aws_secret_access_key` field)
* `external_id` - (`optional`,type `string`)AWS externalID field.
* `role_arn` - (`optional`,type `string`) Assume a role (also requires the `external_id` field)
* `minimum_monitors` - (`optional`,type `string`)Member is down when fewer than minimum monitors report it healthy.
* `update_interval` - (`optional`,type `string`)Update interval for service discovery.

## Attributes Reference

* `aws_sd_json` - The JSON for AWS service discovery block.

