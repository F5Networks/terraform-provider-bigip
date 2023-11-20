package bigip

import (
	"context"
	"log"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceBigipPartition() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBigipPartitionCreate,
		ReadContext:   resourceBigipPartitionRead,
		UpdateContext: resourceBigipPartitionUpdate,
		DeleteContext: resourceBigipPartitionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the partition",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the partition",
			},
			"route_domain_id": {
				Type:        schema.TypeInt,
				Default:     0,
				Optional:    true,
				Description: "The route domain of the partition",
			},
		},
	}
}

func resourceBigipPartitionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*bigip.BigIP)
	name := d.Get("name").(string)
	routeDomain := d.Get("route_domain_id").(int)
	description := d.Get("description").(string)

	partition := &bigip.Partition{
		Name:        name,
		RouteDomain: routeDomain,
	}

	err := client.CreatePartition(partition)

	if err != nil {
		log.Printf("[ERROR] error while creating the partition: %s", name)
		return diag.FromErr(err)
	}

	if description != "" {
		descBody := make(map[string]string)
		descBody["description"] = description

		err := client.ModifyFolderDescription(name, descBody)

		if err != nil {
			log.Printf("[ERROR] error while updating the description of partition: %s", name)
			return diag.FromErr(err)
		}
	}

	d.SetId(name)
	return resourceBigipPartitionRead(ctx, d, m)
}

func resourceBigipPartitionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*bigip.BigIP)
	name := d.Id()
	partition, err := client.GetPartition(name)

	if err != nil {
		log.Printf("[ERROR] error while reading the partition: %s", name)
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Partition: %+v", partition)
	d.Set("name", partition.Name)
	d.Set("route_domain_id", partition.RouteDomain)

	return nil
}

func resourceBigipPartitionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*bigip.BigIP)
	name := d.Id()
	routeDomain := d.Get("route_domain_id").(int)

	oldRD, newRD := d.GetChange("route_domain_id")
	oldDesc, newDesc := d.GetChange("description")

	if oldRD.(int) != newRD.(int) {

		partition := &bigip.Partition{
			RouteDomain: routeDomain,
		}
		err := client.ModifyPartition(name, partition)

		if err != nil {
			log.Printf("[ERROR] error while updating the partition: %s", name)
			return diag.FromErr(err)
		}
	}
	if oldDesc.(string) != newDesc.(string) {
		descBody := make(map[string]string)
		descBody["description"] = newDesc.(string)

		err := client.ModifyFolderDescription(name, descBody)

		if err != nil {
			log.Printf("[ERROR] error while updating the description of partition: %s", name)
			return diag.FromErr(err)
		}
	}

	return resourceBigipPartitionRead(ctx, d, m)
}

func resourceBigipPartitionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	name := d.Id()
	client := m.(*bigip.BigIP)
	err := client.DeletePartition(name)
	if err != nil {
		log.Printf("[ERROR] error while deleting the partition: %s", name)
		return diag.FromErr(err)
	}
	return nil
}
