package bigip

import (
	"context"
	"fmt"
	"log"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceBigipGtmDatacenter() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceBigipGtmDatacenterRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the GTM datacenter",
			},
			"partition": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Partition of the GTM datacenter",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Description of the datacenter",
			},
			"contact": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Contact information for the datacenter",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the datacenter is enabled",
			},
			"location": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Location of the datacenter",
			},
			"prober_fallback": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Type of prober to use for fallback",
			},
			"prober_preference": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Type of prober to prefer",
			},
		},
	}
}

func dataSourceBigipGtmDatacenterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	name := d.Get("name").(string)
	partition := d.Get("partition").(string)
	fullPath := fmt.Sprintf("/%s/%s", partition, name)

	log.Printf("[DEBUG] Reading GTM Datacenter data source: %s", fullPath)

	datacenter, err := client.GetGTMDatacenter(fullPath)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error retrieving GTM Datacenter %s: %v", fullPath, err))
	}
	if datacenter == nil {
		return diag.FromErr(fmt.Errorf("GTM Datacenter %s not found", fullPath))
	}

	d.SetId(fullPath)
	d.Set("name", datacenter.Name)
	d.Set("partition", datacenter.Partition)
	d.Set("description", datacenter.Description)
	d.Set("contact", datacenter.Contact)
	d.Set("location", datacenter.Location)
	d.Set("prober_fallback", datacenter.ProberFallback)
	d.Set("prober_preference", datacenter.ProberPreference)

	if datacenter.Disabled {
		d.Set("enabled", false)
	} else {
		d.Set("enabled", true)
	}

	return nil
}
