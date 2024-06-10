/*
Copyright 2024 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"context"
	"log"
	"os"
	"strings"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/f5devcentral/go-bigip/f5teem"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceBigipLtmProfileRequestLog() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBigipLtmProfileRequestLogCreate,
		ReadContext:   resourceBigipLtmProfileRequestLogRead,
		UpdateContext: resourceBigipLtmProfileRequestLogUpdate,
		DeleteContext: resourceBigipLtmProfileRequestLogDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "Name of the Request Logging profile",
				ValidateFunc: validateF5NameWithDirectory,
			},
			"defaults_from": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "/Common/request-log",
				Description:  "Specifies the profile from which this profile inherits settings. The default is the system-supplied `request-log` profile",
				ValidateFunc: validateF5Name,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "User defined description for Request Logging profile",
			},
			"request_logging": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"disabled",
					"enabled"}, false),
				Description: "Enables or disables request logging. The default is `disabled`",
			},
			"requestlog_protocol": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"mds-udp",
					"mds-tcp"}, false),
				Description: "Specifies the protocol to be used for high-speed logging of requests. The default is `mds-udp`",
			},
			"requestlog_template": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the directives and entries to be logged.",
			},
			"requestlog_error_template": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the directives and entries to be logged.",
			},
			"requestlog_error_protocol": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"mds-udp",
					"mds-tcp"}, false),
				Description: "Defines the protocol to be used for high-speed logging of request errors. The default is `mds-udp`",
			},
			"requestlog_pool": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Defines the pool to send logs to. Typically, the pool will contain one or more syslog servers. It is recommended that you create a pool specifically for logging requests. The default is None",
			},
			"requestlog_error_pool": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Defines the pool associated with logging request errors. The default is None.",
			},
			"proxy_response": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Defines the pool associated with logging request errors. The default is None.",
			},
			"proxyclose_on_error": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Defines the pool associated with logging request errors. The default is None.",
			},
			"proxyrespond_on_loggingerror": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Defines the pool associated with logging request errors. The default is None.",
			},
			"response_logging": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"disabled",
					"enabled"}, false),
				Description: "Enables or disables response logging. The default is `disabled`",
			},
			"responselog_template": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the directives and entries to be logged.",
			},
			"responselog_error_template": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the directives and entries to be logged.",
			},
			"responselog_protocol": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"mds-udp",
					"mds-tcp"}, false),
				Description: "Specifies the protocol to be used for high-speed logging of responses. The default is `mds-udp`",
			},
			"responselog_error_protocol": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"mds-udp",
					"mds-tcp"}, false),
				Description: "Defines the protocol to be used for high-speed logging of responses errors. The default is `mds-udp`",
			},
			"responselog_pool": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Defines the pool to send logs to. Typically, the pool contains one or more syslog servers. It is recommended that you create a pool specifically for logging responses. The default is None",
			},
			"responselog_error_pool": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Defines the pool associated with logging response errors. The default is None.",
			},
		},
	}
}

func resourceBigipLtmProfileRequestLogCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	name := d.Get("name").(string)
	log.Printf("[INFO] Creating Request Log Profile:%+v ", name)

	pss := &bigip.RequestLogProfile{
		Name: name,
	}
	config := getRequestLogProfileConfig(d, pss)

	err := client.AddRequestLogProfile(config)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(name)

	if !client.Teem {
		id := uuid.New()
		uniqueID := id.String()
		assetInfo := f5teem.AssetInfo{
			Name:    "Terraform-provider-bigip",
			Version: client.UserAgent,
			Id:      uniqueID,
		}
		apiKey := os.Getenv("TEEM_API_KEY")
		teemDevice := f5teem.AnonymousClient(assetInfo, apiKey)
		f := map[string]interface{}{
			"Terraform Version": client.UserAgent,
		}
		tsVer := strings.Split(client.UserAgent, "/")
		err = teemDevice.Report(f, "bigip_ltm_request_log_profile", tsVer[3])
		if err != nil {
			log.Printf("[ERROR]Sending Telemetry data failed:%v", err)
		}
	}
	return resourceBigipLtmProfileRequestLogRead(ctx, d, meta)
}

func resourceBigipLtmProfileRequestLogRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	log.Println("[INFO] Fetching HTTP  Profile " + name)

	pp, err := client.GetRequestLogProfile(name)
	if err != nil {
		log.Printf("[ERROR] Unable to retrieve Request Log Profile  (%s) ", err)
		return diag.FromErr(err)
	}
	if pp == nil {
		log.Printf("[WARN] Request Log Profile (%s) not found, removing from state", name)
		d.SetId("")
		return nil
	}
	_ = d.Set("name", name)
	_ = d.Set("defaults_from", pp.DefaultsFrom)

	if _, ok := d.GetOk("request_logging"); ok {
		_ = d.Set("request_logging", pp.RequestLogging)
	}
	if _, ok := d.GetOk("requestlog_pool"); ok {
		_ = d.Set("requestlog_pool", pp.RequestLogPool)
	}
	if _, ok := d.GetOk("description"); ok {
		_ = d.Set("description", pp.Description)
	}
	if _, ok := d.GetOk("requestlog_error_pool"); ok {
		_ = d.Set("requestlog_error_pool", pp.RequestLogErrorPool)
	}
	if _, ok := d.GetOk("requestlog_template"); ok {
		_ = d.Set("requestlog_template", strings.ReplaceAll(pp.RequestLogTemplate, `\"`, `"`))
	}
	if _, ok := d.GetOk("requestlog_protocol"); ok {
		_ = d.Set("requestlog_protocol", pp.RequestLogProtocol)
	}
	if _, ok := d.GetOk("requestlog_error_protocol"); ok {
		_ = d.Set("requestlog_error_protocol", pp.RequestLogErrorProtocol)
	}
	if _, ok := d.GetOk("responselog_protocol"); ok {
		_ = d.Set("responselog_protocol", pp.ResponseLogProtocol)
	}
	if _, ok := d.GetOk("responselog_error_protocol"); ok {
		_ = d.Set("responselog_error_protocol", pp.ResponseLogErrorProtocol)
	}
	if _, ok := d.GetOk("responselog_pool"); ok {
		_ = d.Set("responselog_pool", pp.ResponseLogPool)
	}
	if _, ok := d.GetOk("responselog_error_pool"); ok {
		_ = d.Set("responselog_error_pool", pp.ResponseLogErrorPool)
	}
	if _, ok := d.GetOk("proxy_response"); ok {
		_ = d.Set("proxy_response", pp.ProxyResponse)
	}
	if _, ok := d.GetOk("proxyclose_on_error"); ok {
		_ = d.Set("proxyclose_on_error", pp.ProxyCloseOnError)
	}
	if _, ok := d.GetOk("proxyrespond_on_loggingerror"); ok {
		_ = d.Set("proxyrespond_on_loggingerror", pp.ProxyRespondOnLoggingError)
	}
	if _, ok := d.GetOk("response_logging"); ok {
		_ = d.Set("response_logging", pp.ResponseLogging)
	}
	if _, ok := d.GetOk("responselog_template"); ok {
		_ = d.Set("responselog_template", strings.ReplaceAll(pp.ResponseLogTemplate, `\"`, `"`))
	}
	if _, ok := d.GetOk("requestlog_error_template"); ok {
		_ = d.Set("requestlog_error_template", pp.RequestLogErrorTemplate)
	}
	if _, ok := d.GetOk("responselog_error_template"); ok {
		_ = d.Set("responselog_error_template", pp.ResponseLogErrorTemplate)
	}
	// if _, ok := d.GetOk("request_chunking"); ok {
	// 	_ = d.Set("request_chunking", pp.RequestChunking)
	// }
	// if _, ok := d.GetOk("response_chunking"); ok {
	// 	_ = d.Set("response_chunking", pp.ResponseChunking)
	// }
	// _ = d.Set("response_headers_permitted", pp.ResponseHeadersPermitted)

	// if _, ok := d.GetOk("server_agent_name"); ok {
	// 	_ = d.Set("server_agent_name", pp.ServerAgentName)
	// }
	// if _, ok := d.GetOk("via_host_name"); ok {
	// 	_ = d.Set("via_host_name", pp.ViaHostName)
	// }
	// if _, ok := d.GetOk("via_request"); ok {
	// 	_ = d.Set("via_request", pp.ViaRequest)
	// }
	// if _, ok := d.GetOk("via_response"); ok {
	// 	_ = d.Set("via_response", pp.ViaResponse)
	// }
	// _ = d.Set("xff_alternative_names", pp.XffAlternativeNames)
	return nil
}

func resourceBigipLtmProfileRequestLogUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	name := d.Id()
	log.Printf("[INFO] Updating Request Log Profile:%+v ", name)

	pss := &bigip.RequestLogProfile{
		Name: name,
	}
	config := getRequestLogProfileConfig(d, pss)

	err := client.ModifyRequestLogProfile(name, config)

	if err != nil {
		log.Printf("[ERROR] Unable to Modify Request Log Profile  (%s) (%v)", name, err)
		return diag.FromErr(err)
	}

	return resourceBigipLtmProfileRequestLogRead(ctx, d, meta)
}

func resourceBigipLtmProfileRequestLogDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	name := d.Id()
	log.Println("[INFO] Deleting Request Log Profile " + name)
	err := client.DeleteRequestLogProfile(name)
	if err != nil {
		log.Printf("[ERROR] Unable to Delete Request Log Profile  (%s) (%v) ", name, err)
		return diag.FromErr(err)
	}
	d.SetId("")
	return nil
}

func getRequestLogProfileConfig(d *schema.ResourceData, config *bigip.RequestLogProfile) *bigip.RequestLogProfile {
	config.DefaultsFrom = d.Get("defaults_from").(string)
	config.Description = d.Get("description").(string)
	config.RequestLogPool = d.Get("requestlog_pool").(string)
	config.RequestLogProtocol = d.Get("requestlog_protocol").(string)
	config.RequestLogErrorPool = d.Get("requestlog_error_pool").(string)
	config.RequestLogErrorProtocol = d.Get("requestlog_error_protocol").(string)
	config.ResponseLogPool = d.Get("responselog_pool").(string)
	config.ResponseLogProtocol = d.Get("responselog_protocol").(string)
	config.ResponseLogErrorPool = d.Get("responselog_error_pool").(string)
	config.ResponseLogErrorProtocol = d.Get("responselog_error_protocol").(string)
	config.RequestLogging = d.Get("request_logging").(string)
	config.ResponseLogging = d.Get("response_logging").(string)
	config.RequestLogTemplate = strings.ReplaceAll(d.Get("requestlog_template").(string), `"`, `\"`)
	config.RequestLogErrorTemplate = d.Get("requestlog_error_template").(string)
	config.ResponseLogTemplate = strings.ReplaceAll(d.Get("responselog_template").(string), `"`, `\"`)
	config.ResponseLogErrorTemplate = d.Get("responselog_error_template").(string)
	config.ProxyResponse = d.Get("proxy_response").(string)
	config.ProxyCloseOnError = d.Get("proxyclose_on_error").(string)
	config.ProxyRespondOnLoggingError = d.Get("proxyrespond_on_loggingerror").(string)
	return config
}
