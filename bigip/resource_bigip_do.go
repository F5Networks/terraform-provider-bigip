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

func resourceBigipDo() *schema.Resource {
	return &schema.Resource{
		Create: resourceBigipDoCreate,
		Read:   resourceBigipDoRead,
		Update: resourceBigipDoUpdate,
		Delete: resourceBigipDoDelete,
		Exists: resourceBigipDoExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{

			"do_json": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "DO json",
			},
			"tenant_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "unique identifier for DO resource",
			},
		},
	}
}

func resourceBigipDoCreate(d *schema.ResourceData, meta interface{}) error {
	client_bigip := meta.(*bigip.BigIP)

	do_json := d.Get("do_json").(string)
	name := d.Get("tenant_name").(string)
	log.Printf("[INFO] Creating do config in bigip:%s", do_json)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	client := &http.Client{Transport: tr}
	url := client_bigip.Host + "/mgmt/shared/declarative-onboarding/"
	req, err := http.NewRequest("POST", url, strings.NewReader(do_json))
	if err != nil {
		return fmt.Errorf("Error while creating http request with DO json:%v", err)
	}
	req.SetBasicAuth(client_bigip.User, client_bigip.Password)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	resp, _ := client.Do(req)
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 202 {
		return fmt.Errorf("Error while Sending/Posting http request with DO json :%s  %v", string(body), err)
	}
	if resp.StatusCode == http.StatusAccepted {
		url := client_bigip.Host + "/mgmt/shared/declarative-onboarding"
		for {
			req, _ := http.NewRequest("GET", url, nil)
			req.SetBasicAuth(client_bigip.User, client_bigip.Password)
			req.Header.Set("Accept", "application/json")
			req.Header.Set("Content-Type", "application/json")

			taskResp, _ := client.Do(req)
			body, err := ioutil.ReadAll(taskResp.Body)
			if err != nil {
				return fmt.Errorf("Error while Sending/Posting http request with DO json :%s  %v", string(body), err)
			}
			defer taskResp.Body.Close()
			if taskResp.StatusCode == 202 {
				continue
			} else if taskResp.StatusCode == 200 {
				break
			} else {
				return fmt.Errorf("Error while Sending/Posting http request with DO json :%s  %v", string(body), err)
			}
		}
	}
	d.SetId(name)
	return resourceBigipDoRead(d, meta)
}

func resourceBigipDoRead(d *schema.ResourceData, meta interface{}) error {
	client_bigip := meta.(*bigip.BigIP)
	log.Printf("[INFO] Reading Do config")

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	client := &http.Client{Transport: tr}
	url := client_bigip.Host + "/mgmt/shared/declarative-onboarding"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("Error while creating http request for reading Do config:%v", err)
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

func resourceBigipDoExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	client_bigip := meta.(*bigip.BigIP)
	log.Printf("[INFO] Checking if Do config exists in bigip ")

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	client := &http.Client{Transport: tr}
	url := client_bigip.Host + "/mgmt/shared/declarative-onboarding"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("[ERROR] Error while creating http request for checking Do config : %v", err)
		return false, err
	}
	req.SetBasicAuth(client_bigip.User, client_bigip.Password)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	body, err := ioutil.ReadAll(resp.Body)
	bodyString := string(body)
	if resp.Status == "204 No Content" || err != nil {
		log.Printf("[ERROR] Error while checking doresource present in bigip :%s  %v", bodyString, err)
		defer resp.Body.Close()
		return false, err
	}
	defer resp.Body.Close()
	return true, nil
}

func resourceBigipDoUpdate(d *schema.ResourceData, meta interface{}) error {
	client_bigip := meta.(*bigip.BigIP)

	do_json := d.Get("do_json").(string)
	log.Printf("[INFO] Updating do config in bigip:%s", do_json)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	client := &http.Client{Transport: tr}
	url := client_bigip.Host + "/mgmt/shared/declarative-onboarding/"
	req, err := http.NewRequest("POST", url, strings.NewReader(do_json))
	if err != nil {
		return fmt.Errorf("Error while creating http request with DO json:%v", err)
	}
	req.SetBasicAuth(client_bigip.User, client_bigip.Password)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	resp, _ := client.Do(req)
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 202 {
		return fmt.Errorf("Error while Sending/Posting http request with DO json :%s  %v", string(body), err)
	}
	if resp.StatusCode == http.StatusAccepted {
		url := client_bigip.Host + "/mgmt/shared/declarative-onboarding"
		for {
			req, _ := http.NewRequest("GET", url, nil)
			req.SetBasicAuth(client_bigip.User, client_bigip.Password)
			req.Header.Set("Accept", "application/json")
			req.Header.Set("Content-Type", "application/json")

			taskResp, _ := client.Do(req)
			body, err := ioutil.ReadAll(taskResp.Body)
			if err != nil {
				return fmt.Errorf("Error while Sending/Posting http request with DO json :%s  %v", string(body), err)
			}
			defer taskResp.Body.Close()
			if taskResp.StatusCode == 202 {
				continue
			} else if taskResp.StatusCode == 200 {
				break
			} else {
				return fmt.Errorf("Error while Sending/Posting http request with DO json :%s  %v", string(body), err)
			}
		}
	}
	return resourceBigipDoRead(d, meta)
}

func resourceBigipDoDelete(d *schema.ResourceData, meta interface{}) error {

	return fmt.Errorf("Delete method on DO is not supported")
}
