---
layout: "bigip"
page_title: "BIG-IP: bigip_vwan_config"
sidebar_current: "docs-bigip-datasource-vwan-config-x"
description: |-
    Provides details about bigip_vwan_config data source
---

# bigip_vwan_config

Use this data source (`bigip_vwan_config`) to get the vWAN site config from Azure VWAN Site
 
 
## Example Usage
```hcl

data "bigip_vwan_config" "vwanconfig" {
  azure_vwan_resourcegroup = "azurevwan-bigip-rg-9c8d"
  azure_vwan_name          = "azurevwan-bigip-vwan-9c8d"
  azure_vwan_vpnsite       = "azurevwan-bigip-vsite-9c8d"
  azure_client_id          = "bd28e9c9-ef78-4aac-8517-16384384c80d"
  azure_client_secret      = "6907c06b-3166-404a-9fd4-b288326503f8"
  azure_subsciption_id     = "d31e4e54-7577-4f43-b407-bae6cc0f4f55"
  azure_tenant_id          = "d106871e-7b91-4733-8423-f98586303b68"
}


```      

## Argument Reference

* `azure_vwan_resourcegroup` - (Required) Name of the Azure vWAN resource group

* `azure_vwan_name` - (Required) Name of the Azure vWAN Name

* `azure_vwan_vpnsite` - (Required) Name of the Azure vWAN VPN site from which configuration to be download

* `azure_client_id` - (Required) Specifies the Azure app client ID to use,can be set as Environment Variable `AZURE_CLIENT_ID`.

* `azure_client_secret` - (Required) Specifies the Azure app secret to use,can be set as Environment Variable `AZURE_CLIENT_SECRET`.

* `azure_subsciption_id` - (Required) Specifies the Azure subscription ID to use,can be set as Environment Variable `AZURE_SUBSCRIPTION_ID`.

* `azure_tenant_id` - (Required) Specifies the Tenant to which to authenticate,can be set as Environment Variable `AZURE_TENANT_ID`.

## Attributes Reference

* `type` - The Data Group type (string, ip, integer)"

* `record` - Specifies record of type (string/ip/integer)
