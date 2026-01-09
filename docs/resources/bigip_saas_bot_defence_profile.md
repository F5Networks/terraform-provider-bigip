---
layout: "bigip"
page_title: "BIG-IP: bigip_saas_bot_defense_profile"
subcategory: "Distributed Cloud Services"
description: |-
  Provides details about bigip_saas_bot_defense_profile resource
---

# bigip\_saas\_bot\_defense\_profile

`bigip_saas_bot_defense_profile` Resource used for Configures Distributed Cloud Services Bot Defense profile, for more info [Bot defence](https://techdocs.f5.com/en-us/bigip-17-1-0/big-ip-saas-bot-defense.html)

## Example Usage

```hcl
resource "bigip_saas_bot_defense_profile" "test-bot-defense" {
  name                  = "/Common/test-saas-bot-defense"
  application_id        = "89fb0bfcb4bxxxxx8fad9adb37ce3b19"
  tenant_id             = "a-aaxxxxxvaYOV"
  api_key               = "49840d1dd6faxxxxxxc88762eb398eee"
  shape_protection_pool = "/Common/cs1.pool"
  ssl_profile           = "/Common/cloud-service-default-ssl"
  protected_endpoints {
    name     = "pe1"
    host     = "abc.com"
    endpoint = "/login"
    post     = "enabled"
  }
}
```
## Argument Reference

* `name` - (Required,type `string`) Unique name for the Distributed Cloud Services Bot Defense profile. Full path is the combination of the `partition + profile name`, for example `/Common/test-bot-tc1`.

* `defaults_from` - (Optional,type `string`) Distributed Cloud Services Bot Defense parent profile from which this profile will inherit settings. The default is the system-supplied `bd` profile.

* `description` - (Optional,type `string`) Specifies descriptive text that identifies the BD profile.

* `application_id` - (Required,type `string`) Specifies the Bot Defense API application ID, enter the value provided by F5 Support

* `api_key` - (Required,type `string`) Specifies the API key, enter the value provided by F5 Support.

* `tenant_id` - (Required,type `string`) Specifies the tenant ID, enter the value provided by F5 Support.

* `shape_protection_pool` - (Required,type `string`) Specifies the web hostname to which API requests are made.

* `ssl_profile` - (Required,type `string`) Specifies a server-side SSL profile that is different from what the application pool uses.


* `protected_endpoints` - (Required,`list`) Use these settings to configure which pages on the website will be protected by BD.
It is block `protected_endpoints` block takes input for the protected endpoints. See [protected endpoints](#protected-endpoints) below for more details.

### protected endpoints

This block supports the following arguments:

* `name` - (Required,`string`) Unique name for the protected endpoint.

* `host` - (Optional,`string`) hostname or IP address of the web page to be protected by the Bot Defense.

* `endpoint` - (Optional,`string`) Specifies the path to the web page to be protected by BD. For example, `/login`.

* `mitigation_action` - (Optional,`string`) Specifies whether the BIG-IP or F5 XC Bot Defense handles mitigation of malicious HTTP requests. This field is enabled only if the Service Level field is set to Advanced/Premium.

* `post` - (Optional,`string`) POST field to protect the path when it has a POST method, `enabled` or `disabled`.

* `put` - (Optional,`string`) PUT field to protect the path when it has a PUT method,`enabled` or `disabled`.

## Import

BIG-IP Distributed Cloud Services Bot Defense profile can be imported using the `/<partition>/<profile-name>`, e.g.
```bash
terraform import bigip_saas_bot_defense_profile.test-bot /Common/testbotprofile01
```
