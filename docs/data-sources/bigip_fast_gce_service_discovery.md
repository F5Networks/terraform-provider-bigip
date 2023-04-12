---
layout: "bigip"
page_title: "BIG-IP: bigip_fast_gce_service_discovery"
subcategory: "F5 Automation Tool Chain(ATC)"
description: |-
  Provides details for GCE Service discovery config
---

# bigip\_fast\_gce\_service\_discovery

Use this data source (`bigip_fast_gce_service_discovery`) to get the GCE Service discovery config to be used for `http`/`https` app deployment in FAST.

## Example Usage

```hcl
data "bigip_fast_gce_service_discovery" "TC3" {
  tag_key   = "testgcetag"
  tag_value = "testgcevalue"
  region    = "testgceregion"
}

```      

## Argument Reference

* `tag_key` - (`Required`,type `string`) The tag key associated with the node to add to this pool.
* `tag_value` - (`Required`,type `string`) The tag value associated with the node to add to this pool.
* `region` - (`Required`,type `string`) GCE region in which ADC is running.
* `port` - (`optional`,type `int`)Port to be used for AWS service discovery,default `80`.
* `address_realm` - (`optional`,type `string`)Specifies whether to look for public or private IP addresses,default `private`.
* `undetectable_action` - (`optional`,type `string`)Action to take when node cannot be detected,default `remove`.
* `credential_update` - (`optional`,type `bool`) Specifies whether you are updating your credentials,default `false`.
* `encoded_credentials` - (`optional`,type `string`)Base 64 encoded service account credentials JSON.
* `project_id` - (`optional`,type `string`)For Google Cloud Engine (GCE) only: The ID of the project in which the members are located.
* `minimum_monitors` - (`optional`,type `string`)Member is down when fewer than minimum monitors report it healthy.
* `update_interval` - (`optional`,type `string`)Update interval for service discovery.

## Attributes Reference

* `gce_sd_json` - The JSON for GCE service discovery block.

