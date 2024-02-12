---
layout: "bigip"
page_title: "BIG-IP: bigip_ltm_profile_bot_defense"
subcategory: "Local Traffic Manager(LTM)"
description: |-
  Provides details about bigip_ltm_profile_bot_defense resource
---

# bigip\_ltm\_profile\_bot\_defense

`bigip_ltm_profile_bot_defense` Resource used for Configures a Bot Defense profile.

## Example Usage

```hcl
resource "bigip_ltm_profile_bot_defense" "test-bot-tc1" {
  name          = "/Common/test-bot-tc1"
  defaults_from = "/Common/bot-defense"
  description   = "test-bot"
  template      = "relaxed"
}
```      

## Argument Reference

* `name` (Required,type `string`) Name of the Bot Defense profile,name of Profile should be full path. Full path is the combination of the `partition + profile name`,For example `/Common/test-bot-tc1`.

* `defaults_from` - (optional,type `string`) Specifies the profile from which this profile inherits settings. The default is the system-supplied `bot-defense` profile.

* `description` - (optional,type `string`) Specifies user-defined description.

* `template` - (Optional,type `string`) Profile templates specify Mitigation and Verification Settings default values. possible ptions `balanced`,`relaxed` and `strict`.

* `enforcement_mode` - (Optional,type `string`) Select the enforcement mode, possible values are `transparent` and `blocking`.


## Import

BIG-IP LTM Bot Defense profile can be imported using the `name`, e.g.

```bash
terraform import bigip_ltm_profile_bot_defense.test-bot /Common/testbotprofile01
```
