---
layout: "bigip"
page_title: "BIG-IP: bigip_fast_azure_service_discovery"
subcategory: "F5 Automation Tool Chain(ATC)"
description: |-
  Provides details for Azure Service discovery config
---

# bigip\_fast\_azure\_service\_discovery


Use this data source (`bigip_fast_azure_service_discovery`) to get the Azure Service discovery config to be used for `http`/`https` app deployment in FAST.

## Example Usage

```hcl
data "bigip_fast_azure_service_discovery" "TC3" {
  resource_group  = "testazurerg"
  subscription_id = "testazuresid"
  tag_key         = "testazuretag"
  tag_value       = "testazurevalue"
}
```      

## Argument Reference

* `resource_group` - (`Required`,type `string`) Azure Resource Group name.
* `subscription_id` - (`Required`,type `string`) Azure subscription ID.
* `port` - (`optional`,type `int`)Port to be used for Azure service discovery,default `80`.
* `tag_key` - (`Required`,type `string`) The tag key associated with the node to add to this pool.
* `tag_value` - (`Required`,type `string`) The tag value associated with the node to add to this pool.
* `address_realm` - (`optional`,type `string`)Specifies whether to look for public or private IP addresses,default `private`.
* `undetectable_action` - (`optional`,type `string`)Action to take when node cannot be detected,default `remove`.
* `credential_update` - (`optional`,type `bool`) Specifies whether you are updating your credentials,default `false`.
* `minimum_monitors` - (`optional`,type `string`)Member is down when fewer than minimum monitors report it healthy.
* `update_interval` - (`optional`,type `string`)Update interval for service discovery.

## Attributes Reference

* `azure_sd_json` - The JSON for Azure service discovery block.

