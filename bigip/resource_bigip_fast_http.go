/*
Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"fmt"
	"github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
	//"time"
	"sync"
)

var x1 = 0
var m1 sync.Mutex

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
	m1.Lock()
	defer m1.Unlock()
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
	log.Println("[INFO] Creating Application through FAST template")
	temParameters := &bigip.FastParameters{
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
	err := client.CreateFastTemplate(template)
	if err != nil {
		return fmt.Errorf("Error Creating template %s: %v", name, err)
	}
	d.SetId(name)
	x1 = x1 + 1
	return resourceBigipfasthttpRead(d, meta)
}

func resourceBigipfasthttpRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	name := d.Id()
	tenantName := d.Get("tenant_name").(string)
	applicationName := d.Get("application_name").(string)
	log.Printf("[INFO] Reading Application through FAST :%v\t %v\t %v", name, tenantName, applicationName)
	fastApp, err := client.GetFastTemplate(tenantName, applicationName)
	if err != nil {
		log.Printf("[ERROR] Unable to delete fast template application (%s) (%v) ", applicationName, err)
		return err
	}

	if err := d.Set("name", fastApp.Name); err != nil {
		return fmt.Errorf("[DEBUG] Error saving template name to FAST State  (%s): %s", d.Id(), err)
	}
	if err := d.Set("tenant_name", fastApp.Parameters.TenantName); err != nil {
		return fmt.Errorf("[DEBUG] Error saving Tenant name to FAST State  (%s): %s", d.Id(), err)
	}
	if err := d.Set("application_name", fastApp.Parameters.ApplicationName); err != nil {
		return fmt.Errorf("[DEBUG] Error saving application name to FAST State  (%s): %s", d.Id(), err)
	}
	if err := d.Set("virtual_port", fastApp.Parameters.VirtualPort); err != nil {
		return fmt.Errorf("[DEBUG] Error saving VirtualPort to FAST State  (%s): %s", d.Id(), err)
	}
	if err := d.Set("virtual_address", fastApp.Parameters.VirtualAddress); err != nil {
		return fmt.Errorf("[DEBUG] Error saving VirtualAddress to FAST State  (%s): %s", d.Id(), err)
	}
	if err := d.Set("server_port", fastApp.Parameters.ServerPort); err != nil {
		return fmt.Errorf("[DEBUG] Error saving ServerPort to FAST State  (%s): %s", d.Id(), err)
	}
	log.Printf("[INFO]server_addresses in FAST read:%+v", fastApp.Parameters.ServerAddresses)
	if err := d.Set("server_addresses", fastApp.Parameters.ServerAddresses); err != nil {
		return fmt.Errorf("[DEBUG] Error saving ServerAddresses to FAST State  (%s): %s", d.Id(), err)
	}
	return nil
}

func resourceBigipfasthttpUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	m1.Lock()
	defer m1.Unlock()
	name := d.Id()
	log.Println("Updating Application through Fast Template")
	//name := d.Get("name").(string)
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
	temParameters := &bigip.FastParameters{
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
	err := client.CreateFastTemplate(template)
	if err != nil {
		return fmt.Errorf("Error Creating template %s: %v", name, err)
	}
	x1 = x1 + 1
	return resourceBigipfasthttpRead(d, meta)
}

func resourceBigipfasthttpDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	m1.Lock()
	defer m1.Unlock()
	tenantName := d.Get("tenant_name").(string)
	applicationName := d.Get("application_name").(string)
	log.Printf("[INFO] Deleting Fast application: %v \t in tenant :%v", tenantName, applicationName)

	err := client.DeleteFastTemplate(tenantName, applicationName)
	if err != nil {
		log.Printf("[ERROR] Unable to delete fast template application (%s) (%v) ", applicationName, err)
		return err
	}
	x1 = x1 + 1
	d.SetId("")
	return nil
}
