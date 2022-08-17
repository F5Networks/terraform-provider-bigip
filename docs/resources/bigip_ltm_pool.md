---
layout: "bigip"
page_title: "BIG-IP: bigip_ltm_pool"
subcategory: "Local Traffic Manager(LTM)"
description: |-
  Provides details about bigip_ltm_pool resource
---

# bigip\_ltm\_pool

`bigip_ltm_pool` Manages F5 BIG-IP LTM pools via iControl REST API.

For resources should be named with their `full path`. The full path is the combination of the `partition + name` of the resource or  `partition + directory + name`.
For example `/Common/my-pool`.

## Example Usage

```hcl
resource "bigip_ltm_monitor" "monitor" {
  name   = "/Common/terraform_monitor"
  parent = "/Common/http"
}
resource "bigip_ltm_pool" "pool" {
  name                   = "/Common/Axiom_Environment_APP1_Pool"
  load_balancing_mode    = "round-robin"
  minimum_active_members = 1
  monitors               = [bigip_ltm_monitor.monitor.name]
}
```      

## Argument Reference

* `name` - (Required,type `string`) Name of the pool,it should be `full path`.The full path is the combination of the `partition + name` of the pool.(For example `/Common/my-pool`)

* `monitors` - (Optional,type `list`) List of monitor names to associate with the pool

* `description` - (Optional,type `string`) Specifies descriptive text that identifies the pool. 

* `allow_nat` - (Optional,type `string`) Specifies whether NATs are automatically enabled or disabled for any connections using this pool, [ Default : `yes`, Possible Values `yes` or `no`].

* `allow_snat` - (Optional,type `string`) Specifies whether SNATs are automatically enabled or disabled for any connections using this pool,[ Default : `yes`, Possible Values `yes` or `no`].

* `load_balancing_mode` - (Optional, type `string`) Specifies the load balancing method. The default is Round Robin.

* `minimum_active_members` - (Optional, type `int`) Specifies whether the system load balances traffic according to the priority number assigned to the pool member,Default Value is `0` meaning `disabled`.

* `slow_ramp_time` - (Optional, type `int`) Specifies the duration during which the system sends less traffic to a newly-enabled pool member.

* `service_down_action` - (Optional, type `string`) Specifies how the system should respond when the target pool member becomes unavailable. The default is `None`, Possible values: `[none, reset, reselect, drop]`.

* `reselect_tries` - (Optional, type `int`) Specifies the number of times the system tries to contact a new pool member after a passive failure.

## Importing
An existing pool can be imported into this resource by supplying pool Name in `full path` as `id`.
An example is below:
```sh
$ terraform import bigip_ltm_pool.k8s_prod_import /Common/k8prod_Pool

```