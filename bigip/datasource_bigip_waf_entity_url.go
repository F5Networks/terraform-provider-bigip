/*
Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"encoding/json"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"log"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func datasourceBigipWafEntityUrl() *schema.Resource {
	return &schema.Resource{
		Read: datasourceBigipWafEntityUrlRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies an HTTP URL that the security policy allows.",
				ForceNew:    true,
			},
			"protocol": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "http",
				Description:  "Specifies the name of the template used for the policy creation.",
				ValidateFunc: validation.StringInSlice([]string{"http", "https"}, false),
			},
			"type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "wildcard",
				Description:  "Determines the type of the name attribute. Only when setting the type to wildcard will the special wildcard characters in the name be interpreted as such.",
				ValidateFunc: validation.StringInSlice([]string{"explicit", "wildcard"}, false),
			},
			"method": {
				Type:     schema.TypeString,
				Optional: true,
				//Default:     "utf-8",
				Description: "Unique ID of a URL with a protocol type and name. Select a Method for the URL to create an API endpoint: URL + Method",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Describes the URL (optional).",
			},
			"attacksignatures_check": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Specifies, when true, that you want attack signatures and threat campaigns to be detected on this URL and possibly override the security policy settings of an attack signature or threat campaign specifically for this URL. After you enable this setting, the system displays a list of attack signatures and threat campaigns.",
			},
			"clickjacking_protection": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Specifies that the system adds the X-Frame-Options header to the domain URL’s response header. This is done to protect the web application against clickjacking. Clickjacking occurs when an attacker lures a user to click illegitimate frames and iframes because the attacker hid them on legitimate visible website buttons. Therefore, enabling this option protects the web application from other web sites hiding malicious code behind them. The default is disabled. After you enable this option, you can select whether, and under what conditions, the browser should allow this URL to be rendered in a frame or iframe.",
			},
			"methodsoverrideon_urlcheck": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Specifies, when true, that you want methods to be detected on this URL and possibly override the security policy settings of a method specifically for this URL. After you enable this setting, the system displays a list of methods.",
			},
			"metachars_onurlcheck": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Specifies, when true, that you want meta characters to be detected on this URL and possibly override the security policy settings of a meta character specifically for this URL. After you enable this setting, the system displays a list of meta characters.",
			},
			"wildcard_includes_slash": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Specifies that an asterisk in a wildcard URL matches any number of path segments (separated by slashes); when cleared, specifies that an asterisk matches at most one segment.",
			},
			"wildcard_order": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Specifies the order index for wildcard URLs matching. Wildcard URLs with lower wildcard order will get checked for a match prior to URLs with higher wildcard order.",
			},
			"policy_json": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The payload of the WAF Policy",
			},
		},
	}
}

func datasourceBigipWafEntityUrlRead(d *schema.ResourceData, meta interface{}) error {
	_ = meta.(*bigip.BigIP)
	d.SetId("")
	name := d.Get("name").(string)

	log.Println("[INFO] URL Name " + name)

	policyUrl := &bigip.WafEntityUrl{
		Name: name,
	}
	policyUrl.Protocol = d.Get("protocol").(string)
	policyUrl.Type = d.Get("type").(string)
	policyUrl.Method = d.Get("method").(string)

	data, err := json.Marshal(policyUrl)
	if err != nil {
		return err
	}
	log.Printf("[INFO] AWAF Policy URL Json struct: %+v ", string(data))

	_ = d.Set("policy_json", string(data))
	d.SetId(name)
	return nil
}
