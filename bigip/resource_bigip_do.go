/*
Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"strings"
	"time"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/f5devcentral/go-bigip/f5teem"
	uuid "github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/structure"
)

func resourceBigipDo() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBigipDoCreate,
		ReadContext:   resourceBigipDoRead,
		UpdateContext: resourceBigipDoUpdate,
		DeleteContext: resourceBigipDoDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"do_json": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "DO json",
				StateFunc: func(v interface{}) string {
					jsonString, _ := structure.NormalizeJsonString(v)
					return jsonString
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

func resourceBigipDoCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	clientBigip := meta.(*bigip.BigIP)

	if d.Get("bigip_address").(string) != "" && d.Get("bigip_user").(string) != "" && d.Get("bigip_password").(string) != "" || d.Get("bigip_port").(string) != "" {
		clientBigip2, err := connectBigIP(d)
		if err != nil {
			log.Printf("Connection to BIGIP Failed with :%v", err)
			return diag.FromErr(err)
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
		return diag.FromErr(fmt.Errorf("error while creating http request with DO json:%v", err))
	}
	req.SetBasicAuth(clientBigip.User, clientBigip.Password)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	log.Printf("[INFO] URL:%s", clientBigip.Host)

	resp, err := client.Do(req)

	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("[DEBUG] Could not close the request to %s", url)
		}
	}()

	if err != nil {
		return diag.FromErr(fmt.Errorf("error while receiving  http response with DO json:%v", err))
	}
	// body, err := os.ReadAll(resp.Body)
	var body bytes.Buffer
	_, err = io.Copy(&body, resp.Body)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error while reading http response with DO json:%v", err))
	}

	if resp.StatusCode < 200 || resp.StatusCode > 202 {
		return diag.FromErr(fmt.Errorf("Error while Sending/Posting http request with DO json :%s  %v", body.String(), err))
	}
	respRef := make(map[string]interface{})
	if err := json.Unmarshal(body.Bytes(), &respRef); err != nil {
		return diag.FromErr(err)
	}
	respID := respRef["id"].(string)

	var doSuccess = false

	if resp.StatusCode == 200 {
		log.Printf("[DEBUG] response status is 200 ok and no aysnc flag in declaration")
		doSuccess = true
		d.SetId(respID)
	}

	if resp.StatusCode == http.StatusAccepted {
		start := time.Now()
	forLoop:
		for time.Since(start).Seconds() < float64(timeoutSec) {
			log.Printf("[DEBUG]Value of Timeout counter in seconds :%v", math.Ceil(time.Since(start).Seconds()))
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
				var body bytes.Buffer
				_, err = io.Copy(&body, taskResp.Body)
				// respBody, err := ioutil.ReadAll(taskResp.Body)
				if err != nil {
					d.SetId("")
					return diag.FromErr(fmt.Errorf("error while reading the response body :%v", err))
				}
				respRef1 := make(map[string]interface{})
				if err := json.Unmarshal(body.Bytes(), &respRef1); err != nil {
					return diag.FromErr(err)
				}
				log.Printf("[DEBUG] Got success and setting state id")
				doSuccess = true
				d.SetId(respID)
				break forLoop
			case taskResp.StatusCode == 202:
				var respBody bytes.Buffer
				_, err = io.Copy(&respBody, taskResp.Body)
				// respBody, err := ioutil.ReadAll(taskResp.Body)
				if err != nil {
					d.SetId("")
					return diag.FromErr(fmt.Errorf("error while reading the response body :%v", err))
				}
				respRef1 := make(map[string]interface{})
				if err := json.Unmarshal(respBody.Bytes(), &respRef1); err != nil {
					return diag.FromErr(err)
				}
				resultMap := respRef1["result"]
				if resultMap.(map[string]interface{})["status"] != "RUNNING" {
					return diag.FromErr(fmt.Errorf("error while reading the response body :%v", resultMap))
				}
			default:
				log.Printf("StatusCode:%+v", taskResp.StatusCode)
			}
			time.Sleep(1 * time.Second)
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
			return diag.FromErr(fmt.Errorf("timedout while polling the DO task id with error :%v", err))
		}
		defer taskResp.Body.Close()
		if err != nil {
			d.SetId("")
			return diag.FromErr(fmt.Errorf("timedout while polling the DO task id with error :%v", err))
		}
		var respBody bytes.Buffer
		_, err = io.Copy(&respBody, taskResp.Body)
		if err != nil {
			d.SetId("")
			return diag.FromErr(fmt.Errorf("timedout while polling the DO task id with error :%v", err))
		}
		respRef2 := make(map[string]interface{})
		if err := json.Unmarshal(respBody.Bytes(), &respRef2); err != nil {
			return diag.FromErr(err)
		}
		log.Printf("[DEBUG] timeout resp_body is :%v", respRef2)
		resultMap := respRef2["result"]
		d.SetId("")
		return diag.FromErr(fmt.Errorf("timeout while polling the DO task id with result:%v", resultMap))
	}

	return resourceBigipDoRead(ctx, d, meta)
}

func resourceBigipDoRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	clientBigip := meta.(*bigip.BigIP)
	if d.Get("bigip_address").(string) != "" && d.Get("bigip_user").(string) != "" && d.Get("bigip_password").(string) != "" || d.Get("bigip_port").(string) != "" {
		clientBigip2, err := connectBigIP(d)
		if err != nil {
			log.Printf("Connection to BIGIP Failed with :%v", err)
			return diag.FromErr(err)
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
		return diag.FromErr(fmt.Errorf("error while creating http request for reading Do config:%v", err))
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
		return diag.FromErr(fmt.Errorf("error while receiving http response body in read call :%v ", err))
	}
	var respBody bytes.Buffer
	_, err = io.Copy(&respBody, resp.Body)
	// respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error while reading http response body in read call :%v ", err))
	}
	bodyString := respBody.String()
	if resp.Status != "200 OK" {
		return diag.FromErr(fmt.Errorf("error while Sending/fetching http request :%s ", bodyString))
	}
	respRef1 := make(map[string]interface{})
	if err := json.Unmarshal(respBody.Bytes(), &respRef1); err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[INFO] in read resp_body is :%v", respRef1)
	byteData, _ := json.Marshal(respRef1["declaration"])
	_ = d.Set("do_json", string(byteData))

	return nil

}

func resourceBigipDoUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	clientBigip := meta.(*bigip.BigIP)
	if d.Get("bigip_address").(string) != "" && d.Get("bigip_user").(string) != "" && d.Get("bigip_password").(string) != "" || d.Get("bigip_port").(string) != "" {
		clientBigip2, err := connectBigIP(d)
		if err != nil {
			log.Printf("Connection to BIGIP Failed with :%v", err)
			return diag.FromErr(err)
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
		return diag.FromErr(fmt.Errorf("error while creating http request with DO json:%v ", err))
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
		return diag.FromErr(fmt.Errorf("error while receiving  http response with DO json:%v", err))
	}
	var body bytes.Buffer
	_, err = io.Copy(&body, resp.Body)
	// body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error while reading http response with DO json:%v ", err))
	}

	if resp.StatusCode < 200 || resp.StatusCode > 202 {
		return diag.FromErr(fmt.Errorf("error while Sending/Posting http request with DO json :%s  %v", body.String(), err))
	}
	respRef := make(map[string]interface{})
	if err := json.Unmarshal(body.Bytes(), &respRef); err != nil {
		return diag.FromErr(err)
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
				var respBody bytes.Buffer
				_, err = io.Copy(&respBody, taskResp.Body)
				// respBody, err := ioutil.ReadAll(taskResp.Body)
				if err != nil {
					d.SetId("")
					return diag.FromErr(fmt.Errorf("error while reading the response body :%v", err))
				}
				respRef1 := make(map[string]interface{})
				if err := json.Unmarshal(respBody.Bytes(), &respRef1); err != nil {
					return diag.FromErr(err)
				}
				doSuccess = true
				d.SetId(respID)
				break forLoop
			case taskResp.StatusCode == 202:
				var respBody bytes.Buffer
				_, err = io.Copy(&respBody, taskResp.Body)
				// respBody, err := ioutil.ReadAll(taskResp.Body)
				if err != nil {
					d.SetId("")
					return diag.FromErr(fmt.Errorf("error while reading the response body :%v", err))
				}
				respRef1 := make(map[string]interface{})
				if err := json.Unmarshal(respBody.Bytes(), &respRef1); err != nil {
					return diag.FromErr(err)
				}
				resultMap := respRef1["result"]
				if resultMap.(map[string]interface{})["status"] != "RUNNING" {
					return diag.FromErr(fmt.Errorf("error while reading the response body :%v", resultMap))
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
			return diag.FromErr(fmt.Errorf("Timedout while polling the DO task id with error :%v ", err))
		}
		var respBody bytes.Buffer
		_, err = io.Copy(&respBody, taskResp.Body)
		// respBody, err := ioutil.ReadAll(taskResp.Body)
		if err != nil {
			d.SetId("")
			return diag.FromErr(fmt.Errorf("Timedout while polling the DO task id with error :%v ", err))
		}
		respRef2 := make(map[string]interface{})
		if err := json.Unmarshal(respBody.Bytes(), &respRef2); err != nil {
			return diag.FromErr(err)
		}

		resultMap := respRef2["result"]
		d.SetId("")
		return diag.FromErr(fmt.Errorf("timeout while polling the DO task id with result:%v", resultMap))
	}

	return resourceBigipDoRead(ctx, d, meta)
}

func resourceBigipDoDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	log.Println("[INFO]:Delete Operation is not supported for this resource")
	d.SetId("")
	return nil
}

func connectBigIP(d *schema.ResourceData) (*bigip.BigIP, error) {
	var portVal string
	if _, ok := d.GetOk("bigip_port"); ok {
		portVal = d.Get("bigip_port").(string)
	} else {
		portVal = "443"
	}
	bigipConfig := bigip.Config{
		Address:           d.Get("bigip_address").(string),
		Port:              portVal,
		Username:          d.Get("bigip_user").(string),
		Password:          d.Get("bigip_password").(string),
		CertVerifyDisable: true,
	}

	if d.Get("bigip_token_auth").(bool) {
		bigipConfig.LoginReference = d.Get("bigiq_login_ref").(string)
	}

	return Client(&bigipConfig)
}
