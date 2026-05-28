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

func resourceBigipGtmTopologyRegion() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBigipGtmTopologyRegionCreate,
		ReadContext:   resourceBigipGtmTopologyRegionRead,
		UpdateContext: resourceBigipGtmTopologyRegionUpdate,
		DeleteContext: resourceBigipGtmTopologyRegionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the GTM topology region",
			},
			"partition": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "Common",
				ForceNew:    true,
				Description: "Partition of the GTM topology region",
			},
			"members": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Region member entry. Examples: 'subnet 10.0.0.0/8', 'country US', 'state US/California', 'datacenter /Common/my-dc', 'isp Comcast', 'region /Common/other-region', 'continent NA', 'pool /Common/my-pool'",
						},
					},
				},
				Description: "The members that define this topology region (subnets, countries, states, datacenters, ISPs, etc.)",
			},
		},
	}
}

func resourceBigipGtmTopologyRegionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	name := d.Get("name").(string)
	partition := d.Get("partition").(string)
	fullPath := fmt.Sprintf("/%s/%s", partition, name)

	log.Printf("[INFO] Creating GTM Topology Region: %s", fullPath)

	config := &bigip.GTMRegion{
		Name:      name,
		Partition: partition,
	}

	if v, ok := d.GetOk("members"); ok {
		membersSet := v.(*schema.Set)
		members := make([]bigip.GTMRegionMember, 0, membersSet.Len())
		for _, item := range membersSet.List() {
			memberMap := item.(map[string]interface{})
			members = append(members, bigip.GTMRegionMember{
				Name: memberMap["name"].(string),
			})
		}
		config.Members = members
	}

	err := client.CreateGTMRegion(config)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error creating GTM Topology Region %s: %v", fullPath, err))
	}

	d.SetId(fullPath)

	return resourceBigipGtmTopologyRegionRead(ctx, d, meta)
}

func resourceBigipGtmTopologyRegionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	fullPath := d.Id()

	log.Printf("[INFO] Reading GTM Topology Region: %s", fullPath)

	region, err := client.GetGTMRegion(fullPath)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error retrieving GTM Topology Region %s: %v", fullPath, err))
	}
	if region == nil {
		log.Printf("[WARN] GTM Topology Region %s not found, removing from state", fullPath)
		d.SetId("")
		return nil
	}

	// Parse partition and name from fullPath
	parts := strings.Split(strings.TrimPrefix(fullPath, "/"), "/")
	if len(parts) >= 2 {
		d.Set("partition", parts[0])
		d.Set("name", parts[1])
	} else {
		d.Set("name", region.Name)
		if region.Partition != "" {
			d.Set("partition", region.Partition)
		}
	}

	if len(region.Members) > 0 {
		members := make([]interface{}, 0, len(region.Members))
		for _, member := range region.Members {
			members = append(members, map[string]interface{}{
				"name": member.Name,
			})
		}
		d.Set("members", members)
	}

	return nil
}

func resourceBigipGtmTopologyRegionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	fullPath := d.Id()

	log.Printf("[INFO] Updating GTM Topology Region: %s", fullPath)

	config := &bigip.GTMRegion{
		Name:      d.Get("name").(string),
		Partition: d.Get("partition").(string),
	}

	if v, ok := d.GetOk("members"); ok {
		membersSet := v.(*schema.Set)
		members := make([]bigip.GTMRegionMember, 0, membersSet.Len())
		for _, item := range membersSet.List() {
			memberMap := item.(map[string]interface{})
			members = append(members, bigip.GTMRegionMember{
				Name: memberMap["name"].(string),
			})
		}
		config.Members = members
	}

	err := client.ModifyGTMRegion(fullPath, config)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error updating GTM Topology Region %s: %v", fullPath, err))
	}

	return resourceBigipGtmTopologyRegionRead(ctx, d, meta)
}

func resourceBigipGtmTopologyRegionDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	fullPath := d.Id()

	log.Printf("[INFO] Deleting GTM Topology Region: %s", fullPath)

	err := client.DeleteGTMRegion(fullPath)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error deleting GTM Topology Region %s: %v", fullPath, err))
	}

	d.SetId("")
	return nil
}
