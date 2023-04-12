package bigip

import (
	"context"
	"encoding/json"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceBigipFastGceServiceDiscovery() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataBigipFastGceServiceDiscoveryRead,
		Schema: map[string]*schema.Schema{
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "gce",
			},
			"port": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  80,
			},
			"tag_key": {
				Type:     schema.TypeString,
				Required: true,
			},
			"tag_value": {
				Type:     schema.TypeString,
				Required: true,
			},
			"region": {
				Type:     schema.TypeString,
				Required: true,
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
			"encoded_credentials": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"project_id": {
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
			"gce_sd_json": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The JSON for Consul service discovery block",
			},
		},
	}
}

func dataBigipFastGceServiceDiscoveryRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	config, err := getGceConfig(d)
	if err != nil {
		return diag.FromErr(err)
	}
	_ = d.Set("gce_sd_json", config)
	d.SetId(hashForState(d.Get("gce_sd_json").(string)))
	return nil
}

func getGceConfig(d *schema.ResourceData) (string, error) {
	var sdGce bigip.SDGceObject
	gcePort := d.Get("port").(int)
	sdGce.SdType = d.Get("type").(string)
	sdGce.SdPort = &gcePort
	sdGce.SdTagKey = d.Get("tag_key").(string)
	sdGce.SdTagVal = d.Get("tag_value").(string)
	sdGce.SdRegion = d.Get("region").(string)
	sdGce.SdAddressRealm = d.Get("address_realm").(string)
	sdGce.SdUndetectableAction = d.Get("undetectable_action").(string)
	sdGce.SdEncodedCredentials = d.Get("encoded_credentials").(string)
	sdGce.SdProjectId = d.Get("project_id").(string)
	sdGce.SdMinimumMonitors = d.Get("minimum_monitors").(string)
	sdGce.SdCredentialUpdate = d.Get("credential_update").(bool)
	sdGce.SdUpdateInterval = d.Get("update_interval").(string)

	sdGcestring, err := json.Marshal(&sdGce)
	if err != nil {
		return "", err
	}
	return string(sdGcestring), nil
}
