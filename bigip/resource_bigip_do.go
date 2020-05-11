/*
Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
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
			"timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     20,
				Description: "DO json",
			},
			"tenant_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Deprecated:  "this attribute is no longer in use",
				Description: "unique identifier for DO resource",
			},
		},
	}
}

func resourceBigipDoCreate(d *schema.ResourceData, meta interface{}) error {
	client_bigip := meta.(*bigip.BigIP)

	do_json := d.Get("do_json").(string)
	if ok := bigip.ValidateDOTemplate(do_json); !ok {
		return fmt.Errorf("[DO] Error validating template against DO schema \n")
	}
	//	name := d.Get("tenant_name").(string)
	timeout := d.Get("timeout").(int)
	timeout_sec := timeout * 60
	log.Printf("[DEBUG]timeout_sec is :%d", timeout_sec)
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

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Error while receiving  http response with DO json:%v", err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Errorf("Error while reading http response with DO json:%v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 202 {
		return fmt.Errorf("Error while Sending/Posting http request with DO json :%s  %v", string(body), err)
	}
	respRef := make(map[string]interface{})
	json.Unmarshal(body, &respRef)
	respID := respRef["id"].(string)

	var do_success = false

	if resp.StatusCode == 200 {
		log.Printf("[DEBUG] response status is 200 ok and no aysnc flag in declaration")
		do_success = true
		d.SetId(respID)
	}

	if resp.StatusCode == http.StatusAccepted {
		for i := 0; i <= timeout_sec; i++ {
			log.Printf("[DEBUG]Value of loop counter :%d", i)
			url := client_bigip.Host + "/mgmt/shared/declarative-onboarding/task/" + respID
			req, _ := http.NewRequest("GET", url, nil)
			req.SetBasicAuth(client_bigip.User, client_bigip.Password)
			req.Header.Set("Accept", "application/json")
			req.Header.Set("Content-Type", "application/json")

			taskResp, err := client.Do(req)
			defer taskResp.Body.Close()
			if err != nil {
				log.Printf("[DEBUG]Polling the task id until the timeout")
				time.Sleep(1 * time.Second)
				continue
			}
			if taskResp.StatusCode == 200 {
				resp_body, err := ioutil.ReadAll(taskResp.Body)
				if err != nil {
					d.SetId("")
					return fmt.Errorf("Error while reading the response body :%v", err)
				}
				respRef1 := make(map[string]interface{})
				json.Unmarshal(resp_body, &respRef1)
				log.Printf("[DEBUG] Got success and setting state id")
				do_success = true
				d.SetId(respID)
				break
			} else {
				time.Sleep(1 * time.Second)
				continue
			}
		}
	}

	if do_success == false {
		log.Printf("[DEBUG] Didn't get successful response within timeout")
		url := client_bigip.Host + "/mgmt/shared/declarative-onboarding/task/" + respID
		req, _ := http.NewRequest("GET", url, nil)
		req.SetBasicAuth(client_bigip.User, client_bigip.Password)
		req.Header.Set("Accept", "application/json")
		req.Header.Set("Content-Type", "application/json")
		taskResp, err := client.Do(req)
		defer taskResp.Body.Close()
		if err != nil {
			d.SetId("")
			return fmt.Errorf("Timedout while polling the DO task id with error :%v", err)
		}
		resp_body, err := ioutil.ReadAll(taskResp.Body)
		if err != nil {
			d.SetId("")
			return fmt.Errorf("Timedout while polling the DO task id with error :%v", err)
		}
		respRef2 := make(map[string]interface{})
		json.Unmarshal(resp_body, &respRef2)
		log.Printf("[DEBUG] timeout resp_body is :%v", respRef2)
		result_map := respRef2["result"]
		d.SetId("")
		return fmt.Errorf("Timeout while polling the DO task id with result:%v", result_map)
	}

	return resourceBigipDoRead(d, meta)
}

func resourceBigipDoRead(d *schema.ResourceData, meta interface{}) error {

	client_bigip := meta.(*bigip.BigIP)
	log.Printf("[INFO] Reading Do config")
	ID := d.Id()
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	client := &http.Client{Transport: tr}
	url := client_bigip.Host + "/mgmt/shared/declarative-onboarding/task/" + ID
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("Error while creating http request for reading Do config:%v", err)
	}
	req.SetBasicAuth(client_bigip.User, client_bigip.Password)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)

	if err != nil {
		return fmt.Errorf("Error while receiving http response body in read call :%v", err)
	}
	resp_body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Error while reading http response body in read call :%v", err)
	}
	bodyString := string(resp_body)
	if resp.Status != "200 OK" {
		defer resp.Body.Close()
		return fmt.Errorf("Error while Sending/fetching http request :%s", bodyString)
	}
	respRef1 := make(map[string]interface{})
	json.Unmarshal(resp_body, &respRef1)
	log.Printf("[DEBUG] in read resp_body is :%v", respRef1)

	//dojson := make(map[string]interface{})
	dojson := respRef1["declaration"]
	out, _ := json.Marshal(dojson)
	doString := string(out)
	d.Set("do_json", doString)

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
	//      name := d.Get("tenant_name").(string)
	timeout := d.Get("timeout").(int)
	timeout_sec := timeout * 60
	log.Printf("[DEBUG]timeout_sec is :%d", timeout_sec)
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

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Error while receiving  http response with DO json:%v", err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Errorf("Error while reading http response with DO json:%v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 202 {
		return fmt.Errorf("Error while Sending/Posting http request with DO json :%s  %v", string(body), err)
	}
	respRef := make(map[string]interface{})
	json.Unmarshal(body, &respRef)
	respID := respRef["id"].(string)

	var do_success = false

	if resp.StatusCode == 200 {
		log.Printf("[DEBUG] response status is 200 ok and no aysnc flag in declaration")
		do_success = true
		d.SetId(respID)
	}

	if resp.StatusCode == http.StatusAccepted {
		for i := 0; i <= timeout_sec; i++ {
			log.Printf("[DEBUG]Value of loop counter :%d", i)
			url := client_bigip.Host + "/mgmt/shared/declarative-onboarding/task/" + respID
			req, _ := http.NewRequest("GET", url, nil)
			req.SetBasicAuth(client_bigip.User, client_bigip.Password)
			req.Header.Set("Accept", "application/json")
			req.Header.Set("Content-Type", "application/json")

			taskResp, err := client.Do(req)
			defer taskResp.Body.Close()
			if err != nil {
				log.Printf("[DEBUG]Polling the task id until the timeout")
				time.Sleep(1 * time.Second)
				continue
			}
			if taskResp.StatusCode == 200 {
				resp_body, err := ioutil.ReadAll(taskResp.Body)
				if err != nil {
					d.SetId("")
					return fmt.Errorf("Error while reading the response body :%v", err)
				}
				respRef1 := make(map[string]interface{})
				json.Unmarshal(resp_body, &respRef1)
				do_success = true
				d.SetId(respID)
				break
			} else {
				time.Sleep(1 * time.Second)
				continue
			}
		}
	}

	if do_success == false {
		log.Printf("[DEBUG] Didn't get successful response within timeout")
		url := client_bigip.Host + "/mgmt/shared/declarative-onboarding/task/" + respID
		req, _ := http.NewRequest("GET", url, nil)
		req.SetBasicAuth(client_bigip.User, client_bigip.Password)
		req.Header.Set("Accept", "application/json")
		req.Header.Set("Content-Type", "application/json")
		taskResp, err := client.Do(req)
		defer taskResp.Body.Close()
		if err != nil {
			d.SetId("")
			return fmt.Errorf("Timedout while polling the DO task id with error :%v", err)
		}
		resp_body, err := ioutil.ReadAll(taskResp.Body)
		if err != nil {
			d.SetId("")
			return fmt.Errorf("Timedout while polling the DO task id with error :%v", err)
		}
		respRef2 := make(map[string]interface{})
		json.Unmarshal(resp_body, &respRef2)
		result_map := respRef2["result"]
		d.SetId("")
		return fmt.Errorf("Timeout while polling the DO task id with result:%v", result_map)
	}

	return resourceBigipDoRead(d, meta)
}

func resourceBigipDoDelete(d *schema.ResourceData, meta interface{}) error {

	log.Println("[INFO]:Delete Operation is not supported for this resource")
	d.SetId("")
	return nil
}
