package bigip

import (
	"context"
	"encoding/json"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceBigipFastAzureServiceDiscovery() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataBigipFastAzureServiceDiscoveryRead,
		Schema: map[string]*schema.Schema{
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "azure",
			},
			"port": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  80,
			},
			"resource_group": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of Azure Resource group where Nodes reside",
			},
			"subscription_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Id of Azure Subscription Nodes",
			},
			"address_realm": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies whether to look for public or private IP addresses",
				Default:     "private",
			},
			"undetectable_action": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Action to take when node cannot be detected",
				Default:     "remove",
			},
			"credential_update": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Specifies whether you are updating your credentials",
			},
			"tag_key": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tag_value": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"minimum_monitors": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"update_interval": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"azure_sd_json": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The JSON for Azure service discovery block",
			},
		},
	}
}

func dataBigipFastAzureServiceDiscoveryRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	config, err := getAzureSDConfig(d)
	if err != nil {
		return diag.FromErr(err)
	}
	_ = d.Set("azure_sd_json", config)
	d.SetId(hashForState(d.Get("azure_sd_json").(string)))
	return nil
}

func getAzureSDConfig(d *schema.ResourceData) (string, error) {
	var sdAzure bigip.SDAzureObject
	azurePort := d.Get("port").(int)
	sdAzure.SdType = d.Get("type").(string)
	sdAzure.SdPort = &azurePort
	sdAzure.SdRg = d.Get("resource_group").(string)
	sdAzure.SdSid = d.Get("subscription_id").(string)
	sdAzure.SdRtype = "tag"
	sdAzure.SdUseManagedIdentity = true
	sdAzure.SdAzureTagKey = d.Get("tag_key").(string)
	sdAzure.SdAzureTagVal = d.Get("tag_value").(string)
	sdAzure.SdAddressRealm = d.Get("address_realm").(string)
	sdAzure.SdUndetectableAction = d.Get("undetectable_action").(string)
	sdAzure.SdMinimumMonitors = d.Get("minimum_monitors").(string)
	sdAzure.SdCredentialUpdate = d.Get("credential_update").(bool)
	sdAzure.SdUpdateInterval = d.Get("update_interval").(string)

	sdAzurestring, err := json.Marshal(&sdAzure)
	if err != nil {
		return "", err
	}
	return string(sdAzurestring), nil
}
