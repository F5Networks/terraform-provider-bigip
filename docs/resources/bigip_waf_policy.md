---
layout: "bigip"
page_title: "BIG-IP: bigip_waf_policy"
subcategory: "Web Application Firewall(WAF)"
description: |-
  Provides details about bigip_waf_policy resource
---

# bigip_waf_policy

`bigip_waf_policy` Manages a WAF Policy resource with its adjustments and modifications on a BIG-IP.
It outputs an up-to-date WAF Policy in a JSON format

* [Declarative WAF documentation](https://clouddocs.f5.com/products/waf-declarative-policy/declarative_policy_v16_1.html)

~> **NOTE** This Resource Requires F5 BIG-IP v16.x above version, and ASM need to be provisioned.

## Example Usage 

```hcl

data "bigip_waf_entity_parameter" "Param1" {
  name            = "Param1"
  type            = "explicit"
  data_type       = "alpha-numeric"
  perform_staging = true
}

data "bigip_waf_entity_parameter" "Param2" {
  name            = "Param2"
  type            = "explicit"
  data_type       = "alpha-numeric"
  perform_staging = true
}

data "bigip_waf_entity_url" "URL" {
  name     = "URL1"
  protocol = "http"
}

data "bigip_waf_entity_url" "URL2" {
  name = "URL2"
}

resource "bigip_waf_policy" "test-awaf" {
  name                 = "testpolicyravi"
  partition            = "Common"
  template_name        = "POLICY_TEMPLATE_RAPID_DEPLOYMENT"
  application_language = "utf-8"
  enforcement_mode     = "blocking"
  server_technologies  = ["MySQL", "Unix/Linux", "MongoDB"]
  parameters           = [data.bigip_waf_entity_parameter.Param1.json, data.bigip_waf_entity_parameter.Param2.json]
  urls                 = [data.bigip_waf_entity_url.URL.json, data.bigip_waf_entity_url.URL2.json]
}

```

## Argument Reference

* `name` - (Required,type `string`) The unique user-given name of the policy. Policy names cannot contain spaces or special characters. Allowed characters are a-z, A-Z, 0-9, dot, dash (-), colon (:) and underscore (_).

* `template_name` - (Required,type `string`) Specifies the name of the template used for the policy creation.

* `description` - (Optional,type `string`) Specifies the description of the policy.

* `partition` - (Optional,type `string`) Specifies the partition of the policy. Default is `Common`

* `application_language` - (Optional,type `string`) The character encoding for the web application. The character encoding determines how the policy processes the character sets. The default is `utf-8`

* `case_insensitive` - (Optional,type `bool`) Specifies whether the security policy treats microservice URLs, file types, URLs, and parameters as case sensitive or not. When this setting is enabled, the system stores these security policy elements in lowercase in the security policy configuration

* `enable_passivemode` - (Optional,type `bool`) Passive Mode allows the policy to be associated with a Performance L4 Virtual Server (using a FastL4 profile). With FastL4, traffic is analyzed but is not modified in any way.

* `protocol_independent` - (Optional,type `bool`) When creating a security policy, you can determine whether a security policy differentiates between HTTP and HTTPS URLs. If enabled, the security policy differentiates between HTTP and HTTPS URLs. If disabled, the security policy configures URLs without specifying a specific protocol. This is useful for applications that behave the same for HTTP and HTTPS, and it keeps the security policy from including the same URL twice.

* `enforcement_mode` - (Optional,type `string`) How the system processes a request that triggers a security policy violation

* `type` - (Optional,type `string`) The type of policy you want to create. The default policy type is `security`.

* `server_technologies` - (Optional,type `list`) The server technology is a server-side application, framework, web server or operating system type that is configured in the policy in order to adapt the policy to the checks needed for the respective technology.

* `parameters` - (Optional,type `list`) This section defines parameters that the security policy permits in requests.

* `urls` - (Optional,type `list`) In a security policy, you can manually specify the HTTP URLs that are allowed (or disallowed) in traffic to the web application being protected. If you are using automatic policy building (and the policy includes learning URLs), the system can determine which URLs to add, based on legitimate traffic.

* `signature_sets` - (Optional,type `list`) Defines behavior when signatures found within a signature-set are detected in a request. Settings are culmulative, so if a signature is found in any set with block enabled, that signature will have block enabled.

* `signatures` - (Optional,type `list`) This section defines the properties of a signature on the policy.

* `policy_builder` - (Optional,`set`) `policy_builder` block will provide `learning_mode` options to be used for policy builder.
See [policy builder](#policy-builder) below for more details.

* `graphql_profiles` - (Optional,`list of set`) `graphql_profiles` takes list of graphql profile options to be used for policy builder.
See [graphql profiles](#graphql-profiles) below for more details.

* `file_types` - (Optional,`list of set`) `file_types` takes list of file-types options to be used for policy builder.
See [file types](#file-types) below for more details.

* `open_api_files` - (Optional,type `list`) This section defines the Link for open api files on the policy.

* `policy_import_json` - (Optional,type `string`) The payload of the WAF Policy to be used for IMPORT on to BIG-IP. 

### policy builder
The `policy_builder` block supports the following:

* `learning_mode` - (Optional , `string`) learning mode setting for policy-builder, possible options: [`automatic`,`disabled`, `manual`]

### graphql profiles
The `graphql_profile` block supports the following:

* `name` - (Optional , `string`) name of graphql profile to be used for policy config.

### file types
The `file_types` block supports the following:

* `name` - (Optional , `string`) Specifies the file type name as appearing in the URL extension.

* `type` - (Optional , `string`) Determines the type of the name attribute. Only when setting the type to `wildcard` will the special wildcard characters in the name be interpreted as such

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `policy_id` - The id of the A.WAF Policy as it would be calculated on the BIG-IP.

* `policy_export_json` - Exported WAF policy deployed on BIGIP.


## Import
An existing WAF Policy or if the WAF Policy has been manually created or modified on the BIG-IP WebUI, it can be imported using its `id`.

e.g:

```
terraform import bigip_waf_policy.example <id>
```
