---
layout: "bigip"
page_title: "BIG-IP: bigip_ipsec_profile"
subcategory: "Network"
description: |-
   Provides details about bigip_ipsec_profile resource
---

# bigip_ipsec_profile

`bigip_ipsec_profile` Manage IPSec Profiles on a BIG-IP

## Example Usage

```hcl
resource "bigip_ipsec_profile" "azurevWAN_profile" {
  name             = "/Common/Mytestipsecprofile"
  description      = "mytestipsecprofile"
  traffic_selector = "test-trafficselector"
}

```      

## Argument Reference

* `name` - (Required) Displays the name of the IPsec interface tunnel profile,it should be "full path".The full path is the combination of the partition + name of the IPSec profile.(For example `/Common/test-profile`)

* `description` - (Optional,type `string`) Specifies descriptive text that identifies the IPsec interface tunnel profile.

* `parent_profile` - (Optional,type `string`) Specifies the profile from which this profile inherits settings. The default is the system-supplied `/Common/ipsec` profile

* `traffic_selector` - (Optional,type `string`) Specifies the traffic selector for the IPsec interface tunnel to which the profile is applied 
