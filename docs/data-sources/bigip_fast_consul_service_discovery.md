---
layout: "bigip"
page_title: "BIG-IP: bigip_fast_consul_service_discovery"
subcategory: "F5 Automation Tool Chain(ATC)"
description: |-
  Provides details for Hashicorp Consul Service discovery config
---

# bigip\_fast\_consul\_service\_discovery


Use this data source (`bigip_fast_consul_service_discovery`) to get the Consul Service discovery config to be used for `http`/`https` app deployment in FAST.

## Example Usage

```hcl
data "bigip_fast_consul_service_discovery" "TC2" {
  uri  = "https://192.0.2.100:8500/v1/catalog/nodes"
  port = 8080
}

```      

## Argument Reference

* `uri` - (`Required`,type `string`) The location of the node data.
* `port` - (`optional`,type `int`)Port to be used for AWS service discovery,default `80`.
* `address_realm` - (`optional`,type `string`)Specifies whether to look for public or private IP addresses,default `private`.
* `undetectable_action` - (`optional`,type `string`)Action to take when node cannot be detected,default `remove`.
* `credential_update` - (`optional`,type `bool`) Specifies whether you are updating your credentials,default `false`.
* `encoded_token` - (`optional`,type `string`) Base 64 encoded bearer token to make requests to the Consul API. Will be stored in the declaration in an encrypted format.
* `jmes_path_query` - (`optional`,type `string`)Custom JMESPath Query.
* `reject_unauthorized` - (`optional`,type `bool`)If true, the server certificate is verified against the list of supplied/default CAs when making requests to the Consul API.
* `trust_ca` - (`optional`,type `string`)CA Bundle to validate server certificates.
* `minimum_monitors` - (`optional`,type `string`)Member is down when fewer than minimum monitors report it healthy.
* `update_interval` - (`optional`,type `string`)Update interval for service discovery.

## Attributes Reference

* `consul_sd_json` - The JSON for Hashicorp Consul service discovery block.

