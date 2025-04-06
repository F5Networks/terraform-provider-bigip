/*
Copyright 2024 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"context"
	"log"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceBigipSaasBotDefenseProfile() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBigipSaasBotDefenseProfileCreate,
		ReadContext:   resourceBigipSaasBotDefenseProfileRead,
		UpdateContext: resourceBigipSaasBotDefenseProfileUpdate,
		DeleteContext: resourceBigipSaasBotDefenseProfileDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "Unique name for the Distributed Cloud Services Bot Defense profile",
				ValidateFunc: validateF5NameWithDirectory,
			},
			"defaults_from": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "/Common/bd",
				Description:  "Distributed Cloud Services Bot Defense parent profile from which this profile will inherit settings.",
				ValidateFunc: validateF5Name,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies descriptive text that identifies the BD profile.",
			},
			"application_id": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "Specifies the Bot Defense API application ID, enter the value provided by F5 Support",
			},
			"tenant_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the tenant ID, enter the value provided by F5 Support",
			},
			"api_key": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "Specifies the API key, enter the value provided by F5 Support.",
			},
			"shape_protection_pool": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the web hostname to which API requests are made",
			},
			"ssl_profile": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies a server-side SSL profile that is different from what the application pool uses",
			},
			"protected_endpoints": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Use these settings to configure which pages on the website will be protected by BD",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Unique name for the protected endpoint",
						},
						"host": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "hostname or IP address of the web page to be protected by the Bot Defense",
						},
						"endpoint": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Specifies the path to the web page to be protected by BD. For example, `/login`.",
						},
						"mitigation_action": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Specifies whether the BIG-IP or F5 XC Bot Defense handles mitigation of malicious HTTP requests. This field is enabled only if the Service Level field is set to Advanced/Premium",
						},
						"post": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "POST field to protect the path when it has a POST method, `enabled` or `disabled`",
						},
						"put": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "PUT field to protect the path when it has a PUT method,`enabled` or `disabled`",
						},
					},
				},
			},
		},
	}
}

func resourceBigipSaasBotDefenseProfileCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	name := d.Get("name").(string)
	log.Printf("[INFO] Creating Bot Defense Profile:%+v ", name)
	pss := &bigip.SaasBotDefenseProfile{
		Name: name,
	}
	config := getSaasBotDefenseProfileConfig(d, pss)
	log.Printf("[DEBUG] Bot Defense Profile config :%+v ", config)
	err := client.AddSaasBotDefenseProfile(config)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(name)
	return resourceBigipSaasBotDefenseProfileRead(ctx, d, meta)
}

func resourceBigipSaasBotDefenseProfileRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	log.Printf("[INFO] Reading Bot Defense Profile:%+v ", client)
	name := d.Id()
	log.Printf("[INFO] Reading Bot Defense Profile:%+v ", name)
	botProfile, err := client.GetSaasBotDefenseProfile(name)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] Defense Profile Resp :%+v ", botProfile)
	_ = d.Set("name", botProfile.FullPath)
	_ = d.Set("defaults_from", botProfile.DefaultsFrom)
	_ = d.Set("description", botProfile.Description)
	_ = d.Set("tenant_id", botProfile.TenantId)
	_ = d.Set("api_key", botProfile.ApiKey)
	_ = d.Set("shape_protection_pool", botProfile.ShapeProtectionPool)
	_ = d.Set("ssl_profile", botProfile.SslProfile)
	_ = d.Set("protected_endpoints", flattenProtectedEndpointsReference(botProfile.ProtectedEndpointsReference.Items))
	return nil
}

func resourceBigipSaasBotDefenseProfileUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	name := d.Id()
	log.Printf("[INFO] Updating Bot Defense Profile:%+v ", name)
	pss := &bigip.SaasBotDefenseProfile{
		Name: name,
	}
	config := getSaasBotDefenseProfileConfig(d, pss)

	err := client.ModifySaasBotDefenseProfile(name, config)
	if err != nil {
		return diag.FromErr(err)
	}
	return resourceBigipSaasBotDefenseProfileRead(ctx, d, meta)
}

func resourceBigipSaasBotDefenseProfileDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	name := d.Id()
	log.Println("[INFO] Deleting Bot Defense Profile " + name)
	err := client.DeleteSaasBotDefenseProfile(name)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return nil
}

func getSaasBotDefenseProfileConfig(d *schema.ResourceData, config *bigip.SaasBotDefenseProfile) *bigip.SaasBotDefenseProfile {
	config.Name = d.Get("name").(string)
	config.DefaultsFrom = d.Get("defaults_from").(string)
	config.Description = d.Get("description").(string)
	config.ApplicationId = d.Get("application_id").(string)
	config.TenantId = d.Get("tenant_id").(string)
	config.ApiKey = d.Get("api_key").(string)
	config.ShapeProtectionPool = d.Get("shape_protection_pool").(string)
	config.SslProfile = d.Get("ssl_profile").(string)
	var protectEndpoint []bigip.ProtectedEndpoint
	for _, endpoint := range d.Get("protected_endpoints").([]interface{}) {
		ep := endpoint.(map[string]interface{})
		protectEndpoint = append(protectEndpoint, bigip.ProtectedEndpoint{
			Name:             ep["name"].(string),
			Host:             ep["host"].(string),
			Endpoint:         ep["endpoint"].(string),
			Post:             ep["post"].(string),
			Put:              ep["put"].(string),
			MitigationAction: ep["mitigation_action"].(string),
		})
	}
	config.ProtectedEndpointsReference.Items = protectEndpoint
	log.Printf("[INFO][getSaasBotDefenseProfileConfig] config:%+v ", config)
	return config
}

func flattenProtectedEndpointsReference(data interface{}) []interface{} {
	var prtctEndpt []interface{}
	for _, ep := range data.([]bigip.ProtectedEndpoint) {
		prtctEndpt = append(prtctEndpt, map[string]interface{}{
			"name":              ep.Name,
			"host":              ep.Host,
			"endpoint":          ep.Endpoint,
			"post":              ep.Post,
			"put":               ep.Put,
			"mitigation_action": ep.MitigationAction,
		})
	}
	return prtctEndpt
}

// {
//     "name": "/Common/bd-test",
//     "applicationId": "89fb0bfcb4bf4c578fad9adb37ce3b19",
//     "tenantId": "a-aavN9vaYOV",
//     "apiKey": "49840d1dd6fa4c4d86c88762eb398eee",
//     "shapeProtectionPool": "/Common/cs1.pool",
//     "sslProfile": "/Common/cloud-service-default-ssl",
//     "protectedEndpointsReference": {
//         "items": [
//             {
//                 "name": "pe1",
//                 "host": "abc.com",
//                 "endpoint": "/login",
//                 "post": "enabled"
//             }
//         ]
//     }
// }
