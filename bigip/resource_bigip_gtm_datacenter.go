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
			"prober_pool": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "GTM prober pool to use when prober_preference or prober_fallback is \"pool\". Required in that case; BIG-IP rejects prober_preference/prober_fallback = \"pool\" without a valid prober_pool reference.",
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
	// enabled has a schema Default (true) and is a bool, so it is never
	// truly "unset" -- GetOk cannot distinguish an explicit false from an
	// absent value (both read back as the Go zero value), so it must be
	// read directly via Get rather than GetOk. Using GetOk here meant
	// enabled = false in config was silently dropped (the whole if-ok
	// block was skipped), leaving datacenter.Enabled/Disabled at their own
	// zero values and causing BIG-IP to fall back to its default of
	// enabled.
	enabled := d.Get("enabled").(bool)
	datacenter.Enabled = enabled
	datacenter.Disabled = !enabled
	if v, ok := d.GetOk("location"); ok {
		datacenter.Location = v.(string)
	}
	if v, ok := d.GetOk("prober_fallback"); ok {
		datacenter.ProberFallback = v.(string)
	}
	if v, ok := d.GetOk("prober_preference"); ok {
		datacenter.ProberPreference = v.(string)
	}
	if v, ok := d.GetOk("prober_pool"); ok {
		datacenter.ProberPool = v.(string)
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

	// Set fields unconditionally from the device response so state always
	// reflects the current device configuration, matching the pattern used
	// in resource_bigip_gtm_server.go. There is no GTM API semantic where
	// disabling a datacenter causes contact/description/location to be
	// omitted from the response, so gating these on datacenter.Disabled (as
	// this code previously did) could leave stale values in state if those
	// fields changed on the device while the datacenter was disabled.
	d.Set("contact", datacenter.Contact)
	d.Set("description", datacenter.Description)
	d.Set("location", datacenter.Location)
	d.Set("prober_fallback", datacenter.ProberFallback)
	d.Set("prober_preference", datacenter.ProberPreference)
	d.Set("prober_pool", datacenter.ProberPool)
	d.Set("enabled", datacenter.Enabled)

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
	// enabled/disabled must always be sent on every Update, not just when
	// HasChange("enabled") is true: BIG-IP's GTM datacenter PATCH endpoint
	// resets disabled back to false (i.e. enabled) whenever the request
	// body doesn't include it, rather than leaving the device's existing
	// value alone the way a PATCH normally would for an omitted field.
	// Create always calls Update immediately afterward to apply the
	// remaining properties (see the "Now update with all the additional
	// properties" comment pattern shared across the other GTM resources
	// in this package); HasChange("enabled") is false at that point since
	// d already reflects the fully-planned config with no prior value to
	// diff against, so gating on it here previously meant Create's
	// immediate follow-up Update silently re-enabled every
	// enabled = false datacenter right after Create had correctly
	// disabled it.
	enabled := d.Get("enabled").(bool)
	datacenter.Enabled = enabled
	datacenter.Disabled = !enabled
	if d.HasChange("location") {
		datacenter.Location = d.Get("location").(string)
	}
	if d.HasChange("prober_fallback") {
		datacenter.ProberFallback = d.Get("prober_fallback").(string)
	}
	if d.HasChange("prober_preference") {
		datacenter.ProberPreference = d.Get("prober_preference").(string)
	}
	if d.HasChange("prober_pool") {
		datacenter.ProberPool = d.Get("prober_pool").(string)
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
