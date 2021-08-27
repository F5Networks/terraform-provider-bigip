/*
Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"log"
	"os"
	"strings"

	"github.com/f5devcentral/go-bigip"
	"github.com/f5devcentral/go-bigip/f5teem"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceBigipLtmProfileHttp() *schema.Resource {
	return &schema.Resource{
		Create: resourceBigipLtmProfileHttpCreate,
		Read:   resourceBigipLtmProfileHttpRead,
		Update: resourceBigipLtmProfileHttpUpdate,
		Delete: resourceBigipLtmProfileHttpDelete,
		Exists: resourceBigipLtmProfileHttpExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "Name of the profile",
				ValidateFunc: validateF5NameWithDirectory,
			},
			"proxy_type": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "Specifies the proxy mode for this profile: reverse, explicit, or transparent. The default is Reverse.",
			},
			"defaults_from": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  "Inherit defaults from parent profile",
				ValidateFunc: validateF5Name,
			},
			"app_service": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The application service to which the object belongs.",
			},
			"basic_auth_realm": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies a quoted string for the basic authentication realm. The system sends this string to a client whenever authorization fails. The default value is none",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "User defined description",
			},
			"encrypt_cookies": {
				Type:        schema.TypeSet,
				Set:         schema.HashString,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "Encrypts specified cookies that the BIG-IP system sends to a client system",
			},
			"encrypt_cookie_secret": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies a passphrase for the cookie encryption",
			},
			"fallback_host": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies an HTTP fallback host. HTTP redirection allows you to redirect HTTP traffic to another protocol identifier, host name, port number, or URI path.",
			},
			"fallback_status_codes": {
				Type:        schema.TypeSet,
				Set:         schema.HashString,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "Specifies one or more three-digit status codes that can be returned by an HTTP server,that should trigger a redirection to the fallback host",
			},
			"head_erase": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the header string that you want to erase from an HTTP request. Default is none",
			},
			"head_insert": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies a quoted header string that you want to insert into an HTTP request. Default is none",
			},
			"insert_xforwarded_for": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies, when enabled, that the system inserts an X-Forwarded-For header in an HTTP request with the client IP address, to use with connection pooling. The default is Disabled.",
			},
			"lws_width": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the maximum column width for any given line, when inserting an HTTP header in an HTTP request. The default is 80",
			},
			"lws_separator": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the linear white space (LWS) separator that the system inserts when a header exceeds the maximum width you specify in the LWS Maximum Columns setting.",
			},
			"accept_xff": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "Enables or disables trusting the client IP address, and statistics from the client IP address, based on the request's XFF (X-forwarded-for) headers, if they exist.",
			},
			"oneconnect_transformations": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Enables the system to perform HTTP header transformations for the purpose of keeping server-side connections open. This feature requires configuration of a OneConnect profile.",
			},
			"tm_partition": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Displays the administrative partition within which this profile resides. ",
			},
			"redirect_rewrite": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies whether the system rewrites the URIs that are part of HTTP redirect (3XX) responses. The default is None",
			},
			"response_headers_permitted": {
				Type:        schema.TypeSet,
				Set:         schema.HashString,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Computed:    true,
				Description: "Specifies headers that the BIG-IP system allows in an HTTP response.If you are specifying more than one header, separate the headers with a blank space",
			},
			"request_chunking": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies how the system handles HTTP content that is chunked by a client. The default is Preserve",
			},
			"response_chunking": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies how the system handles HTTP content that is chunked by a server. The default is Selective",
			},
			"server_agent_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the value of the Server header in responses that the BIG-IP itself generates. The default is BigIP. If no string is specified, then no Server header will be added to such responses",
			},
			"via_host_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "Specifies the hostname to include into Via header",
			},
			"via_request": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies whether to append, remove, or preserve a Via header in an HTTP request",
			},
			"via_response": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies whether to append, remove, or preserve a Via header in an HTTP request",
			},
			"xff_alternative_names": {
				Type:        schema.TypeSet,
				Set:         schema.HashString,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Computed:    true,
				Description: "Specifies alternative XFF headers instead of the default X-forwarded-for header",
			},
		},
	}
}

func resourceBigipLtmProfileHttpCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Get("name").(string)
	log.Printf("[INFO] Creating HTTP Profile:%+v ", name)

	pss := &bigip.HttpProfile{
		Name: name,
	}
	config := getHttpProfileConfig(d, pss)

	err := client.AddHttpProfile(config)
	if err != nil {
		return err
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
		err = teemDevice.Report(f, "bigip_ltm_profile_http", tsVer[3])
		if err != nil {
			log.Printf("[ERROR]Sending Telemetry data failed:%v", err)
		}
	}
	return resourceBigipLtmProfileHttpRead(d, meta)
}

func resourceBigipLtmProfileHttpRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	log.Println("[INFO] Fetching HTTP  Profile " + name)

	pp, err := client.GetHttpProfile(name)
	if err != nil {
		log.Printf("[ERROR] Unable to retrieve HTTP Profile  (%s) ", err)
		return err
	}
	if pp == nil {
		log.Printf("[WARN] HTTP  Profile (%s) not found, removing from state", name)
		d.SetId("")
		return nil
	}
	_ = d.Set("name", name)
	_ = d.Set("defaults_from", pp.DefaultsFrom)
	_ = d.Set("proxy_type", pp.ProxyType)

	if _, ok := d.GetOk("accept_xff"); ok {
		_ = d.Set("accept_xff", pp.AcceptXff)
	}
	if _, ok := d.GetOk("basic_auth_realm"); ok {
		_ = d.Set("basic_auth_realm", pp.BasicAuthRealm)
	}
	if _, ok := d.GetOk("description"); ok {
		_ = d.Set("description", pp.Description)
	}
	if _, ok := d.GetOk("encrypt_cookie_secret"); ok {
		_ = d.Set("encrypt_cookie_secret", pp.EncryptCookieSecret)
	}
	if _, ok := d.GetOk("encrypt_cookies"); ok {
		_ = d.Set("encrypt_cookies", pp.EncryptCookies)
	}
	if _, ok := d.GetOk("fallback_host"); ok {
		_ = d.Set("fallback_host", pp.FallbackHost)
	}
	if _, ok := d.GetOk("fallback_status_codes"); ok {
		_ = d.Set("fallback_status_codes", pp.FallbackStatusCodes)
	}
	if _, ok := d.GetOk("head_erase"); ok {
		_ = d.Set("head_erase", pp.HeaderErase)
	}
	if _, ok := d.GetOk("head_insert"); ok {
		_ = d.Set("head_insert", pp.HeaderInsert)
	}
	if _, ok := d.GetOk("insert_xforwarded_for"); ok {
		_ = d.Set("insert_xforwarded_for", pp.InsertXforwardedFor)
	}
	if _, ok := d.GetOk("lws_separator"); ok {
		_ = d.Set("lws_separator", pp.LwsSeparator)
	}
	if _, ok := d.GetOk("oneconnect_transformations"); ok {
		_ = d.Set("oneconnect_transformations", pp.OneconnectTransformations)
	}
	if _, ok := d.GetOk("tm_partition"); ok {
		_ = d.Set("tm_partition", pp.TmPartition)
	}
	if _, ok := d.GetOk("redirect_rewrite"); ok {
		_ = d.Set("redirect_rewrite", pp.RedirectRewrite)
	}
	if _, ok := d.GetOk("request_chunking"); ok {
		_ = d.Set("request_chunking", pp.RequestChunking)
	}
	if _, ok := d.GetOk("response_chunking"); ok {
		_ = d.Set("response_chunking", pp.ResponseChunking)
	}
	log.Printf("[DEBUG] response_headers_permitted:%+v", pp.ResponseHeadersPermitted)
	_ = d.Set("response_headers_permitted", pp.ResponseHeadersPermitted)

	if _, ok := d.GetOk("server_agent_name"); ok {
		_ = d.Set("server_agent_name", pp.ServerAgentName)
	}
	if _, ok := d.GetOk("via_host_name"); ok {
		_ = d.Set("via_host_name", pp.ViaHostName)
	}
	if _, ok := d.GetOk("via_request"); ok {
		_ = d.Set("via_request", pp.ViaRequest)
	}
	if _, ok := d.GetOk("via_response"); ok {
		_ = d.Set("via_response", pp.ViaResponse)
	}
	log.Printf("[DEBUG] xff_alternative_names:%+v", pp.XffAlternativeNames)
	_ = d.Set("xff_alternative_names", pp.XffAlternativeNames)

	return nil
}

func resourceBigipLtmProfileHttpUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	name := d.Id()
	log.Printf("[INFO] Updating HTTP Profile Profile:%+v ", name)

	pss := &bigip.HttpProfile{
		Name: name,
	}
	config := getHttpProfileConfig(d, pss)
	err := client.ModifyHttpProfile(name, config)

	if err != nil {
		log.Printf("[ERROR] Unable to Modify HTTP Profile  (%s) (%v)", name, err)
		return err
	}

	return resourceBigipLtmProfileHttpRead(d, meta)
}

func resourceBigipLtmProfileHttpDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()
	log.Println("[INFO] Deleting HTTPProfile " + name)
	err := client.DeleteHttpProfile(name)
	if err != nil {
		log.Printf("[ERROR] Unable to Delete HTTPProfile  (%s) (%v) ", name, err)
		return err
	}
	d.SetId("")
	return nil
}

func resourceBigipLtmProfileHttpExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	client := meta.(*bigip.BigIP)

	name := d.Id()
	log.Println("[INFO] Fetching HTTPProfile " + name)
	pp, err := client.GetHttpProfile(name)
	if err != nil {
		log.Printf("[ERROR] Unable to retrieve HTTPProfile (%s) (%v) ", name, err)
		return false, err
	}

	if pp == nil {
		log.Printf("[WARN] HTTP Profile  (%s) not found, removing from state", d.Id())
		d.SetId("")
	}

	return pp != nil, nil
}

func getHttpProfileConfig(d *schema.ResourceData, config *bigip.HttpProfile) *bigip.HttpProfile {
	config.AppService = d.Get("app_service").(string)
	config.DefaultsFrom = d.Get("defaults_from").(string)
	config.AcceptXff = d.Get("accept_xff").(string)
	config.BasicAuthRealm = d.Get("basic_auth_realm").(string)
	config.Description = d.Get("description").(string)
	config.EncryptCookieSecret = d.Get("encrypt_cookie_secret").(string)
	config.EncryptCookies = setToStringSlice(d.Get("encrypt_cookies").(*schema.Set))
	config.FallbackHost = d.Get("fallback_host").(string)
	config.FallbackStatusCodes = setToStringSlice(d.Get("fallback_status_codes").(*schema.Set))
	config.HeaderErase = d.Get("head_erase").(string)
	config.HeaderInsert = d.Get("head_insert").(string)
	config.InsertXforwardedFor = d.Get("insert_xforwarded_for").(string)
	config.LwsSeparator = d.Get("lws_separator").(string)
	config.OneconnectTransformations = d.Get("oneconnect_transformations").(string)
	config.TmPartition = d.Get("tm_partition").(string)
	config.ProxyType = d.Get("proxy_type").(string)
	config.RedirectRewrite = d.Get("redirect_rewrite").(string)
	config.RequestChunking = d.Get("request_chunking").(string)
	config.ResponseChunking = d.Get("response_chunking").(string)
	config.ResponseHeadersPermitted = setToInterfaceSlice(d.Get("response_headers_permitted").(*schema.Set))
	config.ServerAgentName = d.Get("server_agent_name").(string)
	config.ViaHostName = d.Get("via_host_name").(string)
	config.ViaRequest = d.Get("via_request").(string)
	config.ViaResponse = d.Get("via_response").(string)
	config.XffAlternativeNames = setToInterfaceSlice(d.Get("xff_alternative_names").(*schema.Set))

	return config
}
