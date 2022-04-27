/*
Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"log"

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
				Description: "The character encoding for the web application. The character encoding determines how the policy processes the character sets. The default is Auto detect",
			},
			"case_insensitive": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Specifies whether the security policy treats microservice URLs, file types, URLs, and parameters as case sensitive or not. When this setting is enabled, the system stores these security policy elements in lowercase in the security policy configuration",
			},
			"enable_passivemode": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Passive Mode allows the policy to be associated with a Performance L4 Virtual Server (using a FastL4 profile). With FastL4, traffic is analyzed but is not modified in any way.",
			},
			"enforcement_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "How the system processes a request that triggers a security policy violation",
				ValidateFunc: validation.StringInSlice([]string{"blocking", "transparent"}, false),
			},
			"server_technologies": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
			"policy_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"content": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
				//ForceNew:    true,
				Description: "Content of certificate on Disk",
			},
			//"partition": {
			//	Type:         schema.TypeString,
			//	Optional:     true,
			//	Default:      "Common",
			//	Description:  "Partition of ssl certificate",
			//	ValidateFunc: validatePartitionName,
			//},
		},
	}
}

func resourceBigipAwafPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	name := d.Get("name").(string)

	log.Println("[INFO] AWAF Name " + name)

	certPath := d.Get("content").(string)

	err := client.ImportAwafJson(name, certPath)
	if err != nil {
		return fmt.Errorf("Error in Importing AWAF json (%s): %s ", name, err)
	}
	d.SetId(name)

	return resourceBigipAwafPolicyRead(d, meta)
}

func resourceBigipAwafPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	name := d.Id()
	//name = "policy-rapid-deployment"
	log.Printf("[DEBUG] Reading AWF Policy: %+v", name)

	wafpolicy, err := client.GetWafPolicy(name)
	if err != nil {
		return fmt.Errorf("error retrieving waf policy %+v: %v", wafpolicy, err)
	}

	_ = d.Set("name", wafpolicy.Name)
	_ = d.Set("description", wafpolicy.Description)
	_ = d.Set("policy_id", wafpolicy.ID)

	return nil
}

func resourceBigipAwafPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	//client := meta.(*bigip.BigIP)
	name := d.Id()
	log.Printf("[DEBUG] Updating AWAF Policy : %+v", name)
	return resourceBigipAwafPolicyRead(d, meta)
}

func resourceBigipAwafPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	name := d.Id()
	log.Printf("[DEBUG] Deleting AWAF Policy : %+v", name)

	err := client.DeleteVlan(name)
	if err != nil {
		return fmt.Errorf(" Error Deleting Vlan : %s", err)
	}

	d.SetId("")
	return nil
}
