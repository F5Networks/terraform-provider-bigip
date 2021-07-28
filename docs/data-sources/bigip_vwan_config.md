---
layout: "bigip"
page_title: "BIG-IP: bigip_vwan_config"
subcategory: "vWAN"
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
}

```      

## Argument Reference

* `azure_vwan_resourcegroup` - (Required, type `string`) Name of the Azure vWAN resource group

* `azure_vwan_name` - (Required,type `string`) Name of the Azure vWAN Name

* `azure_vwan_vpnsite` - (Required,type `string`) Name of the Azure vWAN VPN site from which configuration to be download


## Pre-required Environment Settings:

* `AZURE_CLIENT_ID` - (Required) Set this environment variable with the Azure app client ID to use.

* `AZURE_CLIENT_SECRET` - (Required) Set this environment variable with the Azure app secret to use.

* `AZURE_SUBSCRIPTION_ID` - (Required) Set this environment variable with the Azure subscription ID to use.

* `AZURE_TENANT_ID` - (Required) Set this environment variable with the Tenant ID to which to authenticate.

* `STORAGE_ACCOUNT_NAME` - (Required) Set this environment variable with the storage account for download config.

* `STORAGE_ACCOUNT_KEY` - (Required) Specifies the storage account key to authenticate,set this Environment variable with account key value.

## Attributes Reference

* `bigip_gw_ip` - (type `string`) provides IP address of BIGIP G/W for IPSec Endpoint.

* `preshared_key` - (type `string`) provides pre-shared-key used for IPSec Tunnel creation.

* `hub_address_space` - (type `string`) Provides IP Address space used on vWAN Hub.

* `hub_connected_subnets` - (type `list`) Provides Subnets connected to vWAN Hub.

* `vwan_gw_address` - (type `list`) Provides vWAN Gateway Address for IPSec End point

