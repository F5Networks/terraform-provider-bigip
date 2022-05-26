package bigip

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceBigipWafPb() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBigipWafPbRead,
		Schema: map[string]*schema.Schema{
			"policy_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the WAF policy",
			},
			"partition": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Partition where the WAF policy is deployed",
			},
			"minimum_learning_score": {
				Type:         schema.TypeInt,
				Required:     true,
				Description:  "The minimum learning for suggestions",
				ValidateFunc: validation.IntBetween(1, 100),
			},
			"policy_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "System generated id of the WAF policy",
			},
			"json": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The return payload of the queried PB suggestions",
			},
		},
	}
}

type ExportPb struct {
	PolicyReference map[string]string
	Inline          bool
	Filter          string
}

func dataSourceBigipWafPbRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	d.SetId("")

	policyName := d.Get("policy_name").(string)
	partition := d.Get("partition").(string)
	score := d.Get("minimum_learning_score").(int)
	policyId, err := client.GetWafPolicyId(policyName, partition)
	if err != nil {
		return fmt.Errorf("error retrieving policy %s on partition %s", policyName, partition)
	}
	policyLink := fmt.Sprintf("https://localhost/mgmt/tm/asm/policies/%s", policyId)
	payload := ExportPb{
		PolicyReference: map[string]string{"file": policyLink},
		Inline:          true,
		Filter:          fmt.Sprintf("score gt %d", score),
	}
	export, err := client.PostPbExport(payload)
	if err != nil {
		return fmt.Errorf("error exporting pb suggestions: %v", err)
	}

	task, err := client.GetWafPbExportResult(export.Task_id)
	if err != nil {
		return err
	}
	log.Printf("[DEBUG]Initial response export status %v", task.Status)
	for task.Status != "COMPLETED" && task.Status != "FAILURE" {
		pbtask, err := client.GetWafPbExportResult(export.Task_id)
		if err != nil {
			return err
		}
		task = pbtask
		if task.Status == "FAILURE" || task.Status == "COMPLETED" {
			break
		}
		time.Sleep(3 * time.Second)
	}

	if task.Status == "FAILURE" {
		return fmt.Errorf("export task failed")
	}
	if task.Status == "COMPLETED" {
		pbJson, err := json.Marshal(task.Result)
		if err != nil {
			return err
		}
		_ = d.Set("policy_id", policyId)
		_ = d.Set("json", string(pbJson))
		d.SetId(policyName)
	}
	return nil
}
