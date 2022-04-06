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
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/f5devcentral/go-bigip/f5teem"
	uuid "github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/structure"
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
				StateFunc: func(v interface{}) string {
					json, _ := structure.NormalizeJsonString(v)
					return json
				},
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
			"bigip_address": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "IP Address of BIGIP host to be used for this resource",
			},
			"bigip_user": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "UserName of BIGIP host to be used for this resource",
			},
			"bigip_port": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Port number of BIGIP host to be used for this resource",
			},
			"bigip_password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "Password of  BIGIP host to be used for this resource",
			},
			"bigip_token_auth": {
				Type:        schema.TypeBool,
				Optional:    true,
				Sensitive:   true,
				Description: "Enable to use an external authentication source (LDAP, TACACS, etc)",
				Default:     false,
			},
		},
	}
}

func resourceBigipDoCreate(d *schema.ResourceData, meta interface{}) error {
	clientBigip := meta.(*bigip.BigIP)
	if d.Get("bigip_address").(string) != "" && d.Get("bigip_user").(string) != "" && d.Get("bigip_password").(string) != "" && d.Get("bigip_port").(string) != "" {
		clientBigip2, err := connectBigIP(d)
		if err != nil {
			log.Printf("Connection to BIGIP Failed with :%v", err)
			return err
		}
		clientBigip = clientBigip2
	}
	doJson := d.Get("do_json").(string)
	if !clientBigip.Teem {
		id := uuid.New()
		uniqueID := id.String()
		assetInfo := f5teem.AssetInfo{
			Name:    "Terraform-provider-bigip",
			Version: clientBigip.UserAgent,
			Id:      uniqueID,
		}
		teemDevice := f5teem.AnonymousClient(assetInfo, "")
		f := map[string]interface{}{
			"Terraform Version": clientBigip.UserAgent,
		}
		err := teemDevice.Report(f, "bigip_do", "1")
		if err != nil {
			log.Printf("[ERROR]Sending Telemetry data failed:%v", err)
		}
	}

	timeout := d.Get("timeout").(int)
	timeoutSec := timeout * 60
	log.Printf("[DEBUG]timeout_sec is :%d", timeoutSec)
	log.Printf("[INFO] Creating do config in bigip:%s", doJson)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	client := &http.Client{Transport: tr}
	url := clientBigip.Host + "/mgmt/shared/declarative-onboarding/"
	req, err := http.NewRequest("POST", url, strings.NewReader(doJson))
	if err != nil {
		return fmt.Errorf("Error while creating http request with DO json:%v ", err)
	}
	req.SetBasicAuth(clientBigip.User, clientBigip.Password)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)

	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("[DEBUG] Could not close the request to %s", url)
		}
	}()

	if err != nil {
		return fmt.Errorf("Error while receiving  http response with DO json:%v", err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Error while reading http response with DO json:%v ", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode > 202 {
		return fmt.Errorf("Error while Sending/Posting http request with DO json :%s  %v", string(body), err)
	}
	respRef := make(map[string]interface{})
	if err := json.Unmarshal(body, &respRef); err != nil {
		return err
	}
	respID := respRef["id"].(string)

	var doSuccess = false

	if resp.StatusCode == 200 {
		log.Printf("[DEBUG] response status is 200 ok and no aysnc flag in declaration")
		doSuccess = true
		d.SetId(respID)
	}

	if resp.StatusCode == http.StatusAccepted {
	forLoop:
		for i := 0; i <= timeoutSec; i++ {
			log.Printf("[DEBUG]Value of Timeout counter in seconds :%d", i)
			url := clientBigip.Host + "/mgmt/shared/declarative-onboarding/task/" + respID
			req, _ := http.NewRequest("GET", url, nil)
			req.SetBasicAuth(clientBigip.User, clientBigip.Password)
			req.Header.Set("Accept", "application/json")
			req.Header.Set("Content-Type", "application/json")

			taskResp, err := client.Do(req)
			if taskResp == nil {
				log.Printf("[DEBUG]taskResp of DO is empty,but continue the loop until timeout \n")
				time.Sleep(1 * time.Second)
				continue
			}
			defer taskResp.Body.Close()
			if err != nil {
				log.Printf("[DEBUG]Polling the task id until the timeout")
				time.Sleep(1 * time.Second)
				continue
			}
			switch {
			case taskResp.StatusCode == 200:
				respBody, err := ioutil.ReadAll(taskResp.Body)
				if err != nil {
					d.SetId("")
					return fmt.Errorf("Error while reading the response body :%v ", err)
				}
				respRef1 := make(map[string]interface{})
				if err := json.Unmarshal(respBody, &respRef1); err != nil {
					return err
				}
				log.Printf("[DEBUG] Got success and setting state id")
				doSuccess = true
				d.SetId(respID)
				break forLoop
			case taskResp.StatusCode == 202:
				respBody, err := ioutil.ReadAll(taskResp.Body)
				if err != nil {
					d.SetId("")
					return fmt.Errorf("Error while reading the response body :%v ", err)
				}
				respRef1 := make(map[string]interface{})
				if err := json.Unmarshal(respBody, &respRef1); err != nil {
					return err
				}
				resultMap := respRef1["result"]
				if resultMap.(map[string]interface{})["status"] != "RUNNING" {
					return fmt.Errorf("Error while reading the response body :%v ", resultMap)
				}
			default:
				log.Printf("StatusCode:%+v", taskResp.StatusCode)
				time.Sleep(1 * time.Second)
				continue
			}
		}
	}

	if !doSuccess {
		log.Printf("[DEBUG] Didn't get successful response within timeout")
		url := clientBigip.Host + "/mgmt/shared/declarative-onboarding/task/" + respID
		req, _ := http.NewRequest("GET", url, nil)
		req.SetBasicAuth(clientBigip.User, clientBigip.Password)
		req.Header.Set("Accept", "application/json")
		req.Header.Set("Content-Type", "application/json")
		taskResp, err := client.Do(req)
		if taskResp == nil {
			d.SetId("")
			return fmt.Errorf("Timedout while polling the DO task id with error :%v", err)
		}
		defer taskResp.Body.Close()
		if err != nil {
			d.SetId("")
			return fmt.Errorf("Timedout while polling the DO task id with error :%v", err)
		}
		respBody, err := ioutil.ReadAll(taskResp.Body)
		if err != nil {
			d.SetId("")
			return fmt.Errorf("Timedout while polling the DO task id with error :%v", err)
		}
		respRef2 := make(map[string]interface{})
		if err := json.Unmarshal(respBody, &respRef2); err != nil {
			return err
		}
		log.Printf("[DEBUG] timeout resp_body is :%v", respRef2)
		resultMap := respRef2["result"]
		d.SetId("")
		return fmt.Errorf("Timeout while polling the DO task id with result:%v", resultMap)
	}

	return resourceBigipDoRead(d, meta)
}

