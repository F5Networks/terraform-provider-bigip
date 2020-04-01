/*
Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"crypto/tls"
	"fmt"
	"github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
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
				Type:         schema.TypeString,
				Required:     true,
				Description:  "AS3 json",
				ValidateFunc: validation.ValidateJsonString,
			},
			"tenant_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of Tenant",
			},
		},
	}
}

func resourceBigipAs3Create(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	as3_json := d.Get("as3_json").(string)

	strTrimSpace := strings.TrimSpace(as3_json)
	name := d.Get("tenant_name").(string)
	exmp := client.GetTenantList(as3_json)
	log.Println(exmp)
	log.Printf("[INFO] Creating as3 config in bigip:%s", strTrimSpace)
	err := client.PostAs3Bigip(strTrimSpace)
	if err != nil {
		return fmt.Errorf("Error modifying node %s: %v", name, err)
	}
	d.SetId(name)
	return resourceBigipAs3Read(d, meta)
}
func resourceBigipAs3Read(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	log.Printf("[INFO] Reading As3 config")
	name := d.Id()
	as3exmp, err := client.GetAs3(name)
	if err != nil {
		log.Printf("[ERROR] Unable to retrieve node ")
		return err
	}
	if as3exmp == "" {
		log.Printf("[WARN] Node (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	const s = `{"class":"AS3","action":"deploy","persist":true,"declaration":`
	const s1 = `}`
	as3exmp = s + as3exmp + s1
	d.Set("as3_json", as3exmp)
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
	as3_json := d.Get("as3_json").(string)
	log.Printf("[INFO] Updating As3 Config :%s", as3_json)
	name := d.Id()
	err := client.ModifyAs3(name, as3_json)
	if err != nil {
		return fmt.Errorf("Error modifying node %s: %v", name, err)
	}
	return resourceBigipAs3Read(d, meta)
}

func resourceBigipAs3Delete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	log.Printf("[INFO] Deleting As3 config")
	name := d.Get("tenant_name").(string)
	err := client.DeleteAs3Bigip(name)
	if err != nil {
		log.Printf("[ERROR] Unable to Delete: %v :", err)
		return err
	}
	d.SetId("")
	return nil
}
