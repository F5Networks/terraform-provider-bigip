/*
Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"crypto/tls"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/structure"
	"io/ioutil"
	"log"
	"net/http"
	//"strings"

	"github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceBigipAs3() *schema.Resource {
	return &schema.Resource{
		Create: resourceBigipAs3Create,
		Read:   resourceBigipAs3Read,
		Update: resourceBigipAs3Update,
		Delete: resourceBigipAs3Delete,
		Exists: resourceBigipAs3Exists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"as3_json": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "AS3 json",
				StateFunc: func(v interface{}) string {
					json, _ := structure.NormalizeJsonString(v)
					return json
				},
				ValidateFunc: validation.ValidateJsonString,
			},
			"tenant_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of Tenant",
			},
		},
	}
}

func resourceBigipAs3Create(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	as3Json := d.Get("as3_json").(string)
	if ok := bigip.ValidateAS3Template(as3Json); !ok {
		return fmt.Errorf("[AS3] Error validating template \n")
		//return false
	}
	//strTrimSpace := strings.TrimSpace(as3Json)
	tenantList := client.GetTenantList(as3Json)
	strTrimSpace := client.AddTeemAgent(as3Json)
	var tenants string
	for i := 0; i < len(tenantList)-1; i++ {
		if i == 0 {
			tenants = tenantList[i+1]
			continue
		}
		tenants = tenants + "," + tenantList[i+1]
	}
	log.Printf("[INFO] Tenants in Json:%+v", tenants)
	log.Printf("[INFO] Creating as3 config in bigip:%s", strTrimSpace)
	err := client.PostAs3Bigip(strTrimSpace)
	if err != nil {
		return fmt.Errorf("Error creating json  %s: %v", tenants, err)
	}
	d.SetId(tenants)
	return resourceBigipAs3Read(d, meta)
}
func resourceBigipAs3Read(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	log.Printf("[INFO] Reading As3 config")
	name := d.Id()
	as3Resp, err := client.GetAs3(name)
	if err != nil {
		log.Printf("[ERROR] Unable to retrieve json ")
		return err
	}
	if as3Resp == "" {
		log.Printf("[WARN] Json (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}
	d.Set("as3_json", as3Resp)
	return nil
}

func resourceBigipAs3Exists(d *schema.ResourceData, meta interface{}) (bool, error) {
	client_bigip := meta.(*bigip.BigIP)
	log.Printf("[INFO] Checking if As3 config exists in bigip ")

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	client := &http.Client{Transport: tr}
	url := client_bigip.Host + "/mgmt/shared/appsvcs/declare"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("[ERROR] Error while creating http request for checking As3 config : %v", err)
		return false, err
	}
	req.SetBasicAuth(client_bigip.User, client_bigip.Password)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	body, err := ioutil.ReadAll(resp.Body)
	bodyString := string(body)
	if resp.Status == "204 No Content" || err != nil {
		log.Printf("[ERROR] Error while checking as3resource present in bigip :%s  %v", bodyString, err)
		defer resp.Body.Close()
		return false, err
	}
	defer resp.Body.Close()
	return true, nil
}

func resourceBigipAs3Update(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	as3Json := d.Get("as3_json").(string)
	log.Printf("[INFO] Updating As3 Config :%s", as3Json)
	name := d.Id()
	err := client.ModifyAs3(name, as3Json)
	if err != nil {
		return fmt.Errorf("Error modifying json %s: %v", name, err)
	}
	return resourceBigipAs3Read(d, meta)
}

func resourceBigipAs3Delete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	log.Printf("[INFO] Deleting As3 config")
	//	name := d.Get("tenant_name").(string)
	name := d.Id()
	err := client.DeleteAs3Bigip(name)
	if err != nil {
		log.Printf("[ERROR] Unable to Delete: %v :", err)
		return err
	}
	d.SetId("")
	return nil
}
