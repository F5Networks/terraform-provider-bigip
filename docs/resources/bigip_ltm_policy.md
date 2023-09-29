---
layout: "bigip"
page_title: "BIG-IP: bigip_ltm_policy"
subcategory: "Local Traffic Manager(LTM)"
description: |-
  Provides details about bigip_ltm_policy resource
---

# bigip\_ltm\_policy

`bigip_ltm_policy` Configures ltm policies to manage traffic assigned to a virtual server

For resources should be named with their `full path`. The full path is the combination of the `partition + name` of the resource. For example `/Common/test-policy`.

## Example Usage

```hcl

resource "bigip_ltm_pool" "mypool" {
  name                = "/Common/test-pool"
  allow_nat           = "yes"
  allow_snat          = "yes"
  load_balancing_mode = "round-robin"
}
resource "bigip_ltm_policy" "test-policy" {
  name     = "/Common/test-policy"
  strategy = "first-match"
  requires = ["http"]
  controls = ["forwarding"]
  rule {
    name = "rule6"
    action {
      forward    = true
      connection = false
      pool       = bigip_ltm_pool.mypool.name
    }
  }
  depends_on = [bigip_ltm_pool.mypool]
}
```

## Argument Reference

> [!NOTE]
> The attribute `published_copy` is not required anymore as the resource automatically publishes the policy, hence it's deprecated and will be removed from future release.

* `name`- (Required) Name of the Policy ( policy name should be in full path which is combination of partition and policy name )

* `strategy` - (Optional) Specifies the match strategy

* `description` - (Optional) Specifies descriptive text that identifies the ltm policy.

* `requires` - (Optional) Specifies the protocol

* `published_copy` - (Deprecated) If you want to publish the policy else it will be deployed in Drafts mode. This attribute is deprecated and will be removed in a future release.

*  `controls` - (Optional) Specifies the controls

* `rule` - (Optional,type `list`) List of Rules can be applied using the policy. Each rule is block type with following arguments.
    * `name` -  (Required,type `string`) Name of Rule to be applied in policy.
    * `description` - (Optional) Specifies descriptive text that identifies the irule attached to policy.
    * `condition` - (Optional,type `set`) Block type. See [condition](#condition) block for more details.
    * `action` - (Optional,type `set`) Block type. See [action](#action) block for more details.

* `forward` - (Optional) This action will affect forwarding.

* `pool` - (Optional ) This action will direct the stream to this pool.

* `connection` - (Optional) This action is set to `true` by default, it needs to be explicitly set to `false` for actions it conflicts with.

### condition
* `condition` is block type which defines below attributes, those need to adjusted/provided based on condition match.

    * `http_host`
    * `external`
    * `equals`
    * `address`
    * `all`
    * `app_service`
    * `browser_type`
    * `browser_version`
    * `case_insensitive`
    * `case_sensitive`
    * `cipher`
    * `cipher_bits`
    * `client_ssl`
    * `code`
    * `common_name`
    * `contains`
    * `continent`
    * `country_code`
    * `country_name`
    * `cpu_usage`
    * `device_make`
    * `device_model`
    * `domain`
    * `ends_with`
    * `expiry`
    * `extension`
    * `geoip`
    * `greater`
    * `greater_or_equal`
    * `host`
    * `http_basic_auth`
    * `http_cookie`
    * `http_header`
    * `http_method`
    * `http_referer`
    * `http_set_cookie`
    * `http_status`
    * `http_uri`
    * `http_user_agent`
    * `http_version`
    * `internal`
    * `isp`
    * `last_15secs`
    * `last_1min`
    * `last_5mins`
    * `less`
    * `less_or_equal`
    * `local`
    * `major`
    * `matches`
    * `minor`
    * `missing`
    * `mss`
    * `tm_name`
    * `not`
    * `exists`
    * `org`
    * `password`
    * `path`
    * `path_segment`
    * `port`
    * `present`
    * `protocol`
    * `query_parameter`
    * `query_string`
    * `region_code`
    * `region_name`
    * `remote`
    * `request`
    * `client_accepted`
    * `response`
    * `route_domain`
    * `rtt`
    * `scheme`
    * `server_name`
    * `ssl_cert`
    * `ssl_client_hello`
    * `ssl_extension`
    * `ssl_server_handshake`
    * `ssl_server_hello`
    * `starts_with`
    * `tcp`
    * `text`
    * `unnamed_query_parameter`
    * `user_agent_token`
    * `username`
    * `value`
    * `values`
    * `version`
    * `vlan`
    * `vlan_id`

### action

* `action` is block type which defines below attributes, those need to adjusted/provided based on condition match.

    * `app_service`
    * `application`
    * `asm`
    * `avr`
    * `cache`
    * `carp`
    * `category`
    * `classify`
    * `clone_pool`
    * `code`
    * `compress`
    * `content`
    * `connection`
    * `cookie_hash`
    * `cookie_insert`
    * `cookie_passive`
    * `cookie_rewrite`
    * `decompress`
    * `defer`
    * `destination_address`
    * `disable`
    * `domain`
    * `enable`
    * `expiry`
    * `expiry_secs`
    * `expression`
    * `extension`
    * `facility`
    * `forward`
    * `from_profile`
    * `hash`
    * `host`
    * `http`
    * `http_basic_auth`
    * `http_cookie`
    * `http_header`
    * `http_host`
    * `http_referer`
    * `http_reply`
    * `http_set_cookie`
    * `http_uri`
    * `ifile`
    * `insert`
    * `internal_virtual`
    * `ip_address`
    * `key`
    * `l7dos`
    * `length`
    * `location`
    * `log`
    * `ltm_policy`
    * `member`
    * `message`
    * `tm_name`
    * `netmask`
    * `nexthop`
    * `node`
    * `offset`
    * `path`
    * `pem`
    * `persist`
    * `pin`
    * `policy`
    * `pool`
    * `port`
    * `priority`
    * `profile`
    * `protocol`
    * `query_string`
    * `rateclass`
    * `redirect`
    * `remove`
    * `replace`
    * `request`
    * `request_adapt`
    * `reset`
    * `response`
    * `response_adapt`
    * `scheme`
    * `script`
    * `select`
    * `server_ssl`
    * `set_variable`
    * `shutdown`
    * `snat`
    * `snatpool`
    * `source_address`
    * `ssl_client_hello`
    * `ssl_server_handshake`
    * `ssl_server_hello`
    * `ssl_session_id`
    * `status`
    * `tcl`
    * `tcp_nagle`
    * `text`
    * `timeout`
    * `uie`
    * `universal`
    * `value`
    * `virtual`
    * `vlan`
    * `vlan_id`
    * `wam`
    * `write`

## Importing
An existing policy can be imported into this resource by supplying policy Name in `full path` as `id`.
An example is below:
```sh
$ terraform import bigip_ltm_policy.policy-import-test /Common/policy2
```
