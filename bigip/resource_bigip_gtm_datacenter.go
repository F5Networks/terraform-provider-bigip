package bigip

import (
	"context"
	"fmt"
	"log"
	"strings"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceBigipGtmDatacenter() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBigipGtmDatacenterCreate,
		ReadContext:   resourceBigipGtmDatacenterRead,
		UpdateContext: resourceBigipGtmDatacenterUpdate,
		DeleteContext: resourceBigipGtmDatacenterDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the GTM datacenter",
			},
			"partition": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "Common",
				ForceNew:    true,
				Description: "Partition of the GTM datacenter",
			},
			"contact": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Contact information for the datacenter",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the datacenter",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Enable or disable the datacenter",
			},
			"location": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Location of the datacenter",
			},
			"prober_fallback": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "any-available",
				Description: "Type of prober to use for fallback",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					validOptions := []string{"any-available", "inside-datacenter", "outside-datacenter", "inherit", "pool"}
					for _, opt := range validOptions {
						if v == opt {
							return
						}
					}
					errs = append(errs, fmt.Errorf("%q must be one of %v, got: %s", key, validOptions, v))
					return
				},
			},
			"prober_preference": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "inside-datacenter",
				Description: "Type of prober to prefer",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					validOptions := []string{"inside-datacenter", "outside-datacenter", "inherit", "pool"}
					for _, opt := range validOptions {
						if v == opt {
							return
						}
					}
					errs = append(errs, fmt.Errorf("%q must be one of %v, got: %s", key, validOptions, v))
					return
				},
			},
		},
	}
}

func resourceBigipGtmDatacenterCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	name := d.Get("name").(string)
	partition := d.Get("partition").(string)

	log.Printf("[INFO] Creating GTM Datacenter: %s in partition %s", name, partition)

	datacenter := &bigip.GTMDatacenter{
		Name:      name,
		Partition: partition,
	}

	if v, ok := d.GetOk("contact"); ok {
		datacenter.Contact = v.(string)
	}
	if v, ok := d.GetOk("description"); ok {
		datacenter.Description = v.(string)
	}
	if v, ok := d.GetOk("enabled"); ok {
		enabled := v.(bool)
		datacenter.Enabled = enabled
		datacenter.Disabled = !enabled
	}
	if v, ok := d.GetOk("location"); ok {
		datacenter.Location = v.(string)
	}
	if v, ok := d.GetOk("prober_fallback"); ok {
		datacenter.ProberFallback = v.(string)
	}
	if v, ok := d.GetOk("prober_preference"); ok {
		datacenter.ProberPreference = v.(string)
	}

	err := client.CreateGTMDatacenter(datacenter)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error creating GTM Datacenter %s: %v", name, err))
	}

	fullPath := fmt.Sprintf("/%s/%s", partition, name)
	d.SetId(fullPath)

	return resourceBigipGtmDatacenterUpdate(ctx, d, meta)
}

func resourceBigipGtmDatacenterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	fullPath := d.Id()
	log.Printf("[INFO] Reading GTM Datacenter: %s", fullPath)

	datacenter, err := client.GetGTMDatacenter(fullPath)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error retrieving GTM Datacenter %s: %v", fullPath, err))
	}
	if datacenter == nil {
		log.Printf("[WARN] GTM Datacenter %s not found, removing from state", fullPath)
		d.SetId("")
		return nil
	}

	// Parse partition and name from fullPath
	parts := strings.Split(strings.TrimPrefix(fullPath, "/"), "/")
	if len(parts) >= 2 {
		d.Set("partition", parts[0])
		d.Set("name", parts[1])
	} else {
		d.Set("name", datacenter.Name)
		if datacenter.Partition != "" {
			d.Set("partition", datacenter.Partition)
		}
	}

	// Only set non-empty string fields to preserve user configuration
	// Skip setting contact, description, and location if datacenter is disabled
	if !datacenter.Disabled {
		if _, ok := d.GetOk("contact"); ok || datacenter.Contact != "" {
			d.Set("contact", datacenter.Contact)
		}
		if _, ok := d.GetOk("description"); ok || datacenter.Description != "" {
			d.Set("description", datacenter.Description)
		}
		if _, ok := d.GetOk("location"); ok || datacenter.Location != "" {
			d.Set("location", datacenter.Location)
		}
	}

	if _, ok := d.GetOk("prober_fallback"); ok || datacenter.ProberFallback != "" {
		d.Set("prober_fallback", datacenter.ProberFallback)
	}
	if _, ok := d.GetOk("prober_preference"); ok || datacenter.ProberPreference != "" {
		d.Set("prober_preference", datacenter.ProberPreference)
	}

	// Always set boolean fields
	if datacenter.Enabled {
		d.Set("enabled", true)
	}
	if datacenter.Disabled {
		d.Set("enabled", false)
	}

	return nil
}

func resourceBigipGtmDatacenterUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	fullPath := d.Id()
	log.Printf("[INFO] Updating GTM Datacenter: %s", fullPath)

	datacenter := &bigip.GTMDatacenter{}

	if d.HasChange("contact") {
		datacenter.Contact = d.Get("contact").(string)
	}
	if d.HasChange("description") {
		datacenter.Description = d.Get("description").(string)
	}
	if d.HasChange("enabled") {
		enabled := d.Get("enabled").(bool)
		datacenter.Enabled = enabled
		datacenter.Disabled = !enabled
	}
	if d.HasChange("location") {
		datacenter.Location = d.Get("location").(string)
	}
	if d.HasChange("prober_fallback") {
		datacenter.ProberFallback = d.Get("prober_fallback").(string)
	}
	if d.HasChange("prober_preference") {
		datacenter.ProberPreference = d.Get("prober_preference").(string)
	}

	err := client.ModifyGTMDatacenter(fullPath, datacenter)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error updating GTM Datacenter %s: %v", fullPath, err))
	}

	return resourceBigipGtmDatacenterRead(ctx, d, meta)
}

func resourceBigipGtmDatacenterDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	fullPath := d.Id()
	log.Printf("[INFO] Deleting GTM Datacenter: %s", fullPath)

	err := client.DeleteGTMDatacenter(fullPath)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error deleting GTM Datacenter %s: %v", fullPath, err))
	}

	d.SetId("")
	return nil
}
