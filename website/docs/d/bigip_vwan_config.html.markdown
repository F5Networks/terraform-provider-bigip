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
  azure_client_id          = "XXXXXXXXXXXXXXXXXXXXXXXXXXX"
  azure_client_secret      = "XXXXXXXXXXXXXXXXXXXXXXXXXXX"
  azure_subsciption_id     = "XXXXXXXXXXXXXXXXXXXXXXXXXXX"
  azure_tenant_id          = "XXXXXXXXXXXXXXXXXXXXXXXXXXX"
  storage_accounnt_name    = "XXXXXXXXXXXXXXXX"
  storage_accounnt_key     = "XXXXXXXXXXXXXXXXXXXXX"

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

* `storage_accounnt_name` - (Required) Specifies the storage account for download config,can be set as Environment Variable `STORAGE_ACCOUNT_NAME`.

* `storage_accounnt_key` - (Required) Specifies the storage account key to authenticate,can be set as Environment Variable `STORAGE_ACCOUNT_KEY`.
