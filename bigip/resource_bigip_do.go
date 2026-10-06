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
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"net"
	"net/http"
	"os"
	"strconv"
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
	if clientBigip.Token != "" {
		req.Header.Set("X-F5-Auth-Token", clientBigip.Token)
	} else {
		req.SetBasicAuth(clientBigip.User, clientBigip.Password)
	}
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
			if clientBigip.Token != "" {
				req.Header.Set("X-F5-Auth-Token", clientBigip.Token)
			} else {
				req.SetBasicAuth(clientBigip.User, clientBigip.Password)
			}
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
		if clientBigip.Token != "" {
			req.Header.Set("X-F5-Auth-Token", clientBigip.Token)
		} else {
			req.SetBasicAuth(clientBigip.User, clientBigip.Password)
		}
		req.Header.Set("Accept", "application/json")
		req.Header.Set("Content-Type", "application/json")
		taskResp, err := client.Do(req)
		if err != nil && isRecoverableMgmtPlaneError(err) {
			if waitErr := waitForDeviceManagementPlaneReady(clientBigip, 10*time.Minute); waitErr == nil {
				taskResp, err = client.Do(req)
			}
		}
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
	if err := waitForDeviceManagementPlaneReady(clientBigip, 10*time.Minute); err != nil {
		return diag.FromErr(err)
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
	if clientBigip.Token != "" {
		req.Header.Set("X-F5-Auth-Token", clientBigip.Token)
	} else {
		req.SetBasicAuth(clientBigip.User, clientBigip.Password)
	}
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
	if clientBigip.Token != "" {
		req.Header.Set("X-F5-Auth-Token", clientBigip.Token)
	} else {
		req.SetBasicAuth(clientBigip.User, clientBigip.Password)
	}
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
		start := time.Now()
	forLoop:
		for time.Since(start).Seconds() < float64(timeoutSec) {
			log.Printf("[DEBUG]Value of Timeout counter in seconds :%v", math.Ceil(time.Since(start).Seconds()))
			url := clientBigip.Host + "/mgmt/shared/declarative-onboarding/task/" + respID
			req, _ := http.NewRequest("GET", url, nil)
			if clientBigip.Token != "" {
				req.Header.Set("X-F5-Auth-Token", clientBigip.Token)
			} else {
				req.SetBasicAuth(clientBigip.User, clientBigip.Password)
			}
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
				log.Printf("StatusCode:%+v", taskResp.StatusCode)
			}
			time.Sleep(1 * time.Second)
		}
	}

	if !doSuccess {
		log.Printf("[DEBUG] Didn't get successful response within timeout")
		url := clientBigip.Host + "/mgmt/shared/declarative-onboarding/task/" + respID
		req, _ := http.NewRequest("GET", url, nil)
		if clientBigip.Token != "" {
			req.Header.Set("X-F5-Auth-Token", clientBigip.Token)
		} else {
			req.SetBasicAuth(clientBigip.User, clientBigip.Password)
		}
		req.Header.Set("Accept", "application/json")
		req.Header.Set("Content-Type", "application/json")
		taskResp, err := client.Do(req)
		if err != nil && isRecoverableMgmtPlaneError(err) {
			if waitErr := waitForDeviceManagementPlaneReady(clientBigip, 10*time.Minute); waitErr == nil {
				taskResp, err = client.Do(req)
			}
		}

		defer func() {
			if taskResp != nil && taskResp.Body != nil {
				if err := taskResp.Body.Close(); err != nil {
					log.Printf("[DEBUG] Could not close the request to %s", url)
				}
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
	if err := waitForDeviceManagementPlaneReady(clientBigip, 10*time.Minute); err != nil {
		return diag.FromErr(err)
	}

	return resourceBigipDoRead(ctx, d, meta)
}

func resourceBigipDoDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	log.Println("[INFO]:Delete Operation is not supported for this resource")
	d.SetId("")
	return nil
}

// envIntOrDefault reads name from the environment as an integer, falling
// back to def if unset or unparseable. Mirrors the provider schema's own
// api_timeout/api_retries EnvDefaultFunc behavior (see Provider() in
// provider.go) for connectBigIP's standalone reconnect path below, which
// bypasses the provider's schema-driven Configure lifecycle entirely and so
// has no other way to pick up those same settings.
func envIntOrDefault(name string, def int) int {
	v := os.Getenv(name)
	if v == "" {
		return def
	}
	parsed, err := strconv.Atoi(v)
	if err != nil {
		return def
	}
	return parsed
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
		// Explicit rather than left nil: go-bigip's NewSession only
		// substitutes its own package-level default ConfigOptions (10
		// retries, 10s apart, on every retriable failure including
		// transport-level errors like a refused/timed-out connection)
		// when ConfigOptions is nil -- a genuinely unreachable
		// bigip_address could otherwise hang here, since this
		// reconnect path (unlike the main provider client) isn't
		// driven by the api_timeout/api_retries schema fields.
		// TokenTimeout must still be set to the same 1200s the
		// package default uses even though this path doesn't read
		// token_timeout from anywhere: NewTokenSession compares the
		// device's returned token timeout against
		// ConfigOptions.TokenTimeout and issues an extra PUT to
		// correct it if they differ, so leaving this at its Go zero
		// value would make every bigiq_token_auth reconnect issue a
		// spurious (and, for a device without that PUT permitted,
		// failing) token-timeout adjustment.
		ConfigOptions: &bigip.ConfigOptions{
			APICallTimeout: time.Duration(envIntOrDefault("API_TIMEOUT", 60)) * time.Second,
			APICallRetries: envIntOrDefault("API_RETRIES", 10),
			TokenTimeout:   time.Duration(envIntOrDefault("TOKEN_TIMEOUT", 1200)) * time.Second,
		},
	}

	if d.Get("bigip_token_auth").(bool) {
		bigipConfig.LoginReference = d.Get("bigiq_login_ref").(string)
	}

	return Client(&bigipConfig)
}

func isRecoverableMgmtPlaneError(err error) bool {
	if err == nil {
		return false
	}
	var netErr net.Error
	return errors.As(err, &netErr) || strings.Contains(strings.ToLower(err.Error()), "connection refused")
}

func waitForDeviceManagementPlaneReady(client *bigip.BigIP, timeout time.Duration) error {
	const pollInterval = 5 * time.Second

	var lastErr error
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if err := client.ValidateConnection(); err == nil {
			if freshLoginSucceeds(client) {
				return nil
			}
			lastErr = fmt.Errorf("validate connection succeeded but fresh login is still unavailable")
		} else {
			lastErr = err
		}
		log.Printf("[DEBUG] waiting for device management-plane recovery after DO: %v", lastErr)
		time.Sleep(pollInterval)
	}

	return fmt.Errorf("timed out after %s waiting for device management-plane recovery after DO: %w", timeout, lastErr)
}
