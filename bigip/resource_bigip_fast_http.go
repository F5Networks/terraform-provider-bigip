/*
Original work from https://github.com/DealerDotCom/terraform-provider-bigip
Modifications Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file,You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"fmt"
	"log"

	"github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceBigipfasthttp() *schema.Resource {
	return &schema.Resource{
		Create: resourceBigipfasthttpCreate,
		Read:   resourceBigipfasthttpRead,
		Update: resourceBigipfasthttpUpdate,
		Delete: resourceBigipfasthttpDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Name of the template",
				ValidateFunc: validateF5Name,
			},
                        "tenant_name": {
                                Type:        schema.TypeString,
                                Required:    true,
                                Description: "Name of the tenant",
                        },
			"application_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the application",
			},
			"virtual_port": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Specifies virtual port",
			},
			"virtual_address": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies virtual address",
			},
			"server_port": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Specifies server port ",
			},
			"server_addresses": {
				Type:        schema.TypeList,
				Required:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					},
			},
		},
	}
}

func resourceBigipfasthttpCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Get("name").(string)
	tenant_name := d.Get("tenant_name").(string)
	application_name := d.Get("address").(string)
	virtual_port := d.Get("virtual_port").(int)
	virtual_address := d.Get("virtual_address").(string)
	server_port := d.Get("server_port").(int)
	server_addresses:= d.Get("server_addresses").(string)

	log.Println("[INFO] Creating fast template")
      
        var template *bigip.fathttp
        template = &bigip.fasthttp{
                Name:               name,
		Parameters struct {
                TenantName               string `json:"tenant_name,omitempty"`
                ApplicationName          string `json:"application_name,omitempty"`
                VirtualPort              int `json:"virtual_port,omitempty"`
                VirtualAddress            string `json:"virtual_address,omitempty"`
                ServerPort                int   `json:"server_port,omitempty"`
                ServerAddresses          []string `json:"server_addresses,omitempty"`
        } 

	err = client.CreateFastTemplate(template)

	if err != nil {
		return fmt.Errorf("Error Creating template %s: %v", name, err)
	}

	d.SetId(name)

	return nil
}

func resourceBigipfasthttpRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}


func resourceBigipfasthttpUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceBigipfasthttpDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
