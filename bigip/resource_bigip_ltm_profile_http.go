/*
Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
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
				Description:  "Name of the profile",
				ValidateFunc: validateF5Name,
			},
			"defaults_from": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  "Inherit defaults from parent profile",
				ValidateFunc: validateF5Name,
			},
			"accept_xff": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "Enables or disables trusting the client IP address, and statistics from the client IP address, based on the request's XFF (X-forwarded-for) headers, if they exist.",
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
			"encrypt_cookie_secret": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies a passphrase for the cookie encryption",
			},
			"encrypt_cookies": {
				Type:        schema.TypeSet,
				Set:         schema.HashString,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "Encrypts specified cookies that the BIG-IP system sends to a client system",
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
				Description: "Specifies one or more three-digit status codes that can be returned by an HTTP server.",
			},
			"head_erase": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the header string that you want to erase from an HTTP request. You can also specify none",
			},
			"head_insert": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies a quoted header string that you want to insert into an HTTP request. You can also specify none. ",
			},
			"insert_xforwarded_for": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: "When using connection pooling, which allows clients to make use of other client requests' server-side connections,	you can insert the X-Forwarded-For header and specify a client IP address. ",
			},
			"lws_separator": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies a quoted header string that you want to insert into an HTTP request. You can also specify none. ",
			},
			"oneconnect_transformations": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Enables the system to perform HTTP header transformations for the purpose of  keeping server-side connections open. This feature requires configuration of a OneConnect profile.",
			},
			"tm_partition": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Displays the administrative partition within which this profile resides. ",
			},
			"proxy_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the type of HTTP proxy. ",
			},
			"redirect_rewrite": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies which of the application HTTP redirects the system rewrites to HTTPS.",
			},
			"request_chunking": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies how to handle chunked and unchunked requests.",
			},
			"response_chunking": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies how to handle chunked and unchunked responses.",
			},
			"response_headers_permitted": {
				Type:        schema.TypeSet,
				Set:         schema.HashString,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Computed:    true,
				Description: "Specifies headers that the BIG-IP system allows in an HTTP response.",
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
	parent := d.Get("defaults_from").(string)

	err := client.CreateHttpProfile(
		name,
		parent,
	)
	if err != nil {
		return err
	}
	d.SetId(name)

	err = resourceBigipLtmProfileHttpUpdate(d, meta)
	if err != nil {
		_ = client.DeleteHttpProfile(name)
		return err
	}

	return resourceBigipLtmProfileHttpRead(d, meta)

}

func resourceBigipLtmProfileHttpRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	log.Println("[INFO] Fetching HTTP  Profile " + name)

	pp, err := client.GetHttpProfile(name)
	if err != nil {
		log.Printf("[ERROR] Unable to retrive HTTP Profile  (%s) ", err)
		return err
	}
	if pp == nil {
		log.Printf("[WARN] HTTP  Profile (%s) not found, removing from state", name)
		d.SetId("")
		return nil
	}
	_ = d.Set("name", name)
	_ = d.Set("defaults_from", pp.DefaultsFrom)

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
	if _, ok := d.GetOk("proxy_type"); ok {
		_ = d.Set("proxy_type", pp.ProxyType)
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
	if _, ok := d.GetOk("response_headers_permitted"); ok {
		_ = d.Set("response_headers_permitted", pp.ResponseHeadersPermitted)
	}
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
	if _, ok := d.GetOk("xff_alternative_names"); ok {
		_ = d.Set("xff_alternative_names", pp.XffAlternativeNames)
	}
	return nil
}

func resourceBigipLtmProfileHttpUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	pp := &bigip.HttpProfile{
		AppService:                d.Get("app_service").(string),
		DefaultsFrom:              d.Get("defaults_from").(string),
		AcceptXff:                 d.Get("accept_xff").(string),
		BasicAuthRealm:            d.Get("basic_auth_realm").(string),
		Description:               d.Get("description").(string),
		EncryptCookieSecret:       d.Get("encrypt_cookie_secret").(string),
		EncryptCookies:            setToStringSlice(d.Get("encrypt_cookies").(*schema.Set)),
		FallbackHost:              d.Get("fallback_host").(string),
		FallbackStatusCodes:       setToStringSlice(d.Get("fallback_status_codes").(*schema.Set)),
		HeaderErase:               d.Get("head_erase").(string),
		HeaderInsert:              d.Get("head_insert").(string),
		InsertXforwardedFor:       d.Get("insert_xforwarded_for").(string),
		LwsSeparator:              d.Get("lws_separator").(string),
		OneconnectTransformations: d.Get("oneconnect_transformations").(string),
		TmPartition:               d.Get("tm_partition").(string),
		ProxyType:                 d.Get("proxy_type").(string),
		RedirectRewrite:           d.Get("redirect_rewrite").(string),
		RequestChunking:           d.Get("request_chunking").(string),
		ResponseChunking:          d.Get("response_chunking").(string),
		ResponseHeadersPermitted:  setToStringSlice(d.Get("response_headers_permitted").(*schema.Set)),
		ServerAgentName:           d.Get("server_agent_name").(string),
		ViaHostName:               d.Get("via_host_name").(string),
		ViaRequest:                d.Get("via_request").(string),
		ViaResponse:               d.Get("via_response").(string),
		XffAlternativeNames:       setToStringSlice(d.Get("xff_alternative_names").(*schema.Set)),
	}

	err := client.ModifyHttpProfile(name, pp)
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
		log.Printf("[ERROR] Unable to retrive HTTPProfile (%s) (%v) ", name, err)
		return false, err
	}

	if pp == nil {
		log.Printf("[WARN] HTTP Profile  (%s) not found, removing from state", d.Id())
		d.SetId("")
	}

	return pp != nil, nil
}
