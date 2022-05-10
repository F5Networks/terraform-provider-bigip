/*
Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceBigipAwafPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceBigipAwafPolicyCreate,
		Read:   resourceBigipAwafPolicyRead,
		Update: resourceBigipAwafPolicyUpdate,
		Delete: resourceBigipAwafPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The unique user-given name of the policy. Policy names cannot contain spaces or special characters. Allowed characters are a-z, A-Z, 0-9, dot, dash (-), colon (:) and underscore (_)",
				ForceNew:    true,
			},
			"template_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Specifies the name of the template used for the policy creation.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the description of the policy.",
			},
			"application_language": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "utf-8",
				Description: "The character encoding for the web application. The character encoding determines how the policy processes the character sets. The default is Auto detect",
			},
			"case_insensitive": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Specifies whether the security policy treats microservice URLs, file types, URLs, and parameters as case sensitive or not. When this setting is enabled, the system stores these security policy elements in lowercase in the security policy configuration",
			},
			"enable_passivemode": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Passive Mode allows the policy to be associated with a Performance L4 Virtual Server (using a FastL4 profile). With FastL4, traffic is analyzed but is not modified in any way.",
			},
			"protocol_independent": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "When creating a security policy, you can determine whether a security policy differentiates between HTTP and HTTPS URLs. If enabled, the security policy differentiates between HTTP and HTTPS URLs. If disabled, the security policy configures URLs without specifying a specific protocol. This is useful for applications that behave the same for HTTP and HTTPS, and it keeps the security policy from including the same URL twice.",
			},
			"enforcement_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "How the system processes a request that triggers a security policy violation",
				ValidateFunc: validation.StringInSlice([]string{"blocking", "transparent"}, false),
			},
			"type": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The type of policy you want to create. The default policy type is Security.",
				ValidateFunc: validation.StringInSlice([]string{"parent", "security"}, false),
			},
			"server_technologies": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "The server technology is a server-side application, framework, web server or operating system type that is configured in the policy in order to adapt the policy to the checks needed for the respective technology.",
				Optional:    true,
			},
			"policy_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"parameters": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "This section defines parameters that the security policy permits in requests.",
			},
			"urls": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "In a security policy, you can manually specify the HTTP URLs that are allowed (or disallowed) in traffic to the web application being protected. If you are using automatic policy building (and the policy includes learning URLs), the system can determine which URLs to add, based on legitimate traffic.",
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

func resourceBigipAwafPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	name := d.Get("name").(string)

	log.Println("[INFO] AWAF Name " + name)

	config, err := getpolicyConfig(d)
	if err != nil {
		return fmt.Errorf("error in Json encode for waf policy %+v", err)
	}

	taskId, err := client.ImportAwafJson(name, config)
	log.Printf("[INFO] AWAF Import policy TaskID :%v", taskId)
	if err != nil {
		return fmt.Errorf("Error in Importing AWAF json (%s): %s ", name, err)
	}
	err = client.GetImportStatus(taskId)
	if err != nil {
		return fmt.Errorf("Error in Importing AWAF json (%s): %s ", name, err)
	}
	wafpolicy, err := client.GetWafPolicyQuery(name)
	if err != nil {
		return fmt.Errorf("error retrieving waf policy %+v: %v", wafpolicy, err)
	}
	d.SetId(wafpolicy.ID)
	_ = d.Set("policy_json", config)
	return resourceBigipAwafPolicyRead(d, meta)
}

func resourceBigipAwafPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	policyID := d.Id()
	name := d.Get("name").(string)

	log.Printf("[DEBUG] Reading AWAF Policy : %+v with ID: %+v", name, policyID)

	wafpolicy, err := client.GetWafPolicy(policyID)
	if err != nil {
		return fmt.Errorf("error retrieving waf policy %+v: %v", wafpolicy, err)
	}

	_ = d.Set("name", wafpolicy.FullPath)
	_ = d.Set("policy_id", wafpolicy.ID)

	return nil
}

func resourceBigipAwafPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	_ = d.Id()
	name := d.Get("name").(string)
	log.Printf("[DEBUG] Updating AWAF Policy : %+v", name)

	config, err := getpolicyConfig(d)
	if err != nil {
		return fmt.Errorf("error in Json encode for waf policy %+v", err)
	}

	taskId, err := client.ImportAwafJson(name, config)
	log.Printf("[INFO] AWAF Import policy TaskID :%v", taskId)
	if err != nil {
		return fmt.Errorf("Error in Importing AWAF json (%s): %s ", name, err)
	}
	err = client.GetImportStatus(taskId)
	if err != nil {
		return fmt.Errorf("Error in Importing AWAF json (%s): %s ", name, err)
	}
	wafpolicy, err := client.GetWafPolicyQuery(name)
	if err != nil {
		return fmt.Errorf("error retrieving waf policy %+v: %v", wafpolicy, err)
	}
	return resourceBigipAwafPolicyRead(d, meta)
}

func resourceBigipAwafPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	policyID := d.Id()
	name := d.Get("name").(string)
	log.Printf("[DEBUG] Deleting AWAF Policy : %+v with ID: %+v", name, policyID)

	err := client.DeleteWafPolicy(policyID)
	if err != nil {
		return fmt.Errorf(" Error Deleting AWAF Policy : %s", err)
	}
	d.SetId("")
	return nil
}

func getpolicyConfig(d *schema.ResourceData) (string, error) {
	name := d.Get("name").(string)
	policyWaf := &bigip.WafPolicy{
		Name:                name,
		ApplicationLanguage: d.Get("application_language").(string),
	}
	policyWaf.CaseInsensitive = d.Get("case_insensitive").(bool)
	policyWaf.EnablePassiveMode = d.Get("enable_passivemode").(bool)
	policyWaf.ProtocolIndependent = d.Get("protocol_independent").(bool)
	policyWaf.EnforcementMode = d.Get("enforcement_mode").(string)
	policyWaf.Type = d.Get("type").(string)
	policyWaf.Template = struct {
		Name string `json:"name,omitempty"`
	}{
		Name: d.Get("template_name").(string),
	}
	p := d.Get("server_technologies").([]interface{})

	var sts []struct {
		ServerTechnologyName string `json:"serverTechnologyName,omitempty"`
	}
	for i := 0; i < len(p); i++ {
		st1 := struct {
			ServerTechnologyName string `json:"serverTechnologyName,omitempty"`
		}{
			p[i].(string),
		}
		sts = append(sts, st1)
	}

	log.Printf("[INFO] URLS: %+v ", d.Get("urls"))

	var wafUrls []bigip.WafUrlJson
	urls := d.Get("urls").([]interface{})
	for i := 0; i < len(urls); i++ {
		var wafUrl bigip.WafUrlJson
		_ = json.Unmarshal([]byte(urls[i].(string)), &wafUrl)
		wafUrls = append(wafUrls, wafUrl)
	}
	policyWaf.Urls = wafUrls

	var wafParams []bigip.Parameter
	parmtrs := d.Get("parameters").([]interface{})
	for i := 0; i < len(parmtrs); i++ {
		var wafParam bigip.Parameter
		_ = json.Unmarshal([]byte(parmtrs[i].(string)), &wafParam)
		wafParams = append(wafParams, wafParam)
	}
	policyWaf.Parameters = wafParams

	policyWaf.ServerTechnologies = sts

	policyJson := struct {
		Policy interface{} `json:"policy"`
	}{
		policyWaf,
	}
	data, err := json.Marshal(policyJson)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
