/*
Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/f5devcentral/go-bigip/f5teem"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceBigipCommand() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBigipCommandCreate,
		ReadContext:   resourceBigipCommandRead,
		UpdateContext: resourceBigipCommandUpdate,
		DeleteContext: resourceBigipCommandDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"when": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "apply",
			},
			"commands": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "The commands to send to the remote BIG-IP device over the configured provider",
			},
			"command_result": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Partition of ssl certificate",
			},
		},
	}
}

func resourceBigipCommandCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	var commandList []string
	if d.Get("when").(string) == "apply" {
		if m, ok := d.GetOk("commands"); ok {
			for _, cmd := range m.([]interface{}) {
				commandList = append(commandList, fmt.Sprintf("-c 'tmsh %s'", cmd.(string)))
			}
		}
		log.Printf("[INFO] Running TMSH Command : %v ", commandList)
		var resultList []string
		for _, str := range commandList {
			log.Printf("[INFO] Command to run:%v", str)
			commandConfig := &bigip.BigipCommand{
				Command:     "run",
				UtilCmdArgs: str,
			}
			resultCmd, err := client.RunCommand(commandConfig)
			if err != nil {
				return diag.FromErr(fmt.Errorf("error retrieving Command Result: %v", err))
			}
			resultList = append(resultList, resultCmd.CommandResult)
		}
		_ = d.Set("command_result", resultList)
	}
	d.SetId(d.Get("when").(string))
	if !client.Teem {
		id := uuid.New()
		uniqueID := id.String()
		assetInfo := f5teem.AssetInfo{
			Name:    "Terraform-provider-bigip",
			Version: client.UserAgent,
			Id:      uniqueID,
		}
		apiKey := os.Getenv("TEEM_API_KEY")
		teemDevice := f5teem.AnonymousClient(assetInfo, apiKey)
		f := map[string]interface{}{
			"Terraform Version": client.UserAgent,
		}
		tsVer := strings.Split(client.UserAgent, "/")
		err := teemDevice.Report(f, "bigip_command", tsVer[3])
		if err != nil {
			log.Printf("[ERROR]Sending Telemetry data failed:%v", err)
		}
	}
	return nil
}

func resourceBigipCommandRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Println("[INFO]:Read Operation is not supported for this resource")
	return nil
}
func resourceBigipCommandUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	var commandList []string
	if d.Get("when").(string) == "apply" {
		if m, ok := d.GetOk("commands"); ok {
			for _, cmd := range m.([]interface{}) {
				commandList = append(commandList, fmt.Sprintf("-c 'tmsh %s'", cmd.(string)))
			}
		}
		log.Printf("[INFO] Running TMSH Command : %v ", commandList)
		var resultList []string
		for _, str := range commandList {
			commandConfig := &bigip.BigipCommand{
				Command:     "run",
				UtilCmdArgs: str,
			}
			resultCmd, err := client.RunCommand(commandConfig)
			if err != nil {
				return diag.FromErr(fmt.Errorf("error retrieving Command Result: %v", err))
			}
			resultList = append(resultList, resultCmd.CommandResult)
		}
		_ = d.Set("command_result", resultList)
	}
	return nil
}

func resourceBigipCommandDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	var commandList []string
	if d.Get("when").(string) == "destroy" {
		if m, ok := d.GetOk("commands"); ok {
			for _, cmd := range m.([]interface{}) {
				commandList = append(commandList, fmt.Sprintf("-c 'tmsh %s'", cmd.(string)))
			}
		}
		log.Printf("[INFO] Running Delete TMSH Command: %v ", commandList)

		for _, str := range commandList {
			log.Printf("[INFO] Command to run:%v", str)
			commandConfig := &bigip.BigipCommand{
				Command:     "run",
				UtilCmdArgs: str,
			}
			log.Printf("[INFO] Command struct:%+v", commandConfig)
			resultCmd, err := client.RunCommand(commandConfig)
			if err != nil {
				return diag.FromErr(fmt.Errorf("error retrieving Command Result: %v", err))
			}
			log.Printf("[INFO] Result Command struct:%+v", resultCmd)

		}
	}

	d.SetId("")
	return nil
}
