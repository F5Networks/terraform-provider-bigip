/*
Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	//"context"
	//"encoding/json"
	//"fmt"
	//"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-11-01/network"
	//"github.com/Azure/azure-storage-blob-go/azblob"
	//"github.com/Azure/go-autorest/autorest"
	//"github.com/Azure/go-autorest/autorest/adal"
	//"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	//"io/ioutil"
	//"log"
	//"net/url"
	//"os"
	//"time"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceBigipVwanconfig() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBigipVwanconfigRead,
		Schema: map[string]*schema.Schema{
			"azure_vwan_resourcegroup": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name azure_vwan_resourcegroup",
			},
			"azure_vwan_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name azure_vwan_name",
			},
			"azure_vwan_vpnsite": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name azure_vwan_vpnsite",
			},
			"azure_subsciption_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "partition of resource group",
			},
			"azure_client_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "partition of resource group",
			},
			"azure_client_secret": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "partition of resource group",
			},
			"azure_tenant_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "partition of resource group",
			},
			"full_path": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The full path, if found.",
			},
		},
	}
}
func dataSourceBigipVwanconfigRead(d *schema.ResourceData, meta interface{}) error {
	//client := meta.(*bigip.BigIP)
	d.SetId("")
	//log.Println("[INFO] Reading VWAN Config : " + d.Get("azure_vwan_resourcegroup").(string))
	//DownloadVwanConfig()

	d.SetId("wanconfig")

	return nil

}
