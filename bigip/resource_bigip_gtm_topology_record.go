package bigip

import (
	"context"
	"fmt"
	"log"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceBigipGtmTopologyRecord() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBigipGtmTopologyRecordCreate,
		ReadContext:   resourceBigipGtmTopologyRecordRead,
		UpdateContext: resourceBigipGtmTopologyRecordUpdate,
		DeleteContext: resourceBigipGtmTopologyRecordDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"description": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The topology record description that defines the source and destination match. Format: 'ldns: <source> server: <destination>' (e.g., 'ldns: region /Common/my-region server: datacenter /Common/my-dc')",
			},
			"order": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "The order in which the topology record is evaluated",
			},
			"score": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1,
				Description: "The weight or preference given to this topology record",
			},
		},
	}
}

func resourceBigipGtmTopologyRecordCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	description := d.Get("description").(string)

	log.Printf("[INFO] Creating GTM Topology Record: %s", description)

	config := &bigip.GTMTopologyRecord{
		Description: description,
		Order:       d.Get("order").(int),
		Score:       d.Get("score").(int),
	}

	err := client.CreateGTMTopologyRecord(config)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error creating GTM Topology Record: %v", err))
	}

	d.SetId(description)

	return resourceBigipGtmTopologyRecordRead(ctx, d, meta)
}

func resourceBigipGtmTopologyRecordRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	description := d.Id()

	log.Printf("[INFO] Reading GTM Topology Record: %s", description)

	record, err := client.GetGTMTopologyRecord(description)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error retrieving GTM Topology Record: %v", err))
	}
	if record == nil {
		log.Printf("[WARN] GTM Topology Record not found, removing from state")
		d.SetId("")
		return nil
	}

	d.Set("description", record.Description)
	d.Set("order", record.Order)
	d.Set("score", record.Score)

	return nil
}

func resourceBigipGtmTopologyRecordUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	description := d.Id()

	log.Printf("[INFO] Updating GTM Topology Record: %s", description)

	config := &bigip.GTMTopologyRecord{
		Description: description,
		Order:       d.Get("order").(int),
		Score:       d.Get("score").(int),
	}

	err := client.ModifyGTMTopologyRecord(description, config)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error updating GTM Topology Record: %v", err))
	}

	return resourceBigipGtmTopologyRecordRead(ctx, d, meta)
}

func resourceBigipGtmTopologyRecordDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	description := d.Id()

	log.Printf("[INFO] Deleting GTM Topology Record: %s", description)

	err := client.DeleteGTMTopologyRecord(description)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error deleting GTM Topology Record: %v", err))
	}

	d.SetId("")
	return nil
}
