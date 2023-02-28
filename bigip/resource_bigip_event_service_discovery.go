/*
Original work from https://github.com/DealerDotCom/terraform-provider-bigip
Modifications Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file,You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceServiceDiscovery() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceServiceDiscoveryCreate,
		ReadContext:   resourceServiceDiscoveryRead,
		UpdateContext: resourceServiceDiscoveryUpdate,
		DeleteContext: resourceServiceDiscoveryDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{

			"taskid": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the partition/tenant",
				ForceNew:    true,
			},
			"node": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "name of node",
						},
						"ip": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "ip of nonde",
						},
						"port": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "port",
						},
					},
				},
			},
		},
	}
}

func resourceServiceDiscoveryCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	taskid := d.Get("taskid").(string)
	log.Printf("[INFO]: taskid: %+v", taskid)
	var nodeList []interface{}
	if m, ok := d.GetOk("node"); ok {
		for _, node := range m.(*schema.Set).List() {
			log.Printf("[INFO]: node Value: %+v", node)
			nodeList = append(nodeList, node)
		}
	}
	log.Printf("[INFO]: node Value: %+v", nodeList)
	err := client.AddServiceDiscoveryNodes(taskid, nodeList)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error modifying node %s: %v", nodeList, err))
	}
	d.SetId(taskid)
	return resourceServiceDiscoveryRead(ctx, d, meta)

}

func resourceServiceDiscoveryRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	taskid := d.Id()

	serviceDiscoveryResp, err := client.GetServiceDiscoveryNodes(taskid)
	log.Printf("[DEBUG] serviceDiscoveryResp is :%v", serviceDiscoveryResp)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error Reading node : %v", err))
	}
	nodeList1 := serviceDiscoveryResp.(map[string]interface{})["result"].(map[string]interface{})["providerOptions"].(map[string]interface{})["nodeList"]
	log.Printf("[DEBUG] nodeList1 is :%v", nodeList1)

	if serviceDiscoveryResp == nil {
		d.SetId("")
		return diag.FromErr(fmt.Errorf("[DEBUG]serviceDiscoveryResp is : %s", serviceDiscoveryResp))
	}
	if err := d.Set("node", nodeList1); err != nil {
		return diag.FromErr(fmt.Errorf("error updating nodelist in state: %v", err))
	}
	return nil
}

func resourceServiceDiscoveryUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	taskid := d.Id()
	log.Printf("[INFO]: taskid: %+v", taskid)
	var nodeList []interface{}
	if m, ok := d.GetOk("node"); ok {
		for _, node := range m.(*schema.Set).List() {
			log.Printf("[INFO]: node Value: %+v", node)
			nodeList = append(nodeList, node)
		}
	}
	err := client.AddServiceDiscoveryNodes(taskid, nodeList)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error modifying node %s: %v", nodeList, err))
	}
	return resourceServiceDiscoveryRead(ctx, d, meta)
}

func resourceServiceDiscoveryDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	clientBigip := meta.(*bigip.BigIP)
	taskid := d.Id()
	url := clientBigip.Host + "/mgmt/shared/service-discovery/task/" + taskid + "/nodes/"
	payload := strings.NewReader("[ ]\n")
	log.Printf("[DEBUG] url Complete :%v", url)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	client := &http.Client{Transport: tr}
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error while creating http request for Delete operation:%+v ", err))
	}
	req.SetBasicAuth(clientBigip.User, clientBigip.Password)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return diag.FromErr(err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("[DEBUG] Could not close the request to %s", url)
		}
	}()

	var body bytes.Buffer
	_, err = io.Copy(&body, resp.Body)
	// body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return diag.FromErr(err)
	}

	bodyString := body.String()
	if resp.Status != "200 OK" {
		return diag.FromErr(fmt.Errorf("error while Sending/Posting http request for Delete operation :%s  %v", bodyString, err))
	}

	d.SetId("")
	return nil
}
