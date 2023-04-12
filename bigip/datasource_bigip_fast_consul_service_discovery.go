package bigip

import (
	"context"
	"encoding/json"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceBigipFastConsulServiceDiscovery() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataBigipFastConsulServiceDiscoveryRead,
		Schema: map[string]*schema.Schema{
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "consul",
			},
			"uri": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The location of the node data",
			},
			"port": {
				Type:     schema.TypeInt,
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
			"encoded_token": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Base 64 encoded bearer token to make requests to the Consul API. Will be stored in the declaration in an encrypted format.",
			},
			"jmes_path_query": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Custom JMESPath Query",
			},
			"reject_unauthorized": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "If true, the server certificate is verified against the list of supplied/default CAs when making requests to the Consul API.",
			},
			"trust_ca": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "CA Bundle to validate server certificates",
			},
			"minimum_monitors": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"update_interval": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"consul_sd_json": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The JSON for Consul service discovery block",
			},
		},
	}
}

func dataBigipFastConsulServiceDiscoveryRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	config, err := getConsulSDConfig(d)
	if err != nil {
		return diag.FromErr(err)
	}
	_ = d.Set("consul_sd_json", config)
	d.SetId(d.Get("uri").(string))
	return nil
}

func getConsulSDConfig(d *schema.ResourceData) (string, error) {
	var sdConsul bigip.SDConsulObject
	consulPort := d.Get("port").(int)
	sdConsul.SdType = d.Get("type").(string)
	sdConsul.SdPort = &consulPort
	sdConsul.SdUri = d.Get("uri").(string)
	sdConsul.SdAddressRealm = d.Get("address_realm").(string)
	sdConsul.SdUndetectableAction = d.Get("undetectable_action").(string)
	sdConsul.SdEncodedToken = d.Get("encoded_token").(string)
	sdConsul.SdJmesPathQuery = d.Get("jmes_path_query").(string)
	sdConsul.SdMinimumMonitors = d.Get("minimum_monitors").(string)
	sdConsul.SdCredentialUpdate = d.Get("credential_update").(bool)
	sdConsul.SdRejectUnauthorized = d.Get("reject_unauthorized").(bool)
	sdConsul.SdTrustCA = d.Get("trust_ca").(string)
	sdConsul.SdUpdateInterval = d.Get("update_interval").(string)

	sdConsulstring, err := json.Marshal(&sdConsul)
	if err != nil {
		return "", err
	}
	return string(sdConsulstring), nil
}
