package bigip

import (
	"context"
	"encoding/json"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceBigipFastAwsServiceDiscovery() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataBigipFastCAwsServiceDiscoveryRead,
		Schema: map[string]*schema.Schema{
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "aws",
			},
			"port": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  80,
			},
			"tag_key": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The tag key associated with the node to add to this pool",
			},
			"tag_value": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The tag value associated with the node to add to this pool",
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
			"aws_region": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Empty string (`default`) means region in which ADC is running",
				Optional:    true,
			},
			"aws_access_key": {
				Type:        schema.TypeString,
				Sensitive:   true,
				Description: "Information for discovering AWS nodes that are not in the same region as your BIG-IP (also requires the `aws_secret_access_key` field)",
				Optional:    true,
			},
			"aws_secret_access_key": {
				Type:      schema.TypeString,
				Sensitive: true,
				Optional:  true,
			},
			"external_id": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"role_arn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Assume a role (also requires the externalId field)",
				Optional:    true,
			},
			"minimum_monitors": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"update_interval": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"aws_sd_json": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The JSON for AWS service discovery block",
			},
		},
	}
}

func dataBigipFastCAwsServiceDiscoveryRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	config, err := getAwsSDConfig(d)
	if err != nil {
		return diag.FromErr(err)
	}
	_ = d.Set("aws_sd_json", config)
	d.SetId(hashForState(d.Get("tag_key").(string)))
	return nil
}

func getAwsSDConfig(d *schema.ResourceData) (string, error) {
	var sdAws bigip.SdAwsObj
	awsPort := d.Get("port").(int)
	sdAws.SdType = d.Get("type").(string)
	sdAws.SdPort = &awsPort
	sdAws.SdTagKey = d.Get("tag_key").(string)
	sdAws.SdTagVal = d.Get("tag_value").(string)
	sdAws.SdAddressRealm = d.Get("address_realm").(string)
	sdAws.SdUndetectableAction = d.Get("undetectable_action").(string)
	sdAws.SdAwsRegion = d.Get("aws_region").(string)
	sdAws.SdAccessKeyId = d.Get("aws_access_key").(string)
	sdAws.SdSecretAccessKey = d.Get("aws_secret_access_key").(string)
	sdAws.SdExternalId = d.Get("external_id").(string)
	sdAws.SdRoleARN = d.Get("role_arn").(string)
	sdAws.SdMinimumMonitors = d.Get("minimum_monitors").(string)
	sdAws.SdCredentialUpdate = d.Get("credential_update").(bool)
	sdAws.SdUpdateInterval = d.Get("update_interval").(string)

	sdAwsstring, err := json.Marshal(&sdAws)
	if err != nil {
		return "", err
	}
	return string(sdAwsstring), nil
}
