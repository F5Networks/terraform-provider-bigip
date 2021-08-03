---
layout: "bigip"
page_title: "BIG-IP: bigip_traffic_selector"
subcategory: "Network"
description: |-
  Provides details about bigip_traffic_selector resource
---

# bigip_traffic_selector

`bigip_traffic_selector` Manage IPSec Traffic Selectors on BIG-IP

Resources should be named with their "full path". The full path is the combination of the partition + name (example: /Common/test-selector)


## Example Usage

```hcl
resource bigip_traffic_selector "test-selector" {
  name                = "/Common/test-selector"
  destination_address = "3.10.11.2/32"
  source_address      = "2.10.11.12/32"
}
```      

## Argument Reference

* `name` - (Required) Name of the IPSec traffic-selector,it should be "full path".The full path is the combination of the partition + name of the IPSec traffic-selector.(For example `/Common/test-selector`)

* `description` - (Optional,type `string`) Description of the traffic selector.

* `destination_address` - (Optional,type `string`) Specifies the host or network IP address to which the application traffic is destined.When creating a new traffic selector, this parameter is required. 

* `destination_port` - (Optional,type `int`) Specifies the IP port used by the application. The default value is `All Ports (0)`

* `source_address` - (Optional,type `string`) Specifies the host or network IP address from which the application traffic originates.When creating a new traffic selector, this parameter is required.

* `source_port` - (Optional, type `int`) Specifies the IP port used by the application. The default value is `All Ports (0)`.

* `direction` - (Optional, type `string`) Specifies whether the traffic selector applies to inbound or outbound traffic, or both. The default value is `Both`.

* `ipsec_policy` - (Optional, type `string`) Specifies the IPsec policy that tells the BIG-IP system how to handle the packets.When creating a new traffic selector, if this parameter is not specified, the default is `default-ipsec-policy`.

* `order` - (Optional, type `int`) Specifies the order in which traffic is matched, if traffic can be matched to multiple traffic selectors.Traffic is matched to the traffic selector with the highest priority (lowest order number).
When creating a new traffic selector, if this parameter is not specified, the default is `last`

* `ip_protocol` - (Optional, type `int`) Specifies the network protocol to use for this traffic. The default value is `All Protocols (255)`
