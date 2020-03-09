/*
Original work from https://github.com/DealerDotCom/terraform-provider-bigip
Modifications Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file,You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"fmt"
	"github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
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
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the template",
				//ValidateFunc: validateF5Name,
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
				Type:     schema.TypeList,
				Required: true,
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
	tenantName := d.Get("tenant_name").(string)
	applicationName := d.Get("application_name").(string)
	virtualPort := d.Get("virtual_port").(int)
	virtualAddress := d.Get("virtual_address").(string)
	serverPort := d.Get("server_port").(int)
	var serverAddresses []string
	if m, ok := d.GetOk("server_addresses"); ok {
		for _, serverAddress := range m.([]interface{}) {
			serverAddresses = append(serverAddresses, serverAddress.(string))
		}
	}
	log.Println("[INFO] Creating fast template")
	temParameters := bigip.FastParameters{
		TenantName:      tenantName,
		ApplicationName: applicationName,
		VirtualPort:     virtualPort,
		VirtualAddress:  virtualAddress,
		ServerPort:      serverPort,
		ServerAddresses: serverAddresses,
	}
	template := &bigip.Fasttemplate{
		Name:       name,
		Parameters: temParameters,
	}

	log.Printf("[INFO] Template Before Post Call:%+v", template)
	err := client.CreateFastTemplate(template)
	if err != nil {
		return fmt.Errorf("Error Creating template %s: %v", name, err)
	}
	d.SetId(tenantName)
	return resourceBigipfasthttpRead(d, meta)
	//return nil
}

func resourceBigipfasthttpRead(d *schema.ResourceData, meta interface{}) error {
	//client := meta.(*bigip.BigIP)
	name := d.Id()
	log.Println("[INFO] Fetching node " + name)
	return nil
}

func resourceBigipfasthttpUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceBigipfasthttpCreate(d, meta)
}

func resourceBigipfasthttpDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	tenantName := d.Get("tenant_name").(string)
	applicationName := d.Get("application_name").(string)
	log.Println("[INFO] Deleting fast template application in tenant %s %s ", tenantName, applicationName)

	err := client.DeleteFastTemplate(tenantName, applicationName)
	if err != nil {
		log.Printf("[ERROR] Unable to delete fast template application (%s) (%v) ", applicationName, err)
		return err
	}
	d.SetId("")
	return nil
}
