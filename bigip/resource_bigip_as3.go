/*
Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
			},
			"config_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "unique identifier for AS3 resource",
			},
		},
	}
}

func resourceBigipAs3Create(d *schema.ResourceData, meta interface{}) error {
	client_bigip := meta.(*bigip.BigIP)

	as3_json := d.Get("as3_json").(string)
	name := d.Get("config_name").(string)
	log.Printf("[INFO] Creating as3 config in bigip:%s", as3_json)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	client := &http.Client{Transport: tr}
	url := client_bigip.Host + "/mgmt/shared/appsvcs/declare"
	req, err := http.NewRequest("POST", url, strings.NewReader(as3_json))
	if err != nil {
		return fmt.Errorf("Error while creating http request with AS3 json:%v", err)
	}
	req.SetBasicAuth(client_bigip.User, client_bigip.Password)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	body, err := ioutil.ReadAll(resp.Body)
	bodyString := string(body)
	if resp.Status != "200 OK" || err != nil {
		defer resp.Body.Close()
		return fmt.Errorf("Error while Sending/Posting http request with AS3 json :%s  %v", bodyString, err)
	}

	defer resp.Body.Close()
	d.SetId(name)
	return resourceBigipAs3Read(d, meta)
}
func resourceBigipAs3Read(d *schema.ResourceData, meta interface{}) error {
	client_bigip := meta.(*bigip.BigIP)
	log.Printf("[INFO] Reading As3 config")

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	client := &http.Client{Transport: tr}
	url := client_bigip.Host + "/mgmt/shared/appsvcs/declare"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("Error while creating http request for reading As3 config:%v", err)
	}
	req.SetBasicAuth(client_bigip.User, client_bigip.Password)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	body, err := ioutil.ReadAll(resp.Body)
	bodyString := string(body)
	if resp.Status != "200 OK" || err != nil {
		defer resp.Body.Close()
		return fmt.Errorf("Error while Sending/fetching http request :%s  %v", bodyString, err)
	}

	defer resp.Body.Close()
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
	client_bigip := meta.(*bigip.BigIP)
	as3_json := d.Get("as3_json").(string)
	log.Printf("[INFO] Updating As3 Config :%s", as3_json)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	client := &http.Client{Transport: tr}
	url := client_bigip.Host + "/mgmt/shared/appsvcs/declare"
	req, err := http.NewRequest("PATCH", url, strings.NewReader(as3_json))
	if err != nil {
		return fmt.Errorf("Error while creating http request with AS3 json:%v", err)
	}
	req.SetBasicAuth(client_bigip.User, client_bigip.Password)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	body, err := ioutil.ReadAll(resp.Body)
	bodyString := string(body)
	if resp.Status != "200 OK" || err != nil {
		return fmt.Errorf("Error while Sending/Posting http request with AS3 json :%s  %v", bodyString, err)
	}

	defer resp.Body.Close()
	return resourceBigipAs3Read(d, meta)
}

func resourceBigipAs3Delete(d *schema.ResourceData, meta interface{}) error {
	client_bigip := meta.(*bigip.BigIP)
	log.Printf("[INFO] Deleting As3 config")

	name := d.Get("config_name").(string)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	client := &http.Client{Transport: tr}
	url := client_bigip.Host + "/mgmt/shared/appsvcs/declare/" + name
	req, err := http.NewRequest("DELETE", url, nil)

	if err != nil {
		return fmt.Errorf("Error while creating http request for deleting as3 config:%v", err)
	}
	req.SetBasicAuth(client_bigip.User, client_bigip.Password)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	body, err := ioutil.ReadAll(resp.Body)
	bodyString := string(body)
	if resp.Status != "200 OK" || err != nil {
		return fmt.Errorf("Error while Sending/deleting http request with AS3 json :%s  %v", bodyString, err)
	}

	defer resp.Body.Close()
	d.SetId("")
	return nil
}
