---
layout: "bigip"
page_title: "BIG-IP: bigip_waf signatures"
subcategory: "Web Application Firewall(WAF)"
description: |-
  Provides details about system installed bigip_waf_signatures data source
---

# bigip\_waf\_signatures

Use this data source (`bigip_waf_signatures`) to get the details of attack signatures available on BIG-IP WAF
 
 
## Example Usage

```hcl

data "bigip_waf_signatures" "WAFSIG1" {
  signature_id = 200104004
}

```

## Argument Reference

* `signature_id` - (Required) ID of the signature in the BIG-IP WAF database.


## Attributes Reference

* `name` - Name of the signature as configured on the system.
* `description` - Description of the signature.
* `system_signature_id` - System generated ID of the signature.
* `signature_id` - ID of the signature in the database.
* `type` - Type of the signature.
* `accuracy` - The relative detection accuracy of the signature.
* `risk` - The relative risk level of the attack that matches this signature.