func resourceBigipDoRead(d *schema.ResourceData, meta interface{}) error {
	clientBigip := meta.(*bigip.BigIP)
	if d.Get("bigip_address").(string) != "" && d.Get("bigip_user").(string) != "" && d.Get("bigip_password").(string) != "" && d.Get("bigip_port").(string) != "" {
		clientBigip2, err := connectBigIP(d)
		if err != nil {
			log.Printf("Connection to BIGIP Failed with :%v", err)
			return err
		}
		clientBigip = clientBigip2
	}
	log.Printf("[INFO] Reading Do config")
	ID := d.Id()
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	client := &http.Client{Transport: tr}
	url := clientBigip.Host + "/mgmt/shared/declarative-onboarding/task/" + ID
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("Error while creating http request for reading Do config:%v ", err)
	}
	req.SetBasicAuth(clientBigip.User, clientBigip.Password)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)

	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("[DEBUG] Could not close the request to %s", url)
		}
	}()

	if err != nil {
		return fmt.Errorf("Error while receiving http response body in read call :%v ", err)
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Error while reading http response body in read call :%v ", err)
	}
	bodyString := string(respBody)
	if resp.Status != "200 OK" {
		return fmt.Errorf("Error while Sending/fetching http request :%s ", bodyString)
	}

	respRef1 := make(map[string]interface{})
	if err := json.Unmarshal(respBody, &respRef1); err != nil {
		return err
	}
	log.Printf("[DEBUG] in read resp_body is :%v", respRef1)

	return nil

}

func resourceBigipDoExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	clientBigip := meta.(*bigip.BigIP)
	if d.Get("bigip_address").(string) != "" && d.Get("bigip_user").(string) != "" && d.Get("bigip_password").(string) != "" && d.Get("bigip_port").(string) != "" {
		clientBigip2, err := connectBigIP(d)
		if err != nil {
			log.Printf("Connection to BIGIP Failed with :%v", err)
			return false, err
		}
		clientBigip = clientBigip2
	}
	log.Printf("[INFO] Checking if Do config exists in bigip ")

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	client := &http.Client{Transport: tr}
	url := clientBigip.Host + "/mgmt/shared/declarative-onboarding"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("[ERROR] Error while creating http request for checking Do config : %v", err)
		return false, err
	}
	req.SetBasicAuth(clientBigip.User, clientBigip.Password)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("[DEBUG] Could not close the request to %s", url)
		}
	}()

	body, err := ioutil.ReadAll(resp.Body)
	bodyString := string(body)
	if resp.Status == "204 No Content" || err != nil {
		log.Printf("[ERROR] Error while checking doresource present in bigip :%s  %v", bodyString, err)
		return false, err
	}

	return true, nil
}

func resourceBigipDoUpdate(d *schema.ResourceData, meta interface{}) error {
	clientBigip := meta.(*bigip.BigIP)
	if d.Get("bigip_address").(string) != "" && d.Get("bigip_user").(string) != "" && d.Get("bigip_password").(string) != "" && d.Get("bigip_port").(string) != "" {
		clientBigip2, err := connectBigIP(d)
		if err != nil {
			log.Printf("Connection to BIGIP Failed with :%v", err)
			return err
		}
		clientBigip = clientBigip2
	}

	doJson := d.Get("do_json").(string)
	timeout := d.Get("timeout").(int)
	timeoutSec := timeout * 60
	log.Printf("[DEBUG]timeout_sec is :%d", timeoutSec)
	log.Printf("[INFO] Updating do config in bigip:%s", doJson)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	client := &http.Client{Transport: tr}
	url := clientBigip.Host + "/mgmt/shared/declarative-onboarding/"
	req, err := http.NewRequest("POST", url, strings.NewReader(doJson))
	if err != nil {
		return fmt.Errorf("Error while creating http request with DO json:%v", err)
	}
	req.SetBasicAuth(clientBigip.User, clientBigip.Password)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)

	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("[DEBUG] Could not close the request to %s", url)
		}
	}()

	if err != nil {
		return fmt.Errorf("Error while receiving  http response with DO json:%v", err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Error while reading http response with DO json:%v", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode > 202 {
		return fmt.Errorf("Error while Sending/Posting http request with DO json :%s  %v", string(body), err)
	}
	respRef := make(map[string]interface{})
	if err := json.Unmarshal(body, &respRef); err != nil {
		return err
	}
	respID := respRef["id"].(string)

	var doSuccess = false

	if resp.StatusCode == 200 {
		log.Printf("[DEBUG] response status is 200 ok and no aysnc flag in declaration")
		doSuccess = true
		d.SetId(respID)
	}

	if resp.StatusCode == http.StatusAccepted {
	forLoop:
		for i := 0; i <= timeoutSec; i++ {
			log.Printf("[DEBUG]Value of loop counter :%d", i)
			url := clientBigip.Host + "/mgmt/shared/declarative-onboarding/task/" + respID
			req, _ := http.NewRequest("GET", url, nil)
			req.SetBasicAuth(clientBigip.User, clientBigip.Password)
			req.Header.Set("Accept", "application/json")
			req.Header.Set("Content-Type", "application/json")

			taskResp, err := client.Do(req)

			defer func() {
				if err := taskResp.Body.Close(); err != nil {
					log.Printf("[DEBUG] Could not close the request to %s", url)
				}
			}()

			if err != nil {
				log.Printf("[DEBUG]Polling the task id until the timeout")
				time.Sleep(1 * time.Second)
				continue
			}
			switch {
			case taskResp.StatusCode == 200:
				respBody, err := ioutil.ReadAll(taskResp.Body)
				if err != nil {
					d.SetId("")
					return fmt.Errorf("Error while reading the response body :%v ", err)
				}
				respRef1 := make(map[string]interface{})
				if err := json.Unmarshal(respBody, &respRef1); err != nil {
					return err
				}
				doSuccess = true
				d.SetId(respID)
				break forLoop
			case taskResp.StatusCode == 202:
				respBody, err := ioutil.ReadAll(taskResp.Body)
				if err != nil {
					d.SetId("")
					return fmt.Errorf("Error while reading the response body :%v ", err)
				}
				respRef1 := make(map[string]interface{})
				if err := json.Unmarshal(respBody, &respRef1); err != nil {
					return err
				}
				resultMap := respRef1["result"]
				if resultMap.(map[string]interface{})["status"] != "RUNNING" {
					return fmt.Errorf("Error while reading the response body :%v ", resultMap)
				}
			default:
				time.Sleep(1 * time.Second)
				continue
			}
		}
	}

	if !doSuccess {
		log.Printf("[DEBUG] Didn't get successful response within timeout")
		url := clientBigip.Host + "/mgmt/shared/declarative-onboarding/task/" + respID
		req, _ := http.NewRequest("GET", url, nil)
		req.SetBasicAuth(clientBigip.User, clientBigip.Password)
		req.Header.Set("Accept", "application/json")
		req.Header.Set("Content-Type", "application/json")
		taskResp, err := client.Do(req)

		defer func() {
			if err := taskResp.Body.Close(); err != nil {
				log.Printf("[DEBUG] Could not close the request to %s", url)
			}
		}()

		if err != nil {
			d.SetId("")
			return fmt.Errorf("Timedout while polling the DO task id with error :%v ", err)
		}
		respBody, err := ioutil.ReadAll(taskResp.Body)
		if err != nil {
			d.SetId("")
			return fmt.Errorf("Timedout while polling the DO task id with error :%v ", err)
		}
		respRef2 := make(map[string]interface{})
		if err := json.Unmarshal(respBody, &respRef2); err != nil {
			return err
		}

		resultMap := respRef2["result"]
		d.SetId("")
		return fmt.Errorf("Timeout while polling the DO task id with result:%v", resultMap)
	}

	return resourceBigipDoRead(d, meta)
}

func resourceBigipDoDelete(d *schema.ResourceData, meta interface{}) error {

	log.Println("[INFO]:Delete Operation is not supported for this resource")
	d.SetId("")
	return nil
}

func connectBigIP(d *schema.ResourceData) (*bigip.BigIP, error) {
	bigipConfig := Config{
		Address:  d.Get("bigip_address").(string),
		Port:     d.Get("bigip_port").(string),
		Username: d.Get("bigip_user").(string),
		Password: d.Get("bigip_password").(string),
	}

	if d.Get("bigip_token_auth").(bool) {
		bigipConfig.LoginReference = d.Get("bigiq_login_ref").(string)
	}

	return bigipConfig.Client()
}
