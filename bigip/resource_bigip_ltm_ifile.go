package bigip

import (
	"context"
	"fmt"
	"log"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceBigipLtmIfile() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBigipLtmIfileCreate,
		ReadContext:   resourceBigipLtmIfileRead,
		UpdateContext: resourceBigipLtmIfileUpdate,
		DeleteContext: resourceBigipLtmIfileDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the LTM iFile.",
			},
			"partition": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "Common",
				ForceNew:    true,
				Description: "Partition for the LTM iFile.",
			},
			"sub_path": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Subdirectory within the partition.",
			},
			"file_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The system iFile name to reference (e.g., /Common/my-sys-ifile).",
			},
			"full_path": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Full path of the LTM iFile.",
			},
		},
	}
}

func resourceBigipLtmIfileCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	name := d.Get("name").(string)
	partition := d.Get("partition").(string)
	subPath := d.Get("sub_path").(string)
	fileName := d.Get("file_name").(string)

	ltmIfile := &bigip.LtmIFile{
		Name:      name,
		Partition: partition,
		SubPath:   subPath,
		FileName:  fileName,
	}

	// fullPath := fmt.Sprintf("/%s/%s", partition, name)
	fullPath := buildIFileFullPath(partition, subPath, name)
	log.Printf("[INFO] Creating LTM iFile: %+v", fullPath)

	err := client.CreateLtmIFile(ltmIfile)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error creating LTM iFile: %v", err))
	}

	d.SetId(fullPath)
	return resourceBigipLtmIfileRead(ctx, d, meta)
}

func resourceBigipLtmIfileRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	fullPath := d.Id()

	log.Printf("[DEBUG] Reading LTM iFile: %s", fullPath)

	ltmIfile, err := client.GetLtmIFile(fullPath)

	if err != nil {
		return diag.FromErr(fmt.Errorf("error reading LTM iFile: %v", err))
	}
	if ltmIfile == nil {
		log.Printf("[DEBUG] LTM iFile (%s) not found, removing from state", fullPath)
		d.SetId("")
		return nil
	}

	log.Printf("[resourceBigipLtmIfileRead][INFO] LTM iFile found: %s", ltmIfile.FullPath)

	_ = d.Set("name", ltmIfile.Name)
	_ = d.Set("partition", ltmIfile.Partition)
	_ = d.Set("file_name", ltmIfile.FileName)
	_ = d.Set("full_path", ltmIfile.FullPath)

	return nil
}

func resourceBigipLtmIfileUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	fullPath := d.Id()
	name := d.Get("name").(string)
	partition := d.Get("partition").(string)
	subPath := d.Get("sub_path").(string)
	fileName := d.Get("file_name").(string)

	ltmIfile := &bigip.LtmIFile{
		Name:      name,
		Partition: partition,
		SubPath:   subPath,
		FileName:  fileName,
		FullPath:  fullPath,
	}

	err := client.UpdateLtmIFile(ltmIfile)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error updating LTM iFile: %v", err))
	}

	return resourceBigipLtmIfileRead(ctx, d, meta)
}

func resourceBigipLtmIfileDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	fullPath := d.Id()

	log.Printf("[INFO] Deleting LTM iFile: %+v", fullPath)

	err := client.DeleteLtmIFile(fullPath)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error deleting LTM iFile: %v", err))
	}

	d.SetId("")
	return nil
}
